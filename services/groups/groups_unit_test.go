package groups

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
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
	uAppID       = "12345"
	uAppKey      = "12345"
	uHost        = "idp.host.com"
	uRealm       = "secureauth1"
	uPort        = 443
	uUser1       = "user1"
	uUser2       = "user2"
	uUser3       = "user3"
	uUser4       = "user4"
	uGroup1      = "group1"
	uGroup2      = "group2"
	uGroup3      = "group3"
	uGroup4      = "group4"
	uSpacedGroup = "group 5"
)

func TestGroup_Unit(t *testing.T) {
	client, err := sa.NewClient(uAppID, uAppKey, uHost, uPort, uRealm, true, false)
	if err != nil {
		t.Error(err)
	}
	userToGroupTest, err := userToGroup(client)
	if err != nil {
		t.Error(err)
	}
	if !userToGroupTest {
		t.Error("Add user to group test failed")
	}

	userToGroupsTest, err := userToGroups(client)
	if err != nil {
		t.Error(err)
	}
	if !userToGroupsTest {
		t.Error("Add user to groups test failed")
	}

	groupToUserTest, err := groupToUser(client)
	if err != nil {
		t.Error(err)
	}
	if !groupToUserTest {
		t.Error("Add group to user test failed")
	}

	groupToUsersTest, err := groupToUsers(client)
	if err != nil {
		t.Error(err)
	}
	if !groupToUsersTest {
		t.Error("Add group to users test failed")
	}
}

func userToGroup(client *sa.Client) (bool, error) {
	defer gock.Off()

	responseMock := &Response{
		Status:  "success",
		Message: "",
	}
	bytes, err := json.Marshal(responseMock)
	if err != nil {
		fmt.Println(err)
	}
	responseMockJSON := string(bytes)
	n := time.Now()
	headers := map[string]string{
		"X-SA-DATE":      n.String(),
		"X-SA-SIGNATURE": makeResponseSignature(client, responseMockJSON, n.String()),
	}
	// Set up a test responder for the api.
	gock.New("https://idp.host.com:443").
		Post("/secureauth1/api/v1/users/" + uUser1 + "/groups/" + uGroup1).
		Reply(200).
		BodyString(responseMockJSON).
		SetHeaders(headers)
	groupRequest := new(Request)
	groupResponse, err := groupRequest.AddUserToGroup(client, uUser1, uGroup1)
	if err != nil {
		return false, err
	}
	valid, err := groupResponse.IsSignatureValid(client)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, errors.New("Response signature is invalid")
	}
	return true, nil
}

func userToGroups(client *sa.Client) (bool, error) {
	defer gock.Off()

	responseMock := &Response{
		Status:  "success",
		Message: "",
	}
	bytes, err := json.Marshal(responseMock)
	if err != nil {
		fmt.Println(err)
	}
	responseMockJSON := string(bytes)
	n := time.Now()
	headers := map[string]string{
		"X-SA-DATE":      n.String(),
		"X-SA-SIGNATURE": makeResponseSignature(client, responseMockJSON, n.String()),
	}
	// Set up a test responder for the api.
	gock.New("https://idp.host.com:443").
		Post("/secureauth1/api/v1/users/" + uUser1 + "/groups").
		Reply(200).
		BodyString(responseMockJSON).
		SetHeaders(headers)
	groupRequest := new(Request)
	groups := []string{
		uGroup1,
		uGroup2,
		uGroup3,
		uGroup4,
	}
	groupResponse, err := groupRequest.AddUserToGroups(client, uUser1, groups)
	if err != nil {
		return false, err
	}
	valid, err := groupResponse.IsSignatureValid(client)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, errors.New("Response signature is invalid")
	}
	return true, nil
}

func groupToUser(client *sa.Client) (bool, error) {
	defer gock.Off()

	responseMock := &Response{
		Status:  "success",
		Message: "",
	}
	bytes, err := json.Marshal(responseMock)
	if err != nil {
		fmt.Println(err)
	}
	responseMockJSON := string(bytes)
	n := time.Now()
	headers := map[string]string{
		"X-SA-DATE":      n.String(),
		"X-SA-SIGNATURE": makeResponseSignature(client, responseMockJSON, n.String()),
	}
	// Set up a test responder for the api.
	gock.New("https://idp.host.com:443").
		Post("/secureauth1/api/v1/groups/" + uGroup2 + "/users/" + uUser2).
		Reply(200).
		BodyString(responseMockJSON).
		SetHeaders(headers)
	groupRequest := new(Request)
	groupResponse, err := groupRequest.AddGroupToUser(client, uGroup2, uUser2)
	if err != nil {
		return false, err
	}
	valid, err := groupResponse.IsSignatureValid(client)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, errors.New("Response signature is invalid")
	}
	return true, nil
}

func groupToUsers(client *sa.Client) (bool, error) {
	defer gock.Off()

	responseMock := &Response{
		Status:  "success",
		Message: "",
	}
	bytes, err := json.Marshal(responseMock)
	if err != nil {
		fmt.Println(err)
	}
	responseMockJSON := string(bytes)
	n := time.Now()
	headers := map[string]string{
		"X-SA-DATE":      n.String(),
		"X-SA-SIGNATURE": makeResponseSignature(client, responseMockJSON, n.String()),
	}
	// Set up a test responder for the api.
	gock.New("https://idp.host.com:443").
		Post("/secureauth1/api/v1/groups/" + uGroup3 + "/users").
		Reply(200).
		BodyString(responseMockJSON).
		SetHeaders(headers)
	groupRequest := new(Request)
	users := []string{
		uUser1,
		uUser2,
		uUser3,
		uUser4,
	}
	groupResponse, err := groupRequest.AddGroupToUsers(client, uGroup3, users)
	if err != nil {
		return false, err
	}
	valid, err := groupResponse.IsSignatureValid(client)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, errors.New("Response signature is invalid")
	}
	return true, nil
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
