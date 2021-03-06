package ipeval

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/h2non/gock"
	sa "github.com/secureauthcorp/saidp-sdk-go"
)

/*
**********************************************************************
*   @author jhickman@secureauth.com
*
*  Copyright (c) 2017, SecureAuth
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

const (
	uAppID  = "12345"
	uAppKey = "12345"
	uHost   = "idp.host.com"
	uRealm  = "secureauth1"
	uPort   = 443
	uUser   = "user"
	uUserIP = "192.168.0.1"
)

func TestIpEval_Unit(t *testing.T) {
	defer gock.Off()

	client, err := sa.NewClient(uAppID, uAppKey, uHost, uPort, uRealm, true, false)
	if err != nil {
		t.Error(err)
	}

	n := time.Now()
	headers := map[string]string{
		"X-SA-DATE":      n.String(),
		"X-SA-SIGNATURE": makeResponseSignature(client, generateResponse(), n.String()),
	}
	// Set up a test responder for the api.
	gock.New("https://idp.host.com:443").
		Post("/secureauth1/api/v1/ipeval").
		Reply(200).
		BodyString(generateResponse()).
		SetHeaders(headers)
	ipevalRequest := new(Request)
	ipevalResponse, err := ipevalRequest.EvaluateIP(client, uUser, uUserIP)
	if err != nil {
		t.Error(err)
	}
	valid, err := ipevalResponse.IsSignatureValid(client)
	if err != nil {
		t.Error(err)
	}
	if !valid {
		t.Error("Response signature is invalid")
	}
}

func generateResponse() string {
	var jsonEval = `{"ip_evaluation":{"method":"aggregation","ip":"78.43.156.114","risk_factor":100,"risk_color":"red","risk_desc":"Extreme risk involved","geoloc":{"country":"germany","country_code":"de","region":"baden-wuerttemberg","region_code":"","city":"stuttgart","latitude":"48.7754","longtitude":"9.1818","internet_service_provider":"kabel baden-wuerttemburg gmbh & co. kg","organization":"kabel bw gmbh"},"factoring":{"latitude":48.7754,"longitude":9.1818,"threatType":100,"threatCategory":0},"factor_description":{"geoContinent":"europe","geoCountry":"germany","geoCountryCode":"de","geoCountryCF":"99","geoRegion":"","geoState":"baden-wuerttemberg","geoStateCode":"","geoStateCF":"88","geoCity":"stuttgart","geoCityCF":"77","geoPostalCode":"70173","geoAreaCode":"","geoTimeZone":"1","geoLatitude":"48.7754","geoLongitude":"9.1818","dma":"","msa":"","connectionType":"cable","lineSpeed":"medium","ipRoutingType":"fixed","geoAsn":"29562","sld":"kabel-badenwuerttemberg","tld":"de","organization":"kabel baden-wuerttemburg gmbh & co. kg","carrier":"kabel bw gmbh","anonymizer_status":"","proxyLevel":"","proxyType":"","proxyLastDetected":"","hostingFacility":"false","threatType":"Anonymous Proxy","threatCategory":"Anonymous Proxy"}},"status":"verified","message":""}`
	response := new(Response)
	if err := json.Unmarshal([]byte(jsonEval), &response); err != nil {
		return ""
	}
	bytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	return string(bytes)
}

func makeResponseSignature(c *sa.Client, response string, timeStamp string) string {
	var buffer bytes.Buffer
	buffer.WriteString(timeStamp)
	buffer.WriteString("\n")
	buffer.WriteString(c.AppID)
	buffer.WriteString("\n")
	buffer.WriteString(response)
	raw := buffer.String()
	byteKey, _ := hex.DecodeString(c.AppKey)
	byteData := []byte(raw)
	sig := hmac.New(sha256.New, byteKey)
	sig.Write([]byte(byteData))
	return base64.StdEncoding.EncodeToString(sig.Sum(nil))
}
