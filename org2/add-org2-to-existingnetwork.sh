export CHANNEL_NAME=mychannel

export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=/root/Sample_IOT_BC/create_cert/organizations/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export PEER0_ORG1_CA=/root/Sample_IOT_BC/create_cert/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export PEER0_ORG2_CA=${PWD}/../../artifacts/channel/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt


setGlobalsForOrderer() {
    export CORE_PEER_LOCALMSPID="OrdererMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=/root/Sample_IOT_BC/create_cert/organizations/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem
    export CORE_PEER_MSPCONFIGPATH=/root/Sample_IOT_BC/create_cert/organizations/ordererOrganizations/example.com/users/Admin@example.com/msp

}

setGlobalsForPeer0Org1() {
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
    export CORE_PEER_MSPCONFIGPATH=/root/Sample_IOT_BC/create_cert/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
}


generateOrg3Definition() {
    export FABRIC_CFG_PATH=$PWD/
    configtxgen -printOrg Org2MSP >org2.json
}
# generateOrg3Definition

# this is using on existing org of network
fetchChannelConfig() {
    # setGlobalsForOrderer
    # setGlobalsForPeer0Org1

    # Fetch the config for the channel, writing it to config.json
    # run inside the cli of node
    # echo "Fetching the most recent configuration block for the channel"
    # set -x
    # peer channel fetch config config_block.pb -o orderer.example.com:7050 \
    #     -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
    # set +x
    
    # run this command to copy config_block.pb from the container to localhost
    # docker cp <container id>:/opt/gopath/src/github.com/hyperledger/fabric/peer/config_block.pb /root/Sample_IOT_BC/org2


    # echo "Decoding config block to JSON and isolating config to config.json"
    # set -x
    # configtxlator proto_decode --input config_block.pb --type common.Block | jq .data.data[0].payload.data.config >config.json
    # set +x

    # # Modify the configuration to append the new org
    set -x
    jq -s '.[0] * {"channel_group":{"groups":{"Application":{"groups": {"Org2MSP":.[1]}}}}}' config.json ./org2.json >modified_config.json
    set +x

}
# fetchChannelConfig

createConfigUpdate() {

    CHANNEL="mychannel"
    ORIGINAL="config.json"
    MODIFIED="modified_config.json"
    OUTPUT="org2_update_in_envelope.pb"

    set -x
    configtxlator proto_encode --input "${ORIGINAL}" --type common.Config >original_config.pb

    configtxlator proto_encode --input "${MODIFIED}" --type common.Config >modified_config.pb

    configtxlator compute_update --channel_id "${CHANNEL}" --original original_config.pb --updated modified_config.pb >config_update.pb

    configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate >config_update.json

    echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . >config_update_in_envelope.json

    configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope >"${OUTPUT}"
    set +x

}
# createConfigUpdate