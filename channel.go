package main

import (
	"os"
	"encoding/json"
	"strings"
	"log"
	"fmt"
)

type ChannelOrderer struct {
	OrdererName string `json:"orderer_name"`
	Port string `json:"port"`
}

type ChannelPeer struct {
	PeerName string `json:"peer_name"`
	Port string `json:"port"`
}

type ChannelOrg struct {
	OrgName string `json:"org_name"`
	Peers []ChannelPeer `json:"peers"`
	AnchorPeers []ChannelPeer `json:"anchor_peers"`
}

type ChannelConfig struct {
	Domain string `json:"domain"`
	ChannelName string `json:"channel_name"`
	Orderer ChannelOrderer `json:"orderer"`
	CliName string `json:"cli_name"`
	ChannelOrgs []ChannelOrg `json:"channel_orgs"`
}

func loadChannelConfig(configPath string,config *ChannelConfig) error {
	file,err:=os.Open(configPath)
	if err!=nil {
		return err
	}
	defer file.Close()

	info,err:=os.Stat(configPath)
	if err!=nil {
		return err
	}

	m:=make([]byte,info.Size())
	_,err=file.Read(m)
	if err!=nil {
		return err
	}

	err=json.Unmarshal(m,config)
	if err!=nil {
		return err
	}

	return nil
}

func generateChannelSh(dstPath string,config ChannelConfig) error {
	str:=`
#!/bin/bash

function joinChannel() {
    CHANNEL=$1
    ORDERER=$2
    CLI=$3
    
    v=(`+"`"+`echo "$4" | sed 's/:/\n/g'`+"`"+`)
	host=${v[0]}
	port=${v[1]}
	org=${v[2]}
	hv=(`+"`"+`echo "${host}" | sed 's/\./\n/'`+"`"+`)
	peer=${hv[0]}
	org_domain=${hv[1]}
	echo "docker exec -e "CORE_PEER_ADDRESS=${peer}.${org_domain}:${port}" -e "CORE_PEER_LOCALMSPID=${org}MSP" -e "CORE_PEER_TLS_ENABLED=true" -e "CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/server.crt" -e "CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/server.key" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/users/Admin@${org_domain}/msp" ${CLI} peer channel join -b ${CHANNEL}.block"
	docker exec -e "CORE_PEER_ADDRESS=${peer}.${org_domain}:${port}" -e "CORE_PEER_LOCALMSPID=${org}MSP" -e "CORE_PEER_TLS_ENABLED=true" -e "CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/server.crt" -e "CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/server.key" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/users/Admin@${org_domain}/msp" ${CLI} peer channel join -b ${CHANNEL}.block
	if [ $? -ne 0 ]; then echo "${peer}.${org_domain}:${port} failed to join channel"; exit 1; fi
}

function createAndJoinChannel() {
	CHANNEL=$1
	ORDERER=$2
	CLI=$3

	r=(`+"`"+`echo "${ORDERER}" | sed 's/:/\n/'`+"`"+`)
	ORDERER_HOST=${r[0]}
	echo "docker exec ${CLI} peer channel create -c ${CHANNEL} -f channel-artifacts/channel.tx -o ${ORDERER} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/`+config.Domain+`/orderers/${ORDERER_HOST}/msp/tlscacerts/tlsca.`+config.Domain+`-cert.pem -t 150s"
	docker exec ${CLI} peer channel create -c ${CHANNEL} -f channel-artifacts/channel.tx -o ${ORDERER} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/`+config.Domain+`/orderers/${ORDERER_HOST}/msp/tlscacerts/tlsca.`+config.Domain+`-cert.pem -t 150s
	if [ $? -ne 0 ]; then echo "failed to create channel"; exit 1; fi

	i=1
	for v in $@; do
		if [ ${i} -gt 3 ]; then joinChannel ${CHANNEL} ${ORDERER} ${CLI} ${v}; fi
		i=`+"`"+`expr ${i} + 1`+"`"+`
	done
}

createAndJoinChannel `+config.ChannelName+` `+config.Orderer.OrdererName+`.`+config.Domain+`:`+config.Orderer.Port+` `+config.CliName
	for _,org:=range config.ChannelOrgs{
		for _,peer:=range org.Peers{
			str+=` `+peer.PeerName+`.`+strings.ToLower(org.OrgName)+`.`+config.Domain+`:`+peer.Port+`:`+org.OrgName
		}
	}

	file,err:=os.Create(dstPath)
	if err!=nil {
		return err
	}
	defer file.Close()

	file.Write([]byte(str))
	return nil
}

func generateAnchorSh(dstPath string,config ChannelConfig) error {
	str:=`
#!/bin/bash

function _updateAnchorPeers() {
    CHANNEL=$1
    ORDERER=$2
    CLI=$3

    v=(`+"`"+`echo "$4" | sed 's/:/\n/g'`+"`"+`)
	host=${v[0]}
	port=${v[1]}
	org=${v[2]}
	hv=(`+"`"+`echo "${host}" | sed 's/\./\n/'`+"`"+`)
	peer=${hv[0]}
	org_domain=${hv[1]}

	r=(`+"`"+`echo "${ORDERER}" | sed 's/:/\n/'`+"`"+`)
	ORDERER_HOST=${r[0]}
	echo "docker exec -e "CORE_PEER_ADDRESS=${peer}.${org_domain}:${port}" -e "CORE_PEER_LOCALMSPID=${org}MSP" -e "CORE_PEER_TLS_ENABLED=true" -e "CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/server.crt" -e "CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/server.key" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/users/Admin@${org_domain}/msp" ${CLI} peer channel update -c ${CHANNEL} -f channel-artifacts/channel.tx -o ${ORDERER} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/`+config.Domain+`/orderers/${ORDERER_HOST}/msp/tlscacerts/tlsca.`+config.Domain+`-cert.pem"
	docker exec -e "CORE_PEER_ADDRESS=${peer}.${org_domain}:${port}" -e "CORE_PEER_LOCALMSPID=${org}MSP" -e "CORE_PEER_TLS_ENABLED=true" -e "CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/server.crt" -e "CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/server.key" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/users/Admin@${org_domain}/msp" ${CLI} peer channel update -c ${CHANNEL} -f channel-artifacts/${org}MSPanchors.tx -o ${ORDERER} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/`+config.Domain+`/orderers/${ORDERER_HOST}/msp/tlscacerts/tlsca.`+config.Domain+`-cert.pem
	if [ $? -ne 0 ]; then echo "failed to update anchor peers for ${org}MSP"; exit 1; fi
}

function updateAnchorPeers() {
	CHANNEL=$1
	ORDERER=$2
	CLI=$3

	i=1
	for v in $@; do
		if [ ${i} -gt 3 ]; then _updateAnchorPeers ${CHANNEL} ${ORDERER} ${CLI} ${v}; fi
		i=`+"`"+`expr ${i} + 1`+"`"+`
	done
}

updateAnchorPeers `+config.ChannelName+` `+config.Orderer.OrdererName+`.`+config.Domain+`:`+config.Orderer.Port+` `+config.CliName
	for _,org:=range config.ChannelOrgs{
		for _,peer:=range org.AnchorPeers{
			str+=` `+peer.PeerName+`.`+strings.ToLower(org.OrgName)+`.`+config.Domain+`:`+peer.Port+`:`+org.OrgName
		}
	}

	file,err:=os.Create(dstPath)
	if err!=nil {
		return err
	}
	defer file.Close()

	file.Write([]byte(str))
	return nil
}

func main()  {
	config:=ChannelConfig{}
	err:=loadChannelConfig("channel.json",&config)
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println("generating _channel.sh")
	err=generateChannelSh("_channel.sh",config)
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println("generating _anchor.sh")
	err=generateAnchorSh("_anchor.sh",config)
	if err!=nil {
		log.Fatal(err)
	}
}
