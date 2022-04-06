# join node to network swarm
docker swarm join --token SWMTKN-1-0bqvfage8vqg7klpnhqzvk8o9hfq0yvubwfxw79m6ovfkqv53j-acn062wx3nczjgczjivqbf2r2 165.227.195.156:2377

# change label name of node 
docker node update --label-name name=add-org2 <id of node>

# run ca of node
cd /root/Sample_IOT_BC/org2/docker
docker-compose -f docker-compose-ca.yaml up -d

# check ca run successfull
docker logs <container id of ca> -f

# create certificate of org2
cd /root/Sample_IOT_BC/create_cert
source ./organizations/fabric-ca/registerEnroll.sh 
createOrg2

# run file add-org2-to-existingnetwork.md (modify every path in this file)
run step by step in this file

# copy file org2_update_in_envelope to container cli each node
docker cp /root/Sample_IOT_BC/org2/org2_update_in_envelope.pb <container id>:/opt/gopath/src/github.com/hyperledger/fabric/peer

# run script inside container cli
peer channel signconfigtx -f org2_update_in_envelope.pb

# cmd show "Endorser and orderer connections initialized" that it is success signconfigtx
