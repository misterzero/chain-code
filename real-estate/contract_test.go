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
	"strings"
	//"encoding/json"
)

//TODO clean up error messages
//TODO make generic message for error
//TODO make a single location for expected error message
type TestAttribute struct {
	Key	string
	Value	string
}

const createOwnership = "createOwnership"
const getOwnership = "getOwnership"
const propertyTransaction = "propertyTransaction"
const getProperty = "getProperty"
const ownership_1 = "ownership_1"
const ownership_3 = "ownership_3"
const property_1 = "property_1"
const dateString = `"2017-06-28T21:57:16"`
const emptyOwnershipPropertyJson = `{"properties":[]}`
const errorStatus = int32(500)

const createOwnershipArgError = "Incorrect number of arguments. Expecting ownership id and properties"
const getOwnershipArgError = "Incorrect number of arguments. Expecting ownership id to query"
const getOwnershipNilError = "Nil value for ownershipId:"
const invalidJsonForOwnershipStructError = "Unable to convert json to Ownership struct"

const createPropertyTransactionArgError = "Incorrect number of arguments. Expecting 2"
const createPropertyTransactionMissingSaleDateError = "A sale date is required."
const createPropertyTransactionGreaterThanZeroSalePriceError = "The sale price must be greater than 0"
const createPropertyTransactionTotalPercentageOfOneError = "Total Percentage can not be greater than or less than 1. Your total percentage ="
const createPropertyTransactionNoOwnersError = "At least one owner is required."


func TestCreateOwnership(t *testing.T){

	stub := getStub()

	ownership := TestAttribute{ownership_1, emptyOwnershipPropertyJson}

	checkInvokeOwnership(t, stub, ownership)
	checkOwnershipState(t, stub, ownership)

}

func TestCreateOwnershipExtraArgs(t *testing.T){

	stub := getStub()

	invalidArgument := "invalidArgument"

	invalidArgs := getThreeArgLedgerArray(createOwnership, ownership_1, emptyOwnershipPropertyJson)
	invalidArgs = append(invalidArgs, []byte(invalidArgument))

	message := " | " + createOwnership + " with args: {" + string(invalidArgs[1]) + "," + string(invalidArgs[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, createOwnership, message, createOwnershipArgError, invalidArgs, invalidArgument)

}

func TestCreateOwnershipInvalidJson(t *testing.T){

	stub := getStub()

	invalidJson := `"{"properties":}`

	invalidArgs := getThreeArgLedgerArray(createOwnership, ownership_1, invalidJson)
	message := " | " + createOwnership + " with args: {" + string(invalidArgs[1]) + "," + string(invalidArgs[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, createOwnership, message, invalidJsonForOwnershipStructError, invalidArgs, invalidJson)

}

func TestGetOwnership(t *testing.T){

	stub := getStub()

	ownership := TestAttribute{ownership_1, emptyOwnershipPropertyJson}

	checkInvokeOwnership(t, stub, ownership)
	checkGetOwnership(t, stub, ownership)

}

func TestGetOwnershipExtraArgs(t *testing.T){

	stub := getStub()

	invalidArgument := "invalidArgument"

	invalidArgs := getThreeArgLedgerArray(getOwnership, ownership_1, invalidArgument)
	message := " | " + getOwnership + " with args: {" + string(invalidArgs[1]) + "," + string(invalidArgs[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, getOwnership, message, getOwnershipArgError, invalidArgs, invalidArgument)

}

func TestOwnershipUpdatedDuringPropertyTransaction(t *testing.T){

	stub := getStub()

	originalOwnership := TestAttribute{ownership_3, emptyOwnershipPropertyJson}

	checkInvokeOwnership(t, stub, originalOwnership)
	checkGetOwnership(t, stub, originalOwnership)

	property := getTestProperty(property_1, dateString, 1000, getValidOwners())

	checkPropertyTransaction(t,stub, property)

	updatedOwnership := TestAttribute{ownership_3, `{"properties":[{"id":"property_1","percentage":0.45}]}`}

	checkGetOwnership(t, stub, updatedOwnership)

}

func TestOwnershipCreatedDuringPropertyTransaction(t *testing.T){

	stub := getStub()

	finalPayload := `{"properties":[{"id":"property_1","percentage":0.45}]}`

	ownershipArgs := getTwoArgLedgerArray(getOwnership,ownership_3)
	message := " | " + getOwnership + " with args: {" + string(ownershipArgs[1]) + "}, did not fail. "

	handleExpectedFailures(t, stub, getOwnership, message, getOwnershipNilError, ownershipArgs, finalPayload)

	property := getTestProperty(property_1, dateString, 1000, getValidOwners())

	checkPropertyTransaction(t,stub, property)

	validOwnership := TestAttribute{ownership_3, `{"properties":[{"id":"property_1","percentage":0.45}]}`}

	checkGetOwnership(t, stub, validOwnership)

}

func TestPropertyTransaction(t *testing.T){

	stub := getStub()

	property := getTestProperty(property_1, dateString, 1000, getValidOwners())

	checkPropertyTransaction(t,stub, property)
	checkPropertyState(t, stub, property)

}

func TestPropertyTransactionExtraArgs(t *testing.T) {

	stub := getStub()

	invalidArgument := "invalidArgument"
	validJson := `{"saleDate":"2017-06-28T21:57:16","salePrice":-1,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.55}]}`

	invalidArgs := getThreeArgLedgerArray(propertyTransaction, property_1, validJson)
	invalidArgs = append(invalidArgs, []byte(invalidArgument))

	message := " | " + propertyTransaction + " with args: {" + string(invalidArgs[1]) + "," + string(invalidArgs[2]) + ", " + string(invalidArgs[3]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, message, createPropertyTransactionArgError, invalidArgs, validJson )

}

func TestPropertyTransactionMissingSaleDate(t *testing.T) {

	stub := getStub()

	missingSaleDateJson := `{"salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownership_2","percentage":0.55}]}`

	invalidArgs := getThreeArgLedgerArray(propertyTransaction, property_1, missingSaleDateJson)
	message := " | " + propertyTransaction + " with args: {" + string(invalidArgs[1]) + "," + string(invalidArgs[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, message, createPropertyTransactionMissingSaleDateError, invalidArgs, missingSaleDateJson)

}

func TestPropertyTransactionNegativeSalePrice(t *testing.T) {

	stub := getStub()

	negativeSalePriceJson := `{"saleDate":"2017-06-28T21:57:16","salePrice":-1,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.55}]}`

	invalidArgs := getThreeArgLedgerArray(propertyTransaction, property_1, negativeSalePriceJson)
	message := " | " + propertyTransaction + " with args: {" + string(invalidArgs[1]) + "," + string(invalidArgs[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, message, createPropertyTransactionGreaterThanZeroSalePriceError, invalidArgs, negativeSalePriceJson)

}

func TestPropertyTransactionNoOwners(t *testing.T) {

	stub := getStub()

	noOwnersJson := `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[]}`
	
	invalidArgs := getThreeArgLedgerArray(propertyTransaction, property_1, noOwnersJson)
	message:= " | " + propertyTransaction + " with args: {" + string(invalidArgs[1]) + "," + string(invalidArgs[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, message, createPropertyTransactionNoOwnersError, invalidArgs, noOwnersJson)

}

func TestPropertyTransactionTooLowTotalOwnerPercentage(t *testing.T) {

	stub := getStub()

	tooLowOwnerPercentage := `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.50}]}`

	invalidArgs := getThreeArgLedgerArray(propertyTransaction, property_1, tooLowOwnerPercentage)
	message := " | " + propertyTransaction + " with args: {" + string(invalidArgs[1]) + "," + string(invalidArgs[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, message, createPropertyTransactionTotalPercentageOfOneError, invalidArgs, tooLowOwnerPercentage)

}

func TestPropertyTransactionTooHighTotalOwnerPercentage(t *testing.T) {

	stub := getStub()

	tooHighOwnerPercentage := `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.70}]}`

	invalidArgs := getThreeArgLedgerArray(propertyTransaction, property_1, tooHighOwnerPercentage)
	message := " | " + propertyTransaction + " with args: {" + string(invalidArgs[1]) + "," + string(invalidArgs[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, message, createPropertyTransactionTotalPercentageOfOneError, invalidArgs, tooHighOwnerPercentage)

}

func TestGetProperty(t *testing.T){

	stub := getStub()

	owner1 := `"id":"ownership_3","percentage":0.45`
	owner2 := `"id":"ownership_2","percentage":0.55`
	owners := []string{owner1, owner2}

	property := getTestProperty(property_1, dateString, 1000, owners)

	checkPropertyTransaction(t, stub, property)
	checkGetProperty(t, stub, property)

}

//TODO fix error message
func checkInvokeOwnership(t *testing.T, stub *shim.MockStub, ownership TestAttribute) {

	ownershipArgs := getThreeArgLedgerArray(createOwnership, ownership.Key, ownership.Value)

	res := stub.MockInvoke(createOwnership, ownershipArgs)
	if res.Status != shim.OK {
		fmt.Println(" InvokeOwnership", ownership, "failed", string(res.Message))
		t.FailNow()
	}
}

//TODO fix error message
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

//TODO fix error message
func checkGetOwnership(t *testing.T, stub *shim.MockStub, ownership TestAttribute){

	ownershipArgs := getTwoArgLedgerArray(getOwnership, ownership.Key)

	res := stub.MockInvoke(getOwnership, ownershipArgs)
	if res.Status != shim.OK {
		fmt.Println(getOwnership, ownership, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println(getOwnership, ownership, "failed to get value")
		t.FailNow()
	}
	if string(res.Payload) != ownership.Value {
		fmt.Println(getOwnership, " value", ownership.Key, "was not", ownership.Value, "as expected")
		t.FailNow()
	}

}

//TODO fix error message
func checkPropertyTransaction(t *testing.T, stub *shim.MockStub, property TestAttribute){

	propertyArgs := getThreeArgLedgerArray(propertyTransaction, property.Key, property.Value)

	res := stub.MockInvoke(propertyTransaction, propertyArgs)
	if res.Status != shim.OK {
		fmt.Println("InvokeProperty", property, "failed", string(res.Message))
		t.FailNow()
	}

}

//TODO fix error message
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

//TODO fix error message
func checkGetProperty(t *testing.T, stub *shim.MockStub, property TestAttribute){

	propertyArgs := getTwoArgLedgerArray(getProperty, property.Key)

	res := stub.MockInvoke(getProperty, propertyArgs)
	if res.Status != shim.OK {
		fmt.Println(getProperty, property, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println(getProperty, property, "failed to get value")
		t.FailNow()
	}
	if string(res.Payload) != property.Value {
		fmt.Println(getProperty, " value", property.Key, "was not", property.Value, "as expected")
		t.FailNow()
	}

}

func getTestProperty(propertyId string, saleDate string, salePrice float64, owners []string) (TestAttribute) {

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
	ownersAttribute := getTestOwners(propertyId, owners)
	buffer.WriteString(ownersAttribute.Value)

	buffer.WriteString("}")

	jsonProperty := TestAttribute{propertyId, buffer.String()}

	return jsonProperty

}

func getTestOwners(propertyId string, owners []string) (TestAttribute) {

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

func handleExpectedFailures(t *testing.T, stub *shim.MockStub, argument string, outputMessage string, errorMessage string, invalidArgs [][]byte, attemptedPayload string){

	res := stub.MockInvoke(argument, invalidArgs)
	if res.Status != errorStatus {
		msg := outputMessage +  "[res.Status=" + strconv.FormatInt(int64(res.Status), 10) + "]"
		fmt.Println(msg)
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorMessage) {
		msg := outputMessage + "[res.Message=" + string(res.Message) + "]"
		fmt.Println(msg)
		t.FailNow()
	}
	if string(res.Payload) == attemptedPayload{
		msg := outputMessage + "[res.Payload=" + string(res.Payload) + "]"
		fmt.Println(msg)
		t.FailNow()
	}

}

func getThreeArgLedgerArray(function string, id string, jsonStruct string) ([][]byte) {

	method := []byte(function)
	key := []byte(id)
	value := []byte(jsonStruct)

	args := [][]byte{method, key, value}

	return args
}

func getTwoArgLedgerArray(function string, id string) ([][]byte) {

	method := []byte(function)
	key := []byte(id)

	args := [][]byte{method, key}

	return args
}

func getValidOwners() ([]string) {

	owner1 := `"id":"ownership_3","percentage":0.45`
	owner2 := `"id":"ownership_2","percentage":0.55`

	owners := []string{owner1, owner2}

	return owners

}

func getStub() (*shim.MockStub){

	scc := new(Chaincode)
	stub := shim.NewMockStub("contract", scc)

	return stub

}