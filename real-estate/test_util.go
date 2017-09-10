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

//TODO break this up a little
type TestContext struct {

	t				   *testing.T
	Stub			   *shim.MockStub
	MethodName         string
	Payload            string
	Arguments          [][]byte
	ArgumentBuilder    []string
	Id                 string
	Attributes         []Attribute
	TestFailureMessage string
	ExpectedStatus     int32
	ExpectedResponse   string
	ExpectFailure	   bool
	Response           peer.Response


}

func createProperty(context TestContext) (string) {

	property := Property{}
	property.PropertyId = context.Id
	property.SaleDate = dateString
	property.SalePrice = salePrice
	property.Owners = context.Attributes

	propertyAsBytes, _ := getPropertyAsBytes(property)

	return string(propertyAsBytes)

}

func getChainCodeArgs(context TestContext) ([][]byte) {

	fmt.Println("payloadLength: ", len(context.ArgumentBuilder))

	args := [][]byte{[]byte(context.MethodName), []byte(context.MethodName)}
	for i := 0; i < len(context.ArgumentBuilder); i++ {
		args = append(args, []byte(context.ArgumentBuilder[i]))
	}

	return args

}

//TODO this needs some work
func invokePropertyTransaction(context TestContext){

	context.Arguments = getChainCodeArgs(context)

	if context.ExpectFailure {
		confirmExpectedTestFailure(context)
	}else {

		//TODO this is hacky and needs to be fixed
		context.Payload = ""

		confirmExpectedTestSuccess(context)
	}

}

func confirmExpectedTestSuccess(context TestContext) {

	context.Response = context.Stub.MockInvoke(context.MethodName, context.Arguments)
	context.TestFailureMessage = failureMessageStart + context.MethodName + ", " + context.Payload + "}, "
	context.ExpectedStatus = shim.OK

	verifyExpectedResponseStatus(context)
	verifyNotExpectedPayload(context)

}

func confirmExpectedTestFailure(context TestContext) {

	context.Response = context.Stub.MockInvoke(context.MethodName, context.Arguments)
	context.TestFailureMessage = failureMessageStart + context.MethodName + ", " + context.Payload + "}, "
	context.ExpectedStatus = shim.ERROR

	verifyExpectedResponseStatus(context)
	verifyExpectedResponseMessage(context)
	verifyExpectedPayload(context)

}

func verifyExpectedResponseStatus(context TestContext) {

	if context.Response.Status != context.ExpectedStatus{
		context.TestFailureMessage += responseStatusStart + strconv.FormatInt(int64(context.Response.Status), 10) + "}]"
		displayTestFailure(context)
	}

}

func verifyExpectedResponseMessage(context TestContext) {

	verifyExpectedResponseMessageSet(context)

	if !strings.Contains(context.Response.Message, context.ExpectedResponse) {
		context.TestFailureMessage += responseMessageStart + string(context.Response.Message) + "}]"
		displayTestFailure(context)
	}

}

func verifyExpectedResponseMessageSet(context TestContext) {
	if len(context.ExpectedResponse) == 0 {
		context.TestFailureMessage = "ExpectedResponse is empty in Context"
		displayTestFailure(context)
	}
}

func verifyExpectedPayload(context TestContext) {

	if string(context.Response.Payload) == context.Payload {
		context.TestFailureMessage += responsePayloadStart + string(context.Response.Payload) + "}]"
		displayTestFailure(context)
	}

}

func verifyNotExpectedPayload(context TestContext) {

	if string(context.Response.Payload) != context.Payload {
		context.TestFailureMessage += responsePayloadStart + string(context.Response.Payload) + "}]"
		displayTestFailure(context)
	}

}

func displayTestFailure(context TestContext) {
	fmt.Println(context.TestFailureMessage)
	context.t.FailNow()

}

func getStub() (*shim.MockStub) {

	scc := new(Chaincode)
	stub := shim.NewMockStub("contract", scc)

	return stub

}

func getTestContext(t *testing.T, stub *shim.MockStub) (TestContext){

	context := TestContext{}
	context.t = t
	context.Stub = stub

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