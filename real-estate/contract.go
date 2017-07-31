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
	"fmt"
	"bytes"
	"errors"
	"strings"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

type Chaincode struct {}

type Ownership struct {
	Properties		        []Attribute		`json:"properties"`
}

type Property struct {
	TxId			string			`json:"txid"`
	PropertyId		        string			`json:"id"`
	SaleDate		        string			`json:"saleDate"`
	SalePrice		        float64			`json:"salePrice"`
	Owners			[]Attribute		`json:"owners"`
}

type Attribute struct {
	Id                          string			`json:"id"`
	Percent                     float64			`json:"percent, string"`
}

func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	//No initialization requirements of chain code required at this time
	return shim.Success(nil)

}

func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	var errorMessage string
	function, args := stub.GetFunctionAndParameters()

	if function != "invoke" {
		errorMessage = "Invalid function: " + function
	}

	if args[0] == "getOwnership" {
		return t.getOwnership(stub, args)
	} else if args[0] == "getOwnershipHistory" {
		return t.getOwnershipHistory(stub, args)
	} else if  args[0] == "propertyTransaction" {
		return t.propertyTransaction(stub, args)
	} else if args[0] == "getProperty" {
		return t.getProperty(stub, args)
	}else if args[0] == "getPropertyHistory" {
		return t.getPropertyHistory(stub, args)
	}

	errorMessage = "Invalid method:  " + args[0]

	return shim.Error(errorMessage)

}

func (t *Chaincode) getOwnership(stub shim.ChaincodeStubInterface, args []string) pb.Response{

	if len(args) != 2 {
		return shim.Error("(getOwnership) Incorrect number of arguments: " + strconv.Itoa(len(args)) + ". Expecting 2")
	}

	ownershipId := args[1]

	ownershipPropertiesAsBytes, err := getOwnershipProperties(stub, ownershipId)
	if err != nil {
		return shim.Error(err.Error())
	}

	jsonResp := "{\"OwnershipId\":\"" + ownershipId + "\",\"Ownership Properties Struct\":\"" + string(ownershipPropertiesAsBytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)

	return shim.Success(ownershipPropertiesAsBytes)

}

func (t *Chaincode) getOwnershipHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("(getOwnershipHistory) Incorrect number of arguments: " + strconv.Itoa(len(args)) + ". Expecting 2")
	}

	id := args[1]
	resultsIterator, err := stub.GetHistoryForKey(id)
	if err != nil {
		err = errors.New("Unable to get history for key: " + id + " | "+ err.Error())
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	for resultsIterator.HasNext() {

		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\",")

		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON property)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			ownership := Ownership{}

			err := json.Unmarshal(response.Value, &ownership)
			if err != nil {
				return shim.Error(err.Error())
			}

			buffer.WriteString("\"ownership\":[")
			for i := 0; i < len(ownership.Properties); i++ {
				buffer.WriteString("{")
				buffer.WriteString("\"id\":\"")
				buffer.WriteString(ownership.Properties[i].Id)
				buffer.WriteString("\",\"percent\":")
				percent := strconv.FormatFloat(ownership.Properties[i].Percent, 'f', 2, 64)
				buffer.WriteString(percent)
				buffer.WriteString("}")

				if i != len(ownership.Properties) - 1{
					buffer.WriteString(",")
				}
			}
			buffer.WriteString("]")
		}

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true

	}

	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())

}

func (t *Chaincode) propertyTransaction(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var propertyId string
	var propertyString string
	var err error

	if len(args) != 3 {
		return shim.Error("(propertyTransaction) Incorrect number of arguments: " + strconv.Itoa(len(args)) + ". Expecting 3")
	}

	propertyId = args[1]
	propertyString = args[2]

	property := Property{}
	property.TxId = stub.GetTxID()
	err = json.Unmarshal([]byte(propertyString), &property)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = verifyValidProperty(property)
	if err != nil {
		return shim.Error(err.Error())
	}

	propertyOwnership, err := getOwnershipPropertyUpdateRequirements(stub, property.Owners)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = confirmValidPercentage(property.Owners)
	if err != nil {
		return shim.Error(err.Error())
	}

	propertyAsBytes, err := getPropertyAsBytes(property)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(propertyId, propertyAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = updateOwnershipProperties(stub, propertyId, propertyOwnership)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

}

func (t *Chaincode) getProperty(stub shim.ChaincodeStubInterface, args []string) pb.Response{

	var propertyId string
	var err error

	if len(args) != 2 {
		return shim.Error("(getProperty) Incorrect number of arguments: " + strconv.Itoa(len(args)) + ". Expecting 2")
	}

	propertyId = args[0]

	propertyBytes, err := stub.GetState(propertyId)
	if err != nil {
		return shim.Error(err.Error())
	}

	if propertyBytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + propertyId + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"PropertyId\":\"" + propertyId + "\",\"Property Struct\":\"" + string(propertyBytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)

	return shim.Success(propertyBytes)

}

func (t *Chaincode) getPropertyHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("(getPropertyHistory) Incorrect number of arguments: " + strconv.Itoa(len(args))  + ". Expecting 2")
	}

	id := args[1]
	resultsIterator, err := stub.GetHistoryForKey(id)
	if err != nil {
		err = errors.New("Unable to get history for key: " + id + " | "+ err.Error())
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	for resultsIterator.HasNext() {

		response, err := resultsIterator.Next()
		if err != nil {
			shim.Error(err.Error())
		}

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"txId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\",")

		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON property)

		property := Property{}
		err = json.Unmarshal(response.Value, &property)
		if err != nil {
			err = errors.New("Unable to convert property bytes to Property structure. " + err.Error())

			return shim.Error(err.Error())
		}

		buffer.WriteString("\"saleDate\":\"" + property.SaleDate + "\",")
		buffer.WriteString("\"salePrice\":" + strconv.FormatFloat(property.SalePrice, 'f', -1, 64) + ",")

		propertyOwners := property.Owners

		propertyId := property.PropertyId
		propertyNumber := strings.Replace(propertyId,"property_","",-1)

		property.PropertyId = propertyNumber

		buffer.WriteString("\"propertyId\":\"" + property.PropertyId + "\",")
		buffer.WriteString( "\"owners\":")


		for i := 0; i < len(propertyOwners); i++ {
			ownershipId := propertyOwners[i].Id
			ownershipNumber := strings.Replace(ownershipId,"ownership_","",-1)
			propertyOwners[i].Id = ownershipNumber
		}
		property.Owners = propertyOwners

		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			propertyOwnersAsBytes, err := json.Marshal(property.Owners)
			if err != nil{
				err = errors.New("Unable to convert propertyOwners to json string " + string(propertyOwnersAsBytes))
			}

			buffer.WriteString(string(propertyOwnersAsBytes))
		}

		buffer.WriteString("}")

		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())

}

func getOwnershipFromLedger(stub shim.ChaincodeStubInterface, ownershipId string) ([]byte, error){

	ownershipBytes, err := stub.GetState(ownershipId)
	if err != nil {
		err = errors.New("Unable to retrieve ownershipId: " + ownershipId + ". " + err.Error())
	}
	if ownershipBytes == nil {
		err = errors.New("Nil value for ownershipId: " + ownershipId)
	}

	return ownershipBytes, err
}

func getOwnershipPropertyUpdateRequirements(stub shim.ChaincodeStubInterface, propertyOwners[]Attribute) (map[string][]interface{}, error){

	var ownerAndPercentageData = make(map[string][]interface{})
	var err error

	for i := 0; i < len(propertyOwners); i++ {
		var ownershipId string
		var propertyOwnerBytes []byte
		ownershipId = propertyOwners[i].Id

		propertyOwnerBytes, _ = getOwnershipFromLedger(stub, ownershipId)
		if propertyOwnerBytes == nil {

			ownership := Ownership{}
			propertyOwnerBytes, err = getOwnershipAsBytes(ownership)
			if err != nil {
				return ownerAndPercentageData, err
			}

		}

		var ownerBytesAndPercentage []interface{} = []interface{}{propertyOwnerBytes, propertyOwners[i].Percent}
		ownerAndPercentageData[propertyOwners[i].Id] = ownerBytesAndPercentage

	}

	return ownerAndPercentageData, err

}

func getOwnershipProperties(stub shim.ChaincodeStubInterface, ownershipId string ) ([]byte, error){

	var err error

	ownershipBytes, err := getOwnershipFromLedger(stub, ownershipId)
	if err != nil {
		return ownershipBytes, err
	}

	ownership := Ownership{}
	err = json.Unmarshal(ownershipBytes, &ownership)
	if err != nil {
		return ownershipBytes, err
	}

	ownershipProperties := getOwnershipPropertiesWithProppertyIdOnly(ownership.Properties)

	ownershipPropertiesAsBytes, err := json.Marshal(ownershipProperties)
	if err != nil{
		err = errors.New("Unable to convert ownership properties to json string " + string(ownershipPropertiesAsBytes))
		return ownershipPropertiesAsBytes, err
	}

	return ownershipPropertiesAsBytes, err

}

func getOwnershipPropertiesWithProppertyIdOnly(ownershipProperties []Attribute) ([]Attribute){

	for i := 0; i < len(ownershipProperties); i++ {
		propertyId := ownershipProperties[i].Id
		propertyNumber := strings.Replace(propertyId,"property_","",-1)
		ownershipProperties[i].Id = propertyNumber
	}

	return ownershipProperties

}

func updateOwnershipProperties(stub shim.ChaincodeStubInterface, propertyId string, propertyOwnership map[string][]interface{}) (error){

	var err error

	for k, _ := range propertyOwnership {

		ownership := Ownership{}
		err = json.Unmarshal(propertyOwnership[k][0].([]byte), &ownership)
		if err != nil {
			err = errors.New("Unable to convert byte array to ownership struct,  " + string(propertyOwnership[k][0].([]byte)))
			return err
		}

		ownershipPropertyAttribute := Attribute{}
		ownershipPropertyAttribute.Id = propertyId
		ownershipPropertyAttribute.Percent = propertyOwnership[k][1].(float64)

		ownership.Properties = append(ownership.Properties, ownershipPropertyAttribute)

		ownershipAsBytes, err := getOwnershipAsBytes(ownership)
		if err != nil {
			return err
		}

		err = stub.PutState(k, ownershipAsBytes)
		if err != nil {
			err = errors.New("Unable to add Ownerhsip to Ledger,  " + string(ownershipAsBytes))
			return err
		}

	}

	return err

}

func getPropertyAsBytes(property Property) ([]byte, error){

	var propertyBytes []byte
	var err error

	propertyBytes, err = json.Marshal(property)
	if err != nil{
		err = errors.New("Unable to convert property to json string " + string(propertyBytes))
	}

	return propertyBytes, err

}

func getOwnershipAsBytes(ownership Ownership) ([]byte, error){

	var ownershipBytes []byte
	var err error

	ownershipBytes, err = json.Marshal(ownership)
	if err != nil{
		err = errors.New("Unable to convert ownership to json string " + string(ownershipBytes))
	}

	return ownershipBytes, err

}

func verifyValidProperty(property Property) (error){

	var err error

	if strings.TrimSpace(property.SaleDate) == "" {
		err = errors.New("A sale date is required.")
		return err
	}
	if property.SalePrice < 1 {
		err = errors.New("The sale price must be greater than 0.")
		return err
	}
	if len(property.Owners) < 1 {
		err = errors.New("At least one owner is required.")
	}

	return err
}

func confirmValidPercentage(buyers []Attribute) error{

	var totalPercentage float64
	var err error

	for i := 0; i < len(buyers); i++ {
		totalPercentage += buyers[i].Percent
	}

	if totalPercentage != 1 {
		totalPercentageString := fmt.Sprint(totalPercentage)
		err = errors.New("Total Percentage can not be greater than or less than 1. Your total percentage =" + totalPercentageString)
	}

	return err

}

func main() {

	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}

}
