
echo "peer lifecycle chaincode install <path>"

peer lifecycle chaincode install ${CC_PATH}/$CC_NAME.tar.gz 

peer lifecycle chaincode queryinstalled>&log.txt
