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
	"strconv"
	"bytes"
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

func TestCreateOwnershipIncorrectArgs(t *testing.T){

	stub := getStub()

	function := "createOwnership"
	ownershipId := "ownership_1"
	errorStatus := int32(500)
	errorMessage := "Incorrect number of arguments. Expecting ownership id and properties"

	badArgs := [][]byte{[]byte(function), []byte(ownershipId)}

	res := stub.MockInvoke(function, badArgs)
	output:= " CreateOwnership " + ownershipId + " did not fail |" + string(res.Message)
	if res.Status != errorStatus {
		fmt.Println(output)
		t.FailNow()
	}
	if res.Message != errorMessage {
		fmt.Println(output)
		t.FailNow()
	}

}

func TestGetOwnership(t *testing.T){

	stub := getStub()

	ownership := getOwnership("ownership_1", `{"properties":[]}`)

	checkInvokeOwnership(t, stub, ownership)
	checkGetOwnership(t, stub, ownership)

}

func TestOwnershipUpdatedDuringPropertyTransaction(t *testing.T){

	stub := getStub()

	originalOwnership := getOwnership("ownership_3", `{"properties":[]}`)

	checkInvokeOwnership(t, stub, originalOwnership)
	checkGetOwnership(t, stub, originalOwnership)

	owner1 := `"id":"ownership_3","percentage":0.45`
	owner2 := `"id":"ownership_2","percentage":0.55`

	owners := []string{owner1, owner2}

	property := getProperty("property_1", `"2017-06-28T21:57:16"`, 1000, owners)

	checkPropertyTransaction(t,stub, property)

	updatedOwnership := getOwnership("ownership_3", `{"properties":[{"id":"property_1","percentage":0.45}]}`)

	checkGetOwnership(t, stub, updatedOwnership)

}

func TestOwnershipCreatedDuringPropertyTransaction(t *testing.T){

	stub := getStub()

	invalidOwnership := TestAttribute{"ownership_3", ""}

	verifyOwnershipIsNotInLedger(t,stub, invalidOwnership)

	owner1 := `"id":"ownership_3","percentage":0.45`
	owner2 := `"id":"ownership_2","percentage":0.55`

	owners := []string{owner1, owner2}

	property := getProperty("property_1", `"2017-06-28T21:57:16"`, 1000, owners)

	checkPropertyTransaction(t,stub, property)

	validOwnership := getOwnership("ownership_3", `{"properties":[{"id":"property_1","percentage":0.45}]}`)

	checkGetOwnership(t, stub, validOwnership)

}

func TestPropertyTransaction(t *testing.T){

	stub := getStub()
	owner1 := `"id":"ownership_3","percentage":0.45`
	owner2 := `"id":"ownership_2","percentage":0.55`

	owners := []string{owner1, owner2}

	property := getProperty("property_1", `"2017-06-28T21:57:16"`, 1000, owners)

	checkPropertyTransaction(t,stub, property)
	checkPropertyState(t, stub, property)

}

func TestGetProperty(t *testing.T){

	stub := getStub()
	owner1 := `"id":"ownership_3","percentage":0.45`
	owner2 := `"id":"ownership_2","percentage":0.55`

	owners := []string{owner1, owner2}

	property := getProperty("property_1", `"2017-06-28T21:57:16"`, 1000, owners)

	checkPropertyTransaction(t, stub, property)
	checkGetProperty(t, stub, property)

}

func checkInvokeOwnership(t *testing.T, stub *shim.MockStub, ownership TestAttribute) {

	function := "createOwnership"

	ownershipArgs := create(function, ownership.Key, ownership.Value)

	res := stub.MockInvoke(function, ownershipArgs)
	if res.Status != shim.OK {
		fmt.Println(" InvokeOwnership", ownership, "failed", string(res.Message))
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

	ownershipArgs := get(function, ownership.Key)

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

func checkPropertyTransaction(t *testing.T, stub *shim.MockStub, property TestAttribute){

	function := "propertyTransaction"

	propertyArgs := create(function, property.Key, property.Value)

	res := stub.MockInvoke(function, propertyArgs)
	if res.Status != shim.OK {
		fmt.Println("InvokeProperty", property, "failed", string(res.Message))
		t.FailNow()
	}

}

func checkPropertyState(t *testing.T, stub *shim.MockStub, property TestAttribute) {

	bytes := stub.State[property.Key]
	if bytes == nil {
		fmt.Println("Property", property, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != property.Value {
		fmt.Println("Property value", property.Key, "was not", property.Value, "as expected")
		t.FailNow()
	}

}

func checkGetProperty(t *testing.T, stub *shim.MockStub, property TestAttribute){

	function := "getProperty"

	propertyArgs := get(function, property.Key)

	res := stub.MockInvoke(function, propertyArgs)
	if res.Status != shim.OK {
		fmt.Println(function, property, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println(function, property, "failed to get value")
		t.FailNow()
	}
	if string(res.Payload) != property.Value {
		fmt.Println(function, " value", property.Key, "was not", property.Value, "as expected")
		t.FailNow()
	}

}

func getOwnership(ownershipId string, properties string) (TestAttribute){

	ownership := TestAttribute{ownershipId, properties}

	return ownership

}

func getProperty(propertyId string, saleDate string, salePrice float64, owners []string) (TestAttribute) {

	var buffer bytes.Buffer

	buffer.WriteString("{")

	saleDateKey := `"saleDate":`
	buffer.WriteString(saleDateKey)
	buffer.WriteString(saleDate)
	buffer.WriteString(",")

	salePriceKey := `"salePrice":`
	buffer.WriteString(salePriceKey)
	jsonSalePrice := strconv.FormatFloat(salePrice, 'f', -1, 64)
	buffer.WriteString(jsonSalePrice)
	buffer.WriteString(",")

	ownersKey := `"owners":`
	buffer.WriteString(ownersKey)
	ownersAttribute := getOwners(propertyId, owners)
	buffer.WriteString(ownersAttribute.Value)

	buffer.WriteString("}")

	jsonProperty := TestAttribute{propertyId, buffer.String()}

	return jsonProperty

}

func getOwners(propertyId string, owners []string) (TestAttribute) {

	var buffer bytes.Buffer

	buffer.WriteString("[")

	for i := 0; i < len(owners); i++ {

		buffer.WriteString("{")
		buffer.WriteString(owners[i])
		buffer.WriteString("}")

		if i != (len(owners) - 1){
			buffer.WriteString(",")
		}

	}

	buffer.WriteString("]")

	propertyOwners := TestAttribute{propertyId, buffer.String()}

	return propertyOwners

}

func create(function string, id string, jsonStruct string) ([][]byte) {

	method := []byte(function)
	key := []byte(id)
	value := []byte(jsonStruct)

	args := [][]byte{method, key, value}

	return args
}

func get(function string, id string) ([][]byte) {

	method := []byte(function)
	key := []byte(id)

	args := [][]byte{method, key}

	return args
}

func getStub() (*shim.MockStub){

	scc := new(Chaincode)
	stub := shim.NewMockStub("contract", scc)

	return stub

}

func verifyOwnershipIsNotInLedger(t *testing.T, stub *shim.MockStub, ownership TestAttribute){

	function := "getOwnership"
	errorStatus := int32(500)

	ownershipArgs := get(function, ownership.Key)

	res := stub.MockInvoke(function, ownershipArgs)
	if res.Status != errorStatus {
		fmt.Println(function, ownership, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload != nil {
		fmt.Println(function, ownership, "value returned")
		t.FailNow()
	}
	if string(res.Payload) != ownership.Value {
		fmt.Println(function, " value", ownership.Key, "was not", ownership.Value, "as expected")
		t.FailNow()
	}

}