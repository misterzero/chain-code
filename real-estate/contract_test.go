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
	"testing"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"fmt"
	"strings"
)

func TestGetOwnershipMissingOwnership(t *testing.T){

	stub := getStub()

	payload := ownership_1
	args := [][]byte{
		[]byte(getOwnership),
		[]byte(getOwnership),
		[]byte(payload)}

	failedTestMessage := " | " + getOwnership + " with args: {" + string(args[1]) + "}, did not fail. "

	handleExpectedFailures(t, stub, getOwnership, failedTestMessage, args, payload, nilValueForOwnershipId)

}

func TestGetOwnershipExtraArgs(t *testing.T){

	stub := getStub()

	payload := invalidArgument
	args := [][]byte{
		[]byte(getOwnership),
		[]byte(getOwnership),
		[]byte(ownership_1),
		[]byte(payload)}

	failedTestMessage := " | " + getOwnership + " with args: {" + string(args[0]) + ", " + string(args[1]) + ", " + string(args[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, getOwnership, failedTestMessage, args, invalidArgument, incorrectNumberOfArgs)

}

func TestOwnershipCreatedDuringPropertyTransaction(t *testing.T){

	stub := getStub()

	var property = Property{}

	owners := getValidOwnersList()

	property, propertyAsString := createProperty(property_1, dateString, 1000, owners)

	propertyAsBytes := stub.State[property.PropertyId]
	if propertyAsBytes != nil{
		fmt.Println(property_1 + " should not exist.")
		t.Fail()
	}

	invokePropertyTransaction(t, stub, property.PropertyId, propertyAsString)

	propertyId := "1"
	propertyOwnersList := getValidPropertyOwnersList(propertyId)

	ownershipPropertyAsString := getAttributesAsString([]Attribute{propertyOwnersList[0]})

	invokeGetOwnership(t, stub, owners[0].Id, ownershipPropertyAsString)

}

//TODO next
func TestPropertyTransaction(t *testing.T){

	stub := getStub()

	property, propertyString := createProperty(property_1, dateString, 1000, getValidOwnersList())

	invokePropertyTransaction(t,stub, property.PropertyId, propertyString)
	checkPropertyState(t, stub, property, propertyString)

}

func TestPropertyTransactionExtraArgs(t *testing.T) {

	stub := getStub()

	_, propertyAsString := createProperty(property_1, dateString, 1000, getValidOwnersList())

	args := [][]byte{
		[]byte(propertyTransaction),
		[]byte(propertyTransaction),
		[]byte(property_1),
		[]byte(propertyAsString),
		[]byte("extraArg")}

	failedTestMessage := " | " + propertyTransaction + " with args: {" + string(args[1]) + ", " + string(args[2]) + ", " + string(args[3]) + ", " + string(args[4]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, failedTestMessage, args, propertyAsString, incorrectNumberOfArgs)

}

func TestPropertyTransactionStringAsSalePrice(t *testing.T) {

	stub := getStub()

	validJson := `{"saleDate":"2017-06-28T21:57:16","salePrice":"1000","owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.55}]}`

	args := [][]byte{
		[]byte(propertyTransaction),
		[]byte(propertyTransaction),
		[]byte(property_1),
		[]byte(validJson)}

	failedTestMessage := " | " + propertyTransaction + " with args: {" + string(args[1]) + "," + string(args[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, failedTestMessage, args, validJson, cannotUnmarshalStringIntoFloat64)

}

func TestPropertyTransactionMissingSaleDate(t *testing.T) {

	stub := getStub()

	missingSaleDateJson := `{"salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownership_2","percentage":0.55}]}`

	args := [][]byte{
		[]byte(propertyTransaction),
		[]byte(propertyTransaction),
		[]byte(property_1),
		[]byte(missingSaleDateJson)}

	failedTestMessage := " | " + propertyTransaction + " with args: {" + string(args[1]) + "," + string(args[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, failedTestMessage, args, missingSaleDateJson, saleDateRequired)

}

func TestPropertyTransactionNegativeSalePrice(t *testing.T) {

	stub := getStub()

	negativeSalePriceJson := `{"saleDate":"2017-06-28T21:57:16","salePrice":-1,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.55}]}`

	args := [][]byte{
		[]byte(propertyTransaction),
		[]byte(propertyTransaction),
		[]byte(property_1),
		[]byte(negativeSalePriceJson)}

	failedTestMessage := " | " + propertyTransaction + " with args: {" + string(args[1]) + "," + string(args[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, failedTestMessage, args, negativeSalePriceJson, salePriceMustBeGreaterThan0)

}

func TestPropertyTransactionNoOwners(t *testing.T) {

	stub := getStub()

	noOwnersJson := `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[]}`

	args := [][]byte{
		[]byte(propertyTransaction),
		[]byte(propertyTransaction),
		[]byte(property_1),
		[]byte(noOwnersJson)}

	failedTestMessage := " | " + propertyTransaction + " with args: {" + string(args[1]) + "," + string(args[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, failedTestMessage, args, noOwnersJson, atLeastOneOwnerIsRequired)

}

func TestPropertyTransactionTooLowTotalOwnerPercentage(t *testing.T) {

	stub := getStub()

	tooLowOwnerPercentage := `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.50}]}`

	args := [][]byte{
		[]byte(propertyTransaction),
		[]byte(propertyTransaction),
		[]byte(property_1),
		[]byte(tooLowOwnerPercentage)}

	failedTestMessage := " | " + propertyTransaction + " with args: {" + string(args[1]) + "," + string(args[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, failedTestMessage, args, tooLowOwnerPercentage, totalPercentageCanNotBeGreaterThan1)

}

func TestPropertyTransactionTooHighTotalOwnerPercentage(t *testing.T) {

	stub := getStub()

	tooHighOwnerPercentage := `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.70}]}`

	args := [][]byte{
		[]byte(propertyTransaction),
		[]byte(propertyTransaction),
		[]byte(property_1),
		[]byte(tooHighOwnerPercentage)}

	failedTestMessage := " | " + propertyTransaction + " with args: {" + string(args[1]) + "," + string(args[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, failedTestMessage, args, tooHighOwnerPercentage, totalPercentageCanNotBeGreaterThan1)

}

func TestGetProperty(t *testing.T){

	stub := getStub()

	property, propertyString := createProperty(property_1, dateString, 1000, getValidOwnersList())

	invokePropertyTransaction(t, stub, property.PropertyId, propertyString)
	checkGetProperty(t, stub, property, propertyString)

}

func TestGetPropertyExtraArgs(t *testing.T){

	stub := getStub()

	property, propertyString := createProperty(property_1, dateString, 1000, getValidOwnersList())

	invokePropertyTransaction(t, stub, property.PropertyId, propertyString)

	args := [][]byte{
		[]byte(getProperty),
		[]byte(getProperty),
		[]byte(property.PropertyId),
		[]byte(propertyString)}

	failedTestMessage := " | " + getProperty + " with args: {" + string(args[1]) + "," + string(args[2]) + "," + string(args[3]) + "}, did not fail. "

	handleExpectedFailures(t, stub, getProperty, failedTestMessage, args, propertyString, incorrectNumberOfArgs)

}

func TestGetPropertyMissingProperty(t *testing.T){

	stub := getStub()

	args := [][]byte{
		[]byte(getProperty),
		[]byte(getProperty),
		[]byte(property_1)}

	failedTestMessage := " | " + getProperty + " with args: {" + property_1 + "}, did not fail. "

	handleExpectedFailures(t, stub, getProperty, failedTestMessage, args, emptyPropertyJson, nilAmountFor)

}

func handleExpectedSuccess(t *testing.T, stub *shim.MockStub, argument string, outputMessage string, args [][]byte, attemptedPayload string){

	res := stub.MockInvoke(argument, args)

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

func handleExpectedFailures(t *testing.T, stub *shim.MockStub, argument string, failedTestMessage string, args [][]byte, payload string, expectedResponseMessage string){

	response := stub.MockInvoke(argument, args)

	if response.Status != 500 {
		msg := failedTestMessage +  "[response.Status=" + strconv.FormatInt(int64(response.Status), 10) + "]"
		fmt.Println(msg)
		t.FailNow()
	}
	if !strings.Contains(response.Message, expectedResponseMessage) {
		msg := failedTestMessage + "[response.Message=" + string(response.Message) + "]"
		fmt.Println(msg)
		t.FailNow()
	}
	if string(response.Payload) == payload {
		msg := failedTestMessage + "[response.Payload=" + string(response.Payload) + "]"
		fmt.Println(msg)
		t.FailNow()
	}

}

func checkGetProperty(t *testing.T, stub *shim.MockStub, property Property, propertyString string){

	args := [][]byte{
		[]byte(getProperty),
		[]byte(getProperty),
		[]byte(property.PropertyId)}

	failedTestMessage := " | " + getProperty + " with args: {" + string(args[1]) + "}, failed. "

	handleExpectedSuccess(t, stub, getProperty, failedTestMessage, args, propertyString)

}

func checkPropertyState(t *testing.T, stub *shim.MockStub, property Property, propertyString string) []byte{

	bytes := stub.State[property.PropertyId]
	if bytes == nil {
		fmt.Println("Property", string(bytes), "failed to get value")
		t.FailNow()
	}
	if string(bytes) != propertyString {
		fmt.Println("Property value", property.PropertyId, "was not", propertyString, "as expected")
		t.FailNow()
	}

	return bytes

}

func invokeGetOwnership(t *testing.T, stub *shim.MockStub,ownershipId string, payload string){

	args := [][]byte{
		[]byte(getOwnership),
		[]byte(getOwnership),
		[]byte(ownershipId)}

	failedTestMessage := " | " + getOwnership + " with args: {" + string(args[0]) + ", "+ string(args[1]) + ", " + string(args[2]) + "}, failed. "

	handleExpectedSuccess(t, stub, getOwnership, failedTestMessage, args, payload)

}

func invokePropertyTransaction(t *testing.T, stub *shim.MockStub, propertyId string, payload string ){

	args := [][]byte{
		[]byte(propertyTransaction),
		[]byte(propertyTransaction),
		[]byte(propertyId),
		[]byte(payload)}

	failedTestMessage := " | " + propertyTransaction + " with args: {" + string(args[1]) + ", " + string(args[2]) + ", " + string(args[3]) + "}, failed. "

	response := stub.MockInvoke(propertyTransaction, args)
	if response.Status != shim.OK {
		message := failedTestMessage +  "[response.Status=" + strconv.FormatInt(int64(response.Status), 10) + "]"
		fmt.Println(message)
		t.FailNow()
	}

}