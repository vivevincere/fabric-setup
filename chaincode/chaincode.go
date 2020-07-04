package main

import (
    "encoding/hex"
    "fmt"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const (
    indexName = "GooseEvent"
)

// DltAudit provides functions for logging GOOSE events to the ledger
type DltAudit struct {
    contractapi.Contract
}

// InitLedger provides initialization
func (d *DltAudit) InitLedger(ctx contractapi.TransactionContextInterface) error {
    return nil
}

// LogEvent will log a new GOOSE event to the ledger
func (d *DltAudit) LogEvent(ctx contractapi.TransactionContextInterface, id string, timestamp string, goosePacket string) error {
    key, err := ctx.GetStub().CreateCompositeKey(indexName, []string{id, string(timestamp)})
    if err != nil {
        return fmt.Errorf("Error creating composite key for GooseEvent")
    }

    packetBytes, err := hex.DecodeString(goosePacket)
    if err != nil {
        return fmt.Errorf("Error decoding the GOOSE packet (expected valid hex string)")
    }

    err = ctx.GetStub().PutState(key, packetBytes)
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

    return hex.EncodeToString(eventAsBytes), nil
}

func main(){
    dltAudit := new(DltAudit)
    contract,err := contractapi.NewChaincode(dltAudit)

    if err != nil {
        panic(err.Error())
    }

    if err := contract.Start(); err != nil{
        panic(err.Error())
    }
}