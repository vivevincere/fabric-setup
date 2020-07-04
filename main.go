package main

import(
	"fmt"

	"github.com/vivevincere/goosefabric"
	"time"
	"log"

)

func main(){

	objects := goosefabric.Init("/Users/vincentchung/cleanslate/thefix/config.yaml", "mychannel", "Admin", "Org1")
	err := objects.LogEvent("base", "faucet","1234", "6d61697961686565")
	if err != nil{
		log.Fatal("LogEvent failed: %v", err)
	}
	time.Sleep(2 * time.Second)
	k, err := objects.Get("base", "faucet","1234")
	if err != nil{
		log.Fatal("QueryEvent failed: %v", err)
	}
	fmt.Printf(k)	

}