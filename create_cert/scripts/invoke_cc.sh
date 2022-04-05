
echo "Invoke chaincode "
fcn_call='{"function":"InitLedger","Args":[]}'
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n ${CC_NAME} $PEER_CONN_PARMS --isInit -c ${fcn_call}


# fcn_call='{"function":"CreateAsset","Args":["1","2","2","2","2"]}'
# peer chaincode invoke -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n ${CC_NAME} $PEER_CONN_PARMS -c ${fcn_call}
