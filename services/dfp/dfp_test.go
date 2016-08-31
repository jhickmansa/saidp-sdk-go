package dfp

import (
	"testing"
	"fmt"
	sa "github.com/SecureAuthCorp/saidp-sdk-go"
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

const (
	appId = ""
	appKey = ""
	host = "host.company.com"
	realm = "secureauth1"
	port = 443
	user = "user"
	host_addr = "192.168.0.1"
	fingerprintJson = ``
	accept = "";
	accept_encode = ""
	accept_lang = ""
	accept_charset = ""
)

func TestDFPRequest (t *testing.T){
	client, err := sa.NewClient(appId, appKey, host, port, realm, true, false)
	if err != nil {
		fmt.Println(err)
	}
	dfpRequest := new(Request)
	jsResponse, err := dfpRequest.GetDfpJs(client)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Response Struct for JavaScript Source:")
	fmt.Println(jsResponse)
	dfpValResponse, err := dfpRequest.ValidateDfp(client, user, host_addr, "", fingerprintJson, accept, accept_charset, accept_encode, accept_lang)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Response Struct for DFP Validate:")
	fmt.Println(dfpValResponse)
	dfpConResponse, err := dfpRequest.ConfirmDfp(client, user, dfpValResponse.FingerprintId)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Response Struct for DFP Confirm:")
	fmt.Println(dfpConResponse)
}