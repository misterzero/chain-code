package main

import (
	"bytes"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"errors"
	"testing"
	"strconv"
	"strings"
	"fmt"
)

const getOwnership = "getOwnership"
const propertyTransaction = "propertyTransaction"
const getProperty = "getProperty"
const ownership_1 = "ownership_1"
const ownership_2 = "ownership_2"
const ownership_3 = "ownership_3"
const ownership_4 = "ownership_4"
const property_1 = "property_1"
const dateString = `"2017-06-28T21:57:16"`
const emptyPropertyJson = `{"saleDate":"","salePrice":0,"owners":[]}`
const invalidArgument = "invalidArgument"

const incorrectNumberOfArgs = "Incorrect number of arguments: "
const nilValueForOwnershipId = "Nil value for ownershipId:"
const nilAmountFor = "Nil amount for"
const cannotUnmarshalStringIntoFloat64 = "cannot unmarshal string into Go struct field Property.salePrice of type float64"
const saleDateRequired = "A sale date is required."
const salePriceMustBeGreaterThan0 = "The sale price must be greater than 0"
const totalPercentageCanNotBeGreaterThan1 = "Total Percentage can not be greater than or less than 1. Your total percentage ="
const atLeastOneOwnerIsRequired = "At least one owner is required."

const failureMessageStart = "| FAIL [{args}, {<response-failure>}] | [{"
const responseMessageStart = "{response.Message="
const responseStatusStart = "{response.Status="
const responsePayloadStart = "{response.Payload="


func createProperty(propertyId string, saleDate string, salePrice float64, owners []Attribute) (Property, string) {

	property := Property{}
	property.PropertyId = propertyId
	property.SaleDate = saleDate
	property.SalePrice = salePrice
	property.Owners = owners

	propertyAsBytes, _ := getPropertyAsBytes(property)

	return property, string(propertyAsBytes)

}

func getListOfOwnersForProperty(owner1 string, owner1Percentage float64, owner2 string, owner2Percentage float64, date string) []Attribute {

	ownerAttribute1 := Attribute{}
	ownerAttribute2 := Attribute{}
	ownerAttribute1.Id = owner1
	ownerAttribute1.Percent = owner1Percentage
	ownerAttribute1.SaleDate = dateString

	ownerAttribute2.Id = owner2
	ownerAttribute2.Percent = owner2Percentage
	ownerAttribute2.SaleDate = dateString

	ownershipInputList := []Attribute{ownerAttribute1, ownerAttribute2}

	return ownershipInputList

}

func getPropertyListForOwner(propertyId string) []Attribute{

	property1 := Attribute{}
	property2 := Attribute{}

	property1.Id = propertyId
	property1.Percent = 0.45
	property1.SaleDate = dateString

	property2.Id = propertyId
	property2.Percent = 0.55
	property2.SaleDate = dateString

	ownershipInputList := []Attribute{property1, property2}

	return ownershipInputList

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

func getAttributeAsString(attribute Attribute) (string, error){

	var attributeBytes []byte
	var err error

	attributeBytes, err = json.Marshal(attribute)
	if err != nil{
		err = errors.New("Unable to convert property to json string " + string(attributeBytes))
	}

	return string(attributeBytes), err

}

func getChainCodeArgs(chainCodeMethodName string, payload ...string) ([][]byte){

	args := [][]byte{[]byte(chainCodeMethodName), []byte(chainCodeMethodName)}
	for i := 0; i < len(payload); i++ {
		args = append(args, []byte(payload[i]))
	}
	return args

}

func confirmAdditionOfPropertyTransactionToLedger(t *testing.T, stub *shim.MockStub, ownershipId1 string, ownerPercentage1 float64, ownershipId2 string, ownershipPercentage2 float64) {

	ownerList := getListOfOwnersForProperty(ownershipId1, ownerPercentage1, ownershipId2, ownershipPercentage2, dateString)
	property, propertyString := createProperty(property_1, dateString, 1000, ownerList)

	invokePropertyTransaction(t, stub, property.PropertyId, propertyString)
	checkGetProperty(t, stub, property, propertyString)

}

func checkGetProperty(t *testing.T, stub *shim.MockStub, property Property, propertyString string){

	args := getChainCodeArgs(getProperty, property.PropertyId)

	handleExpectedSuccess(t, stub, getProperty, args, propertyString)

}

func invokeGetOwnership(t *testing.T, stub *shim.MockStub,ownershipId string, payload string){

	args := getChainCodeArgs(getOwnership, ownershipId)

	handleExpectedSuccess(t, stub, getOwnership, args, payload)

}

func invokePropertyTransaction(t *testing.T, stub *shim.MockStub, propertyId string, payload string ){

	args := getChainCodeArgs(propertyTransaction, propertyId, payload)

	handleExpectedSuccess(t, stub, propertyTransaction, args, "")

}

func handleExpectedSuccess(t *testing.T, stub *shim.MockStub, argument string, args [][]byte, payload string){

	response := stub.MockInvoke(argument, args)

	failureMessage := failureMessageStart + argument + ", " + payload + "}, "

	verifyExpectedResponseStatus(t, response, failureMessage, shim.OK)
	verifyExpectedInvalidPayload(t, response, failureMessage, payload)

}

func handleExpectedFailures(t *testing.T, stub *shim.MockStub, args [][]byte, payload string, argument string, expectedResponseMessage string){

	response := stub.MockInvoke(argument, args)

	failureMessage := failureMessageStart + argument + ", " + payload + "}, "

	verifyExpectedResponseStatus(t, response, failureMessage, shim.ERROR)
	verifyExpectedResponseMessage(t, response, failureMessage, expectedResponseMessage)
	verifyExpectedValidPayload(t, response, failureMessage, payload)

}

func verifyExpectedResponseStatus(t *testing.T, response peer.Response, failureMessage string, statusValue int32) {

	if response.Status != statusValue {
		failureMessage += responseStatusStart + strconv.FormatInt(int64(response.Status), 10) + "}]"
		displayFailure(t, failureMessage)
	}

}

func verifyExpectedResponseMessage(t *testing.T, response peer.Response, failureMessage string, expectedResponseMessage string) {

	if !strings.Contains(response.Message, expectedResponseMessage) {
		failureMessage += responseMessageStart + string(response.Message) + "}]"
		displayFailure(t, failureMessage)
	}

}

func verifyExpectedValidPayload(t *testing.T, response peer.Response, failureMessage string, payload string) {

	if string(response.Payload) == payload {
		failureMessage += responsePayloadStart + string(response.Payload) + "}]"
		displayFailure(t, failureMessage)
	}

}

func verifyExpectedInvalidPayload(t *testing.T, response peer.Response, failureMessage string, payload string) {

	if string(response.Payload) != payload {
		failureMessage += responsePayloadStart + string(response.Payload) + "}]"
		displayFailure(t, failureMessage)
	}

}

func displayFailure(t *testing.T, failureMessage string) {
	fmt.Println(failureMessage)
	t.FailNow()

}

func getStub() (*shim.MockStub){

	scc := new(Chaincode)
	stub := shim.NewMockStub("contract", scc)

	return stub

}
