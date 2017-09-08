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

	handleExpectedFailures(t, stub, getOwnership, args, payload, nilValueForOwnershipId)

}

func TestGetOwnershipExtraArgs(t *testing.T){

	stub := getStub()

	payload := invalidArgument
	args := [][]byte{
		[]byte(getOwnership),
		[]byte(getOwnership),
		[]byte(ownership_1),
		[]byte(payload)}

	handleExpectedFailures(t, stub, getOwnership, args, invalidArgument, incorrectNumberOfArgs)

}

func TestOwnershipCreatedDuringPropertyTransaction(t *testing.T){

	stub := getStub()

	var property = Property{}

	owners := getValidOwnersList(ownership_1, 0.45, ownership_2, 0.55, dateString)

	property, propertyAsString := createProperty(property_1, dateString, 1000, owners)

	propertyAsBytes := stub.State[property.PropertyId]
	if propertyAsBytes != nil{
		fmt.Println(property_1 + " should not exist.")
		t.Fail()
	}

	invokePropertyTransaction(t, stub, property.PropertyId, propertyAsString)

	propertyId := "1"
	propertyOwnersList := getValidPropertyListForOwner(propertyId)

	ownershipPropertyAsString := getAttributesAsString([]Attribute{propertyOwnersList[0]})

	invokeGetOwnership(t, stub, owners[0].Id, ownershipPropertyAsString)

}

//TODO next
func TestPropertyTransaction(t *testing.T){

	stub := getStub()

	ownerList := getValidOwnersList(ownership_1, 0.45, ownership_2, 0.55, dateString)
	property, propertyString := createProperty(property_1, dateString, 1000, ownerList)

	invokePropertyTransaction(t,stub, property.PropertyId, propertyString)
	checkPropertyState(t, stub, property, propertyString)

}

func TestMultiplePropertyTransactions(t *testing.T){

	stub := getStub()

	ownerList := getValidOwnersList(ownership_1, 0.45, ownership_2, 0.55, dateString)
	property, propertyString := createProperty(property_1, dateString, 1000, ownerList)

	invokePropertyTransaction(t,stub, property.PropertyId, propertyString)
	checkPropertyState(t, stub, property, propertyString)

	ownerList = getValidOwnersList(ownership_3, 0.35, ownership_4, 0.65, dateString)
	property, propertyString = createProperty(property_1, dateString, 1000, ownerList)

	invokePropertyTransaction(t,stub, property.PropertyId, propertyString)
	checkPropertyState(t, stub, property, propertyString)

}

func TestMultiplePropertyTransactionsWithRepeatOwners(t *testing.T){

	stub := getStub()

	ownerList := getValidOwnersList(ownership_1, 0.45, ownership_2, 0.55, dateString)
	property, propertyString := createProperty(property_1, dateString, 1000, ownerList)

	invokePropertyTransaction(t,stub, property.PropertyId, propertyString)
	checkPropertyState(t, stub, property, propertyString)

	ownerList = getValidOwnersList(ownership_1, 0.35, ownership_3, 0.65, dateString)
	property, propertyString = createProperty(property_1, dateString, 1000, ownerList)

	invokePropertyTransaction(t,stub, property.PropertyId, propertyString)
	checkPropertyState(t, stub, property, propertyString)

}

func TestPropertyTransactionExtraArgs(t *testing.T) {

	stub := getStub()

	ownerList := getValidOwnersList(ownership_1, 0.35, ownership_3, 0.65, dateString)
	_, propertyAsString := createProperty(property_1, dateString, 1000, ownerList)

	args := [][]byte{
		[]byte(propertyTransaction),
		[]byte(propertyTransaction),
		[]byte(property_1),
		[]byte(propertyAsString),
		[]byte("extraArg")}

	handleExpectedFailures(t, stub, propertyTransaction, args, propertyAsString, incorrectNumberOfArgs)

}

func TestPropertyTransactionStringAsSalePrice(t *testing.T) {

	stub := getStub()

	stringAsSalePrice := `{"saleDate":"2017-06-28T21:57:16","salePrice":"1000","owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.55}]}`

	args := [][]byte{
		[]byte(propertyTransaction),
		[]byte(propertyTransaction),
		[]byte(property_1),
		[]byte(stringAsSalePrice)}

	handleExpectedFailures(t, stub, propertyTransaction, args, stringAsSalePrice, cannotUnmarshalStringIntoFloat64)

}

func TestPropertyTransactionMissingSaleDate(t *testing.T) {

	stub := getStub()

	missingSaleDateJson := `{"salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownership_2","percentage":0.55}]}`

	args := [][]byte{
		[]byte(propertyTransaction),
		[]byte(propertyTransaction),
		[]byte(property_1),
		[]byte(missingSaleDateJson)}

	handleExpectedFailures(t, stub, propertyTransaction, args, missingSaleDateJson, saleDateRequired)

}

func TestPropertyTransactionNegativeSalePrice(t *testing.T) {

	stub := getStub()

	negativeSalePriceJson := `{"saleDate":"2017-06-28T21:57:16","salePrice":-1,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.55}]}`

	args := [][]byte{
		[]byte(propertyTransaction),
		[]byte(propertyTransaction),
		[]byte(property_1),
		[]byte(negativeSalePriceJson)}

	handleExpectedFailures(t, stub, propertyTransaction, args, negativeSalePriceJson, salePriceMustBeGreaterThan0)

}

func TestPropertyTransactionNoOwners(t *testing.T) {

	stub := getStub()

	noOwnersJson := `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[]}`

	args := [][]byte{
		[]byte(propertyTransaction),
		[]byte(propertyTransaction),
		[]byte(property_1),
		[]byte(noOwnersJson)}

	handleExpectedFailures(t, stub, propertyTransaction, args, noOwnersJson, atLeastOneOwnerIsRequired)

}

func TestPropertyTransactionTooLowTotalOwnerPercentage(t *testing.T) {

	stub := getStub()

	tooLowOwnerPercentage := `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.50}]}`

	args := [][]byte{
		[]byte(propertyTransaction),
		[]byte(propertyTransaction),
		[]byte(property_1),
		[]byte(tooLowOwnerPercentage)}

	handleExpectedFailures(t, stub, propertyTransaction, args, tooLowOwnerPercentage, totalPercentageCanNotBeGreaterThan1)

}

func TestPropertyTransactionTooHighTotalOwnerPercentage(t *testing.T) {

	stub := getStub()

	tooHighOwnerPercentage := `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.70}]}`

	args := [][]byte{
		[]byte(propertyTransaction),
		[]byte(propertyTransaction),
		[]byte(property_1),
		[]byte(tooHighOwnerPercentage)}

	handleExpectedFailures(t, stub, propertyTransaction, args, tooHighOwnerPercentage, totalPercentageCanNotBeGreaterThan1)

}

func TestGetProperty(t *testing.T){

	stub := getStub()

	ownerList := getValidOwnersList(ownership_1, 0.35, ownership_3, 0.65, dateString)
	property, propertyString := createProperty(property_1, dateString, 1000, ownerList)

	invokePropertyTransaction(t, stub, property.PropertyId, propertyString)
	checkGetProperty(t, stub, property, propertyString)

}

func TestGetPropertyExtraArgs(t *testing.T){

	stub := getStub()

	ownerList := getValidOwnersList(ownership_1, 0.35, ownership_3, 0.65, dateString)
	property, propertyString := createProperty(property_1, dateString, 1000, ownerList)

	invokePropertyTransaction(t, stub, property.PropertyId, propertyString)

	args := [][]byte{
		[]byte(getProperty),
		[]byte(getProperty),
		[]byte(property.PropertyId),
		[]byte(propertyString)}

	handleExpectedFailures(t, stub, getProperty, args, propertyString, incorrectNumberOfArgs)

}

func TestGetPropertyMissingProperty(t *testing.T){

	stub := getStub()

	args := [][]byte{
		[]byte(getProperty),
		[]byte(getProperty),
		[]byte(property_1)}

	handleExpectedFailures(t, stub, getProperty, args, emptyPropertyJson, nilAmountFor)

}

func checkGetProperty(t *testing.T, stub *shim.MockStub, property Property, propertyString string){

	args := [][]byte{
		[]byte(getProperty),
		[]byte(getProperty),
		[]byte(property.PropertyId)}

	handleExpectedSuccess(t, stub, getProperty, args, propertyString)

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

	handleExpectedSuccess(t, stub, getOwnership, args, payload)

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

func handleExpectedSuccess(t *testing.T, stub *shim.MockStub, argument string, args [][]byte, payload string){

	response := stub.MockInvoke(argument, args)

	msg := "| FAIL [{args}, {<response-failure>}] | [{" + argument + ", " + payload + "}, "

	if response.Status != shim.OK {
		msg += "{response.Status=" + strconv.FormatInt(int64(response.Status), 10) + "}]"
		fmt.Println(msg)
		t.FailNow()
	}
	if response.Payload == nil {
		msg += "{response.Message=" + string(response.Message) + "}]"
		fmt.Println(msg)
		t.FailNow()
	}
	if string(response.Payload) != payload {
		msg += "{response.Payload=" + string(response.Payload) + "}]"
		fmt.Println(msg)
		t.FailNow()
	}

}

func handleExpectedFailures(t *testing.T, stub *shim.MockStub, argument string, args [][]byte, payload string, expectedResponseMessage string){

	response := stub.MockInvoke(argument, args)

	msg := "| FAIL [{args}, {<response-failure>}] | [{" + argument + ", " + payload + "}, "

	if response.Status != 500 {
		msg += "{response.Status=" + strconv.FormatInt(int64(response.Status), 10) + "}]"
		fmt.Println(msg)
		t.FailNow()
	}
	if !strings.Contains(response.Message, expectedResponseMessage) {
		msg += "{response.Message=" + string(response.Message) + "}]"
		fmt.Println(msg)
		t.FailNow()
	}
	if string(response.Payload) == payload {
		msg += "{response.Payload=" + string(response.Payload) + "}]"
		fmt.Println(msg)
		t.FailNow()
	}

}


