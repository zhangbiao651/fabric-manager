cryptogen generate --config=./config/crypto/crypto.yaml --output=orgs
cd ./config/configtx
configtxgen -profile OrgsOrdererGenesis -outputBlock ../../data/genesis.block -channelID fabric-channel
configtxgen -profile OrgsChannel -outputCreateChannelTx ../../data/mychannel.tx -channelID mychannel
configtxgen -profile OrgsChannel -outputAnchorPeersUpdate ../../data/mychannelOrg1anchors.tx -channelID mychannel -asOrg Org1MSP

configtxgen -profile OrgsChannel -outputAnchorPeersUpdate ../../data/mychannelOrg2anchors.tx -channelID mychannel -asOrg Org2MSP

configtxgen -profile OrgsChannel -outputCreateChannelTx ../../data/testchannel.tx -channelID testchannel
configtxgen -profile OrgsChannel -outputAnchorPeersUpdate ../../data/testchannelOrg1anchors.tx -channelID testchannel -asOrg Org1MSP

configtxgen -profile OrgsChannel -outputAnchorPeersUpdate ../../data/testchannelOrg2anchors.tx -channelID testchannel -asOrg Org2MSP

