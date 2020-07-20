package main

import (
	"encoding/hex"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"bytes"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"encoding/json"
	"strconv"
)
const (
	indexName = "id~timestamp"
)

// DltAudit provides functions for logging GOOSE events to the ledger
type DltAudit struct {
	contractapi.Contract
}


type EventMessage struct {
	Timestamp int `json:"timestamp"`
	DeviceID string `json:"id"`
	Message string `json:"message"`
}



// Helper function to construct response from SDK iterator
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false

	
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}

// InitLedger provides initialization
func (d *DltAudit) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

// LogEvent will log a new GOOSE event to the ledger
func (d *DltAudit) LogEvent(ctx contractapi.TransactionContextInterface, id string, timestamp string, goosePacket string) error {
	packetBytes, err := hex.DecodeString(goosePacket)
	if err != nil {
		return fmt.Errorf("Error decoding the GOOSE packet (expected valid hex string)")
	}
	timeInt, _ := strconv.Atoi(timestamp)
	theinput := &EventMessage{timeInt,id,goosePacket}
	EventAsJSONBytes , err := json.Marshal(theinput)
	if err != nil {
		return fmt.Errorf("Error marshalling into json object: %v", err)
	}

	key, err := ctx.GetStub().CreateCompositeKey(indexName, []string{id, string(timestamp)})
	if err != nil {
		return fmt.Errorf("Error creating composite key for GooseEvent")
	}

	err = ctx.GetStub().PutState(key, EventAsJSONBytes)
	if err != nil {
		return fmt.Errorf("Error putting key '%s' and value '%s'", key, packetBytes)
	}

	return ctx.GetStub().SetEvent("logEvent", packetBytes)
}

// QueryEvent gets the event for a given IED ID and timestamp
func (d *DltAudit) QueryEvent(ctx contractapi.TransactionContextInterface, id string, timestamp string) (string, error) {

	key, err := ctx.GetStub().CreateCompositeKey(indexName, []string{id, string(timestamp)})
	if err != nil {
		return "", fmt.Errorf("Error creating composite key for GooseEvent")
	}

	eventAsBytes, err := ctx.GetStub().GetState(key)

	if err != nil {
		return "", fmt.Errorf("Failed to read: %s", err.Error())
	}

	if eventAsBytes == nil {
		return "", fmt.Errorf("%s does not exist", id)
	}

	return string(eventAsBytes), nil
}


// Pass start and end timestamps of desired range as strings 
// If start is empty string, will search with no lower bound
// If end is empty string, will search with no upper bound
// If both empty string, will return all elements
func (d *DltAudit) QueryAllByDateRange(ctx contractapi.TransactionContextInterface, start string, end string) (string, error) {
	var queryString string
	startInt,_ := strconv.Atoi(start)
	endInt,_ := strconv.Atoi(end)
	// both empty string means query all
	if start == "" && end == ""{
		queryString = fmt.Sprintf("{\"selector\":{\"timestamp\":{\"$gte\": 0}}}")
	} else if start == "" {
		queryString = fmt.Sprintf("{\"selector\":{\"timestamp\":{ \"$lte\": %d}}}", endInt)
		//query up to end
	} else if end == ""{
		queryString = fmt.Sprintf("{\"selector\":{\"timestamp\":{\"$gte\": %d}}}", startInt)
		// query from start
	} else{
		queryString = fmt.Sprintf("{\"selector\":{\"timestamp\":{\"$gte\": %d,\"$lte\": %d}}}", startInt, endInt)
		//query startid to end
	}


	// SDK Rich query function
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil{
		return "" , err
	}

	defer resultsIterator.Close()


	// Helper function to construct a response from returned iterator
	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil{
		return "", err
	}

	return buffer.String(), nil
}


// Pass start and end timestamps of desired range and deviceID as strings
// If start is empty string, will search with no lower bound
// If end is empty string, will search with no upper bound
// If both empty string, will return all elements
func (d *DltAudit) QueryDeviceByDateRange(ctx contractapi.TransactionContextInterface, start string, end string, deviceID string) (string, error) {
	var queryString string
	startInt,_ := strconv.Atoi(start)
	endInt,_ := strconv.Atoi(end)

	//both empty string means query all logs from this device
	if start == "" && end == ""{
		queryString = fmt.Sprintf("{\"selector\":{\"timestamp\":{\"$gte\": 0} , \"id\": \"%s\"}}", deviceID)
	} else if start == "" {
		queryString = fmt.Sprintf("{\"selector\":{\"timestamp\":{ \"$lte\": %d}, \"id\": \"%s\"}}", endInt, deviceID)
		//query up to end
	} else if end == ""{
		queryString = fmt.Sprintf("{\"selector\":{\"timestamp\":{\"$gte\": %d}, \"id\": \"%s\"}}", startInt, deviceID)
		// query from start
	} else{
		queryString = fmt.Sprintf("{\"selector\":{\"timestamp\":{\"$gte\": %d,\"$lte\": %d}, \"id\": \"%s\"}}", startInt, endInt, deviceID)
		//query startid to end
	}

	// SDK rich query function
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil{
		return "" , err
	}

	defer resultsIterator.Close()

	// Helper function to construct response from iterator
	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil{
		return "", err
	}

	return buffer.String(), nil
}



func main() {
	dltAudit := new(DltAudit)
	contract, err := contractapi.NewChaincode(dltAudit)

	if err != nil {
		panic(err.Error())
	}

	if err := contract.Start(); err != nil {
		panic(err.Error())
	}
}
