package main

import (
	"os"
	"bufio"
	"fmt"
	"github.com/alecthomas/jsonschema"
	"encoding/json"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func buildPropertyJson() {
	propertyFileName, err := os.Create("Property.json")
	check(err)

	defer propertyFileName.Close()

	bufferedWriter := bufio.NewWriter(propertyFileName)

	jsonProperty := jsonschema.Reflect(&Property{})
	property, _ := json.MarshalIndent(jsonProperty, "", "    ")
	fmt.Fprintf(bufferedWriter, string(property))

	bufferedWriter.Flush()
}

func buildOwnershipJson(){
	ownershipFileName, err := os.Create("Ownership.json")
	check(err)

	defer ownershipFileName.Close()

	bufferedWriter := bufio.NewWriter(ownershipFileName)

	jsonOwnership := jsonschema.Reflect(&Ownership{})
	ownership, _ := json.MarshalIndent(jsonOwnership, "", "    ")
	fmt.Fprintf(bufferedWriter, string(ownership))

	bufferedWriter.Flush()
}


