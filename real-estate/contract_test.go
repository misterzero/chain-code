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

// TODO make context class level var that will be reused
// TODO the context should not be updated as it is passed around, just used ... refactor opportunity
func TestGetOwnershipMissingOwnership(t *testing.T) {

	stub := getStub()

	context := SessionContext{}
	context.MethodName = getOwnership
	context.Payload = ownership_1
	context.Arguments = getChainCodeArgs(context.MethodName, context.Payload)
	context.ExpectedResponse = nilValueForOwnershipId

	handleExpectedFailures(t, stub, context)

}

func TestGetOwnershipExtraArgs(t *testing.T) {

	stub := getStub()

	context := SessionContext{}
	context.MethodName = getOwnership
	context.Payload = invalidArgument
	context.Id = ownership_1
	context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload)
	context.ExpectedResponse = incorrectNumberOfArgs

	handleExpectedFailures(t, stub, context)

}

func TestOwnershipCreatedDuringPropertyTransaction(t *testing.T) {

	stub := getStub()

	context := SessionContext{}
	context.MethodName = propertyTransaction
	context.Id = property_1

	context.Attributes = []Attribute{
		{Id: ownership_1, Percent:0.45},
		{Id: ownership_2, Percent:0.55}}

	confirmPropertyTransaction(t,stub, context)

	context = SessionContext{}
	context.MethodName = getOwnership
	context.Id = ownership_1

	property := []Attribute{{Id:"1", Percent:0.45, SaleDate:dateString}}
	context.Payload = getAttributesAsString(property)

	invokeGetOwnership(t, stub, context)

}

func TestPropertyTransaction(t *testing.T) {

	stub := getStub()

	context := SessionContext{}
	context.MethodName = propertyTransaction
	context.Id = property_1

	context.Attributes = []Attribute{
		{Id: ownership_1, Percent:0.45},
		{Id: ownership_2, Percent:0.55}}

	confirmPropertyTransaction(t, stub, context)
}

func TestMultiplePropertyTransactions(t *testing.T) {

	stub := getStub()

	context := SessionContext{}
	context.MethodName = propertyTransaction
	context.Id = property_1

	context.Attributes = []Attribute{
		{Id: ownership_1, Percent:0.45},
		{Id: ownership_2, Percent:0.55}}

	confirmPropertyTransaction(t, stub, context)

	context = SessionContext{}
	context.MethodName = propertyTransaction
	context.Id = property_1

	context.Attributes = []Attribute{
		{Id: ownership_3, Percent:0.35},
		{Id: ownership_4, Percent:0.65}}

	confirmPropertyTransaction(t, stub, context)

}

func TestMultiplePropertyTransactionsWithRepeatOwners(t *testing.T) {

	stub := getStub()

	context := SessionContext{}
	context.MethodName = propertyTransaction
	context.Id = property_1

	//owners := []Attribute{
	//	{Id: ownership_1, Percent:0.45},
	//	{Id: ownership_2, Percent:0.55}}

	context.Attributes = []Attribute{
			{Id: ownership_1, Percent:0.45},
			{Id: ownership_2, Percent:0.55}}

	//confirmPropertyTransaction(t, stub, owners)
	confirmPropertyTransaction(t, stub, context)

	context = SessionContext{}
	context.MethodName = propertyTransaction
	context.Id = property_1

	//owners = []Attribute{
	//	{Id: ownership_1, Percent:0.35},
	//	{Id: ownership_3, Percent:0.65}}

	context.Attributes = []Attribute{
		{Id: ownership_1, Percent:0.35},
		{Id: ownership_3, Percent:0.65}}

	//confirmPropertyTransaction(t, stub, owners)
	confirmPropertyTransaction(t, stub, context)

}

func TestPropertyTransactionExtraArgs(t *testing.T) {

	stub := getStub()

	context := SessionContext{}
	context.MethodName = getOwnership
	context.Id = property_1

	context.Attributes = []Attribute{
		{Id: ownership_1, Percent:0.35},
		{Id: ownership_3, Percent:0.65}}

	_, context.Payload = createProperty(context)

	context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload, "extraArgs")
	context.ExpectedResponse = incorrectNumberOfArgs

	handleExpectedFailures(t, stub, context)

}

func TestPropertyTransactionStringAsSalePrice(t *testing.T) {

	stub := getStub()

	stringAsSalePrice := `{"saleDate":"2017-06-28T21:57:16","salePrice":"1000","owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.55}]}`

	context := SessionContext{}
	context.MethodName = propertyTransaction
	context.Payload = stringAsSalePrice
	context.Id = property_1
	context.ExpectedResponse = cannotUnmarshalStringIntoFloat64
	context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload)

	handleExpectedFailures(t, stub, context)

}

func TestPropertyTransactionMissingSaleDate(t *testing.T) {

	stub := getStub()

	missingSaleDate := `{"salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownership_2","percentage":0.55}]}`

	context := SessionContext{}
	context.MethodName = propertyTransaction
	context.Payload = missingSaleDate
	context.Id = property_1
	context.ExpectedResponse = saleDateRequired
	context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload)

	handleExpectedFailures(t, stub, context)

}

func TestPropertyTransactionNegativeSalePrice(t *testing.T) {

	stub := getStub()

	negativeSalePrice := `{"saleDate":"2017-06-28T21:57:16","salePrice":-1,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.55}]}`

	context := SessionContext{}
	context.MethodName = propertyTransaction
	context.Payload = negativeSalePrice
	context.Id = property_1
	context.ExpectedResponse = salePriceMustBeGreaterThan0
	context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload)

	handleExpectedFailures(t, stub, context)

}

func TestPropertyTransactionNoOwners(t *testing.T) {

	stub := getStub()

	noOwners := `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[]}`

	context := SessionContext{}
	context.MethodName = propertyTransaction
	context.Payload = noOwners
	context.Id = property_1
	context.ExpectedResponse = atLeastOneOwnerIsRequired
	context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload)

	handleExpectedFailures(t, stub, context)

}

func TestPropertyTransactionTooLowTotalOwnerPercentage(t *testing.T) {

	stub := getStub()

	tooLowOwnerPercentage := `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.50}]}`

	context := SessionContext{}
	context.MethodName = propertyTransaction
	context.Payload = tooLowOwnerPercentage
	context.Id = property_1
	context.ExpectedResponse = totalPercentageCanNotBeGreaterThan1
	context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload)

	handleExpectedFailures(t, stub, context)

}

func TestPropertyTransactionTooHighTotalOwnerPercentage(t *testing.T) {

	stub := getStub()

	tooHighOwnerPercentage := `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.70}]}`

	context := SessionContext{}
	context.MethodName = propertyTransaction
	context.Payload = tooHighOwnerPercentage
	context.Id = property_1
	context.ExpectedResponse = totalPercentageCanNotBeGreaterThan1
	context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload)

	handleExpectedFailures(t, stub, context)

}

func TestGetProperty(t *testing.T) {

	stub := getStub()

	context := SessionContext{}
	context.MethodName = propertyTransaction
	context.Id = property_1

	context.Attributes = []Attribute{
		{Id: ownership_1, Percent:0.35},
		{Id: ownership_3, Percent:0.65}}

	_, context.Payload = createProperty(context)

	 confirmPropertyTransaction(t, stub, context)

}

func TestGetPropertyExtraArgs(t *testing.T) {

	stub := getStub()

	context := SessionContext{}
	context.MethodName = propertyTransaction
	context.Id = property_1

	context.Attributes = []Attribute{
		{Id: ownership_1, Percent:0.35},
		{Id: ownership_3, Percent:0.65}}

	_, context.Payload = createProperty(context)

	confirmPropertyTransaction(t, stub, context)

	context = SessionContext{}
	context.MethodName = getProperty
	context.Payload = "invalidProperty"
	context.Id = property_1
	context.ExpectedResponse = incorrectNumberOfArgs
	context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload)

	handleExpectedFailures(t, stub, context)

}

func TestGetPropertyMissingProperty(t *testing.T) {

	stub := getStub()

	context := SessionContext{}
	context.MethodName = getProperty
	context.Payload = emptyPropertyJson
	context.Id = property_1
	context.ExpectedResponse = incorrectNumberOfArgs
	context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload)

	handleExpectedFailures(t, stub, context)

}
