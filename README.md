# fabric-setup

Hyperledger Fabric Network for storing and querying GOOSE messages.

HF Network: 
Use run.sh script to start the network containers with default settings, install the chaincode and instantiate all the nodes. 

Use same script with -clean option to take down the network.

Default configs set up a single orderer, single peer & couchDB network. Edit the relevant config files to change settings.

Chaincode can be found in chaincode folder. Supports basic query/store & various query by date range. Custom couchDB queries can be made.

HFLib:

Golang library to enable a golang client to interact with the Hyperledger Fabric Network and invoke the various chaincode functions.

Uses the Hyperledger Fabric golang SDK.



