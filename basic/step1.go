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
	Options 		[]Option 		`json:"options"`
	Status 			int 			`json:"status"`
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
	} else if function == "addUser" {
		return t.addUser(stub, args)
	} else if function == "newQuery" {
		return t.newQuery(stub, args)
	} else if function == "addNewActivePollToUser" {
		return t.addNewActivePollToUser(stub, args)
	} else if function == "activeToInactivePoll" {
		return t.activeToInactivePoll(stub, args)
	} else if function == "addNewPoll" {
		return t.addNewPoll(stub, args)
	} else if function == "vote" {
		return t.vote(stub, args)
	} else if function == "changeStatusToZero" {
		return t.changeStatusToZero(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"move\" \"delete\" \"query\" \"findAll\"")
}

//TODO modify
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	//No need

	return shim.Success(nil)
}


//peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n mycc -c '{"Args":["addUser","c"]}'
func (t *SimpleChaincode) addUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var user User 
	var userAsJsonByteArray []byte
    var userAsJsonString string
    var userExists bool
    var err error

	if len(args) != 1 {
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
	fmt.Println("New user: "+ key + " added to the ledger")

	return shim.Success(nil)
}

//peer chaincode query -C mychannel -n mycc -c '{"Args":["newQuery","c"]}'
//peer chaincode query -C mychannel -n mycc -c '{"Args":["newQuery","my_poll0"]}'
func (t *SimpleChaincode) newQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {

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

//peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n mycc -c '{"Args":["addNewActivePollToUser","c","my_poll0"]}'
//peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n mycc -c '{"Args":["addNewActivePollToUser","c","my_poll1"]}'
//peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n mycc -c '{"Args":["addNewActivePollToUser","c","my_poll2"]}'
func (t *SimpleChaincode) addNewActivePollToUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var user User 
	var userAsJsonByteArray []byte
    var userAsJsonString string
    var userExists, pollExists, pollAlreadyAddToUser bool
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

    pollExists, err = isExistingPoll(stub, args[1])
    if pollExists==false{
    	return shim.Error("No poll with name:" + args[1])
    }

	userAsBytes, err := stub.GetState(key)

    if err != nil{

        return shim.Error(err.Error())

    }

    user, err = getUserFromJsonByteArray(userAsBytes)
    if err != nil{
        return shim.Error(err.Error())
    }

    pollAlreadyAddToUser, err = isPollAlreadyAddToUser(stub, key, args[1])
    if pollAlreadyAddToUser==true {
    	return shim.Error("this poll: " + args[1] + " has already been add to the user: " +key)
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

//peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n mycc -c '{"Args":["activeToInactivePoll","c","my_poll0"]}'
func (t *SimpleChaincode) activeToInactivePoll(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var user User 
	var userAsJsonByteArray []byte
    var userAsJsonString string
    var userExists, pollContained bool
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
    if userAsBytes == nil {
		return shim.Error("Entity not found")
	}


    user, err = getUserFromJsonByteArray(userAsBytes)
	

    if err != nil{

        return shim.Error(err.Error())

    }

    pollContained = false
    indexOfThePoll = 0
    for i := 0; i < len(user.Active); i++ {
        if user.Active[i].Name==pollName {
        	indexOfThePoll = i
        	pollContained = true
        }
    }

    if pollContained==false {
    	return shim.Error("No such poll for this user")
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

//peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n mycc -c '{"Args":["addNewPoll","my_poll0","{\"Options\":[{\"Name\":\"opt1\",\"Count\":0},{\"Name\":\"opt2\",\"Count\":0}],\"status\":1}"]}'
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
		return shim.Error("Expecting integer value for user holding")
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
        return shim.Error("A poll already exists with id:" + key)
    }
	

	err = stub.PutState(key, []byte(pollAsJsonString))
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("New poll: "+ key + " has been added to the ledger")
	
	return shim.Success(nil)
}

//peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n mycc -c '{"Args":["vote","c","my_poll0","opt1"]}'
func (t *SimpleChaincode) vote(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var userAsJsonBytes, pollAsJsonBytes []byte
	var userAsJsonString, pollAsJsonString string
	var UserAllowed, userExists, pollExists, optionExists, statusOne bool
	var pollKey, userKey, option string
	var err error
	var user User
	var poll Poll
	var indexOfThePoll, indexOfTheOption int

	if len(args) < 3 {
		return shim.Error("put operation must include three arguments, a userKey, a pollKey and an option ")
	}

	userKey = args[0]
	pollKey = args[1]
	option = args[2]


	userExists, err = isExistingUser(stub, userKey)
    if err != nil {
        return shim.Error(err.Error())
    }
    if userExists==false {
        return shim.Error("No user with id:" + userKey)
    }

    pollExists, err = isExistingPoll(stub, pollKey)
    if err != nil {
        return shim.Error(err.Error())
    }
    if pollExists==false {
        return shim.Error("No poll named: " + pollKey + " has been created")
    }

    UserAllowed, indexOfThePoll, err = isUserAllowed(stub,userKey,pollKey)
    if UserAllowed==false {
    	return shim.Error("The user: " + userKey + " is not allowed to participate to this poll")
    }

    optionExists, indexOfTheOption, err = isExistingOption(stub,pollKey,option)
    if optionExists==false {
    	return shim.Error("The option: " + option + " does not exist")
    }

    userAsJsonBytes, err = stub.GetState(userKey)
    if err != nil{
        return shim.Error(err.Error())
    }
    if userAsJsonBytes == nil {
		return shim.Error("Entity not found")
	}
    user, err = getUserFromJsonByteArray(userAsJsonBytes)
	if err != nil{
        return shim.Error(err.Error())
    }

    pollAsJsonBytes, err = stub.GetState(pollKey)
    if err != nil{
        return shim.Error(err.Error())
    }
    if pollAsJsonBytes == nil {
		return shim.Error("Entity not found")
	}
    poll, err = getPollFromJsonByteArray(pollAsJsonBytes)
	if err != nil{
        return shim.Error(err.Error())
    }

    statusOne, err = isStatusOne(stub, poll)
    if statusOne == false {
    	return shim.Error("the current poll is closed")
    }

    user.Active[indexOfThePoll].Token = user.Active[indexOfThePoll].Token - 1
    poll.Options[indexOfTheOption].Count = poll.Options[indexOfTheOption].Count + 1

    userAsJsonBytes, err = getUserAsJsonByteArray(user)
    if err != nil{
        return shim.Error(err.Error())
    }
    userAsJsonString = string(userAsJsonBytes)
    err = stub.PutState(userKey, []byte(userAsJsonString))
	if err != nil {
		return shim.Error(err.Error())
	}

	pollAsJsonBytes, err = getPollAsJsonByteArray(poll)
	if err != nil{
        return shim.Error(err.Error())
    }
    pollAsJsonString = string(pollAsJsonBytes)
	err = stub.PutState(pollKey, []byte(pollAsJsonString))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
//peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n mycc -c '{"Args":["changeStatusToZero","my_poll0"]}'
func (t *SimpleChaincode) changeStatusToZero(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	var pollKey, pollAsJsonString string
	var pollAsBytes []byte
	var poll Poll
	var statusOne bool
	var err error

	if len(args) != 1 {
		return shim.Error("put operation must include three arguments, a userKey, a pollKey and an option ")
	}
	
	pollKey = args[0]
	pollAsBytes, err = stub.GetState(pollKey)
    poll, err = getPollFromJsonByteArray(pollAsBytes)

    statusOne, err = isStatusOne(stub, poll)
    if statusOne == false{
    	return shim.Error("The status of the poll: "+ pollKey +" has already been changed")
    }

    poll.Status = poll.Status - 1

    pollAsBytes, err = getPollAsJsonByteArray(poll)
	if err != nil{
        return shim.Error(err.Error())
    }
    pollAsJsonString = string(pollAsBytes)
	err = stub.PutState(pollKey, []byte(pollAsJsonString))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
//TODO delete
func (t *SimpleChaincode) move(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//no need of this function
	return shim.Success(nil)
}

// Deletes an entity from state
//TODO modify
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
//TODO delete
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//no need

	return shim.Success(nil)
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

func isUserAllowed(stub shim.ChaincodeStubInterface, userKey string, pollKey string) (bool, int, error){

    var err error
    index :=0
    result := false

    userAsBytes, err := stub.GetState(userKey)

    user, err := getUserFromJsonByteArray(userAsBytes)

    for i := 0; i < len(user.Active); i++ {
        if user.Active[i].Name==pollKey {
        	if user.Active[i].Token >0{
	        	result = true
	        	index = i
	        }
        }
    }

    return result, index, err
}

func isPollAlreadyAddToUser(stub shim.ChaincodeStubInterface, userKey string, pollKey string) (bool, error){

    var err error
    result := false

    userAsBytes, err := stub.GetState(userKey)

    user, err := getUserFromJsonByteArray(userAsBytes)

    for i := 0; i < len(user.Active); i++ {
        if user.Active[i].Name==pollKey {
        	if user.Active[i].Token >0{
	        	result = true
	        }
        }
    }
    for i := 0; i < len(user.Inactive); i++ {
        if user.Inactive[i]==pollKey {
	        result = true
        }
    }
    return result, err
}

func isExistingOption(stub shim.ChaincodeStubInterface, pollKey string, option string) (bool, int, error){

    var err error
    indexOfTheOption := 0
    result := false
    pollAsBytes, err := stub.GetState(pollKey)
    poll, err := getPollFromJsonByteArray(pollAsBytes)

    for i := 0; i < len(poll.Options); i++ {
        if poll.Options[i].Name==option {
        	result = true
        	indexOfTheOption=i
        }
    }

    return result,indexOfTheOption, err
}

func isStatusOne(stub shim.ChaincodeStubInterface, poll Poll) (bool,error){
	var err error
	var result bool

	result = false
    if poll.Status == 1 {
    	result = true
    }

    return result, err
}

func getUserAsJsonByteArray(user User) ([]byte, error){

    var jsonUser []byte

    var err error

    jsonUser, err = json.Marshal(user)

    if err != nil{

        errors.New("Unable to convert user to json string")

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

        errors.New("Unable to convert poll to json string")

    }

    return jsonPoll, err

}


func getPollFromJsonByteArray(pollAsJsonByte []byte) (Poll, error){
    var poll Poll
    var err error

    err = json.Unmarshal(pollAsJsonByte, &poll)
    if err != nil{
        errors.New("Unable to convert poll to json string")
    }
    return poll, err
}

