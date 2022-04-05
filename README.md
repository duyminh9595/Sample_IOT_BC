# Sample_IOT_BC

# Crate Swarm network
docker swarm init --advertise-addr 165.227.195.156
<!-- copy the command and replace by your ip -->
docker swarm join --token SWMTKN-1-0bqvfage8vqg7klpnhqzvk8o9hfq0yvubwfxw79m6ovfkqv53j-acn062wx3nczjgczjivqbf2r2 165.227.195.156:2377 --advertise-addr <Your IP>
docker network create --driver=overlay --attachable test
# install portainer for docker
https://docs.portainer.io/v/ce-2.9/start/install/server/swarm/linux

# update host name of Node
docker node list
docker node update --label-add name=<Name> <ID of node>


# run file docker-compose-ca in org1 and orderer for creating ca
docker-compose -f docker-compose-ca.yaml up -d



# create certificate
cd create_cert/
source ./organizations/fabric-ca/registerEnroll.sh 
createOrg1

# create genesis block
./scripts/createGenesis.sh 

# create channel file
./scripts/createChannelTx.sh

# deploy peer and orderer
cd /root/Sample_IOT_BC/org1/docker
docker-compose -f docker-compose-test-net.yaml up -d

cd /root/Sample_IOT_BC/orderer/docker
docker-compose -f docker-compose-test-net.yaml up -d

cd /root/Sample_IOT_BC/org1/docker
docker-compose -f docker-compose-couch.yaml up -d

# create cli
cd /root/Sample_IOT_BC/org1/docker
docker-compose -f docker-compose-cli.yaml up -d

# inside cli 
# create channel
docker exec -it <container id of cli> bash
export CHANNEL_NAME=mychannel
echo $CHANNEL_NAME
./scripts/create_app_channel.sh 

# join channel
peer channel join -b ./channel-artifacts/mychannel.block 

# update anchor peer
./scripts/updateAnchorPeer.sh mychannel Org1MSP

# chaincode packing and install
export CC_NAME=basic
./scripts/package_cc.sh 
./scripts/install_cc.sh 

# approve chaincode
export CHANNEL_NAME=mychannel
export CC_NAME=basic
export PACKAGE_ID=basic_1:9820659c595e662a849033ca23b4424e87a126e8f40b5f81ace59820b81fe8e7
peer lifecycle chaincode approveformyorg -o orderer.example.com:7050 --tls --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name $CC_NAME --version 1 --package-id $PACKAGE_ID --sequence 1 --init-required --signature-policy "OR ('Org1MSP.peer')"

peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name $CC_NAME --version 1 --sequence 1 --init-required --signature-policy "OR ('Org1MSP.peer')" 

peer lifecycle chaincode commit -o orderer.example.com:7050 --tls --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name $CC_NAME $PEER_CONN_PARMS --version 1 --sequence 1 --init-required --signature-policy "OR ('Org1MSP.peer')" 

peer lifecycle chaincode querycommitted --channelID mychannel --name basic 

fcn_call='{"function":"InitLedger","Args":[]}'
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n ${CC_NAME} $PEER_CONN_PARMS --isInit -c ${fcn_call}

fcn_call='{"function":"CreateAsset","Args":["1","2","2","2","2"]}'
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n ${CC_NAME} $PEER_CONN_PARMS -c ${fcn_call}
