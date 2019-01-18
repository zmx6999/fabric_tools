package main

import (
	"os"
	"encoding/json"
	"strings"
	"strconv"
	"log"
	"fmt"
)

type Config struct {
	Domain string `json:"domain"`
	Orderers []string `json:"orderers"`
	Kafkas []string `json:"kafkas"`
	PeerOrgs []PeerOrg `json:"peer_orgs"`
	GenesisProfile string `json:"genesis_profile"`
	ChannelProfile string `json:"channel_profile"`
	Channel string `json:"channel"`
}

type PeerOrg struct {
	OrgName string `json:"org_name"`
	PeerCount int `json:"peer_count"`
	UserCount int `json:"user_count"`
}

func main()  {
	config:=Config{}
	err:=loadConfig("config.json",&config)
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println("generating crypto-config.yaml")
	err=generateCryptoConfig("crypto-config.yaml",config)
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println("generating configtx.yaml")
	err=generateConfigtx("configtx.yaml",config)
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println("generating _generate.sh")
	err=generateShell("_generate.sh",config)
	if err!=nil {
		log.Fatal(err)
	}
}

func loadConfig(configPath string,config *Config) error {
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

func generateCryptoConfig(dstPath string,config Config) error {
	str:= `
OrdererOrgs:
  - Name: Orderer
    Domain: `+config.Domain+`
    Specs:`
	for _,orderer:=range config.Orderers{
		_str:=`
      - Hostname: `+orderer
		str+=_str
	}

	str+=`

PeerOrgs:`

	for _,peerOrg:=range config.PeerOrgs{
		_str:=`
  - Name: `+peerOrg.OrgName+`
    Domain: `+strings.ToLower(peerOrg.OrgName)+`.`+config.Domain+`
    EnableNodeOUs: true
    Template:
      Count: `+strconv.Itoa(peerOrg.PeerCount)+`
    Users:
      Count: `+strconv.Itoa(peerOrg.UserCount)
		str+=_str
	}

	file,err:=os.Create(dstPath)
	if err!=nil {
		return err
	}
	defer file.Close()

	file.Write([]byte(str))
	return nil
}

func generateConfigtx(dstPath string,config Config) error {
	str:=`
Organizations:`

	str+=`
    - &OrdererOrg
        Name: OrdererOrg
        ID: OrdererMSP
        MSPDir: crypto-config/ordererOrganizations/`+config.Domain+`/msp`

	for _,peerOrg:=range config.PeerOrgs{
		_str:=`
    - &`+peerOrg.OrgName+`
        Name: `+peerOrg.OrgName+`MSP
        ID: `+peerOrg.OrgName+`MSP
        MSPDir: crypto-config/peerOrganizations/`+strings.ToLower(peerOrg.OrgName)+`.`+config.Domain+`/msp
        AnchorPeers:
            - Host: peer0.`+strings.ToLower(peerOrg.OrgName)+`.`+config.Domain+`
              Port: 7051`
		str+=_str
	}

	str+=`

Capabilities:
    Global: &ChannelCapabilities
        V1_1: true
    Orderer: &OrdererCapabilities
        V1_1: true
    Application: &ApplicationCapabilities
        V1_2: true
`

	str+=`
Application: &ApplicationDefaults
    Organizations:
`

	str+=`
Orderer: &OrdererDefaults
    OrdererType: kafka
    Addresses:`
	for _,orderer:=range config.Orderers{
		_str:=`
        - `+orderer+`.`+config.Domain+`:7050`
		str+=_str
	}
	str+=`
    BatchTimeout: 2s
    BatchSize:
        MaxMessageCount: 10
        AbsoluteMaxBytes: 99 MB
        PreferredMaxBytes: 512 KB
    Kafka:
        Brokers:`
	for _,kafka:=range config.Kafkas{
		_str:=`
            - `+kafka+`:9092`
		str+=_str
	}
	str+=`
    Organizations:
`

	str+=`
Profiles:
    `+config.GenesisProfile+`:
        Capabilities:
            <<: *ChannelCapabilities
        Orderer:
            <<: *OrdererDefaults
            Organizations:
                - *OrdererOrg
            Capabilities:
                <<: *OrdererCapabilities
        Consortiums:
            SampleConsortium:
                Organizations:`
	for _,peerOrg:=range config.PeerOrgs{
		_str:=`
                    - *`+peerOrg.OrgName
		str+=_str
	}
	str+=`
    `+config.ChannelProfile+`:
        Consortium: SampleConsortium
        Application:
            <<: *ApplicationDefaults
            Organizations:`
	for _,peerOrg:=range config.PeerOrgs{
		_str:=`
                - *`+peerOrg.OrgName
		str+=_str
	}
    str+=`
            Capabilities:
                <<: *ApplicationCapabilities
`

	file,err:=os.Create(dstPath)
	if err!=nil {
		return err
	}
	defer file.Close()

	file.Write([]byte(str))
	return nil
}

func generateShell(dstPath string,config Config) error {
	str:=`
#!/bin/bash

function updateAnchorPeer() {
    CHANNEL=$1
    GENESIS_PROFILE=$2
    CHANNEL_PROFILE=$3
    v=$4
    echo "configtxgen -profile ${CHANNEL_PROFILE} -outputAnchorPeersUpdate channel-artifacts/${v}anchors.tx -channelID ${CHANNEL} -asOrg ${v}"
    configtxgen -profile ${CHANNEL_PROFILE} -outputAnchorPeersUpdate channel-artifacts/${v}anchors.tx -channelID ${CHANNEL} -asOrg ${v}
    if [ $? -ne 0 ]; then echo "failed to generate ${v}anchors.tx"; exit 1; fi
}

function generate() {
    CHANNEL=$1
    GENESIS_PROFILE=$2
    CHANNEL_PROFILE=$3

    if [ -d crypto-config ]; then rm -rf crypto-config/*; fi
    if [ -d channel-artifacts ]; then rm -rf channel-artifacts/*; else mkdir channel-artifacts; fi
    
    echo "cryptogen generate --config=crypto-config.yaml"
    cryptogen generate --config=crypto-config.yaml
    if [ $? -ne 0 ]; then echo "failed to generate crypto"; exit 1; fi
    
    echo "configtxgen -profile ${GENESIS_PROFILE} -outputBlock channel-artifacts/genesis.block"
    configtxgen -profile ${GENESIS_PROFILE} -outputBlock channel-artifacts/genesis.block
    if [ $? -ne 0 ]; then echo "failed to generate genesis.block"; exit 1; fi
    
    echo "configtxgen -profile ${CHANNEL_PROFILE} -outputCreateChannelTx channel-artifacts/channel.tx -channelID ${CHANNEL}"
    configtxgen -profile ${CHANNEL_PROFILE} -outputCreateChannelTx channel-artifacts/channel.tx -channelID ${CHANNEL}
    if [ $? -ne 0 ]; then echo "failed to generate channel.tx"; exit 1; fi
    
    i=1
    for v in $@; do
        if [ $i -gt 3 ]; then
            updateAnchorPeer ${CHANNEL} ${GENESIS_PROFILE} ${CHANNEL_PROFILE} ${v}
        fi
        i=`+"`"+`expr ${i} + 1`+"`"+`
    done
}

export PATH=../../bin:$PATH

generate `+config.Channel+` `+config.GenesisProfile+` `+config.ChannelProfile
	for _,peerOrg:=range config.PeerOrgs{
		str+=` `+peerOrg.OrgName+`MSP`
	}

	file,err:=os.Create(dstPath)
	if err!=nil {
		return err
	}
	defer file.Close()

	file.Write([]byte(str))
	return nil
}
