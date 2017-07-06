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
)

type Chaincode struct {}

type Ownership struct {
	Properties		[]Attribute		`json:"properties"`
}

//implement two new fields?
//when return history, need to be returning a []Property?

//The PropertyTransaction will be used for the history of a property (only)
//	- still using Property as before
type Property struct {
	//TxID		not used in ledger
	//PropertyId	not used in ledger
	SaleDate		string			`json:"saleDate"`
	SalePrice	        float64 		`json:"salePrice"`
	Owners 			[]Attribute 		`json:"owners"`
}

//TODO use PropertyTransaction unpackage history of Property into it
type PropertyTransaction struct {
	TxId			string			`json:"txid"`
	PropertyId		string			`json:"id"`
	SaleDate		string			`json:"saleDate"`
	SalePrice		float64			`json:"percentage"`
	Owners			[]Attribute		`json:"owners"`
}

//TODO use OwnershipTransaction, unpackage history of Ownership
type OwnershipTransaction struct{
	//TxId			string			`json:"txid"`
	Properties		[]OwnershipAttribute	`json:"properties"`
}

//TODO only use with Property
type Attribute struct {
	Id			string			`json:"id"`
	Percentage 		float64			`json:"percentage, string"`
}

//TODO, this is really the ownershipTransactions
type OwnershipAttribute struct {
	TxId			string			`json:txid`
	Id			string			`json:"id"`
	Percentage 		float64			`json:"percentage, string"`
}

func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	//No initialization requirements of chain code required at this time
	return shim.Success(nil)

}

func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	fmt.Printf("Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "getOwnership" {
		return t.getOwnership(stub, args)
	} else if function == "getOwnershipHistory" {
		return t.getOwnershipHistory(stub, args)
	} else if  function == "propertyTransaction" {
		return t.propertyTransaction(stub, args)
	} else if function == "getProperty" {
		return t.getProperty(stub, args)
	}else if function == "getPropertyHistory" {
		return t.getPropertyHistory(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"createOwnership\"  \"getOwnership\" \"getOwnershipHistory\" \"propertyTransaction\" \"getProperty\" \"getPropertyHistory\"")

}

func (t *Chaincode) getOwnership(stub shim.ChaincodeStubInterface, args []string) pb.Response{

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting ownership id to query")
	}

	ownershipId := args[0]

	ownershipBytes, err := getOwnershipFromLedger(stub, ownershipId)
	if err != nil {
		return shim.Error(err.Error())
	}

	jsonResp := "{\"OwnershipId\":\"" + ownershipId + "\",\"Ownership Struct\":\"" + string(ownershipBytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)

	return shim.Success(ownershipBytes)

}

func (t *Chaincode) getOwnershipHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	buffer, err := getHistory(stub, args[0], "ownership")
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(buffer)

}

func (t *Chaincode) propertyTransaction(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var propertyId string
	var propertyString string
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	propertyId = args[0]
	propertyString = args[1]

	property := Property{}
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

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting property id to query")
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

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	buffer, err := getHistory(stub, args[0], "property")
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(buffer)

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

		var ownerBytesAndPercentage []interface{} = []interface{}{propertyOwnerBytes, propertyOwners[i].Percentage}
		ownerAndPercentageData[propertyOwners[i].Id] = ownerBytesAndPercentage

	}

	return ownerAndPercentageData, err

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
		ownershipPropertyAttribute.Percentage = propertyOwnership[k][1].(float64)

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

func getHistory(stub shim.ChaincodeStubInterface, id string, historyType string) ([]byte, error){

	var empty []byte
	resultsIterator, err := stub.GetHistoryForKey(id)
	if err != nil {
		err = errors.New("Unable to get history for key: " + id + " | "+ err.Error())
		return empty, err
	}

	defer resultsIterator.Close()

	var buffer bytes.Buffer
	var historyMessage string
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	for resultsIterator.HasNext() {

		response, err := resultsIterator.Next()
		if err != nil {
			return empty, err
		}

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		if historyType == "ownership"{
			historyMessage = ", \"Ownership\":"
		} else if historyType == "property" {
			historyMessage = ", \"Property\":"
		}

		buffer.WriteString(historyMessage)

		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON property)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")

	return buffer.Bytes(), err

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
		totalPercentage += buyers[i].Percentage
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

	buildPropertyJson()
	buildOwnershipJson()

}
