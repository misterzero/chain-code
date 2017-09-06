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
)

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


	ownershipInputList := getMockValidOwners()

	property, propertyString := getMockProperty(property_1, dateString, 1000, ownershipInputList)

	propertyAsBytes := stub.State[property.PropertyId]
	if propertyAsBytes != nil{
		t.Fail()
	}

	checkPropertyTransaction(t, stub, property.PropertyId, propertyString)

	ownershipProperty := getAttributesAsString([]Attribute{getMockValidOwnersProperty()[0]})

	checkGetOwnership(t, stub, ownershipInputList[0].Id, ownershipProperty)

}

func TestPropertyTransaction(t *testing.T){

	stub := getStub()

	property, propertyString := getMockProperty(property_1, dateString, 1000, getMockValidOwners())

	checkPropertyTransaction(t,stub, property.PropertyId, propertyString)
	checkPropertyState(t, stub, property, propertyString)

}

func TestPropertyTransactionExtraArgs(t *testing.T) {

	stub := getStub()

	invalidArgument := "invalidArgument"
	_, propertyString := getMockProperty(property_1, dateString, 1000, getMockValidOwners())

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

	property, propertyString := getMockProperty(property_1, dateString, 1000, getMockValidOwners())

	checkPropertyTransaction(t, stub, property.PropertyId, propertyString)
	checkGetProperty(t, stub, property, propertyString)

}

func TestGetPropertyExtraArgs(t *testing.T){

	stub := getStub()

	property, propertyString := getMockProperty(property_1, dateString, 1000, getMockValidOwners())

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
