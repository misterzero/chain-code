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
	Properties                  []Attribute   	`json:"properties"`
}

type Property struct {
	TxId                        string      	`json:"txid"`
	PropertyId                  string      	`json:"id"`
	SaleDate                    string      	`json:"saleDate"`
	SalePrice                   float64     	`json:"salePrice"`
	Owners                      []Attribute   	`json:"owners"`
}

type Attribute struct {
	Id                          string      	`json:"id"`
	SaleDate 					string			`json:"saleDate"`
	Name						string			`json:"name"`
	Percent                     float64     	`json:"percent"`
}

//chaincode methods
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

	ownershipPropertiesAsBytes, err := getOwnershipPropertiesAsBytes(stub, ownershipId)
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

			ownershipProperties := getOwnershipPropertiesIdValues(ownership.Properties)

			for i := 0; i < len(ownershipProperties); i++ {
				buffer.WriteString("{")
				buffer.WriteString("\"id\":\"")
				buffer.WriteString(ownershipProperties[i].Id)
				buffer.WriteString("\",\"percent\":")
				percent := strconv.FormatFloat(ownershipProperties[i].Percent, 'f', 2, 64)
				buffer.WriteString(percent)
				buffer.WriteString(",\"saleDate\":\"")

				propertyAsBytes, err := getPropertyFromLedger(stub, "property_" + ownershipProperties[i].Id)

				property := Property{}
				json.Unmarshal(propertyAsBytes, &property)
				if err != nil {
					return shim.Error(err.Error())
				}

				buffer.WriteString(ownershipProperties[i].SaleDate)
				buffer.WriteString("\"")

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

	err = confirmValidPercentage(property.Owners)
	if err != nil {
		return shim.Error(err.Error())
	}

	propertyBytes, err := stub.GetState(propertyId)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = updatePropertyOwnership(stub, property, propertyBytes, propertyId)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = addPropertyToLedger(stub, property, propertyId)
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

	propertyId = args[1]

	propertyBytes, err := getPropertyFromLedger(stub, propertyId)
	if err != nil {
		return shim.Error(err.Error())
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

//property transaction methods
func updatePropertyOwnership(stub shim.ChaincodeStubInterface, newProperty Property, originalPropertyBytes []byte, propertyId string) error{
	var err error
	var sameOwnersList = []Attribute{}
	var updateNewOwnersList = []Attribute{}
	var updatedOldOwnersList = []Attribute{}

	originalProperty := Property{}
	if originalPropertyBytes != nil {

		err =json.Unmarshal(originalPropertyBytes, &originalProperty)
		if err != nil {
			err = errors.New("Unable to create originalPropertyBytes: " + string(originalPropertyBytes) + ". " + err.Error())
			return err
		}

	}

	sameOwnersList, updateNewOwnersList, updatedOldOwnersList = getOwnershipLists(newProperty.Owners, originalProperty.Owners)

	err = removePropertyFromOwnership(stub, updatedOldOwnersList, propertyId)
	if err != nil {
		return err
	}

	err = addPropertyToOwnership(stub, updateNewOwnersList, newProperty)
	if err != nil {
		return err
	}

	err = updatePropertyForSameOwnership(stub, sameOwnersList, newProperty)
	if err != nil {
		return err
	}

	return err

}

func updatePropertyForSameOwnership(stub shim.ChaincodeStubInterface, sameOwnersList []Attribute, newProperty Property) error{

	var err error
	var propertyAttribute = Attribute{}

	for i := 0; i < len(sameOwnersList); i++ {

		ownershipAsBytes, err := getOwnershipFromLedger(stub, sameOwnersList[i].Id)

		ownership := Ownership{}
		if ownershipAsBytes != nil {

			err = json.Unmarshal(ownershipAsBytes, &ownership)
			if err != nil {
				return err
			}

			for _, v := range ownership.Properties {

				if v.Id == newProperty.PropertyId {

					ownership.Properties = removePropertyFromOwnershipList(ownership.Properties, newProperty.PropertyId)

					err = addOwnershipToLedger(stub,ownership, sameOwnersList[i].Id)
					if err != nil {
						return err
					}

					propertyAttribute.Id = newProperty.PropertyId
					propertyAttribute.SaleDate = newProperty.SaleDate
					propertyAttribute.Percent = sameOwnersList[i].Percent
					propertyAttribute.Name = sameOwnersList[i].Name

					for j:= 0; j < len(newProperty.Owners); j++ {

						if newProperty.Owners[j].Id == sameOwnersList[i].Id {
							propertyAttribute.Percent = newProperty.Owners[j].Percent
						}

					}

					ownership.Properties = append(ownership.Properties, propertyAttribute)

					err = addOwnershipToLedger(stub, ownership, sameOwnersList[i].Id)
					if err != nil {
						return err
					}

				}

			}

		}else {
			err = nil
		}

	}

	return err

}

func addPropertyToOwnership(stub shim.ChaincodeStubInterface, newOwnersList []Attribute, newProperty Property) error{

	var err error
	var propertyAttribute = Attribute{}

	for i := 0; i < len(newOwnersList); i++ {

		ownershipAsBytes, err := getOwnershipFromLedger(stub, newOwnersList[i].Id)

		propertyAttribute.Id = newProperty.PropertyId
		propertyAttribute.SaleDate = newProperty.SaleDate
		propertyAttribute.Percent = newOwnersList[i].Percent
		propertyAttribute.Name = newOwnersList[i].Name

		ownership := Ownership{}
		if ownershipAsBytes != nil {

			err = json.Unmarshal(ownershipAsBytes, &ownership)
			if err != nil {
				return err
			}

		} else {
			err = nil
		}

		ownership.Properties = append(ownership.Properties, propertyAttribute)

		err = addOwnershipToLedger(stub, ownership, newOwnersList[i].Id)
		if err != nil {
			return err
		}

	}

	return err

}

func removePropertyFromOwnership(stub shim.ChaincodeStubInterface, oldOwners []Attribute, propertyId string) error{

	var err error

	for i := 0; i < len(oldOwners); i++ {

		ownershipAsBytes, err := getOwnershipFromLedger(stub, oldOwners[i].Id)

		ownership := Ownership{}
		if ownershipAsBytes != nil {

			err = json.Unmarshal(ownershipAsBytes, &ownership)
			if err != nil {
				return err
			}

			ownership.Properties = removePropertyFromOwnershipList(ownership.Properties, propertyId)

		}else {
			err = nil
		}

		err = addOwnershipToLedger(stub, ownership, oldOwners[i].Id)
		if err != nil {
			return err
		}

	}

	return err

}

func removePropertyFromOwnershipList(properties []Attribute, propertyId string) []Attribute{

	for i, v := range properties {
		if v.Id == propertyId {
			properties[i] = properties[len(properties) - 1]
			properties = properties[: len(properties) - 1]
		}
	}

	return properties
}

func getOwnershipLists(newOwners []Attribute, oldOwners []Attribute) ([]Attribute, []Attribute, []Attribute){
	var sameOwners = []Attribute{}
	var updatedNewOwners = []Attribute{}
	var updatedOldOwners = []Attribute{}

	if len(newOwners) >= len(oldOwners) {
		sameOwners, updatedNewOwners, updatedOldOwners = buildOwnershipLists(newOwners, oldOwners)
	} else {
		sameOwners, updatedOldOwners, updatedNewOwners = buildOwnershipLists(oldOwners, newOwners)
	}

	return sameOwners, updatedNewOwners, updatedOldOwners

}

func buildOwnershipLists(longestOwnerList []Attribute, shortestOwnerList []Attribute) ([]Attribute, []Attribute, []Attribute){

	var sameOwnersList = []Attribute{}
	var updatedLongestOwnerList = []Attribute{}
	var updatedShortestOwnerList = []Attribute{}

	sameOwnersList, updatedLongestOwnerList = getLongestAndSameOwnershipLists(longestOwnerList, shortestOwnerList)

	updatedShortestOwnerList = getShortestOwnershipList(shortestOwnerList, sameOwnersList)

	return sameOwnersList, updatedLongestOwnerList, updatedShortestOwnerList

}

func getLongestAndSameOwnershipLists(longestOwnersList []Attribute, shortestOwnersList []Attribute) ([]Attribute, []Attribute){

	var sameOwners = []Attribute{}
	var updatedLongestOwnersList = []Attribute{}

	for i := 0; i < len(longestOwnersList); i++ {
		var foundMatch = false

		for j := 0; j< len(shortestOwnersList); j++ {

			if longestOwnersList[i].Id == shortestOwnersList[j].Id {

				sameOwners = append(sameOwners, longestOwnersList[i])
				foundMatch = true
				break

			} else {
				foundMatch = false
			}

		}

		if !foundMatch{
			updatedLongestOwnersList = append(updatedLongestOwnersList, longestOwnersList[i])
		}

	}

	return sameOwners, updatedLongestOwnersList

}

func getShortestOwnershipList(shortestOwnersList []Attribute, sameOwnersList []Attribute) ([]Attribute){

	var updatedShortestOwnersList = []Attribute{}

	for k := 0; k <len(shortestOwnersList); k++ {

		var foundMatch = false

		for m:= 0; m < len(sameOwnersList); m++ {

			if shortestOwnersList[k].Id == sameOwnersList[m].Id {
				foundMatch = true
				break
			} else {
				foundMatch = false
			}

		}

		if !foundMatch {
			updatedShortestOwnersList = append(updatedShortestOwnersList, shortestOwnersList[k])
		}

	}

	return updatedShortestOwnersList

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

//helper methods
func addOwnershipToLedger(stub shim.ChaincodeStubInterface, ownership Ownership, ownershipId string) error{

	updatedOwnershipAsBytes, err := getOwnershipAsBytes(ownership)
	if err != nil {
		return err
	}

	err = stub.PutState(ownershipId, updatedOwnershipAsBytes)
	if err != nil {
		err = errors.New("Unable to add property for new Owners,  " + string(updatedOwnershipAsBytes))
		return err
	}

	return err
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

func getOwnershipAsBytes(ownership Ownership) ([]byte, error){

	var ownershipBytes []byte
	var err error

	ownershipBytes, err = json.Marshal(ownership)
	if err != nil{
		err = errors.New("Unable to convert ownership to json string " + string(ownershipBytes))
	}

	return ownershipBytes, err

}

func getOwnershipPropertiesAsBytes(stub shim.ChaincodeStubInterface, ownershipId string ) ([]byte, error){

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

	ownershipProperties := getOwnershipPropertiesIdValues(ownership.Properties)

	ownershipPropertiesAsBytes, err := json.Marshal(ownershipProperties)
	if err != nil{
		err = errors.New("Unable to convert ownership properties to json string " + string(ownershipPropertiesAsBytes))
		return ownershipPropertiesAsBytes, err
	}

	return ownershipPropertiesAsBytes, err

}

func getOwnershipPropertiesIdValues(ownershipProperties []Attribute) ([]Attribute){

	for i := 0; i < len(ownershipProperties); i++ {
		propertyId := ownershipProperties[i].Id
		propertyNumber := strings.Replace(propertyId,"property_","",-1)
		ownershipProperties[i].Id = propertyNumber
	}

	return ownershipProperties

}

func addPropertyToLedger(stub shim.ChaincodeStubInterface, property Property, propertyId string) error{

	propertyAsBytes, err := getPropertyAsBytes(property)
	if err != nil {
		return err
	}

	err = stub.PutState(propertyId, propertyAsBytes)
	if err != nil {
		return err
	}

	return err

}

func getPropertyFromLedger(stub shim.ChaincodeStubInterface, propertyId string) ([]byte, error){

	var err error

	propertyBytes, err := stub.GetState(propertyId)
	if err != nil {
		return propertyBytes, err
	}

	if propertyBytes == nil {
		err = errors.New("{\"Error\":\"Nil amount for " + propertyId + "\"}")
		return propertyBytes, err
	}

	return propertyBytes, err

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

func main() {

	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}

}