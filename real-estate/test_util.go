package main

import (
	"bytes"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"errors"
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

func createProperty(propertyId string, saleDate string, salePrice float64, owners []Attribute) (Property, string) {

	property := Property{}
	property.PropertyId = propertyId
	property.SaleDate = saleDate
	property.SalePrice = salePrice
	property.Owners = owners

	propertyAsBytes, _ := getPropertyAsBytes(property)

	return property, string(propertyAsBytes)

}

func getValidOwnersList(owner1 string, owner1Percentage float64, owner2 string, owner2Percentage float64, date string) []Attribute {

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

//TODO this should include a name as well (chaincode needs to be tightened up)
func getValidPropertyListForOwner(propertyId string) []Attribute{

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

func getStub() (*shim.MockStub){

	scc := new(Chaincode)
	stub := shim.NewMockStub("contract", scc)

	return stub

}
