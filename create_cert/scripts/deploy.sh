
. ./scripts/envVar.sh



peer channel create -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/mychannel.tx --outputBlock ./channel-artifacts/mychannel.block --tls --cafile $ORDERER_CA

peer channel join -b ./channel-artifacts/mychannel.block

setGlobals 2

peer channel join -b ./channel-artifacts/mychannel.block

setGlobals 3

peer channel join -b ./channel-artifacts/mychannel.block



CC_NAME=basic

peer lifecycle chaincode package ${CC_PATH}/${CC_NAME}.tar.gz --path ${CC_PATH}/$CC_NAME --label ${CC_NAME}_1
sleep 5
setGlobals 1
peer lifecycle chaincode install ${CC_PATH}/$CC_NAME.tar.gz 
sleep 5
setGlobals 2
peer lifecycle chaincode install ${CC_PATH}/$CC_NAME.tar.gz 
sleep 5
setGlobals 3
peer lifecycle chaincode install ${CC_PATH}/$CC_NAME.tar.gz 
sleep 5



setGlobals 1
peer lifecycle chaincode queryinstalled>&log.txt

echo "Arrpove for org for all orgs"
PACKAGE_ID=$(sed -n "/${CC_NAME}_${CC_VERSION}/{s/^Package ID: //; s/, Label:.*$//; p;}" log.txt)
peer lifecycle chaincode approveformyorg -o orderer.example.com:7050 --tls --cafile $ORDERER_CA --channelID mychannel --name basic --version 1 --package-id $PACKAGE_ID --sequence 1 --init-required --signature-policy "OR ('Org1MSP.peer','Org2MSP.peer','Org3MSP.peer')" 

sleep 5
setGlobals 2
peer lifecycle chaincode approveformyorg -o orderer.example.com:7050 --tls --cafile $ORDERER_CA --channelID mychannel --name basic --version 1 --package-id $PACKAGE_ID --sequence 1 --init-required --signature-policy "OR ('Org1MSP.peer','Org2MSP.peer','Org3MSP.peer')" 

sleep 5
setGlobals 3
peer lifecycle chaincode approveformyorg -o orderer.example.com:7050 --tls --cafile $ORDERER_CA --channelID mychannel --name basic --version 1 --package-id $PACKAGE_ID --sequence 1 --init-required --signature-policy "OR ('Org1MSP.peer','Org2MSP.peer','Org3MSP.peer')" 


setGlobals 1
echo "check commit readiness"
peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name basic --version 1 --sequence 1 --init-required --signature-policy "OR ('Org1MSP.peer','Org2MSP.peer','Org3MSP.peer')" 

parsePeerConnectionParameters 1 2 3

echo "chaincode commit "
sleep 5
peer lifecycle chaincode commit -o orderer.example.com:7050 --tls --cafile $ORDERER_CA --channelID mychannel --name basic $PEER_CONN_PARMS --version 1 --sequence 1 --init-required --signature-policy "OR ('Org1MSP.peer','Org2MSP.peer','Org3MSP.peer')" 

echo "query commited"
peer lifecycle chaincode querycommitted --channelID mychannel --name basic 


echo "Invoke chaincode "
fcn_call='{"function":"InitLedger","Args":[]}'
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n ${CC_NAME} $PEER_CONN_PARMS --isInit -c ${fcn_call} >&log.txt

sleep 5
fcn_call='{"function":"CreateAsset","Args":["1","2","2","2","2"]}'
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n ${CC_NAME} $PEER_CONN_PARMS -c ${fcn_call}
