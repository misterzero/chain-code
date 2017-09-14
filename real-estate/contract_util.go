package main

import (
	"strings"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"errors"
)

//TODO find better way to process errors

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
