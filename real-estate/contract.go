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
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Chaincode struct {
}

//TODO
// - make sure errors are all handled (custom responses where needed)
// - remove comments that are not needed
// - make layout of methods uniform
type Ownership struct {
	Properties		[]Attribute		`json:"properties"`
}

type Property struct {
	SaleDate		string			`json:"saleDate"`
	SalePrice	        float64 		`json:"salePrice"`
	Owners 			[]Attribute 		`json:"owners"`
}

type Attribute struct {
	Id			string			`json:"id"`
	Percentage 		float64			`json:"percentage"`
}

func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	//No initialization requirements of chain code required at this time
	return shim.Success(nil)
}

func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Printf("Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "createOwnership" {
		return t.createOwnership(stub, args)
	} else if function == "getOwnership" {
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

	return shim.Error("Invalid invoke function name. Expecting \"delete\" \"createOwnership\"  \"getOwnership\" \"getOwnershipHistory\" \"propertyTransaction\" \"getProperty\" \"getPropertyHistory\"")
}

//peer chaincode invoke -C mychannel -n mycc -c '{"Args":["createOwnership","ownership_1","{\"properties\":[{\"id\":\"genesis\",\"percentage\":0}]}"]}'
//peer chaincode invoke -C mychannel -n mycc -c '{"Args":["createOwnership","ownership_1","{\"properties\":[]}"]}'
func (t *Chaincode) createOwnership(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var ownershipId string
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	ownershipId = args[0]
	ownershipString := args[1]

	ownership := Ownership{}
	err = json.Unmarshal([]byte(ownershipString), &ownership)
	if err != nil {
		return shim.Error(err.Error())
	}

	ownershipAsBytes, err := getOwnershipAsBytes(ownership)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(ownershipId, ownershipAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

}

//peer chaincode query -C mychannel -n mycc -c '{"Args":["getOwnership","ownership_1"]}'
func (t *Chaincode) getOwnership(stub shim.ChaincodeStubInterface, args []string) pb.Response{

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting ownership id to query")
	}

	ownershipId := args[0]

	ownershipBytes, err := queryOwnershipInLedger(stub, ownershipId)
	if err != nil {
		return shim.Error(err.Error())
	}

	jsonResp := "{\"OwnershipId\":\"" + ownershipId + "\",\"Ownership Struct\":\"" + string(ownershipBytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)

	return shim.Success(ownershipBytes)

}

//peer chaincode query -C mychannel -n mycc -c '{"Args":["getOwnershipHistory","ownership_1"]}'
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

//peer chaincode invoke -C mychannel -n mycc -c '{"Args":["propertyTransaction","property_1","{\"saleDate\": \"2017-06-28T21:57:16\", \"salePrice\": 1000, \"owners\": [{\"id\":\"ownership_3\",\"percentage\":0.45},{\"id\":\"ownership_2\",\"percentage\":0.55}]}"]}'
//peer chaincode invoke -C mychannel -n mycc -c '{"Args":["propertyTransaction","property_1","{\"saleDate\": \"2017-06-28T21:57:16\", \"salePrice\": 1000, \"owners\": [{\"id\":\"ownership_1\",\"percentage\":0.32},{\"id\":\"ownership_4\",\"percentage\":0.68}]}"]}'
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

	propertyOwnership := map[string][]byte{}
	propertyOwnershipPercentage := map[string]float64{}

	for i := 0; i < len(property.Owners); i++ {
		propertyOwnersBytes, err := queryOwnershipInLedger(stub, property.Owners[i].Id)
		if err != nil {
			return shim.Error("Unable to find ownershipId: " + property.Owners[i].Id + "| " + err.Error())
		}

		propertyOwnership[property.Owners[i].Id] = propertyOwnersBytes
		propertyOwnershipPercentage[property.Owners[i].Id] = property.Owners[i].Percentage
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

	for k, v := range propertyOwnership {
		ownership := Ownership{}
		err = json.Unmarshal(v, &ownership)
		if err != nil {
			return shim.Error(err.Error())
		}

		ownershipPropertyAttribute := Attribute{}
		ownershipPropertyAttribute.Id = propertyId
		ownershipPropertyAttribute.Percentage = propertyOwnershipPercentage[k]

		ownership.Properties = append(ownership.Properties, ownershipPropertyAttribute)

		ownershipAsBytes, err := getOwnershipAsBytes(ownership)
		if err != nil {
			return shim.Error(err.Error())
		}

		err = stub.PutState(k, ownershipAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

	}

	return shim.Success(nil)
}

//peer chaincode query -C mychannel -n mycc -c '{"Args":["getProperty","property_1"]}'
func (t *Chaincode) getProperty(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	var propertyId string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting property id to query")
	}

	propertyId = args[0]

	// Get the state from the ledger
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

//peer chaincode query -C mychannel -n mycc -c '{"Args":["getPropertyHistory","property_1"]}'
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

func queryOwnershipInLedger(stub shim.ChaincodeStubInterface, ownershipId string) ([]byte, error){

	ownershipBytes, err := stub.GetState(ownershipId)
	if err != nil {
		err = errors.New("Unable to retrieve ownershipId: " + ownershipId + ". " + err.Error())
	}
	if ownershipBytes == nil {
		err = errors.New("Nil value for ownershipId: " + ownershipId)
	}

	return ownershipBytes, err
}

func getAttributeAsBytes(attribute Attribute) ([]byte, error){

	var attributeBytes []byte
	var err error

	attributeBytes, err = json.Marshal(attribute)
	if err != nil{
		err = errors.New("Unable to convert list of attributes to json string")
	}

	return attributeBytes, err

}

func getPropertyAsBytes(property Property) ([]byte, error){

	var propertyBytes []byte
	var err error

	propertyBytes, err = json.Marshal(property)
	if err != nil{
		err = errors.New("Unable to convert property to json string")
	}

	return propertyBytes, err

}

func getOwnershipAsBytes(ownership Ownership) ([]byte, error){

	var ownershipBytes []byte
	var err error

	ownershipBytes, err = json.Marshal(ownership)
	if err != nil{
		err = errors.New("Unable to convert ownership to json string")
	}

	return ownershipBytes, err

}

func getHistory(stub shim.ChaincodeStubInterface, id string, historyType string) ([]byte, error){

	var empty []byte
	resultsIterator, err := stub.GetHistoryForKey(id)
	if err != nil {
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

func confirmValidPercentage(buyers []Attribute) error{

	var totalPercentage float64
	var err error

	for i := 0; i < len(buyers); i++ {
		totalPercentage += buyers[i].Percentage
	}

	if totalPercentage != 1 {
		totalPercentageString := fmt.Sprint(totalPercentage)
		err = errors.New("Total Percentage is not correct: " + totalPercentageString)
	}

	return err

}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
