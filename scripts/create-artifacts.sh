
chmod -R 0755 ./crypto-config
# Delete existing artifacts
rm -rf ./crypto-config
rm genesis.block mychannel.tx
rm -rf ../../channel-artifacts/*

#Generate Crypto artifactes for organizations
cryptogen generate --config=../configs/crypto-config.yaml --output=../artifacts/channel/crypto-config/



# System channel
SYS_CHANNEL="sys-channel"

# channel name defaults to "mychannel"
CHANNEL_NAME="mychannel"

echo $CHANNEL_NAME

# Generate System Genesis block
configtxgen -profile OrdererGenesis -configPath ../artifacts/channel -channelID $SYS_CHANNEL  -outputBlock ../artifacts/channel/genesis.block


# Generate channel configuration block
configtxgen -profile BasicChannel -configPath ../artifacts/channel -outputCreateChannelTx ../artifacts/channel/mychannel.tx -channelID $CHANNEL_NAME

# echo "#######    Generating anchor peer update for Org1MSP  ##########"
configtxgen -profile BasicChannel -configPath ../artifacts/channel -outputAnchorPeersUpdate ../artifacts/channel/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP

