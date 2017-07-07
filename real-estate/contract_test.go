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
)

//TODO clean up error messages
type TestAttribute struct {
	Key	string
	Value	string
}

const getOwnership = "getOwnership"
const propertyTransaction = "propertyTransaction"
const getProperty = "getProperty"
const ownership_1 = "ownership_1"
const ownership_3 = "ownership_3"
const property_1 = "property_1"
const txId = `"b4fa043e2f2d0db269c37cf415635a529dc4f2fb7f0baf64d59e23536c60cda4"`
const dateString = `"2017-06-28T21:57:16"`
const emptyOwnershipPropertyJson = `{"properties":[]}`
const emptyPropertyJson = `{"saleDate":"","salePrice":0,"owners":[]}`
const errorStatus = int32(500)

const getOwnershipArgError = "Incorrect number of arguments. Expecting ownership id to query"
const getOwnershipMissingError = "Nil value for ownershipId:"
const getOwnershipNilError = "Nil value for ownershipId:"
const getPropertyNilError = "Nil amount for"
const createPropertyTransactionArgError = "Incorrect number of arguments. Expecting 2"
const createPropertyTransactionUnmarshalWrongTypeError = "cannot unmarshal string into Go struct field Property.salePrice of type float64"
const createPropertyTransactionMissingSaleDateError = "A sale date is required."
const createPropertyTransactionGreaterThanZeroSalePriceError = "The sale price must be greater than 0"
const createPropertyTransactionTotalPercentageOfOneError = "Total Percentage can not be greater than or less than 1. Your total percentage ="
const createPropertyTransactionNoOwnersError = "At least one owner is required."
const getPropertyArgError = "Incorrect number of arguments. Expecting property id to query"

func TestGetOwnershipMissingOwnership(t *testing.T){

	stub := getStub()

	invalidArgs := getTwoArgLedgerArray(getOwnership, ownership_1)
	message := " | " + getOwnership + " with args: {" + string(invalidArgs[1]) + "}, did not fail. "

	handleExpectedFailures(t, stub, getOwnership, message, getOwnershipMissingError, invalidArgs, emptyOwnershipPropertyJson)
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

	var property = TestAttribute{}

	property = getTestProperty(txId, property_1, dateString, 1000, getValidOwners())
	checkPropertyTransaction(t,stub, property)

	initialOwnership := TestAttribute{ownership_3, `{"properties":[{"id":"property_1","percentage":0.45}]}`}
	checkGetOwnership(t, stub, initialOwnership)

	owner3 := `"id":"ownership_3","percentage":0.45`
	owner2 := `"id":"ownership_2","percentage":0.55`
	owners := []string{owner3, owner2}

	property = getTestProperty(txId, "property_2", dateString, 1000, owners)
	checkPropertyTransaction(t,stub, property)

	updatedOwnership := TestAttribute{ownership_3, `{"properties":[{"id":"property_1","percentage":0.45},{"id":"property_2","percentage":0.45}]}`}
	checkGetOwnership(t, stub, updatedOwnership)

}

func TestOwnershipCreatedDuringPropertyTransaction(t *testing.T){

	stub := getStub()

	ownershipArgs := getTwoArgLedgerArray(getOwnership, ownership_3)

	message := " | " + getOwnership + " with args: {" + string(ownershipArgs[1]) + "}, did not fail. "
	handleExpectedFailures(t, stub, getOwnership, message, getOwnershipNilError, ownershipArgs, emptyOwnershipPropertyJson)

	property := getTestProperty(txId, property_1, dateString, 1000, getValidOwners())
	checkPropertyTransaction(t,stub, property)

	validOwnership := TestAttribute{ownership_3, `{"properties":[{"id":"property_1","percentage":0.45}]}`}

	checkGetOwnership(t, stub, validOwnership)

}

func TestPropertyTransaction(t *testing.T){

	stub := getStub()

	property := getTestProperty(txId, property_1, dateString, 1000, getValidOwners())

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

func TestPropertyTransactionStringAsSalePrice(t *testing.T) {

	stub := getStub()

	validJson := `{"saleDate":"2017-06-28T21:57:16","salePrice":"1000","owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.55}]}`

	invalidArgs := getThreeArgLedgerArray(propertyTransaction, property_1, validJson)

	message := " | " + propertyTransaction + " with args: {" + string(invalidArgs[1]) + "," + string(invalidArgs[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, message, createPropertyTransactionUnmarshalWrongTypeError, invalidArgs, validJson )

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

	property := getTestProperty(txId, property_1, dateString, 1000, getValidOwners())

	checkPropertyTransaction(t, stub, property)
	checkGetProperty(t, stub, property)

}

func TestGetPropertyExtraArgs(t *testing.T){

	stub := getStub()

	property := getTestProperty(txId, property_1, dateString, 1000, getValidOwners())

	checkPropertyTransaction(t, stub, property)

	invalidArgs := getThreeArgLedgerArray(getProperty, property.Key, property.Value)

	message:= " | " + getProperty + " with args: {" + string(invalidArgs[1]) + "," + string(invalidArgs[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, getProperty, message, getPropertyArgError, invalidArgs, property.Value)

}

func TestGetPropertyMissingProperty(t *testing.T){

	stub := getStub()

	invalidArgs := getTwoArgLedgerArray(getProperty, property_1)
	message := " | " + getProperty + " with args: {" + property_1 + "}, did not fail. "

	handleExpectedFailures(t, stub, getProperty, message, getPropertyNilError, invalidArgs, emptyPropertyJson)

}

func checkGetOwnership(t *testing.T, stub *shim.MockStub, ownership TestAttribute){

	ownershipArgs := getTwoArgLedgerArray(getOwnership, ownership.Key)

	message:= " | " + getOwnership + " with args: {" + string(ownershipArgs[1]) + "}, failed. "

	handleExpectedSuccess(t, stub, getOwnership, message, ownershipArgs, ownership.Value)

}

func checkPropertyTransaction(t *testing.T, stub *shim.MockStub, property TestAttribute){

	propertyArgs := getThreeArgLedgerArray(propertyTransaction, property.Key, property.Value)

	message:= " | " + propertyTransaction + " with args: {" + string(propertyArgs[1]) + "," + string(propertyArgs[2]) + "}, failed. "

	res := stub.MockInvoke(propertyTransaction, propertyArgs)
	if res.Status != shim.OK {
		message := message +  "[res.Status=" + strconv.FormatInt(int64(res.Status), 10) + "]"
		fmt.Println(message)
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

	propertyArgs := getTwoArgLedgerArray(getProperty, property.Key)

	message:= " | " + getProperty + " with args: {" + string(propertyArgs[1]) + "}, failed. "

	handleExpectedSuccess(t, stub, getProperty, message, propertyArgs, property.Value)

}

//func getTestProperty(propertyId string, saleDate string, salePrice float64, owners []string) (TestAttribute) {
//
//	var buffer bytes.Buffer
//
//	buffer.WriteString("{")
//
//	saleDateKey := `"saleDate":`
//	buffer.WriteString(saleDateKey)
//	buffer.WriteString(saleDate)
//	buffer.WriteString(",")
//
//	salePriceKey := `"salePrice":`
//	buffer.WriteString(salePriceKey)
//	jsonSalePrice := strconv.FormatFloat(salePrice, 'f', -1, 64)
//	buffer.WriteString(jsonSalePrice)
//	buffer.WriteString(",")
//
//	ownersKey := `"owners":`
//	buffer.WriteString(ownersKey)
//	ownersAttribute := getTestOwners(propertyId, owners)
//	buffer.WriteString(ownersAttribute.Value)
//
//	buffer.WriteString("}")
//
//	jsonProperty := TestAttribute{propertyId, buffer.String()}
//
//	return jsonProperty
//
//}


//{"txid":"358b5b558c6f1e3fcb09af9510570f748e8e5a0bb16934471f98b0fe0a1514a9","id":"property_1","saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownership_2","percentage":0.55}]}
func getTestProperty(txid string, propertyId string, saleDate string, salePrice float64, owners []string) (TestAttribute) {

	var buffer bytes.Buffer

	buffer.WriteString("{")

	txIdKey := `"txid":`
	buffer.WriteString(txIdKey)
	buffer.WriteString(txid)
	buffer.WriteString(",")

	propertyIdKey := `"id":`
	buffer.WriteString(propertyIdKey)
	buffer.WriteString("\"")
	buffer.WriteString(propertyId)
	buffer.WriteString("\",")

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

	fmt.Println(jsonProperty)

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

func handleExpectedSuccess(t *testing.T, stub *shim.MockStub, argument string, outputMessage string, invalidArgs [][]byte, attemptedPayload string){

	res := stub.MockInvoke(argument, invalidArgs)
	if res.Status != shim.OK {
		msg := outputMessage +  "[res.Status=" + strconv.FormatInt(int64(res.Status), 10) + "]"
		fmt.Println(msg)
		t.FailNow()
	}
	if res.Payload == nil {
		msg := outputMessage + "[res.Message=" + string(res.Message) + "]"
		fmt.Println(msg)
		t.FailNow()
	}
	if string(res.Payload) != attemptedPayload {
		msg := outputMessage + "[res.Payload=" + string(res.Payload) + "]"
		fmt.Println(msg)
		t.FailNow()
	}

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