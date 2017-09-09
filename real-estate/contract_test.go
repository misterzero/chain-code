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

//TODO code needs to be refactored to deal with the name field of Attribute better
func TestGetOwnershipMissingOwnership(t *testing.T){

	stub := getStub()

	payload := ownership_1
	args := getChainCodeArgs(getOwnership, ownership_1)

	handleExpectedFailures(t, stub, args, payload, getOwnership, nilValueForOwnershipId)

}

func TestGetOwnershipExtraArgs(t *testing.T){

	stub := getStub()

	payload := invalidArgument
	args := getChainCodeArgs(getOwnership, ownership_1, payload)

	handleExpectedFailures(t, stub, args, invalidArgument, getOwnership, incorrectNumberOfArgs)

}

func TestOwnershipCreatedDuringPropertyTransaction(t *testing.T){

	stub := getStub()

	owners := []Attribute{
		{Id: ownership_1, Percent:0.45},
		{Id: ownership_2, Percent:0.55}}

	confirmPropertyTransaction(t, stub, owners)

	property := []Attribute{{Id:"1", Percent:0.45, SaleDate:dateString}}
	ownershipPropertyAsString := getAttributesAsString(property)

	invokeGetOwnership(t, stub, ownership_1, ownershipPropertyAsString)

}

func TestPropertyTransaction(t *testing.T){

	stub := getStub()

	owners := []Attribute{
		{Id: ownership_1, Percent:0.45},
		{Id: ownership_2, Percent:0.55}}

	confirmPropertyTransaction(t, stub, owners)

}

func TestMultiplePropertyTransactions(t *testing.T){

	stub := getStub()

	owners := []Attribute{
		{Id: ownership_1, Percent:0.45},
		{Id: ownership_2, Percent:0.55}}

	confirmPropertyTransaction(t, stub, owners)

	owners = []Attribute{
		{Id: ownership_3, Percent:0.35},
		{Id: ownership_4, Percent:0.65}}

	confirmPropertyTransaction(t, stub, owners)

}

func TestMultiplePropertyTransactionsWithRepeatOwners(t *testing.T){

	stub := getStub()

	owners := []Attribute{
		{Id: ownership_1, Percent:0.45},
		{Id: ownership_2, Percent:0.55}}

	confirmPropertyTransaction(t, stub, owners)

	owners = []Attribute{
		{Id: ownership_1, Percent:0.35},
		{Id: ownership_3, Percent:0.65}}

	confirmPropertyTransaction(t, stub, owners)

}

func TestPropertyTransactionExtraArgs(t *testing.T) {

	stub := getStub()

	owners := []Attribute{
		{Id: ownership_1, Percent:0.35},
		{Id: ownership_3, Percent:0.65}}

	_, propertyAsString := createProperty(property_1, owners)

	args := getChainCodeArgs(getOwnership, property_1, propertyAsString, "extraArg")

	handleExpectedFailures(t, stub, args, propertyAsString, propertyTransaction, incorrectNumberOfArgs)

}

func TestPropertyTransactionStringAsSalePrice(t *testing.T) {

	stub := getStub()

	stringAsSalePrice := `{"saleDate":"2017-06-28T21:57:16","salePrice":"1000","owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.55}]}`
	args := getChainCodeArgs(propertyTransaction, property_1, stringAsSalePrice)

	handleExpectedFailures(t, stub, args, stringAsSalePrice, propertyTransaction, cannotUnmarshalStringIntoFloat64)

}

func TestPropertyTransactionMissingSaleDate(t *testing.T) {

	stub := getStub()

	missingSaleDate := `{"salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownership_2","percentage":0.55}]}`
	args := getChainCodeArgs(propertyTransaction, property_1, missingSaleDate)

	handleExpectedFailures(t, stub, args, missingSaleDate, propertyTransaction, saleDateRequired)

}

func TestPropertyTransactionNegativeSalePrice(t *testing.T) {

	stub := getStub()

	negativeSalePrice := `{"saleDate":"2017-06-28T21:57:16","salePrice":-1,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.55}]}`
	args := getChainCodeArgs(propertyTransaction, property_1, negativeSalePrice)

	handleExpectedFailures(t, stub, args, negativeSalePrice, propertyTransaction, salePriceMustBeGreaterThan0)

}

func TestPropertyTransactionNoOwners(t *testing.T) {

	stub := getStub()

	noOwners := `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[]}`
	args := getChainCodeArgs(propertyTransaction, property_1, noOwners)

	handleExpectedFailures(t, stub, args, noOwners, propertyTransaction, atLeastOneOwnerIsRequired)

}

func TestPropertyTransactionTooLowTotalOwnerPercentage(t *testing.T) {

	stub := getStub()

	tooLowOwnerPercentage := `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.50}]}`
	args := getChainCodeArgs(propertyTransaction, property_1, tooLowOwnerPercentage)

	handleExpectedFailures(t, stub, args, tooLowOwnerPercentage, propertyTransaction, totalPercentageCanNotBeGreaterThan1)

}

func TestPropertyTransactionTooHighTotalOwnerPercentage(t *testing.T) {

	stub := getStub()

	tooHighOwnerPercentage := `{"saleDate":"2017-06-28T21:57:16","salePrice":1000,"owners":[{"id":"ownership_3","percentage":0.45},{"id":"ownerhip_2","percentage":0.70}]}`
	args := getChainCodeArgs(propertyTransaction, property_1, tooHighOwnerPercentage)

	handleExpectedFailures(t, stub, args, tooHighOwnerPercentage, propertyTransaction, totalPercentageCanNotBeGreaterThan1)

}

func TestGetProperty(t *testing.T){

	stub := getStub()

	owners := []Attribute{
		{Id: ownership_1, Percent:0.35},
		{Id: ownership_3, Percent:0.65}}

	property, propertyString := createProperty(property_1, owners)

	invokePropertyTransaction(t, stub, property.PropertyId, propertyString)
	invokeGetProperty(t, stub, property, propertyString)

}

func TestGetPropertyExtraArgs(t *testing.T){

	stub := getStub()

	owners := []Attribute{
		{Id: ownership_1, Percent:0.35},
		{Id: ownership_3, Percent:0.65}}

	property, propertyString := createProperty(property_1, owners)

	invokePropertyTransaction(t, stub, property.PropertyId, propertyString)

	args := getChainCodeArgs(getProperty, property.PropertyId, propertyString)

	handleExpectedFailures(t, stub, args, propertyString, getProperty, incorrectNumberOfArgs)

}

func TestGetPropertyMissingProperty(t *testing.T){

	stub := getStub()

	args := getChainCodeArgs(getProperty, property_1)

	handleExpectedFailures(t, stub, args, emptyPropertyJson, getProperty, nilAmountFor)

}
