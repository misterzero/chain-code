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
const salePrice = 125000

const incorrectNumberOfArgs = "Incorrect number of arguments: "
const nilValueForOwnershipId = "Nil value for ownershipId:"
const cannotUnmarshalStringIntoFloat64 = "cannot unmarshal string into Go struct field Property.salePrice of type float64"
const saleDateRequired = "A sale date is required."
const salePriceMustBeGreaterThan0 = "The sale price must be greater than 0"
const totalPercentageCanNotBeGreaterThan1 = "Total Percentage can not be greater than or less than 1. Your total percentage ="
const atLeastOneOwnerIsRequired = "At least one owner is required."

const failureMessageStart = "| FAIL [{args}, {<response-failure>}] | [{"
const responseMessageStart = "{response.Message="
const responseStatusStart = "{response.Status="
const responsePayloadStart = "{response.Payload="

type SessionContext struct {
	Test      					TestContext
	Arguments 					ArgumentContext
}

type TestContext struct {
	t                       	*testing.T
	Stub                    	*shim.MockStub
	TestFailureMessage      	string
	ExpectedStatus          	int32
	ExpectedResponseMessage 	string
	ExpectedResponsePayload 	string
	ShouldFailTest          	bool
	Response                	peer.Response
}

type ArgumentContext struct{
	MethodName    				string
	Payload       				string
	ChainCodeArgs 				[][]byte
	Builder       				[]string
	Id            				string
	Attributes    				[]Attribute
}

func createProperty(context SessionContext) (string) {

	property := Property{}
	property.PropertyId = context.Arguments.Id
	property.SaleDate = dateString
	property.SalePrice = salePrice
	property.Owners = context.Arguments.Attributes

	propertyAsBytes, _ := getPropertyAsBytes(property)

	return string(propertyAsBytes)

}

func getChainCodeArgs(context SessionContext) ([][]byte) {

	args := [][]byte{[]byte(context.Arguments.MethodName), []byte(context.Arguments.MethodName)}
	for i := 0; i < len(context.Arguments.Builder); i++ {
		args = append(args, []byte(context.Arguments.Builder[i]))
	}

	return args

}

func invokePropertyTransaction(context SessionContext) {
	invoke(context)
}

func invokeGetOwnership(context SessionContext) {
	invoke(context)
}

func invokeGetProperty(context SessionContext) {
	invoke(context)
}

func invoke(context SessionContext){

	context.Arguments.ChainCodeArgs = getChainCodeArgs(context)

	if context.Test.ShouldFailTest {
		confirmExpectedTestFailure(context)
	}else {
		confirmExpectedTestSuccess(context)
	}

}

func confirmExpectedTestSuccess(context SessionContext) {

	context.Test.Response = context.Test.Stub.MockInvoke(context.Arguments.MethodName, context.Arguments.ChainCodeArgs)
	context.Test.TestFailureMessage = failureMessageStart + context.Arguments.MethodName + ", " + context.Arguments.Payload + "}, "
	context.Test.ExpectedStatus = shim.OK

	verifyExpectedResponseStatus(context)
	verifyExpectedPayload(context)

}

func confirmExpectedTestFailure(context SessionContext) {

	context.Test.Response = context.Test.Stub.MockInvoke(context.Arguments.MethodName, context.Arguments.ChainCodeArgs)
	context.Test.TestFailureMessage = failureMessageStart + context.Arguments.MethodName + ", " + context.Arguments.Payload + "}, "
	context.Test.ExpectedStatus = shim.ERROR

	verifyExpectedResponseStatus(context)
	verifyExpectedResponseMessage(context)
	verifyExpectedPayload(context)

}

func verifyExpectedResponseStatus(context SessionContext) {

	if context.Test.Response.Status != context.Test.ExpectedStatus{
		context.Test.TestFailureMessage += responseStatusStart + strconv.FormatInt(int64(context.Test.Response.Status), 10) + "}]"
		displayTestFailure(context)
	}

}

func verifyExpectedResponseMessage(context SessionContext) {

	verifyExpectedResponseMessageSet(context)

	if !strings.Contains(context.Test.Response.Message, context.Test.ExpectedResponseMessage) {
		context.Test.TestFailureMessage += responseMessageStart + string(context.Test.Response.Message) + "}]"
		displayTestFailure(context)
	}

}

func verifyExpectedResponseMessageSet(context SessionContext) {

	if len(context.Test.ExpectedResponseMessage) == 0 {
		context.Test.TestFailureMessage = "ExpectedResponseMessage is empty in Context"
		displayTestFailure(context)
	}

}

func verifyExpectedPayload(context SessionContext) {

	if string(context.Test.Response.Payload) != context.Test.ExpectedResponsePayload{
		context.Test.TestFailureMessage += responsePayloadStart + string(context.Test.Response.Payload) + "}]"
		displayTestFailure(context)
	}

}

func displayTestFailure(context SessionContext) {
	fmt.Println(context.Test.TestFailureMessage)
	context.Test.t.FailNow()

}

func getStub() (*shim.MockStub) {

	scc := new(Chaincode)
	stub := shim.NewMockStub("contract", scc)

	return stub

}

func getTestContext(t *testing.T, stub *shim.MockStub) (SessionContext){

	context := SessionContext{}
	context.Test.Stub = stub
	context.Test.t = t

	return context

}

func getAttributesAsString(attributes []Attribute) (string) {

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

func getAttributeAsString(attribute Attribute) (string, error) {

	var attributeBytes []byte
	var err error

	attributeBytes, err = json.Marshal(attribute)
	if err != nil{
		err = errors.New("Unable to convert property to json string " + string(attributeBytes))
	}

	return string(attributeBytes), err

}