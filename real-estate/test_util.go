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

//TODO set this up to accept only context
func createProperty(propertyId string, owners []Attribute) (Property, string) {

	property := Property{}
	property.PropertyId = propertyId
	property.SaleDate = dateString
	property.SalePrice = salePrice
	property.Owners = owners

	propertyAsBytes, _ := getPropertyAsBytes(property)

	return property, string(propertyAsBytes)

}

func createProperty2(context SessionContext) (Property, string) {

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

//TODO working
func confirmPropertyTransaction2(t *testing.T, stub *shim.MockStub, context SessionContext) {

	_, context.Payload = createProperty2(context)

	invokePropertyTransaction2(t, stub, context)

	context.MethodName = getProperty
	invokeGetProperty2(t, stub, context)

}

func invokeGetProperty(t *testing.T, stub *shim.MockStub, propertyId, propertyString string) {

	args := getChainCodeArgs(getProperty, propertyId)

	handleExpectedSuccess(t, stub, getProperty, args, propertyString)

}

func invokeGetProperty2(t *testing.T, stub *shim.MockStub, context SessionContext) {

	context.Arguments = getChainCodeArgs(context.MethodName, context.Id)

	handleExpectedSuccess2(t, stub, context)

}

func invokeGetOwnership2(t *testing.T, stub *shim.MockStub, context SessionContext) {

	context.Arguments = getChainCodeArgs(context.MethodName, context.Id)

	handleExpectedSuccess2(t, stub, context)

}

func invokePropertyTransaction(t *testing.T, stub *shim.MockStub, propertyId string, payload string ) {

	args := getChainCodeArgs(propertyTransaction, propertyId, payload)
	handleExpectedSuccess(t, stub, propertyTransaction, args, "")

}

func invokePropertyTransaction2(t *testing.T, stub *shim.MockStub, context SessionContext) {

	context.Arguments = getChainCodeArgs(context.MethodName, context.Id, context.Payload)

	context.Payload = ""

	handleExpectedSuccess2(t, stub, context)

}


//TODO less parameters
func handleExpectedSuccess(t *testing.T, stub *shim.MockStub, argument string, args [][]byte, payload string) {

	response := stub.MockInvoke(argument, args)

	fmt.Println("Response: ", response)

	failureMessage := failureMessageStart + argument + ", " + payload + "}, "

	verifyExpectedResponseStatus(t, response, failureMessage, shim.OK)
	verifyNotExpectedPayload(t, response, failureMessage, payload)

}

func handleExpectedSuccess2(t *testing.T, stub *shim.MockStub, context SessionContext) {

	context.Response = stub.MockInvoke(context.MethodName, context.Arguments)

	context.TestFailureMessage = failureMessageStart + context.MethodName + ", " + context.Payload + "}, "

	context.ExpectedStatus = shim.OK
	verifyExpectedResponseStatus2(t, context)

	verifyNotExpectedPayload2(t, context)

}

func handleExpectedFailures(t *testing.T, stub *shim.MockStub, context SessionContext) {

	context.Response = stub.MockInvoke(context.MethodName, context.Arguments)
	context.TestFailureMessage = failureMessageStart + context.MethodName + ", " + context.Payload + "}, "
	context.ExpectedStatus = shim.ERROR

	verifyExpectedResponseStatus2(t, context)
	verifyExpectedResponseMessage2(t, context)
	verifyExpectedPayload2(t, context)

}
//TODO less parameters
func verifyExpectedResponseStatus(t *testing.T, response peer.Response, failureMessage string, statusValue int32) {

	if response.Status != statusValue {
		failureMessage += responseStatusStart + strconv.FormatInt(int64(response.Status), 10) + "}]"
		displayTestFailure(t, failureMessage)
	}

}

func verifyExpectedResponseStatus2(t *testing.T, context SessionContext) {

	if context.Response.Status != context.ExpectedStatus{
		context.TestFailureMessage += responseStatusStart + strconv.FormatInt(int64(context.Response.Status), 10) + "}]"
		displayTestFailure(t, context.TestFailureMessage)
	}

}

func verifyExpectedResponseMessage2(t *testing.T, context SessionContext) {

	verifyExpectedResponseMessageSet2(t, context)

	if !strings.Contains(context.Response.Message, context.ExpectedResponse) {
		context.TestFailureMessage += responseMessageStart + string(context.Response.Message) + "}]"
		displayTestFailure(t, context.TestFailureMessage)
	}

}

func verifyExpectedResponseMessageSet2(t *testing.T, context SessionContext) {
	if len(context.ExpectedResponse) == 0 {
		failureMessage := "ExpectedResponse is empty in Context"
		displayTestFailure(t, failureMessage)
	}
}

func verifyExpectedPayload2(t *testing.T, context SessionContext) {

	if string(context.Response.Payload) == context.Payload {
		context.TestFailureMessage += responsePayloadStart + string(context.Response.Payload) + "}]"
		displayTestFailure(t, context.TestFailureMessage)
	}

}

//TODO less parameters
//TODO needs new name
func verifyNotExpectedPayload(t *testing.T, response peer.Response, failureMessage string, payload string) {

	if string(response.Payload) != payload {
		failureMessage += responsePayloadStart + string(response.Payload) + "}]"
		displayTestFailure(t, failureMessage)
	}

}

func verifyNotExpectedPayload2(t *testing.T, context SessionContext) {

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