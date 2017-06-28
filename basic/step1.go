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

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"fmt"
	"strconv"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	"encoding/json"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type ActivePoll struct{
	Name 		string 				`json:"name"`
	Token		int 				`json:"token"`
}

type Option struct{
	Name 		string 				`json:"name"`
	Count 		int 				`json:"count"`
}

type User struct{
	Active 		[]ActivePoll 		`json:"active"`
	Inactive	[]string 			`json:"inactive"`
}

type Poll struct{
	Options 		[]Option 				`json:"options"`
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
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
	} else if function == "addUser" {
		return t.addUser(stub, args)
	} else if function == "queryUser" {
		return t.queryUser(stub, args)
	} else if function == "addNewActivePollToUser" {
		return t.addNewActivePollToUser(stub, args)
	} else if function == "ActiveToInactivePoll" {
		return t.ActiveToInactivePoll(stub, args)
	} else if function == "addNewPoll" {
		return t.addNewPoll(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"move\" \"delete\" \"query\" \"findAll\"")
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Init")
	_, args := stub.GetFunctionAndParameters()
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state to the ledger
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

//peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n mycc -c '{"Args":["addUser","c"]}'
func (t *SimpleChaincode) addUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var user User 
	var userAsJsonByteArray []byte
    var userAsJsonString string
    var userExists bool
    var err error

	if len(args) < 1 {
		return shim.Error("put operation must include one arguments, a key")
	}
	key := args[0]
	user, err = createUser()

	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}

	userAsJsonByteArray, err = getUserAsJsonByteArray(user)

    if err != nil{

        return shim.Error(err.Error())

    }

    userAsJsonString = string(userAsJsonByteArray)
    userExists, err = isExistingUser(stub, key)

    if err != nil {

        return shim.Error(err.Error())

    }

    if userExists {

        return shim.Error("An asset already exists with id:" + key)

    }
	

	err = stub.PutState(key, []byte(userAsJsonString))
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("######################## New user ! ##################################")

	return shim.Success(nil)
}

//peer chaincode query -C mychannel -n mycc -c '{"Args":["queryUser","c"]}'
//peer chaincode query -C mychannel -n mycc -c '{"Args":["queryUser","my_poll0"]}'
func (t *SimpleChaincode) queryUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    var key string // Entities

    var err error

    if len(args) != 1 {

        return shim.Error("Incorrect number of arguments. Expecting name of the person to query")

    }

    key = args[0]

    // Get the state from the ledger

    keyValBytes, err := stub.GetState(key)

    if err != nil {

        jsonResp := "{\"Error\":\"Failed to get state for " + key + "\"}"

        return shim.Error(jsonResp)

    }

    if keyValBytes == nil {

        jsonResp := "{\"Error\":\"Nil amount for " + key + "\"}"

        return shim.Error(jsonResp)

    }

    jsonResp := "{\"Id\":\"" + key + "\",\"Value\":\"" + string(keyValBytes) + "\"}"

    fmt.Printf("Query Response:%s\n", jsonResp)

    return shim.Success(keyValBytes)

}

//peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n mycc -c '{"Args":["addNewActivePollToUser","c","my new poll"]}'
func (t *SimpleChaincode) addNewActivePollToUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var user User 
	var userAsJsonByteArray []byte
    var userAsJsonString string
    var userExists bool
    var err error
    var pollToAdd ActivePoll

	if len(args) < 2 {
		return shim.Error("put operation must include one arguments, a key")
	}
	key := args[0]

    userExists, err = isExistingUser(stub, key)

    if err != nil {

        return shim.Error(err.Error())

    }

    if userExists==false {

        return shim.Error("No user with id:" + key)

    }

	userAsBytes, err := stub.GetState(key)

    if err != nil{

        return shim.Error(err.Error())

    }

    user, err = getUserFromJsonByteArray(userAsBytes)
	

    if err != nil{

        return shim.Error(err.Error())

    }

    pollToAdd,err = createActivePoll(args[1])

    if err != nil{

        return shim.Error(err.Error())

    }

    user.Active = append(user.Active,pollToAdd)

    userAsJsonByteArray, err = getUserAsJsonByteArray(user)

    if err != nil{

        return shim.Error(err.Error())

    }

    userAsJsonString = string(userAsJsonByteArray)

	err = stub.PutState(key, []byte(userAsJsonString))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n mycc -c '{"Args":["ActiveToInactivePoll","c","my new poll1"]}'
func (t *SimpleChaincode) ActiveToInactivePoll(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var user User 
	var userAsJsonByteArray []byte
    var userAsJsonString string
    var userExists bool
    var err error
   	var pollName string
   	var indexOfThePoll int

	if len(args) < 2 {
		return shim.Error("put operation must include one arguments, a key")
	}
	key := args[0]
	pollName = args[1]

    userExists, err = isExistingUser(stub, key)

    if err != nil {

        return shim.Error(err.Error())

    }

    if userExists==false {

        return shim.Error("No user with id:" + key)

    }

	userAsBytes, err := stub.GetState(key)

    if err != nil{

        return shim.Error(err.Error())

    }

    user, err = getUserFromJsonByteArray(userAsBytes)
	

    if err != nil{

        return shim.Error(err.Error())

    }

    indexOfThePoll = 0
    for i := 0; i < len(user.Active); i++ {
        if user.Active[i].Name==pollName {
        	indexOfThePoll = i
        }
    }

    user.Active = append(user.Active[:indexOfThePoll], user.Active[indexOfThePoll+1:]...)

    user.Inactive = append(user.Inactive,pollName)

    userAsJsonByteArray, err = getUserAsJsonByteArray(user)

    if err != nil{

        return shim.Error(err.Error())

    }


    userAsJsonString = string(userAsJsonByteArray)

	err = stub.PutState(key, []byte(userAsJsonString))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n mycc -c '{"Args":["addNewPoll","my_poll0","{\"Options\":[{\"Name\":\"opt1\",\"Count\":0},{\"Name\":\"opt2\",\"Count\":0}]}"]}'
func (t *SimpleChaincode) addNewPoll(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var poll Poll
	var pollAsJsonByteArray []byte
    var pollAsJsonString string
	var pollExists bool
	var err error

	if len(args) < 2 {
		return shim.Error("put operation must include one arguments, a key")
	}

	key := args[0]
	s := args[1]

	poll, err = createNewPoll(s)
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}

	pollAsJsonByteArray, err = getPollAsJsonByteArray(poll)
    if err != nil{
        return shim.Error(err.Error())
    }

    pollAsJsonString = string(pollAsJsonByteArray)
    pollExists, err = isExistingPoll(stub, key)
    if err != nil {
        return shim.Error(err.Error())
    }
    if pollExists {
        return shim.Error("An asset already exists with id:" + key)
    }
	

	err = stub.PutState(key, []byte(pollAsJsonString))
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("######################## New poll ! ##################################")
	
	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) move(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
func (t *SimpleChaincode) findAll(stub shim.ChaincodeStubInterface) pb.Response {
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

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func createUser() (User, error){
	emptyActive := []ActivePoll{}
	emptyInactive := []string{}
	var err error

    return User{Active: emptyActive, Inactive: emptyInactive}, err
}

func createActivePoll(pollName string) (ActivePoll, error){
	tok := 1
	var err error

    return ActivePoll{Name: pollName, Token: tok}, err
}

func createNewPoll(s string) (Poll, error){
	var poll Poll
	var err error
	err = json.Unmarshal([]byte(s), &poll)

	return poll,err

}

func isExistingUser(stub shim.ChaincodeStubInterface, key string) (bool, error){

    var err error

    result := false

    userAsBytes, err := stub.GetState(key)

    if(len(userAsBytes) != 0){

        result = true

    }

    return result, err
}

func isExistingPoll(stub shim.ChaincodeStubInterface, key string) (bool, error){

    var err error

    result := false

    pollAsBytes, err := stub.GetState(key)

    if(len(pollAsBytes) != 0){

        result = true

    }

    return result, err
}

func getUserAsJsonByteArray(user User) ([]byte, error){

    var jsonUser []byte

    var err error

    jsonUser, err = json.Marshal(user)

    if err != nil{

        errors.New("Unable to convert asset to json string")

    }

    return jsonUser, err

}

func getUserFromJsonByteArray(userAsJsonByte []byte) (User, error) {

    var user User

    var err error

    err = json.Unmarshal(userAsJsonByte, &user)

    if err != nil {

        errors.New("Unable to convert json []byte to an User")

    }

    return user, err

}

func getPollAsJsonByteArray(poll Poll) ([]byte, error){

    var jsonPoll []byte

    var err error

    jsonPoll, err = json.Marshal(poll)

    if err != nil{

        errors.New("Unable to convert asset to json string")

    }

    return jsonPoll, err

}
