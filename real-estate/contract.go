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

//TODO create contractContext wrapper
//TODO change t.* to something more descriptive

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

type ContractContext struct {
	Chaincode					*Chaincode
	Stub						shim.ChaincodeStubInterface
	Arguments					[]string

	//TODO wip
	OwnershipContext			OwnershipContext
}

type OwnershipContext struct {
	Ownership					Ownership
}

type PropertyContext struct {
	Property					Property
}

var contractContext = ContractContext{}

func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	//No initialization requirements of chain code required at this time
	return shim.Success(nil)

}

//func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
//
//	contractContext := ContractContext{}
//	contractContext.Stub = stub
//
//	var errorMessage string
//	function, args := stub.GetFunctionAndParameters()
//
//	if function != "invoke" {
//		errorMessage = "Invalid function: " + function
//	}
//
//	if args[0] == "getOwnership" {
//		return t.getOwnership(stub, args)
//	} else if args[0] == "getOwnershipHistory" {
//		return t.getOwnershipHistory(stub, args)
//	} else if  args[0] == "propertyTransaction" {
//		return t.propertyTransaction(stub, args)
//	} else if args[0] == "getProperty" {
//		return t.getProperty(stub, args)
//	}else if args[0] == "getPropertyHistory" {
//		return t.getPropertyHistory(stub, args)
//	}
//
//	errorMessage = "Invalid method:  " + args[0]
//
//	return shim.Error(errorMessage)
//
//}

//TODO make one function call to handle arg requirement size
func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	contractContext := ContractContext{}
	contractContext.Stub = stub

	var errorMessage string

	function, args := stub.GetFunctionAndParameters()
	contractContext.Arguments = args

	if function != "invoke" {
		errorMessage = "Invalid function: " + function
	}

	if args[0] == "getOwnership" {
		return t.getOwnership(contractContext)
	} else if args[0] == "getOwnershipHistory" {
		return t.getOwnershipHistory(contractContext)
	} else if  args[0] == "propertyTransaction" {
		return t.propertyTransaction(contractContext)
	} else if args[0] == "getProperty" {
		return t.getProperty(contractContext)
	}else if args[0] == "getPropertyHistory" {
		return t.getPropertyHistory(contractContext)
	}

	errorMessage = "Invalid method:  " + args[0]

	return shim.Error(errorMessage)

}

func (t *Chaincode) getOwnership(contractContext ContractContext) pb.Response{

	if len(contractContext.Arguments) != 2 {
		return shim.Error("(getOwnership) Incorrect number of arguments: " + strconv.Itoa(len(contractContext.Arguments)) + ". Expecting 2")
	}

	ownershipId := contractContext.Arguments[1]

	ownershipPropertiesAsBytes, err := getOwnershipPropertiesAsBytes(contractContext.Stub, ownershipId)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(ownershipPropertiesAsBytes)

}

func (t *Chaincode) getOwnershipHistory(contractContext ContractContext) pb.Response {

	if len(contractContext.Arguments) != 2 {
		return shim.Error("(getOwnershipHistory) Incorrect number of arguments: " + strconv.Itoa(len(contractContext.Arguments)) + ". Expecting 2")
	}

	id := contractContext.Arguments[1]
	resultsIterator, err := contractContext.Stub.GetHistoryForKey(id)
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

				propertyAsBytes, err := getPropertyFromLedger(contractContext.Stub, "property_" + ownershipProperties[i].Id)

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

func (t *Chaincode) propertyTransaction(contractContext ContractContext) pb.Response {

	var propertyId string
	var propertyString string
	var err error

	if len(contractContext.Arguments) != 3 {
		return shim.Error("(propertyTransaction) Incorrect number of arguments: " + strconv.Itoa(len(contractContext.Arguments)) + ". Expecting 3")
	}

	propertyId = contractContext.Arguments[1]
	propertyString = contractContext.Arguments[2]

	property := Property{}
	property.TxId = contractContext.Stub.GetTxID()
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

	propertyBytes, err := contractContext.Stub.GetState(propertyId)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = updatePropertyOwnership(contractContext.Stub, property, propertyBytes, propertyId)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = addPropertyToLedger(contractContext.Stub, property, propertyId)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

}

func (t *Chaincode) getProperty(contractContext ContractContext) pb.Response{

	var propertyId string
	var err error

	if len(contractContext.Arguments) != 2 {
		return shim.Error("(getProperty) Incorrect number of arguments: " + strconv.Itoa(len(contractContext.Arguments)) + ". Expecting 2")
	}

	propertyId = contractContext.Arguments[1]

	propertyBytes, err := getPropertyFromLedger(contractContext.Stub, propertyId)
	if err != nil {
		return shim.Error(err.Error())
	}

	jsonResp := "{\"PropertyId\":\"" + propertyId + "\",\"Property Struct\":\"" + string(propertyBytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)

	return shim.Success(propertyBytes)

}

func (t *Chaincode) getPropertyHistory(contractContext ContractContext) pb.Response {

	if len(contractContext.Arguments) != 2 {
		return shim.Error("(getPropertyHistory) Incorrect number of arguments: " + strconv.Itoa(len(contractContext.Arguments))  + ". Expecting 2")
	}

	id := contractContext.Arguments[1]
	resultsIterator, err := contractContext.Stub.GetHistoryForKey(id)
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

func main() {

	contractContext.Chaincode = new(Chaincode)

	err := shim.Start(contractContext.Chaincode)
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}



}