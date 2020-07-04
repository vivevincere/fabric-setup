package goosefabric

import(
	"log"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"

)

type SdkObjects struct{
	sdk *fabsdk.FabricSDK
	client *channel.Client
}

// Creates a new key variable to pass into SDK query function
// id and timestamp forms a composite key
func NewKey(id string, timestamp string)[][]byte{
return [][]byte{[]byte(id),[]byte(timestamp)}
}

// Creates a new key + payload/goosepacket variable to pass into SDK execute function
// id and timestamp for composite key, goosepacket is the payload i.e. value part of the key/value pair
func NewPayload(id string, timestamp string, goosePacket string) [][]byte{
return [][]byte{[]byte(id),[]byte(timestamp), []byte(goosePacket)}
}


// Initializes the fabric sdk objects with given parameters and returns a SdkObjects struct
func Init(configFile string, channelID string, user string, org string) SdkObjects{
	var toReturn SdkObjects
	var err error
	toReturn.sdk, err = fabsdk.New(config.FromFile(configFile))		// init fabsdk by loading from config.yaml location
	if err != nil {
		log.Fatal("failed to create sdk: %v", err)
	}
	clientContext := toReturn.sdk.ChannelContext(channelID,fabsdk.WithUser(user), fabsdk.WithOrg(org))	// loads from preexisting channel with the relevant user and org
	toReturn.client, err = channel.New(clientContext)
	if err != nil {
		log.Fatal("Failed to create new channel: %v", err)
	}
	return toReturn
}


// Creates or updates a key/value pair
// chaincodeID is the name of the chaincode that is installed
func (thesdk SdkObjects) LogEvent( chaincodeID string, id string, timestamp string, goosePacket string) error{
	defArgs := NewPayload(id, timestamp, goosePacket)	// creates id+timestamp composite key and goosePacket payload
	_, err := thesdk.client.Execute(channel.Request{ChaincodeID: chaincodeID, Fcn: "LogEvent", Args: defArgs})	// SDK invoke chaincode function
	if err != nil{
		return err	
	}
	return  nil
}

//Queries an existing key/value pair
// chaincodeID is the name of the chaincode that is installed
func (thesdk SdkObjects) Get(chaincodeID string, id string, timestamp string) (string , error){
	defArgs := NewKey(id, timestamp)	// Creates composite key with id+timestamp
	response, err := thesdk.client.Query(channel.Request{ChaincodeID: chaincodeID, Fcn: "QueryEvent", Args: defArgs})	// SDK invoke chaincode function for queries
	if err != nil{
		return "", err
	}
	k := string(response.Payload)
	return k, nil
}