infoln "Vendoring Go dependencies at ${CC_PATH}/${CC_NAME}"
pushd ${CC_PATH}/${CC_NAME}
GO111MODULE=on go mod vendor
popd
successln "Finished vendoring Go dependencies"


peer lifecycle chaincode package ${CC_PATH}/${CC_NAME}.tar.gz --path ${CC_PATH}/$CC_NAME --label ${CC_NAME}_1