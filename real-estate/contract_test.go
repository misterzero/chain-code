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

func TestGetOwnershipMissingOwnership(t *testing.T) {

	stub := getStub()

	context := getSessionContext(t, stub)
	context.Arguments.MethodName = getOwnership
	context.Arguments.Payload = ownership_1
	context.Arguments.Builder = []string{context.Arguments.Payload}

	context.Test.ExpectedResponseMessage = nilValueForOwnershipId
	context.Test.ShouldFailTest = true

	invokeGetOwnership(context)

}

func TestGetOwnershipExtraArgs(t *testing.T) {

	stub := getStub()

	context := getSessionContext(t, stub)
	context.Arguments.MethodName = getOwnership
	context.Arguments.Payload = invalidArgument
	context.Arguments.Id = ownership_1
	context.Arguments.Builder = []string{context.Arguments.Id, context.Arguments.Payload}

	context.Test.ExpectedResponseMessage = incorrectNumberOfArgs
	context.Test.ShouldFailTest = true

	invokeGetOwnership(context)

}

func TestOwnershipCreatedDuringPropertyTransaction(t *testing.T) {

	stub := getStub()

	context := getSessionContext(t, stub)
	context.Arguments.MethodName = propertyTransaction
	context.Arguments.Id = property_1

	context.Arguments.Attributes = []Attribute{
		{Id: ownership_1, Percent:0.45},
		{Id: ownership_2, Percent:0.55}}

	context.Arguments.Payload = createProperty(context)
	context.Arguments.Builder = []string{context.Arguments.Id, context.Arguments.Payload}

	invokePropertyTransaction(context)

	context = getSessionContext(t, stub)
	context.Arguments.Id = ownership_1
	context.Arguments.MethodName = getOwnership

	property := []Attribute{{Id:"1", Percent:0.45, SaleDate:dateString}}
	context.Arguments.Payload = getAttributesAsString(property)
	context.Arguments.Builder = []string{context.Arguments.Id}

	context.Test.ExpectedResponsePayload = context.Arguments.Payload

	invokeGetOwnership(context)

}

func TestPropertyTransaction(t *testing.T) {

	stub := getStub()

	context := getSessionContext(t, stub)
	context.Arguments.MethodName = propertyTransaction
	context.Arguments.Id = property_1

	context.Arguments.Attributes = []Attribute{
		{Id: ownership_1, Percent:0.45},
		{Id: ownership_2, Percent:0.55}}

	context.Arguments.Payload = createProperty(context)
	context.Arguments.Builder = []string{context.Arguments.Id, context.Arguments.Payload}

	invokePropertyTransaction(context)

}

func TestMultiplePropertyTransactionsDifferentOwners(t *testing.T) {

	stub := getStub()

	context := getSessionContext(t, stub)
	context.Arguments.MethodName = propertyTransaction
	context.Arguments.Id = property_1

	context.Arguments.Attributes = []Attribute{
		{Id: ownership_1, Percent:0.45},
		{Id: ownership_2, Percent:0.55}}

	context.Arguments.Payload = createProperty(context)
	context.Arguments.Builder = []string{context.Arguments.Id, context.Arguments.Payload}

	invokePropertyTransaction(context)

	context = getSessionContext(t, stub)
	context.Arguments.MethodName = propertyTransaction
	context.Arguments.Id = property_1

	context.Arguments.Attributes = []Attribute{
		{Id: ownership_3, Percent:0.35},
		{Id: ownership_4, Percent:0.65}}

	context.Arguments.Payload = createProperty(context)
	context.Arguments.Builder = []string{context.Arguments.Id, context.Arguments.Payload}

	invokePropertyTransaction(context)
}

func TestMultiplePropertyTransactionsWithRepeatOwners(t *testing.T) {

	stub := getStub()

	context := getSessionContext(t, stub)
	context.Arguments.MethodName = propertyTransaction
	context.Arguments.Id = property_1

	context.Arguments.Attributes = []Attribute{
			{Id: ownership_1, Percent:0.45},
			{Id: ownership_2, Percent:0.55}}

	context.Arguments.Payload = createProperty(context)
	context.Arguments.Builder = []string{context.Arguments.Id, context.Arguments.Payload}

	invokePropertyTransaction(context)

	context = getSessionContext(t, stub)
	context.Arguments.MethodName = propertyTransaction
	context.Arguments.Id = property_1

	context.Arguments.Attributes = []Attribute{
		{Id: ownership_1, Percent:0.35},
		{Id: ownership_3, Percent:0.65}}

	context.Arguments.Payload = createProperty(context)
	context.Arguments.Builder = []string{context.Arguments.Id, context.Arguments.Payload}

	invokePropertyTransaction(context)

}

func TestPropertyTransactionExtraArgs(t *testing.T) {

	stub := getStub()

	context := getSessionContext(t, stub)
	context.Arguments.MethodName = getOwnership
	context.Arguments.Id = property_1

	context.Arguments.Attributes = []Attribute{
		{Id: ownership_1, Percent:0.35},
		{Id: ownership_3, Percent:0.65}}

	context.Arguments.Payload = createProperty(context)
	context.Arguments.Builder = []string{context.Arguments.Payload, context.Arguments.Id, context.Arguments.Payload, "extraArgs"}
	context.Test.ExpectedResponseMessage = incorrectNumberOfArgs
	context.Test.ShouldFailTest = true

	invokePropertyTransaction(context)
}

func TestPropertyTransactionStringAsSalePrice(t *testing.T) {

	stub := getStub()

	context := getSessionContext(t, stub)
	context.Arguments.MethodName = propertyTransaction
	context.Arguments.Payload = `{"saleDate":"2017-06-28T21:57:16","salePrice":"1000","owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.55}]}`
	context.Arguments.Id = property_1
	context.Test.ExpectedResponseMessage = cannotUnmarshalStringIntoFloat64
	context.Arguments.Builder = []string{context.Arguments.Id, context.Arguments.Payload}
	context.Test.ShouldFailTest = true

	invokePropertyTransaction(context)
}

func TestPropertyTransactionMissingSaleDate(t *testing.T) {

	stub := getStub()

	context := getSessionContext(t, stub)
	context.Arguments.MethodName = propertyTransaction
	context.Arguments.Payload = `{"salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownership_2","percentage":0.55}]}`
	context.Arguments.Id = property_1
	context.Test.ExpectedResponseMessage = saleDateRequired
	context.Arguments.Builder = []string{context.Arguments.Id, context.Arguments.Payload}
	context.Test.ShouldFailTest = true

	invokePropertyTransaction(context)

}

func TestPropertyTransactionNegativeSalePrice(t *testing.T) {

	stub := getStub()

	context := getSessionContext(t, stub)
	context.Arguments.MethodName = propertyTransaction
	context.Arguments.Payload = `{"saleDate":"2017-06-28T21:57:16","salePrice":-1,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.55}]}`
	context.Arguments.Id = property_1
	context.Test.ExpectedResponseMessage = salePriceMustBeGreaterThan0
	context.Arguments.Builder = []string{context.Arguments.Id, context.Arguments.Payload}
	context.Test.ShouldFailTest = true

	invokePropertyTransaction(context)
}

func TestPropertyTransactionNoOwners(t *testing.T) {

	stub := getStub()

	context := getSessionContext(t, stub)
	context.Arguments.MethodName = propertyTransaction
	context.Arguments.Payload = `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[]}`
	context.Arguments.Id = property_1
	context.Test.ExpectedResponseMessage = atLeastOneOwnerIsRequired
	context.Arguments.Builder = []string{context.Arguments.Id, context.Arguments.Payload}
	context.Test.ShouldFailTest = true

	invokePropertyTransaction(context)
}

func TestPropertyTransactionTooLowTotalOwnerPercentage(t *testing.T) {

	stub := getStub()

	context := getSessionContext(t, stub)
	context.Arguments.MethodName = propertyTransaction
	context.Arguments.Payload = `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.50}]}`
	context.Arguments.Id = property_1
	context.Test.ExpectedResponseMessage = totalPercentageCanNotBeGreaterThan1
	context.Arguments.Builder =[]string{context.Arguments.Id, context.Arguments.Payload}
	context.Test.ShouldFailTest = true

	invokePropertyTransaction(context)

}

func TestPropertyTransactionTooHighTotalOwnerPercentage(t *testing.T) {

	stub := getStub()

	context := getSessionContext(t, stub)
	context.Arguments.MethodName = propertyTransaction
	context.Arguments.Payload = `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.70}]}`
	context.Arguments.Id = property_1
	context.Arguments.Builder = []string{context.Arguments.Id, context.Arguments.Payload}

	context.Test.ExpectedResponseMessage = totalPercentageCanNotBeGreaterThan1
	context.Test.ShouldFailTest =true

	invokePropertyTransaction(context)

}

func TestGetProperty(t *testing.T) {

	stub := getStub()

	context := getSessionContext(t, stub)
	context.Arguments.MethodName = propertyTransaction
	context.Arguments.Id = property_1

	context.Arguments.Attributes = []Attribute{
		{Id: ownership_1, Percent:0.35},
		{Id: ownership_3, Percent:0.65}}

	propertyAsString := createProperty(context)

	context.Arguments.Payload = propertyAsString
	context.Arguments.Builder = []string{context.Arguments.Id, context.Arguments.Payload}

	invokePropertyTransaction(context)

	context = getSessionContext(t, stub)
	context.Arguments.Id = property_1
	context.Arguments.MethodName = getProperty
	context.Arguments.Payload = propertyAsString
	context.Arguments.Builder = []string{context.Arguments.Id}

	context.Test.ExpectedResponsePayload = context.Arguments.Payload

	invokeGetProperty(context)

}

func TestGetPropertyExtraArgs(t *testing.T) {

	stub := getStub()

	context := getSessionContext(t, stub)
	context.Arguments.Id = property_1
	context.Arguments.MethodName = getProperty
	context.Arguments.Payload = createProperty(context)
	context.Arguments.Builder = []string{context.Arguments.Id, context.Arguments.Payload}

	context.Test.ExpectedResponseMessage = incorrectNumberOfArgs
	context.Test.ShouldFailTest = true

	invokeGetProperty(context)

}

func TestGetPropertyMissingProperty(t *testing.T) {

	stub := getStub()

	context := getSessionContext(t, stub)
	context.Arguments.MethodName = getProperty
	context.Arguments.Payload = emptyPropertyJson
	context.Arguments.Id = property_1
	context.Arguments.Builder = []string{context.Arguments.Id, context.Arguments.Payload}

	context.Test.ExpectedResponseMessage = incorrectNumberOfArgs
	context.Test.ShouldFailTest = true

	invokeGetProperty(context)

}
