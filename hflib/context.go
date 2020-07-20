package hflib

import (
	"log"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// Context contains the relevant go structs for interacting with Hyperledger Fabric. It
// currently includes:
// - The chaincode ID
// - A pointer to the fabric SDK object
// - A pointer to the channel client
type Context struct {
	ChaincodeID string
	SDK         *fabsdk.FabricSDK
	Client      *channel.Client
}

// Config is a helper struct that holds the configuration required to create a Context
// using the Init function
type Config struct {
	ConfigFile  string
	ChaincodeID string
	ChannelID   string
	User        string
	Org         string
}

// Args is a helper type that matches the expected type when passing arguments to
// the fabric SDK Args paraemeter when executing or querying chaincode.
type Args [][]byte

func ThreeArgs(first string, second string, third string) Args {
	return [][]byte{[]byte(first),[]byte(second),[]byte(third)}
}
// NewQueryEventArgs takes an id and timestamp for an event then converts it into the
// expected format for arguments when executing the "QueryEvent" chaincode.
func NewQueryEventArgs(id string, timestamp string) Args {
	return [][]byte{[]byte(id), []byte(timestamp)}
}

// NewLogEventArgs takes and id, timestamp, and data for an eventthen converts into
// the expected format for arguments when executing the "LogEvent" chaincode
func NewLogEventArgs(id string, timestamp string, data string) Args {
	return [][]byte{[]byte(id), []byte(timestamp), []byte(data)}
}

// Init will create the fabric sdk context with given parameters. A context contains the
// following values:
// - A reference to the Fabric SDK
// - A reference to the channel client
func Init(c *Config) Context {
	// init fabsdk by loading from config.yaml location
	sdk, err := fabsdk.New(config.FromFile(c.ConfigFile))
	if err != nil {
		log.Fatalf("failed to create sdk: %v", err)
	}

	// loads from preexisting channel with the relevant user and org
	clientContext := sdk.ChannelContext(c.ChannelID, fabsdk.WithUser(c.User), fabsdk.WithOrg(c.Org))
	client, err := channel.New(clientContext)
	if err != nil {
		log.Fatalf("Failed to create new channel: %v", err)
	}

	return Context{
		ChaincodeID: c.ChaincodeID,
		SDK:         sdk,
		Client:      client,
	}
}

// LogEvent creates or updates a key/value pair
// chaincodeID is the name of the chaincode that is installed
func (c *Context) LogEvent(id string, timestamp string, goosePacket string) error {
	// Creates arguments to pass to "LogEvent" chaincode
	args := NewLogEventArgs(id, timestamp, goosePacket)

	// Create request
	req := channel.Request{
		ChaincodeID: c.ChaincodeID,
		Fcn:         "LogEvent",
		Args:        args,
	}

	// SDK invoke chaincode function
	_, err := c.Client.Execute(req)
	if err != nil {
		return err
	}

	return nil
}

// QueryEvent queries an existing key/value pair
// chaincodeID is the name of the chaincode that is installed
func (c *Context) QueryEvent(id string, timestamp string) (string, error) {
	// Creates arguments to pass to "QueryEvent" chaincode
	args := NewQueryEventArgs(id, timestamp)

	// Create request
	req := channel.Request{
		ChaincodeID: c.ChaincodeID,
		Fcn:         "QueryEvent",
		Args:        args,
	}

	// SDK invoke chaincode function for queries
	response, err := c.Client.Query(req)
	if err != nil {
		return "", err
	}

	// Convert to string
	k := string(response.Payload)
	return k, nil
}

// Pass start and end timestamps of desired range as strings
// If start is empty string, will search with no lower bound
// If end is empty string, will search with no upper bound
// If both empty string, will return all elements
func (c *Context) QueryAllByDateRange(start string, end string) (string, error) {
	// Creates arguments to pass to "QueryAllByDateRange" chaincode
	args := NewQueryEventArgs(start, end)

	// Create request
	req := channel.Request{
		ChaincodeID: c.ChaincodeID,
		Fcn:         "QueryAllByDateRange",
		Args:        args,
	}

	// SDK invoke chaincode function for queries
	response, err := c.Client.Query(req)
	if err != nil {
		return "", err
	}

	// Convert to string
	k := string(response.Payload)
	return k, nil
}

// Pass start and end timestamps of desired range as strings and deviceID
// If start is empty string, will search with no lower bound
// If end is empty string, will search with no upper bound
// If both empty string, will return all elements
func (c *Context) QueryDeviceByDateRange(start string, end string, deviceID string) (string, error) {
	// Creates arguments to pass to "QueryDeviceByDateRange" chaincode
	args := ThreeArgs(start, end, deviceID)

	// Create request
	req := channel.Request{
		ChaincodeID: c.ChaincodeID,
		Fcn:         "QueryDeviceByDateRange",
		Args:        args,
	}

	// SDK invoke chaincode function for queries
	response, err := c.Client.Query(req)
	if err != nil {
		return "", err
	}

	// Convert to string
	k := string(response.Payload)
	return k, nil
}