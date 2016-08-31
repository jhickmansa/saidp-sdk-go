package ipeval

import (
	"net/http"
	sa "github.com/SecureAuthCorp/saidp-sdk-go"
	"encoding/json"
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

const endpoint = "/api/v1/ipeval"

// Summary:
//	Response struct that will be populated after the post request.

type Response struct {
	IpEvaluation	IpEvaluation	`json:"ip_evaluation,omitempty"`
	HttpResponse	*http.Response  `json:"-,omitempty"`
}

// Summary:
//	Request struct to build the required post parameters.
// Fields:
//	[Required] UserId: the username that you want to evaluate.
//	[Required] EvalType: currently, only 'risk' is supported. Sets the eval type.
//	[Required] IpAddress: the IP Address of the user to be evaluated.

type Request struct {
	UserId		string		`json:"user_id"`
	EvalType	string		`json:"type"`
	IpAddress	string		`json:"ip_address"`
}

// Summary:
//	Struct providing data from the post request.

type IpEvaluation struct {
	Method		string		`json:"method,omitempty"`
	Ip		string		`json:"ip,omitempty"`
	RiskFactor	float32		`json:"risk_factor,omitempty"`
	RiskColor	string		`json:"risk_color,omitempty"`
	RiskDesc	string		`json:"risk_desc,omitempty"`
	GeoLoc		GeoLoc		`json:"geoloc,omitempty"`
	Factoring	Factoring	`json:"factoring,omitempty"`
	Status		string		`json:"status,omitempty"`
	Message		string		`json:"message,omitempty"`
}

// Summary:
//	Struct providing data from the post request.

type GeoLoc struct {
	Country		string		`json:"country,omitempty"`
	CountryCode	string		`json:"country_code,omitempty"`
	Region		string		`json:"region,omitempty"`
	RegionCode	string		`json:"region_code,omitempty"`
	City		string		`json:"city,omitempty"`
	Latitude	string		`json:"latitude,omitempty"`
	Longtitude	string		`json:"longtitude,omitempty"`
	Isp		string		`json:"internet_service_provider,omitempty"`
	Organization	string		`json:"organization,omitempty"`
}

// Summary:
//	Struct providing data from the post request.

type Factoring struct {
	ThreatType	float32		`json:"threatType,omitempty"`
	ThreatCategory	float32		`json:"threatCategory,omitempty"`
}

// Summary:
//	Executes a post to the adaptauth endpoint.
// Parameters:
// 	[Required] r: should have all the required fields of the struct populated before using.
//	[Required] c: passing in the client containing authorization and host information.
// Returns:
//	Response: Struct marshaled from the Json response from the API endpoints.
//	Error: If an error is encountered, response will be nil and the error must be handled.

func (r *Request) Post(c *sa.Client)(*Response, error){
	jsonRequest, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	httpRequest, err := c.BuildPostRequest(endpoint, string(jsonRequest))
	if err != nil {
		return nil, err
	}
	httpResponse, err := c.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	ipEvalResponse := new(Response)
	if err := json.NewDecoder(httpResponse.Body).Decode(ipEvalResponse); err != nil {
		return nil, err
	}
	ipEvalResponse.HttpResponse = httpResponse
	httpResponse.Body.Close()
	return ipEvalResponse, nil
}

// Summary:
//	Helper function for making IpEval Posts
// Parameters:
//	[Required] c: passing in the client containing authorization and host information.
//	[Required] userId: the user you wish to evaluate via adaptive auth.
//	[Required] ipAddress: the ip address of the user being evaluated.
// Returns:
//	Response: Struct marshaled from the Json response from the API endpoints.
//	Error: If an error is encountered, response will be nil and the error must be handled.

func (r *Request) EvaluateIp(c *sa.Client, userId string, ipAddress string)(*Response, error){
	r.UserId = userId
	r.EvalType = "risk"
	r.IpAddress = ipAddress
	ipEvalResponse, err := r.Post(c)
	if err != nil {
		return nil, err
	}
	return ipEvalResponse, nil
}