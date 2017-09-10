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

// TODO the context should not be updated as it is passed around, just used ... refactor opportunity
func TestGetOwnershipMissingOwnership(t *testing.T) {

	stub := getStub()

	context := getTestContext(t, stub)
	context.MethodName = getOwnership
	context.Payload = ownership_1
	//context.Arguments = getChainCodeArgs(context.MethodName, context.Payload)
	context.ArgumentBuilder = []string{context.Payload}
	context.Arguments = getChainCodeArgs2(context)
	//context.Arguments = getChainCodeArgs(context.MethodName, context.Payload)
	context.ExpectedResponse = nilValueForOwnershipId

	handleExpectedFailures(context)

}

func TestGetOwnershipExtraArgs(t *testing.T) {

	stub := getStub()

	context := getTestContext(t, stub)
	context.MethodName = getOwnership
	context.Payload = invalidArgument
	context.Id = ownership_1
	//context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload)
	context.ArgumentBuilder = []string{context.Id, context.Payload}
	context.Arguments = getChainCodeArgs2(context)
	context.ExpectedResponse = incorrectNumberOfArgs

	handleExpectedFailures(context)

}

func TestOwnershipCreatedDuringPropertyTransaction(t *testing.T) {

	stub := getStub()

	context := getTestContext(t, stub)
	context.MethodName = propertyTransaction
	context.Id = property_1

	context.Attributes = []Attribute{
		{Id: ownership_1, Percent:0.45},
		{Id: ownership_2, Percent:0.55}}

	context.Payload = createProperty(context)
	context.ArgumentBuilder = []string{context.Id, context.Payload}

	context.Arguments = getChainCodeArgs2(context)

	context.Payload = ""

	handleExpectedSuccess(context)

	context = getTestContext(t, stub)
	context.Id = ownership_1
	context.MethodName = getOwnership

	property := []Attribute{{Id:"1", Percent:0.45, SaleDate:dateString}}
	context.Payload = getAttributesAsString(property)

	context.ArgumentBuilder = []string{context.Id}
	context.Arguments = getChainCodeArgs2(context)

	handleExpectedSuccess(context)

}

func TestPropertyTransaction(t *testing.T) {

	stub := getStub()

	context := getTestContext(t, stub)
	context.MethodName = propertyTransaction
	context.Id = property_1

	context.Attributes = []Attribute{
		{Id: ownership_1, Percent:0.45},
		{Id: ownership_2, Percent:0.55}}

	context.Payload = createProperty(context)
	context.ArgumentBuilder = []string{context.Id, context.Payload}
	context.Arguments = getChainCodeArgs2(context)

	//TODO this is hacky and needs to be fixed
	context.Payload = ""

	handleExpectedSuccess(context)
}

func TestMultiplePropertyTransactions(t *testing.T) {

	stub := getStub()

	context := getTestContext(t, stub)
	context.MethodName = propertyTransaction
	context.Id = property_1

	context.Attributes = []Attribute{
		{Id: ownership_1, Percent:0.45},
		{Id: ownership_2, Percent:0.55}}

	//confirmPropertyTransaction(context)
	context.Payload = createProperty(context)
	context.ArgumentBuilder = []string{context.Id, context.Payload}
	context.Arguments = getChainCodeArgs2(context)

	//TODO this is hacky and needs to be fixed
	context.Payload = ""

	handleExpectedSuccess(context)

	context = getTestContext(t, stub)
	context.MethodName = propertyTransaction
	context.Id = property_1

	context.Attributes = []Attribute{
		{Id: ownership_3, Percent:0.35},
		{Id: ownership_4, Percent:0.65}}

	//confirmPropertyTransaction(context)
	context.Payload = createProperty(context)
	context.ArgumentBuilder = []string{context.Id, context.Payload}
	context.Arguments = getChainCodeArgs2(context)

	//TODO this is hacky and needs to be fixed
	context.Payload = ""

	handleExpectedSuccess(context)

}

func TestMultiplePropertyTransactionsWithRepeatOwners(t *testing.T) {

	stub := getStub()

	context := getTestContext(t, stub)
	context.MethodName = propertyTransaction
	context.Id = property_1

	context.Attributes = []Attribute{
			{Id: ownership_1, Percent:0.45},
			{Id: ownership_2, Percent:0.55}}

	//confirmPropertyTransaction(context)
	context.Payload = createProperty(context)
	context.ArgumentBuilder = []string{context.Id, context.Payload}
	context.Arguments = getChainCodeArgs2(context)

	//TODO this is hacky and needs to be fixed
	context.Payload = ""

	handleExpectedSuccess(context)

	context = getTestContext(t, stub)
	context.MethodName = propertyTransaction
	context.Id = property_1

	context.Attributes = []Attribute{
		{Id: ownership_1, Percent:0.35},
		{Id: ownership_3, Percent:0.65}}

	//confirmPropertyTransaction(context)
	context.Payload = createProperty(context)
	context.ArgumentBuilder = []string{context.Id, context.Payload}
	context.Arguments = getChainCodeArgs2(context)

	//TODO this is hacky and needs to be fixed
	context.Payload = ""

	handleExpectedSuccess(context)

}

func TestPropertyTransactionExtraArgs(t *testing.T) {

	stub := getStub()

	context := getTestContext(t, stub)
	context.MethodName = getOwnership
	context.Id = property_1

	context.Attributes = []Attribute{
		{Id: ownership_1, Percent:0.35},
		{Id: ownership_3, Percent:0.65}}

	context.Payload = createProperty(context)
	//context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload, "extraArgs")
	context.ArgumentBuilder = []string{context.Payload, context.Id, context.Payload, "extraArgs"}
	context.Arguments = getChainCodeArgs2(context)
	context.ExpectedResponse = incorrectNumberOfArgs


	handleExpectedFailures(context)

}

func TestPropertyTransactionStringAsSalePrice(t *testing.T) {

	stub := getStub()

	context := getTestContext(t, stub)
	context.MethodName = propertyTransaction
	context.Payload = `{"saleDate":"2017-06-28T21:57:16","salePrice":"1000","owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.55}]}`
	context.Id = property_1
	context.ExpectedResponse = cannotUnmarshalStringIntoFloat64
	//context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload)
	context.ArgumentBuilder = []string{context.Id, context.Payload}
	context.Arguments = getChainCodeArgs2(context)

	handleExpectedFailures(context)

}

func TestPropertyTransactionMissingSaleDate(t *testing.T) {

	stub := getStub()

	context := getTestContext(t, stub)
	context.MethodName = propertyTransaction
	context.Payload = `{"salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownership_2","percentage":0.55}]}`
	context.Id = property_1
	context.ExpectedResponse = saleDateRequired
	//context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload)
	context.ArgumentBuilder = []string{context.Id, context.Payload}
	context.Arguments = getChainCodeArgs2(context)

	handleExpectedFailures(context)

}

func TestPropertyTransactionNegativeSalePrice(t *testing.T) {

	stub := getStub()

	context := getTestContext(t, stub)
	context.MethodName = propertyTransaction
	context.Payload = `{"saleDate":"2017-06-28T21:57:16","salePrice":-1,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.55}]}`
	context.Id = property_1
	context.ExpectedResponse = salePriceMustBeGreaterThan0
	//context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload)
	context.ArgumentBuilder = []string{context.Id, context.Payload}
	context.Arguments = getChainCodeArgs2(context)

	handleExpectedFailures(context)

}

func TestPropertyTransactionNoOwners(t *testing.T) {

	stub := getStub()

	context := getTestContext(t, stub)
	context.MethodName = propertyTransaction
	context.Payload = `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[]}`
	context.Id = property_1
	context.ExpectedResponse = atLeastOneOwnerIsRequired
	//context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload)
	context.ArgumentBuilder = []string{context.Id, context.Payload}
	context.Arguments = getChainCodeArgs2(context)

	handleExpectedFailures(context)

}

func TestPropertyTransactionTooLowTotalOwnerPercentage(t *testing.T) {

	stub := getStub()

	context := getTestContext(t, stub)
	context.MethodName = propertyTransaction
	context.Payload = `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.50}]}`
	context.Id = property_1
	context.ExpectedResponse = totalPercentageCanNotBeGreaterThan1
	//context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload)
	context.ArgumentBuilder =[]string{context.Id, context.Payload}
	context.Arguments = getChainCodeArgs2(context)

	handleExpectedFailures(context)

}

func TestPropertyTransactionTooHighTotalOwnerPercentage(t *testing.T) {

	stub := getStub()

	context := getTestContext(t, stub)
	context.MethodName = propertyTransaction
	context.Payload = `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.70}]}`
	context.Id = property_1
	context.ExpectedResponse = totalPercentageCanNotBeGreaterThan1
	//context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload)
	context.ArgumentBuilder = []string{context.Id, context.Payload}
	context.Arguments = getChainCodeArgs2(context)

	handleExpectedFailures(context)

}

func TestGetProperty(t *testing.T) {

	stub := getStub()

	context := getTestContext(t, stub)
	context.MethodName = propertyTransaction
	context.Id = property_1

	context.Attributes = []Attribute{
		{Id: ownership_1, Percent:0.35},
		{Id: ownership_3, Percent:0.65}}

	propertyAsString := createProperty(context)

	context.Payload = propertyAsString
	context.ArgumentBuilder = []string{context.Id, context.Payload}
	context.Arguments = getChainCodeArgs2(context)

	//TODO this is hacky and needs to be fixed
	context.Payload = ""

	handleExpectedSuccess(context)

	context = getTestContext(t, stub)
	context.Id = property_1
	context.MethodName = getProperty
	context.Payload = propertyAsString
	context.ArgumentBuilder = []string{context.Id}
	context.Arguments = getChainCodeArgs2(context)

	handleExpectedSuccess(context)

}

func TestGetPropertyExtraArgs(t *testing.T) {

	stub := getStub()

	context := getTestContext(t, stub)
	context.MethodName = propertyTransaction
	context.Id = property_1

	context.Attributes = []Attribute{
		{Id: ownership_1, Percent:0.35},
		{Id: ownership_3, Percent:0.65}}

	propertyAsString := createProperty(context)

	context.Payload = propertyAsString
	context.ArgumentBuilder = []string{context.Id, context.Payload}
	context.Arguments = getChainCodeArgs2(context)

	//confirmPropertyTransaction(context)

	//TODO this is hacky and needs to be fixed
	context.Payload = ""

	handleExpectedSuccess(context)

	context = getTestContext(t, stub)
	context.Id = property_1
	context.MethodName = getProperty
	context.Payload = propertyAsString
	context.ExpectedResponse = incorrectNumberOfArgs
	context.ArgumentBuilder = []string{context.Id, context.Payload}
	context.Arguments = getChainCodeArgs2(context)

	handleExpectedFailures(context)

}

func TestGetPropertyMissingProperty(t *testing.T) {

	stub := getStub()

	context := getTestContext(t, stub)
	context.MethodName = getProperty
	context.Payload = emptyPropertyJson
	context.Id = property_1
	context.ExpectedResponse = incorrectNumberOfArgs
	//context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload)
	context.ArgumentBuilder = []string{context.Id, context.Payload}
	context.Arguments = getChainCodeArgs2(context)

	handleExpectedFailures(context)

}
