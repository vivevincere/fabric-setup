package main

import(
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	// "github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"

)

func main(){
	sdk, err := fabsdk.New(config.FromFile("/Users/vincentchung/cleanslate/basic-setup-fabric/config.yaml"))
	if err != nil {
		fmt.Printf("failed to create sdk: %v", err)
	}
	clientContext := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org1"))
	orgResMgmt, err := resmgmt.New(clientContext)
	if err != nil {
		fmt.Printf("Creation of resmgmt client failed: %s", err)
	}

	ccPkg, err := packager.NewCCPackage("github.com/vivevincere/chaincodes/chaincode.go","/Users/vincentchung/go")
	if err != nil{
		fmt.Printf("CC Package finder error: %v",err)
	}
	installCCReq := resmgmt.InstallCCRequest{Name: "c", Path: "github.com/vivevincere/chaincodes/chaincode.go", Version: "0", Package: ccPkg}
	_, err = orgResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		fmt.Printf("Install CC failed error: %v", err)
	}



}