
#!/bin/bash

cleanup () {
    rm -rf ../artifacts/channel/crypto-config
    rm  ../artifacts/channel/*.tx
    rm ../artifacts/channel/*.block
    rm ../channel-artifacts/*.block
    docker-compose down
}



if [ "$1" = "cleanup" ] 
then
	cleanup
else
	./create-artifacts.sh
	docker-compose up -d
	sleep 5s
	./createChannel.sh
	sudo ./deployChaincode.sh
fi
