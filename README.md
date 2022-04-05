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
