package main

import( 
	"testing"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkAddUser(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("addUser", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkAddUser_Exists(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.ERROR {
		fmt.Println("addUser_Exists", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkNewQuery(t *testing.T, stub *shim.MockStub, name string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("newQuery"), []byte(name)})
	if res.Status != shim.OK {
		fmt.Println("newQuery", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("newQuery", name, "failed to get value")
		t.FailNow()
	}
	if string(res.Payload) != value {
		fmt.Println("newQuery value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkNewQuery_NoKey(t *testing.T, stub *shim.MockStub, name string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("newQuery"), []byte(name)})
	if res.Status != shim.ERROR {
		fmt.Println("newQuery_NoKey", name, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkAddNewPoll(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("addNewPoll", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkAddNewPoll_Exists(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.ERROR {
		fmt.Println("addNewPoll", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkAddNewActivePollToUser(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("addNewActivePollToUser", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkAddNewActivePollToUser_NoPoll(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.ERROR {
		fmt.Println("addNewActivePollToUser_NoPoll", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkAddNewActivePollToUser_NoUser(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.ERROR {
		fmt.Println("addNewActivePollToUser_NoUser", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkAddNewActivePollToUser_AlreadyAdded(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.ERROR {
		fmt.Println("checkAddNewActivePollToUser_AlreadyAdded", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkActiveToInactivePoll(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("ActiveToInactivePoll", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkActiveToInactivePoll_NoUser(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.ERROR {
		fmt.Println("ActiveToInactivePoll_NoUser", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkActiveToInactivePoll_NoPollForUser(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.ERROR {
		fmt.Println("ActiveToInactivePoll_NoPollForUser", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkVote(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("vote", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkVote_NoUser(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.ERROR {
		fmt.Println("vote_NoUser", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkVote_NoPoll(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.ERROR {
		fmt.Println("vote_NoPoll", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkVote_NotAllowed(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.ERROR {
		fmt.Println("vote_NotAllowed", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkVote_VoteClosed(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.ERROR {
		fmt.Println("vote_VoteClosed", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkVote_HasAlreadyVoted(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.ERROR {
		fmt.Println("vote_HasAlreadyVotes", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkVote_NoOption(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.ERROR {
		fmt.Println("vote_NoOption", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkChangeStatusToZero(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("changeStatusToZero", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkChangeStatusToZero_NoPoll(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.ERROR {
		fmt.Println("changeStatusToZero", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkChangeStatusToZero_StatusAlready0(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.ERROR {
		fmt.Println("changeStatusToZero", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func TestExample02_AddUser(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddUser(t, stub, [][]byte{[]byte("addUser"), []byte("c")})
	checkState(t, stub, "c", "{\"active\":[],\"inactive\":[]}")
}

func TestExample02_AddUser_Exists(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddUser(t, stub, [][]byte{[]byte("addUser"), []byte("c")})
	checkAddUser_Exists(t, stub, [][]byte{[]byte("addUser"), []byte("c")})
	checkState(t, stub, "c", "{\"active\":[],\"inactive\":[]}")
}

func TestExample02_NewQuery(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddUser(t, stub, [][]byte{[]byte("addUser"), []byte("c")})
	// Query A
	checkNewQuery(t, stub, "c", "{\"active\":[],\"inactive\":[]}")
}

func TestExample02_NewQuery_NoKey(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddUser(t, stub, [][]byte{[]byte("addUser"), []byte("d")})
	// Query A
	checkNewQuery_NoKey(t, stub, "c", "{\"active\":[],\"inactive\":[]}")
}

func TestExample02_AddNewPoll(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddNewPoll(t, stub, [][]byte{[]byte("addNewPoll"), []byte("mypoll"), []byte("{\"Options\":[{\"Name\":\"opt1\",\"Count\":0},{\"Name\":\"opt2\",\"Count\":0}],\"status\":1}")})
	checkNewQuery(t, stub, "mypoll", "{\"options\":[{\"name\":\"opt1\",\"count\":0},{\"name\":\"opt2\",\"count\":0}],\"status\":1}")
}

func TestExample02_AddNewPoll_Exists(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddNewPoll(t, stub, [][]byte{[]byte("addNewPoll"), []byte("mypoll"), []byte("{\"Options\":[{\"Name\":\"opt1\",\"Count\":0},{\"Name\":\"opt2\",\"Count\":0}],\"status\":1}")})
	checkAddNewPoll_Exists(t, stub, [][]byte{[]byte("addNewPoll"), []byte("mypoll"), []byte("{\"Options\":[{\"Name\":\"opt3\",\"Count\":0},{\"Name\":\"opt4\",\"Count\":0}],\"status\":1}")})
	checkNewQuery(t, stub, "mypoll", "{\"options\":[{\"name\":\"opt1\",\"count\":0},{\"name\":\"opt2\",\"count\":0}],\"status\":1}")
}

func TestExample02_AddNewActivePollToUser(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddNewPoll(t, stub, [][]byte{[]byte("addNewPoll"), []byte("mypoll"), []byte("{\"Options\":[{\"Name\":\"opt1\",\"Count\":0},{\"Name\":\"opt2\",\"Count\":0}],\"status\":1}")})
	checkAddUser(t, stub, [][]byte{[]byte("addUser"), []byte("c")})
	checkAddNewActivePollToUser(t, stub, [][]byte{[]byte("addNewActivePollToUser"),[]byte("c"), []byte("mypoll")})
	checkNewQuery(t, stub, "c", "{\"active\":[{\"name\":\"mypoll\",\"token\":1}],\"inactive\":[]}")
}

func TestExample02_AddNewActivePollToUser_NoPoll(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddUser(t, stub, [][]byte{[]byte("addUser"), []byte("c")})
	checkAddNewActivePollToUser_NoPoll(t, stub, [][]byte{[]byte("addNewActivePollToUser"),[]byte("c"), []byte("mypoll")})
	checkNewQuery(t, stub, "c", "{\"active\":[],\"inactive\":[]}")
}

func TestExample02_AddNewActivePollToUser_NoUser(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddNewPoll(t, stub, [][]byte{[]byte("addNewPoll"), []byte("mypoll"), []byte("{\"Options\":[{\"Name\":\"opt1\",\"Count\":0},{\"Name\":\"opt2\",\"Count\":0}],\"status\":1}")})
	checkAddNewActivePollToUser_NoUser(t, stub, [][]byte{[]byte("addNewActivePollToUser"),[]byte("c"), []byte("mypoll0")})
}

func TestExample02_AddNewActivePollToUser_AlreadyAdded(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddNewPoll(t, stub, [][]byte{[]byte("addNewPoll"), []byte("mypoll"), []byte("{\"Options\":[{\"Name\":\"opt1\",\"Count\":0},{\"Name\":\"opt2\",\"Count\":0}],\"status\":1}")})
	checkAddUser(t, stub, [][]byte{[]byte("addUser"), []byte("c")})
	checkAddNewActivePollToUser(t, stub, [][]byte{[]byte("addNewActivePollToUser"),[]byte("c"), []byte("mypoll")})
	checkAddNewActivePollToUser_AlreadyAdded(t, stub, [][]byte{[]byte("addNewActivePollToUser"),[]byte("c"), []byte("mypoll")})
}

func TestExample02_ActiveToInactivePoll(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddNewPoll(t, stub, [][]byte{[]byte("addNewPoll"), []byte("mypoll"), []byte("{\"Options\":[{\"Name\":\"opt1\",\"Count\":0},{\"Name\":\"opt2\",\"Count\":0}],\"status\":1}")})
	checkAddUser(t, stub, [][]byte{[]byte("addUser"), []byte("c")})
	checkAddNewActivePollToUser(t, stub, [][]byte{[]byte("addNewActivePollToUser"),[]byte("c"), []byte("mypoll")})
	checkActiveToInactivePoll(t, stub, [][]byte{[]byte("activeToInactivePoll"),[]byte("c"), []byte("mypoll")})
	checkNewQuery(t, stub, "c", "{\"active\":[],\"inactive\":[\"mypoll\"]}")
}

func TestExample02_ActiveToInactivePoll_NoUser(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddNewPoll(t, stub, [][]byte{[]byte("addNewPoll"), []byte("mypoll"), []byte("{\"Options\":[{\"Name\":\"opt1\",\"Count\":0},{\"Name\":\"opt2\",\"Count\":0}],\"status\":1}")})
	checkAddUser(t, stub, [][]byte{[]byte("addUser"), []byte("c")})
	checkAddNewActivePollToUser(t, stub, [][]byte{[]byte("addNewActivePollToUser"),[]byte("c"), []byte("mypoll")})
	checkActiveToInactivePoll_NoUser(t, stub, [][]byte{[]byte("activeToInactivePoll"),[]byte("d"), []byte("mypoll")})
}

func TestExample02_ActiveToInactivePoll_NoPollForUser(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddNewPoll(t, stub, [][]byte{[]byte("addNewPoll"), []byte("mypoll"), []byte("{\"Options\":[{\"Name\":\"opt1\",\"Count\":0},{\"Name\":\"opt2\",\"Count\":0}],\"status\":1}")})
	checkAddUser(t, stub, [][]byte{[]byte("addUser"), []byte("c")})
	checkAddNewActivePollToUser(t, stub, [][]byte{[]byte("addNewActivePollToUser"),[]byte("c"), []byte("mypoll")})
	checkActiveToInactivePoll_NoPollForUser(t, stub, [][]byte{[]byte("activeToInactivePoll"),[]byte("c"), []byte("mypoll0")})
	checkNewQuery(t, stub, "c", "{\"active\":[{\"name\":\"mypoll\",\"token\":1}],\"inactive\":[]}")
}

func TestExample02_Vote(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddNewPoll(t, stub, [][]byte{[]byte("addNewPoll"), []byte("mypoll"), []byte("{\"Options\":[{\"Name\":\"opt1\",\"Count\":0},{\"Name\":\"opt2\",\"Count\":0}],\"status\":1}")})
	checkAddUser(t, stub, [][]byte{[]byte("addUser"), []byte("c")})
	checkAddNewActivePollToUser(t, stub, [][]byte{[]byte("addNewActivePollToUser"),[]byte("c"), []byte("mypoll")})
	checkVote(t, stub, [][]byte{[]byte("vote"),[]byte("c"), []byte("mypoll"), []byte("opt1")})
	checkChangeStatusToZero(t, stub, [][]byte{[]byte("changeStatusToZero"), []byte("mypoll")})
	checkNewQuery(t, stub, "c", "{\"active\":[{\"name\":\"mypoll\",\"token\":0}],\"inactive\":[]}")
	checkNewQuery(t, stub, "mypoll", "{\"options\":[{\"name\":\"opt1\",\"count\":1},{\"name\":\"opt2\",\"count\":0}],\"status\":0}")
}

func TestExample02_Vote_NoUser(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddNewPoll(t, stub, [][]byte{[]byte("addNewPoll"), []byte("mypoll"), []byte("{\"Options\":[{\"Name\":\"opt1\",\"Count\":0},{\"Name\":\"opt2\",\"Count\":0}],\"status\":1}")})
	checkAddUser(t, stub, [][]byte{[]byte("addUser"), []byte("c")})
	checkAddNewActivePollToUser(t, stub, [][]byte{[]byte("addNewActivePollToUser"),[]byte("c"), []byte("mypoll")})
	checkVote_NoUser(t, stub, [][]byte{[]byte("vote"),[]byte("d"), []byte("mypoll"), []byte("opt1")})

	checkNewQuery(t, stub, "c", "{\"active\":[{\"name\":\"mypoll\",\"token\":1}],\"inactive\":[]}")
	checkNewQuery(t, stub, "mypoll", "{\"options\":[{\"name\":\"opt1\",\"count\":0},{\"name\":\"opt2\",\"count\":0}],\"status\":1}")
}

func TestExample02_Vote_NoPoll(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddNewPoll(t, stub, [][]byte{[]byte("addNewPoll"), []byte("mypoll"), []byte("{\"Options\":[{\"Name\":\"opt1\",\"Count\":0},{\"Name\":\"opt2\",\"Count\":0}],\"status\":1}")})
	checkAddUser(t, stub, [][]byte{[]byte("addUser"), []byte("c")})
	checkAddNewActivePollToUser(t, stub, [][]byte{[]byte("addNewActivePollToUser"),[]byte("c"), []byte("mypoll")})
	checkVote_NoPoll(t, stub, [][]byte{[]byte("vote"),[]byte("c"), []byte("mypoll0"), []byte("opt1")})

	checkNewQuery(t, stub, "c", "{\"active\":[{\"name\":\"mypoll\",\"token\":1}],\"inactive\":[]}")
	checkNewQuery(t, stub, "mypoll", "{\"options\":[{\"name\":\"opt1\",\"count\":0},{\"name\":\"opt2\",\"count\":0}],\"status\":1}")
}

func TestExample02_Vote_NotAllowed(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddNewPoll(t, stub, [][]byte{[]byte("addNewPoll"), []byte("mypoll"), []byte("{\"Options\":[{\"Name\":\"opt1\",\"Count\":0},{\"Name\":\"opt2\",\"Count\":0}],\"status\":1}")})
	checkAddUser(t, stub, [][]byte{[]byte("addUser"), []byte("c")})
	checkAddNewActivePollToUser(t, stub, [][]byte{[]byte("addNewActivePollToUser"),[]byte("c"), []byte("mypoll")})

	checkVote_NotAllowed(t, stub, [][]byte{[]byte("vote"),[]byte("c"), []byte("mypoll0"), []byte("opt1")})
	checkNewQuery(t, stub, "mypoll", "{\"options\":[{\"name\":\"opt1\",\"count\":0},{\"name\":\"opt2\",\"count\":0}],\"status\":1}")
}

func TestExample02_Vote_VoteClosed(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddNewPoll(t, stub, [][]byte{[]byte("addNewPoll"), []byte("mypoll"), []byte("{\"Options\":[{\"Name\":\"opt1\",\"Count\":0},{\"Name\":\"opt2\",\"Count\":0}],\"status\":1}")})
	checkAddUser(t, stub, [][]byte{[]byte("addUser"), []byte("c")})
	checkAddNewActivePollToUser(t, stub, [][]byte{[]byte("addNewActivePollToUser"),[]byte("c"), []byte("mypoll")})
	checkChangeStatusToZero(t, stub, [][]byte{[]byte("changeStatusToZero"), []byte("mypoll")})
	checkVote_VoteClosed(t, stub, [][]byte{[]byte("vote"),[]byte("c"), []byte("mypoll"), []byte("opt1")})
	checkNewQuery(t, stub, "c", "{\"active\":[{\"name\":\"mypoll\",\"token\":1}],\"inactive\":[]}")
	checkNewQuery(t, stub, "mypoll", "{\"options\":[{\"name\":\"opt1\",\"count\":0},{\"name\":\"opt2\",\"count\":0}],\"status\":0}")
}

func TestExample02_Vote_HasAlreadyVoted(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddNewPoll(t, stub, [][]byte{[]byte("addNewPoll"), []byte("mypoll"), []byte("{\"Options\":[{\"Name\":\"opt1\",\"Count\":0},{\"Name\":\"opt2\",\"Count\":0}],\"status\":1}")})
	checkAddUser(t, stub, [][]byte{[]byte("addUser"), []byte("c")})
	checkAddNewActivePollToUser(t, stub, [][]byte{[]byte("addNewActivePollToUser"),[]byte("c"), []byte("mypoll")})
	checkVote(t, stub, [][]byte{[]byte("vote"),[]byte("c"), []byte("mypoll"), []byte("opt1")})
	checkVote_HasAlreadyVoted(t, stub, [][]byte{[]byte("vote"),[]byte("c"), []byte("mypoll"), []byte("opt1")})
	checkNewQuery(t, stub, "c", "{\"active\":[{\"name\":\"mypoll\",\"token\":0}],\"inactive\":[]}")
	checkChangeStatusToZero(t, stub, [][]byte{[]byte("changeStatusToZero"), []byte("mypoll")})
	checkNewQuery(t, stub, "mypoll", "{\"options\":[{\"name\":\"opt1\",\"count\":1},{\"name\":\"opt2\",\"count\":0}],\"status\":0}")
}

func TestExample02_Vote_NoOption(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddNewPoll(t, stub, [][]byte{[]byte("addNewPoll"), []byte("mypoll"), []byte("{\"Options\":[{\"Name\":\"opt1\",\"Count\":0},{\"Name\":\"opt2\",\"Count\":0}],\"status\":1}")})
	checkAddUser(t, stub, [][]byte{[]byte("addUser"), []byte("c")})
	checkAddNewActivePollToUser(t, stub, [][]byte{[]byte("addNewActivePollToUser"),[]byte("c"), []byte("mypoll")})
	checkVote_NoOption(t, stub, [][]byte{[]byte("vote"),[]byte("c"), []byte("mypoll"), []byte("opt3")})
	checkNewQuery(t, stub, "c", "{\"active\":[{\"name\":\"mypoll\",\"token\":1}],\"inactive\":[]}")
	checkNewQuery(t, stub, "mypoll", "{\"options\":[{\"name\":\"opt1\",\"count\":0},{\"name\":\"opt2\",\"count\":0}],\"status\":1}")
}

func TestExample02_ChangeStatusToZero(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddNewPoll(t, stub, [][]byte{[]byte("addNewPoll"), []byte("mypoll"), []byte("{\"Options\":[{\"Name\":\"opt1\",\"Count\":0},{\"Name\":\"opt2\",\"Count\":0}],\"status\":1}")})
	checkAddUser(t, stub, [][]byte{[]byte("addUser"), []byte("c")})
	checkAddNewActivePollToUser(t, stub, [][]byte{[]byte("addNewActivePollToUser"),[]byte("c"), []byte("mypoll")})
	checkChangeStatusToZero(t, stub, [][]byte{[]byte("changeStatusToZero"), []byte("mypoll")})
	checkNewQuery(t, stub, "mypoll", "{\"options\":[{\"name\":\"opt1\",\"count\":0},{\"name\":\"opt2\",\"count\":0}],\"status\":0}")
}

func TestExample02_ChangeStatusToZero_NoPoll(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddNewPoll(t, stub, [][]byte{[]byte("addNewPoll"), []byte("mypoll"), []byte("{\"Options\":[{\"Name\":\"opt1\",\"Count\":0},{\"Name\":\"opt2\",\"Count\":0}],\"status\":1}")})
	checkAddUser(t, stub, [][]byte{[]byte("addUser"), []byte("c")})
	checkAddNewActivePollToUser(t, stub, [][]byte{[]byte("addNewActivePollToUser"),[]byte("c"), []byte("mypoll")})
	checkChangeStatusToZero_NoPoll(t, stub, [][]byte{[]byte("changeStatusToZero"), []byte("mypoll0")})
	checkNewQuery(t, stub, "mypoll", "{\"options\":[{\"name\":\"opt1\",\"count\":0},{\"name\":\"opt2\",\"count\":0}],\"status\":1}")
}

func TestExample02_ChangeStatusToZero_StatusAlready0(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	checkAddNewPoll(t, stub, [][]byte{[]byte("addNewPoll"), []byte("mypoll"), []byte("{\"Options\":[{\"Name\":\"opt1\",\"Count\":0},{\"Name\":\"opt2\",\"Count\":0}],\"status\":1}")})
	checkAddUser(t, stub, [][]byte{[]byte("addUser"), []byte("c")})
	checkAddNewActivePollToUser(t, stub, [][]byte{[]byte("addNewActivePollToUser"),[]byte("c"), []byte("mypoll")})
	checkChangeStatusToZero(t, stub, [][]byte{[]byte("changeStatusToZero"), []byte("mypoll")})
	checkChangeStatusToZero_StatusAlready0(t, stub, [][]byte{[]byte("changeStatusToZero"), []byte("mypoll")})
	checkNewQuery(t, stub, "mypoll", "{\"options\":[{\"name\":\"opt1\",\"count\":0},{\"name\":\"opt2\",\"count\":0}],\"status\":0}")
}