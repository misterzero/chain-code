package main

import (
	"bytes"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"errors"

	"testing"
	"strconv"
	"fmt"
	"strings"
)

const getOwnership = "getOwnership"
const propertyTransaction = "propertyTransaction"
const getProperty = "getProperty"
const ownership_1 = "ownership_1"
const property_1 = "property_1"
const dateString = `"2017-06-28T21:57:16"`
const emptyOwnershipPropertyJson = `{"properties":[]}`
const emptyPropertyJson = `{"saleDate":"","salePrice":0,"owners":[]}`
const errorStatus = int32(500)

const incorrectArgCountError = "Incorrect number of arguments: "
const getOwnershipMissingError = "Nil value for ownershipId:"
const getPropertyNilError = "Nil amount for"
const createPropertyTransactionUnmarshalWrongTypeError = "cannot unmarshal string into Go struct field Property.salePrice of type float64"
const createPropertyTransactionMissingSaleDateError = "A sale date is required."
const createPropertyTransactionGreaterThanZeroSalePriceError = "The sale price must be greater than 0"
const createPropertyTransactionTotalPercentageOfOneError = "Total Percentage can not be greater than or less than 1. Your total percentage ="
const createPropertyTransactionNoOwnersError = "At least one owner is required."

func checkGetOwnership(t *testing.T, stub *shim.MockStub,ownershipId string, ownershipString string){

	ownershipArgs := getThreeArgs(getOwnership, ownershipId)

	message:= " | " + getOwnership + " with args: {" + string(ownershipArgs[0]) + ", "+ string(ownershipArgs[1]) + ", " + string(ownershipArgs[2]) + "}, failed. "

	handleExpectedSuccess(t, stub, getOwnership, message, ownershipArgs, ownershipString)

}

func checkPropertyTransaction(t *testing.T, stub *shim.MockStub, propertyId string, propertyString string ){

	propertyArgs := getFourArgs(propertyTransaction, propertyId, propertyString)

	message:= " | " + propertyTransaction + " with args: {" + string(propertyArgs[1]) + ", " + string(propertyArgs[2]) + ", " + string(propertyArgs[3]) + "}, failed. "

	res := stub.MockInvoke(propertyTransaction, propertyArgs)
	if res.Status != shim.OK {
		message := message +  "[res.Status=" + strconv.FormatInt(int64(res.Status), 10) + "]"
		fmt.Println(message)
		t.FailNow()
	}

}

func checkPropertyState(t *testing.T, stub *shim.MockStub, property Property, propertyString string) []byte{

	bytes := stub.State[property.PropertyId]
	if bytes == nil {
		fmt.Println("Property", string(bytes), "failed to get value")
		t.FailNow()
	}
	if string(bytes) != propertyString {
		fmt.Println("Property value", property.PropertyId, "was not", propertyString, "as expected")
		t.FailNow()
	}

	return bytes

}

func checkGetProperty(t *testing.T, stub *shim.MockStub, property Property, propertyString string){

	propertyArgs := getThreeArgs(getProperty, property.PropertyId)

	message:= " | " + getProperty + " with args: {" + string(propertyArgs[1]) + "}, failed. "

	handleExpectedSuccess(t, stub, getProperty, message, propertyArgs, propertyString)

}


func handleExpectedSuccess(t *testing.T, stub *shim.MockStub, argument string, outputMessage string, invalidArgs [][]byte, attemptedPayload string){

	res := stub.MockInvoke(argument, invalidArgs)
	if res.Status != shim.OK {
		msg := outputMessage +  "[res.Status=" + strconv.FormatInt(int64(res.Status), 10) + "]"
		fmt.Println(msg)
		t.FailNow()
	}
	if res.Payload == nil {
		msg := outputMessage + "[res.Message=" + string(res.Message) + "]"
		fmt.Println(msg)
		t.FailNow()
	}
	if string(res.Payload) != attemptedPayload {
		msg := outputMessage + "[res.Payload=" + string(res.Payload) + "]"
		fmt.Println(msg)
		t.FailNow()
	}

}

func handleExpectedFailures(t *testing.T, stub *shim.MockStub, argument string, outputMessage string, errorMessage string, invalidArgs [][]byte, attemptedPayload string){

	res := stub.MockInvoke(argument, invalidArgs)

	if res.Status != errorStatus {
		msg := outputMessage +  "[res.Status=" + strconv.FormatInt(int64(res.Status), 10) + "]"
		fmt.Println(msg)
		t.FailNow()
	}
	if !strings.Contains(res.Message, errorMessage) {
		msg := outputMessage + "[res.Message=" + string(res.Message) + "]"
		fmt.Println(msg)
		t.FailNow()
	}
	if string(res.Payload) == attemptedPayload{
		msg := outputMessage + "[res.Payload=" + string(res.Payload) + "]"
		fmt.Println(msg)
		t.FailNow()
	}

}

func getFourArgs(function string, id string, jsonStruct string) ([][]byte) {

	method := []byte(function)
	key := []byte(id)
	value := []byte(jsonStruct)

	args := [][]byte{method, method, key, value}

	return args
}

func getThreeArgs(function string, id string) ([][]byte) {

	method := []byte(function)
	key := []byte(id)

	args := [][]byte{method, method, key}

	return args
}

func getMockProperty(propertyId string, saleDate string, salePrice float64, owners []Attribute) (Property, string) {

	property := Property{}
	property.PropertyId = propertyId
	property.SaleDate = saleDate
	property.SalePrice = salePrice
	property.Owners = owners

	propertyAsBytes, _ := getPropertyAsBytes(property)

	return property, string(propertyAsBytes)

}

func getMockValidOwners() []Attribute {

	owner1 := Attribute{}
	owner2 := Attribute{}
	owner1.Id = "ownership_3"
	owner1.Percent = 0.45
	owner1.SaleDate = dateString

	owner2.Id = "ownership_2"
	owner2.Percent = 0.55
	owner2.SaleDate = dateString

	ownershipInputList := []Attribute{owner1, owner2}

	return ownershipInputList

}

func getMockValidOwnersProperty() []Attribute{

	property1 := Attribute{}
	property2 := Attribute{}

	property1.Id = "1"
	property1.Percent = 0.45
	property1.SaleDate = dateString

	property2.Id = "1"
	property2.Percent = 0.55
	property2.SaleDate = dateString

	ownershipInputList := []Attribute{property1, property2}

	return ownershipInputList

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

func getStub() (*shim.MockStub){

	scc := new(Chaincode)
	stub := shim.NewMockStub("contract", scc)

	return stub

}
