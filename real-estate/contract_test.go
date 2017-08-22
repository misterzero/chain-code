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
	"encoding/json"
	"errors"
)

const getOwnership = "getOwnership"
const propertyTransaction = "propertyTransaction"
const getProperty = "getProperty"
const ownership_1 = "ownership_1"
const property_1 = "property_1"
const dateString = `"2017-06-28T21:57:16"`
const emptyOwnershipPropertyJson = `{"properties":[]}`
const emptyPropertyJson = `{"saleDate":"","salePrice":0,"owners":[]}`
const errorStatus = int32(500)

const incorrectArgCountError = "Incorrect number of arguments: "
const getOwnershipMissingError = "Nil value for ownershipId:"
const getPropertyNilError = "Nil amount for"
const createPropertyTransactionUnmarshalWrongTypeError = "cannot unmarshal string into Go struct field Property.salePrice of type float64"
const createPropertyTransactionMissingSaleDateError = "A sale date is required."
const createPropertyTransactionGreaterThanZeroSalePriceError = "The sale price must be greater than 0"
const createPropertyTransactionTotalPercentageOfOneError = "Total Percentage can not be greater than or less than 1. Your total percentage ="
const createPropertyTransactionNoOwnersError = "At least one owner is required."

func TestGetOwnershipMissingOwnership(t *testing.T){

	stub := getStub()

	invalidArgs := getThreeArgs(getOwnership, ownership_1)
	message := " | " + getOwnership + " with args: {" + string(invalidArgs[1]) + "}, did not fail. "

	handleExpectedFailures(t, stub, getOwnership, message, getOwnershipMissingError, invalidArgs, emptyOwnershipPropertyJson)
}

func TestGetOwnershipExtraArgs(t *testing.T){

	stub := getStub()

	invalidArgument := "invalidArgument"

	invalidArgs := getFourArgs(getOwnership, ownership_1, invalidArgument)
	message := " | " + getOwnership + " with args: {" + string(invalidArgs[0]) + ", " + string(invalidArgs[1]) + ", " + string(invalidArgs[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, getOwnership, message, incorrectArgCountError, invalidArgs, invalidArgument)

}

func TestOwnershipCreatedDuringPropertyTransaction(t *testing.T){

	stub := getStub()

	var property = Property{}

	ownershipInputList := getValidOwners()

	property, propertyString := getTestProperty(property_1, dateString, 1000, ownershipInputList)

	propertyAsBytes := stub.State[property.PropertyId]
	if propertyAsBytes != nil{
		t.Fail()
	}

	checkPropertyTransaction(t, stub, property.PropertyId, propertyString)

	ownershipProperty := getAttributesAsString([]Attribute{getValidOwnersProperty1()[0]})

	checkGetOwnership(t, stub, ownershipInputList[0].Id, ownershipProperty)

}

func TestPropertyTransaction(t *testing.T){

	stub := getStub()

	property, propertyString := getTestProperty(property_1, dateString, 1000, getValidOwners())

	checkPropertyTransaction(t,stub, property.PropertyId, propertyString)
	checkPropertyState(t, stub, property, propertyString)

}

func TestPropertyTransactionExtraArgs(t *testing.T) {

	stub := getStub()

	invalidArgument := "invalidArgument"
	_, propertyString := getTestProperty(property_1, dateString, 1000, getValidOwners())

	invalidArgs := getFourArgs(propertyTransaction, property_1, propertyString)
	invalidArgs = append(invalidArgs, []byte(invalidArgument))

	message := " | " + propertyTransaction + " with args: {" + string(invalidArgs[1]) + ", " + string(invalidArgs[2]) + ", " + string(invalidArgs[3]) + ", " + string(invalidArgs[4]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, message, incorrectArgCountError, invalidArgs, propertyString)

}

func TestPropertyTransactionStringAsSalePrice(t *testing.T) {

	stub := getStub()

	validJson := `{"saleDate":"2017-06-28T21:57:16","salePrice":"1000","owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.55}]}`

	invalidArgs := getFourArgs(propertyTransaction, property_1, validJson)

	message := " | " + propertyTransaction + " with args: {" + string(invalidArgs[1]) + "," + string(invalidArgs[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, message, createPropertyTransactionUnmarshalWrongTypeError, invalidArgs, validJson )

}

func TestPropertyTransactionMissingSaleDate(t *testing.T) {

	stub := getStub()

	missingSaleDateJson := `{"salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownership_2","percentage":0.55}]}`

	invalidArgs := getFourArgs(propertyTransaction, property_1, missingSaleDateJson)
	message := " | " + propertyTransaction + " with args: {" + string(invalidArgs[1]) + "," + string(invalidArgs[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, message, createPropertyTransactionMissingSaleDateError, invalidArgs, missingSaleDateJson)

}

func TestPropertyTransactionNegativeSalePrice(t *testing.T) {

	stub := getStub()

	negativeSalePriceJson := `{"saleDate":"2017-06-28T21:57:16","salePrice":-1,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.55}]}`

	invalidArgs := getFourArgs(propertyTransaction, property_1, negativeSalePriceJson)
	message := " | " + propertyTransaction + " with args: {" + string(invalidArgs[1]) + "," + string(invalidArgs[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, message, createPropertyTransactionGreaterThanZeroSalePriceError, invalidArgs, negativeSalePriceJson)

}

func TestPropertyTransactionNoOwners(t *testing.T) {

	stub := getStub()

	noOwnersJson := `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[]}`
	
	invalidArgs := getFourArgs(propertyTransaction, property_1, noOwnersJson)
	message:= " | " + propertyTransaction + " with args: {" + string(invalidArgs[1]) + "," + string(invalidArgs[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, message, createPropertyTransactionNoOwnersError, invalidArgs, noOwnersJson)

}

func TestPropertyTransactionTooLowTotalOwnerPercentage(t *testing.T) {

	stub := getStub()

	tooLowOwnerPercentage := `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.50}]}`

	invalidArgs := getFourArgs(propertyTransaction, property_1, tooLowOwnerPercentage)
	message := " | " + propertyTransaction + " with args: {" + string(invalidArgs[1]) + "," + string(invalidArgs[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, message, createPropertyTransactionTotalPercentageOfOneError, invalidArgs, tooLowOwnerPercentage)

}

func TestPropertyTransactionTooHighTotalOwnerPercentage(t *testing.T) {

	stub := getStub()

	tooHighOwnerPercentage := `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.70}]}`

	invalidArgs := getFourArgs(propertyTransaction, property_1, tooHighOwnerPercentage)
	message := " | " + propertyTransaction + " with args: {" + string(invalidArgs[1]) + "," + string(invalidArgs[2]) + "}, did not fail. "

	handleExpectedFailures(t, stub, propertyTransaction, message, createPropertyTransactionTotalPercentageOfOneError, invalidArgs, tooHighOwnerPercentage)

}

func TestGetProperty(t *testing.T){

	stub := getStub()

	property, propertyString := getTestProperty(property_1, dateString, 1000, getValidOwners())

	checkPropertyTransaction(t, stub, property.PropertyId, propertyString)
	checkGetProperty(t, stub, property, propertyString)

}

func TestGetPropertyExtraArgs(t *testing.T){

	stub := getStub()

	property, propertyString := getTestProperty(property_1, dateString, 1000, getValidOwners())

	checkPropertyTransaction(t, stub, property.PropertyId, propertyString)

	invalidArgs := getFourArgs(getProperty, property.PropertyId, propertyString)

	message:= " | " + getProperty + " with args: {" + string(invalidArgs[1]) + "," + string(invalidArgs[2]) + "," + string(invalidArgs[3]) + "}, did not fail. "

	handleExpectedFailures(t, stub, getProperty, message, incorrectArgCountError, invalidArgs, propertyString)

}

func TestGetPropertyMissingProperty(t *testing.T){

	stub := getStub()

	invalidArgs := getThreeArgs(getProperty, property_1)
	message := " | " + getProperty + " with args: {" + property_1 + "}, did not fail. "

	handleExpectedFailures(t, stub, getProperty, message, getPropertyNilError, invalidArgs, emptyPropertyJson)

}

//====================================================

func checkGetOwnership(t *testing.T, stub *shim.MockStub,ownershipId string, ownershipString string){

	ownershipArgs := getThreeArgs(getOwnership, ownershipId)

	message:= " | " + getOwnership + " with args: {" + string(ownershipArgs[0]) + ", "+ string(ownershipArgs[1]) + ", " + string(ownershipArgs[2]) + "}, failed. "

	handleExpectedSuccess(t, stub, getOwnership, message, ownershipArgs, ownershipString)

}

func checkPropertyTransaction(t *testing.T, stub *shim.MockStub, propertyId string, propertyString string ){

	propertyArgs := getFourArgs(propertyTransaction, propertyId, propertyString)

	message:= " | " + propertyTransaction + " with args: {" + string(propertyArgs[1]) + ", " + string(propertyArgs[2]) + ", " + string(propertyArgs[3]) + "}, failed. "

	res := stub.MockInvoke(propertyTransaction, propertyArgs)
	if res.Status != shim.OK {
		message := message +  "[res.Status=" + strconv.FormatInt(int64(res.Status), 10) + "]"
		fmt.Println(message)
		t.FailNow()
	}

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

func checkGetProperty(t *testing.T, stub *shim.MockStub, property Property, propertyString string){

	propertyArgs := getThreeArgs(getProperty, property.PropertyId)

	message:= " | " + getProperty + " with args: {" + string(propertyArgs[1]) + "}, failed. "

	handleExpectedSuccess(t, stub, getProperty, message, propertyArgs, propertyString)

}

func getTestProperty(propertyId string, saleDate string, salePrice float64, owners []Attribute) (Property, string) {

	property := Property{}
	property.PropertyId = propertyId
	property.SaleDate = saleDate
	property.SalePrice = salePrice
	property.Owners = owners

	propertyAsBytes, _ := getPropertyAsBytes(property)

	return property, string(propertyAsBytes)

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

func getFourArgs(function string, id string, jsonStruct string) ([][]byte) {

	method := []byte(function)
	key := []byte(id)
	value := []byte(jsonStruct)

	args := [][]byte{method, method, key, value}

	return args
}

func getThreeArgs(function string, id string) ([][]byte) {

	method := []byte(function)
	key := []byte(id)

	args := [][]byte{method, method, key}

	return args
}

func getValidOwners() []Attribute {

	owner1 := Attribute{}
	owner2 := Attribute{}
	owner1.Id = "ownership_3"
	owner1.Percent = 0.45
	owner1.SaleDate = dateString

	owner2.Id = "ownership_2"
	owner2.Percent = 0.55
	owner2.SaleDate = dateString

	ownershipInputList := []Attribute{owner1, owner2}

	return ownershipInputList

}

func getValidOwnersProperty1() []Attribute{

	property1 := Attribute{}
	property2 := Attribute{}

	property1.Id = "1"
	property1.Percent = 0.45
	property1.SaleDate = dateString

	property2.Id = "1"
	property2.Percent = 0.55
	property2.SaleDate = dateString

	ownershipInputList := []Attribute{property1, property2}

	return ownershipInputList

}

func getAttributeAsString(attribute Attribute) (string, error){

	var attributeBytes []byte
	var err error

	attributeBytes, err = json.Marshal(attribute)
	if err != nil{
		err = errors.New("Unable to convert property to json string " + string(attributeBytes))
	}

	return string(attributeBytes), err

}

func getAttributesAsString(attributes []Attribute) string{

	var buffer bytes.Buffer

	buffer.WriteString("[")

	for i := 0; i< len(attributes); i++ {

		currentAttribute, _ := getAttributeAsString(attributes[i])
		buffer.WriteString(currentAttribute)

		if i != (len(attributes) - 1){
			buffer.WriteString(",")
		}

	}

	buffer.WriteString("]")

	return buffer.String()

}

func getStub() (*shim.MockStub){

	scc := new(Chaincode)
	stub := shim.NewMockStub("contract", scc)

	return stub

}