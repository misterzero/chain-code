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

	MethodName         string
	Payload            string
	Arguments          [][]byte
	Id                 string
	Attributes         []Attribute
	TestFailureMessage string
	ExpectedStatus     int32
	ExpectedResponse   string
	Response           peer.Response

}

func createProperty(context SessionContext) (Property, string) {

	property := Property{}
	property.PropertyId = context.Id
	property.SaleDate = dateString
	property.SalePrice = salePrice
	property.Owners = context.Attributes

	propertyAsBytes, _ := getPropertyAsBytes(property)

	return property, string(propertyAsBytes)

}

//TODO convert to use just the context object as parameter
func getChainCodeArgs(chainCodeMethodName string, payload ...string) ([][]byte) {

	args := [][]byte{[]byte(chainCodeMethodName), []byte(chainCodeMethodName)}
	for i := 0; i < len(payload); i++ {
		args = append(args, []byte(payload[i]))
	}

	return args

}

func confirmPropertyTransaction(t *testing.T, stub *shim.MockStub, context SessionContext) {

	_, context.Payload = createProperty(context)

	invokePropertyTransaction(t, stub, context)

	context.MethodName = getProperty
	invokeGetProperty(t, stub, context)

}

func invokeGetProperty(t *testing.T, stub *shim.MockStub, context SessionContext) {

	context.Arguments = getChainCodeArgs(context.MethodName, context.Id)

	handleExpectedSuccess(t, stub, context)

}

func invokeGetOwnership(t *testing.T, stub *shim.MockStub, context SessionContext) {

	context.Arguments = getChainCodeArgs(context.MethodName, context.Id)

	handleExpectedSuccess(t, stub, context)

}

func invokePropertyTransaction(t *testing.T, stub *shim.MockStub, context SessionContext) {

	context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload)

	context.Payload = ""

	handleExpectedSuccess(t, stub, context)

}

func handleExpectedSuccess(t *testing.T, stub *shim.MockStub, context SessionContext) {

	context.Response = stub.MockInvoke(context.MethodName, context.Arguments)

	context.TestFailureMessage = failureMessageStart + context.MethodName + ", " + context.Payload + "}, "

	context.ExpectedStatus = shim.OK
	verifyExpectedResponseStatus(t, context)

	verifyNotExpectedPayload(t, context)

}

func handleExpectedFailures(t *testing.T, stub *shim.MockStub, context SessionContext) {

	context.Response = stub.MockInvoke(context.MethodName, context.Arguments)
	context.TestFailureMessage = failureMessageStart + context.MethodName + ", " + context.Payload + "}, "
	context.ExpectedStatus = shim.ERROR

	verifyExpectedResponseStatus(t, context)
	verifyExpectedResponseMessage(t, context)
	verifyExpectedPayload(t, context)

}

func verifyExpectedResponseStatus(t *testing.T, context SessionContext) {

	if context.Response.Status != context.ExpectedStatus{
		context.TestFailureMessage += responseStatusStart + strconv.FormatInt(int64(context.Response.Status), 10) + "}]"
		displayTestFailure(t, context.TestFailureMessage)
	}

}

func verifyExpectedResponseMessage(t *testing.T, context SessionContext) {

	verifyExpectedResponseMessageSet(t, context)

	if !strings.Contains(context.Response.Message, context.ExpectedResponse) {
		context.TestFailureMessage += responseMessageStart + string(context.Response.Message) + "}]"
		displayTestFailure(t, context.TestFailureMessage)
	}

}

func verifyExpectedResponseMessageSet(t *testing.T, context SessionContext) {
	if len(context.ExpectedResponse) == 0 {
		failureMessage := "ExpectedResponse is empty in Context"
		displayTestFailure(t, failureMessage)
	}
}

func verifyExpectedPayload(t *testing.T, context SessionContext) {

	fmt.Println("Full Response: ", context.Response)
	fmt.Println("Payload: ", context.Response.Payload)

	if string(context.Response.Payload) == context.Payload {
		context.TestFailureMessage += responsePayloadStart + string(context.Response.Payload) + "}]"
		displayTestFailure(t, context.TestFailureMessage)
	}

}

func verifyNotExpectedPayload(t *testing.T, context SessionContext) {

	if string(context.Response.Payload) != context.Payload {
		context.TestFailureMessage += responsePayloadStart + string(context.Response.Payload) + "}]"
		displayTestFailure(t, context.TestFailureMessage)
	}

}

func displayTestFailure(t *testing.T, failureMessage string) {
	fmt.Println(failureMessage)
	t.FailNow()

}

func getStub() (*shim.MockStub) {

	scc := new(Chaincode)
	stub := shim.NewMockStub("contract", scc)

	return stub

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