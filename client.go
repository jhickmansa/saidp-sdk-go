package saidp_sdk_go

import (
	"net/http"
	"errors"
	"fmt"
	"bytes"
	"encoding/hex"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"time"
	"strconv"
	"net/url"
	"crypto/tls"
)

/*
**********************************************************************
*   @author jhickman@secureauth.com
*
*  Copyright (c) 2016, SecureAuth
*  All rights reserved.
*
*    Redistribution and use in source and binary forms, with or without modification,
*    are permitted provided that the following conditions are met:
*
*    1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
*
*    2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer
*    in the documentation and/or other materials provided with the distribution.
*
*    3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived
*    from this software without specific prior written permission.
*
*    THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO,
*    THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR
*    CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
*    PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
*    LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
*    EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
**********************************************************************
*/

var (
	list = []string{ http.MethodGet, http.MethodPost, http.MethodPut }
	hdrContentTypeKey   	= http.CanonicalHeaderKey("Content-Type")
	hdrAcceptKey        	= http.CanonicalHeaderKey("Accept")
	hdrDateKey	    	= http.CanonicalHeaderKey("Date")
	hdrAuthorizationKey 	= http.CanonicalHeaderKey("Authorization")
	jsonContentType		= "application/json; charset=utf-8"
)

// Summary:
//	Client struct to hold configuration information for connecting to the SecureAuth APIs.

type Client struct {
	AppId			string
	AppKey			string
	Host			string
	Port			int
	Realm			string
	SSL			bool
	BypassCertValidation	bool
}

// Summary:
//	Function supporting the building of get requests for each service package. Will handle signing and creation of the auth header as well as timestamp and other headers needed.
// Parameters:
//	[Required] c: passing in the client containing authorization and host information.
//	[Required] endpoint: the api endpoint (after the SecureAuth# realm) that the get will be performed against.
// Returns:
//	http.Request: Http Request struct that can be used via Http Client to make the request.
//	Error: If an error is encountered, response will be nil and the error must be handled.

func (c Client) BuildGetRequest(endpoint string)(*http.Request, error) {
	url, err := buildEndpointUrl(c.SSL, c.Host, c.Port, c.Realm, endpoint)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	sig, err := c.Sign(http.MethodGet, endpoint, "")
	if err != nil {
		return nil, err
	}
	req.Header.Set(hdrAcceptKey, jsonContentType)
	req.Header.Set(hdrContentTypeKey, jsonContentType)
	req.Header.Set(hdrDateKey, getGMTTimestamp())
	req.Header.Set(hdrAuthorizationKey, sig)
	return req, nil
}

// Summary:
//	Function supporting the building of post requests for each service package. Will handle signing and creation of the auth header as well as timestamp and other headers needed.
// Parameters:
//	[Required] c: passing in the client containing authorization and host information.
//	[Required] endpoint: the api endpoint (after the SecureAuth# realm) that the get will be performed against.
//	[Required] content: the json content to be posted to the api endpoint.
// Returns:
//	http.Request: Http Request struct that can be used via Http Client to make the request.
//	Error: If an error is encountered, response will be nil and the error must be handled.

func (c Client) BuildPostRequest(endpoint string, content string)(*http.Request, error) {
	url, err := buildEndpointUrl(c.SSL, c.Host, c.Port, c.Realm, endpoint)
	if err != nil {
		return nil, err
	}
	jsonContent := []byte(content)
	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewBuffer(jsonContent))
	if err != nil {
		return nil, err
	}
	sig, err := c.Sign(http.MethodPost, endpoint, content)
	if err != nil {
		return nil, err
	}
	req.Header.Set(hdrAcceptKey, jsonContentType)
	req.Header.Set(hdrContentTypeKey, jsonContentType)
	req.Header.Set(hdrDateKey, getGMTTimestamp())
	req.Header.Set(hdrAuthorizationKey, sig)
	return req, nil
}

// Summary:
//	Function supporting the building of put requests for each service package. Will handle signing and creation of the auth header as well as timestamp and other headers needed.
// Parameters:
//	[Required] c: passing in the client containing authorization and host information.
//	[Required] endpoint: the api endpoint (after the SecureAuth# realm) that the get will be performed against.
//	[Required] content: the json content to be put to the api endpoint.
// Returns:
//	http.Request: Http Request struct that can be used via Http Client to make the request.
//	Error: If an error is encountered, response will be nil and the error must be handled.

func (c Client) BuildPutRequest(endpoint string, content string)(*http.Request, error) {
	url, err := buildEndpointUrl(c.SSL, c.Host, c.Port, c.Realm, endpoint)
	if err != nil {
		return nil, err
	}
	jsonContent := []byte(content)
	req, err := http.NewRequest(http.MethodPut, url.String(), bytes.NewBuffer(jsonContent))
	if err != nil {
		return nil, err
	}
	sig, err := c.Sign(http.MethodPut, endpoint, content)
	if err != nil {
		return nil, err
	}
	req.Header.Set(hdrAcceptKey, jsonContentType)
	req.Header.Set(hdrContentTypeKey, jsonContentType)
	req.Header.Set(hdrDateKey, getGMTTimestamp())
	req.Header.Set(hdrAuthorizationKey, sig)
	return req, nil
}

// Summary:
//	Function to execute a Http Request.
// Parameters:
//	[Required] c: passing in the client containing authorization and host information.
//	[Required] req: http.Request struct to execute.
// Returns:
//	http.Response: Http Response struct that can be used to get the body for the api response.
//	Error: If an error is encountered, response will be nil and the error must be handled.

func (c Client) Do(req *http.Request)(*http.Response, error) {
	httpClient := new(http.Client)
	transport := new(http.Transport)
	if c.BypassCertValidation {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		httpClient.Transport = transport
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Summary:
//	Function to create the Authorization headed needed to perform API calls to SecureAuth.
// Parameters:
//	[Required] c: passing in the client containing authorization and host information.
//	[Required] method: the http verb of the method being used (GET, POST, PUT)
//	[Required] endpoint: the api endpoint (after the SecureAuth# realm) that the get will be performed against.
//	content: the json content to be put to the api endpoint.
// Returns:
//	string: Authorization header string.
//	Error: If an error is encountered, response will be nil and the error must be handled.

func (c Client) Sign(method string, endpoint string, content string)(string, error) {
	if !validateMethod(method) {
		return "", errors.New("Method invalid. Try using: POST, GET, or PUT")
	}
	if len(endpoint) <= 0 {
		return "", errors.New("A API endpoint is required.")
	}
	if strings.Contains(method, "PUT") || strings.Contains(method, "POST") {
		if len(content) <= 0 {
			return "", errors.New("PUT and POST methods require content.")
		}
	}
	var buffer bytes.Buffer
	timestamp := getGMTTimestamp()
	payload := buildAuthPayload(method, timestamp, c.AppId, c.Realm, endpoint, content)
	encryptStr := makeHmac(payload, c.AppKey)
	buffer.WriteString(c.AppId)
	buffer.WriteString(":")
	buffer.WriteString(string(encryptStr))
	authStr := base64.StdEncoding.EncodeToString(buffer.Bytes())
	buffer.Reset()
	buffer.WriteString("Basic ")
	buffer.WriteString(authStr)
	return buffer.String(), nil
}

// Summary:
//	Helper function to create a Client struct.
// Parameters:
//	[Required] appId: SecureAuth API AppId.
//	[Required] appKey: SecureAuth API AppKey.
//	[Required] host: the host name (fully qualified/dns route-able) of the SecureAuth server.
//	[Required] port: the port that SecureAuth's web service is running on.
//	[Required] realm: the SecureAuth realm that will be serving the APIs.
//	[Required] ssl: if the SecureAuth realm/web service is running over ssl, set to true.
//	[Required] bypassCert: bypass certificate validation.
// Returns:
//	Client: a pointer to a client struct with the supplied values.
//	Error: If an error is encountered, response will be nil and the error must be handled

func NewClient(appId string, appKey string, host string, port int, realm string, ssl bool, bypassCert bool) (*Client, error) {
	params := []string{appId,appKey,host,realm}
	for _, v := range params {
		if isNil(v) {
			return nil, errors.New(fmt.Sprintf("%v is required to create a new client.", v))
		}
	}
	c := new(Client)
	c.AppId = appId
	c.AppKey = appKey
	c.Host = host
	if port == 0 {
		c.Port = 443
	} else {
		c.Port = port
	}
	c.Realm = realm
	c.SSL = ssl
	c.BypassCertValidation = bypassCert

	return c, nil
}

// Summary:
//	non-exportable helper to nil check a string.

func isNil(s string) (bool) {
	if len(s) <= 0 {
		return true
	} else {
		return false
	}
}

// Summary:
//	non-exportable helper to build the GMT timestamp used in authorization and http headers.

func getGMTTimestamp() string {
	location, _ := time.LoadLocation("Etc/GMT")
	return time.Now().In(location).Format(time.RFC1123)
}

// Summary:
//	non-exportable helper to build the auth payload for building the authorization header.

func buildAuthPayload (method string, timestamp string, appid string, realm string, endpoint string, content string) string {
	var buffer bytes.Buffer
	switch method {
	case "GET" :
		buffer.WriteString(method)
		buffer.WriteString("\n")
		buffer.WriteString(timestamp)
		buffer.WriteString("\n")
		buffer.WriteString(appid)
		buffer.WriteString("\n")
		buffer.WriteString("/")
		buffer.WriteString(realm)
		buffer.WriteString(endpoint)
	case "POST", "PUT" :
		buffer.WriteString(method)
		buffer.WriteString("\n")
		buffer.WriteString(timestamp)
		buffer.WriteString("\n")
		buffer.WriteString(appid)
		buffer.WriteString("\n")
		buffer.WriteString("/")
		buffer.WriteString(realm)
		buffer.WriteString(endpoint)
		buffer.WriteString("\n")
		buffer.WriteString(content)
	}
	return buffer.String()
}

// Summary:
//	non-exportable helper to build the full url to the api endpoint.

func buildEndpointUrl (ssl bool,host string, port int, realm string, endpoint string) (*url.URL, error) {
	var buffer bytes.Buffer
	if ssl {
		buffer.WriteString("https://")
	} else {
		buffer.WriteString("http://")
	}
	buffer.WriteString(host)
	buffer.WriteString(":")
	buffer.WriteString(strconv.Itoa(port))
	buffer.WriteString("/")
	buffer.WriteString(realm)
	buffer.WriteString(endpoint)
	return url.Parse(buffer.String())
}

// Summary:
//	non-exportable helper to do SHA256 HMAC

func makeHmac (data string, key string) (string) {
	byteKey,_ := hex.DecodeString(key)
	byteData := []byte(data)
	sig := hmac.New(sha256.New, byteKey)
	sig.Write([]byte(byteData))
	return base64.StdEncoding.EncodeToString(sig.Sum(nil))
}

// Summary:
//	non-exportable helper to validate expected http method verbs.

func validateMethod (str string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}