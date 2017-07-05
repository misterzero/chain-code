/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"testing"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type TestAttribute struct {
	Key	string
	Value	string
}

func TestCreateOwnership(t *testing.T){

	stub := getStub()

	ownership := getOwnership("ownership_1", `{"properties":[]}`)

	checkInvokeOwnership(t, stub, ownership)
	checkOwnershipState(t, stub, ownership)

}

func TestGetOwnership(t *testing.T){

	stub := getStub()

	ownership := getOwnership("ownership_1", `{"properties":[]}`)

	checkInvokeOwnership(t, stub, ownership)
	checkGetOwnership(t, stub, ownership)

}


func checkInvokeOwnership(t *testing.T, stub *shim.MockStub, ownership TestAttribute) {

	function := "createOwnership"

	ownershipArgs := getOwnershipCall(function, ownership.Key, ownership.Value)

	res := stub.MockInvoke(function, ownershipArgs)
	if res.Status != shim.OK {
		fmt.Println("InvokeOwnership", ownership, "failed", string(res.Message))
		t.FailNow()
	}

}

func checkOwnershipState(t *testing.T, stub *shim.MockStub, ownership TestAttribute) {

	bytes := stub.State[ownership.Key]
	if bytes == nil {
		fmt.Println("Properties", ownership, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != ownership.Value {
		fmt.Println("Properties value", ownership.Key, "was not", ownership.Value, "as expected")
		t.FailNow()
	}

}

func checkGetOwnership(t *testing.T, stub *shim.MockStub, ownership TestAttribute){

	function := "getOwnership"

	ownershipArgs := [][]byte{[]byte(function), []byte(ownership.Key)}

	res := stub.MockInvoke(function, ownershipArgs)
	if res.Status != shim.OK {
		fmt.Println(function, ownership, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println(function, ownership, "failed to get value")
		t.FailNow()
	}
	if string(res.Payload) != ownership.Value {
		fmt.Println(function, " value", ownership.Key, "was not", ownership.Value, "as expected")
		t.FailNow()
	}

}

func getOwnership(ownershipId string, properties string) (TestAttribute){

	ownership := TestAttribute{ownershipId, properties}

	return ownership

}

func getOwnershipCall(function string, ownershipId string, properties string) ([][]byte) {

	method := []byte(function)
	key := []byte(ownershipId)
	value := []byte(properties)

	args := [][]byte{method, key, value}

	return args

}

func getStub() (*shim.MockStub){

	scc := new(Chaincode)
	stub := shim.NewMockStub("contract", scc)

	return stub

}