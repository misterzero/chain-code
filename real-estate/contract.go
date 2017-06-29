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

//TODO
// - getOwnership(id) Ownership (CurrentState)
// - getAllOwnership(list<id>) list<Ownership>
// - getOwnershipHistory(id) Ownership (History)
// - getProperty(id) Property (CurrentState)
// - getAllProperty(list<id>) list<Property> (CurrentState)
// - getPropertyHistory(id) Property (History)
// - propertyTransaction(Date, SalePrice, list<Attribute>) nil
package main

import (
	"fmt"
	"strconv"
	"errors"
	"time"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Chaincode struct {
}

type Ownership struct {
	Properties	[]Attribute	`json:"properties"`
}

type Property struct {
	Date		time.Time 	`json:"dateOfSale"`
	SalePrice	float64		`json:"salePrice"`
	Owners		[]Attribute	`json:"owners"`
}

type Attribute struct {
	Id		int64		`json:"id"`
	Percentage 	float64		`json:"percentage"`
}

func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	//No initialization requirements of chain code required at this time
	return shim.Success(nil)
}

func (t *Chaincode) processTransaction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string    // Entities
	var Aval int // Asset holdings
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	fmt.Printf("Aval = %d", Aval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}


func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Printf("Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "move" {
		// Make payment of X units from A to B
		return t.move(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	} else if function == "findAll" {
		return t.findAll(stub)
	} else if  function == "processTransasction" {
		return t.processTransaction(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"move\" \"delete\" \"query\" \"findAll\"")
}

// Transaction makes payment of X units from A to B
func (t *Chaincode) move(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Deletes an entity from state
func (t *Chaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *Chaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)

	jsonAvalBytes, err := getFullByteArray(A, Avalbytes)

	return shim.Success(jsonAvalBytes)
}

// query callback representing the query of a chaincode
func (t *Chaincode) findAll(stub shim.ChaincodeStubInterface) pb.Response {
	var err error

	aKey := "a"
	bKey := "b"

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(aKey)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + aKey + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + aKey + "\"}"
		return shim.Error(jsonResp)
	}

	aJsonResp := "{\"Name\":\"" + aKey + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", aJsonResp)

	Bvalbytes, err := stub.GetState(bKey)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + bKey + "\"}"
		return shim.Error(jsonResp)
	}

	if Bvalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + bKey + "\"}"
		return shim.Error(jsonResp)
	}

	bJsonResp := "{\"Name\":\"" + bKey + "\",\"Amount\":\"" + string(Bvalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", bJsonResp)

	aAndBvalueResponse, err := getAllFullByteArrays(Avalbytes, Bvalbytes)

	return shim.Success(aAndBvalueResponse)
}

func getAllFullByteArrays(aValBytes []byte, bValBytes []byte)([]byte, error){
	aResult, err := getFullByteArray("a", aValBytes)
	bResult, err := getFullByteArray("b", bValBytes)

	resultString := `{"ledger":[`+ string(aResult) + `,` + string(bResult) + `]}`
	result := []byte(resultString)

	return result, err
}

func getFullByteArray(id string, byteArray []byte) ([]byte, error){
	byteString := `{"id":"` + id + `", "value":`+ string(byteArray) + `}`

	fullStruct := []byte(byteString)

	return fullStruct, nil
}

//====================================================================================================================
func createAttributeFromArgs(args []string) (*Attribute, error){

	var attribute Attribute
	var id int64
	var percentage float64
	var err error

	if len(args) != 2 {
		err = errors.New("expected args length of 2, but received " + string(len(args)))
	}

	id, err = strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		err = errors.New("unable to parse " + string(args[0]) + " as int")
	}

	percentage, err = strconv.ParseFloat(args[1], 64)
	if err != nil {
		err = errors.New("unable to parse " + string(args[1]) + " as float")
	}

	attribute.Id = id
	attribute.Percentage = percentage

	return &attribute, err

}

func getAttributeListAsBytes(attribute []Attribute) ([]byte, error){

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

func getFormattedTimeAsString(time time.Time, format string)(string, error){

	var err error
	var formattedTimeString string

	formattedTimeString = time.Format(format)
	if len(formattedTimeString) == 0 {
		err = errors.New("Unable to format time " + string(time.String()))
	}

	return formattedTimeString, err

}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}