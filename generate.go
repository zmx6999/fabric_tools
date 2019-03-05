package main

import (
	"os"
	"encoding/json"
	"strings"
	"strconv"
	"fmt"
		)

type Config struct {
	Domain string `json:"domain"`
	Orderers []string `json:"orderers"`
	Kafkas []string `json:"kafkas"`
	PeerOrgs []PeerOrg `json:"peer_orgs"`
	GenesisProfile string `json:"genesis_profile"`
	Channels []ConfigChannel `json:"channels"`
}

type PeerOrg struct {
	OrgName string `json:"org_name"`
	PeerCount int `json:"peer_count"`
	UserCount int `json:"user_count"`
	AnchorPeers []string `json:"anchor_peers"`
}

type ConfigChannel struct {
	ChannelName string `json:"channel_name"`
	Orgs []string `json:"orgs"`
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
        AnchorPeers:`
		for _,anchorPeer:=range peerOrg.AnchorPeers{
			_str+=`
            - Host: `+anchorPeer+`.`+strings.ToLower(peerOrg.OrgName)+`.`+config.Domain+`
              Port: 7051`
		}
		str+=_str
	}

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
        Orderer:
            <<: *OrdererDefaults
            Organizations:
                - *OrdererOrg
        Consortiums:`
	for _,channel:=range config.Channels{
		str+=`
            `+channel.ChannelName+`Consortium:
                Organizations:`
		for _,org:=range channel.Orgs{
			_str:=`
                    - *`+org
			str+=_str
		}
	}
	for _,channel:=range config.Channels{
		str+=`
    `+channel.ChannelName+`Channel:
        Consortium: `+channel.ChannelName+`Consortium
        Application:
            <<: *ApplicationDefaults
            Organizations:`
		for _,org:=range channel.Orgs{
			_str:=`
                - *`+org
			str+=_str
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

func generateShell(dstPath string,config Config) error {
	str:=`
#!/bin/bash

function generate() {
    if [ -d crypto-config ]; then rm -rf crypto-config/*; fi
    if [ -d channel-artifacts ]; then rm -rf channel-artifacts/*; else mkdir channel-artifacts; fi
    
    echo "cryptogen generate --config=crypto-config.yaml"
    cryptogen generate --config=crypto-config.yaml
    if [ $? -ne 0 ]; then echo "failed to generate crypto"; exit 1; fi
    
    echo "configtxgen -profile `+config.GenesisProfile+` -outputBlock channel-artifacts/genesis.block"
    configtxgen -profile `+config.GenesisProfile+` -outputBlock channel-artifacts/genesis.block
    if [ $? -ne 0 ]; then echo "failed to generate genesis.block"; exit 1; fi`
	for _,channel:=range config.Channels{
		str+=`

	echo "configtxgen -profile `+channel.ChannelName+`Channel -outputCreateChannelTx channel-artifacts/`+channel.ChannelName+`.tx -channelID `+channel.ChannelName+`"
    configtxgen -profile `+channel.ChannelName+`Channel -outputCreateChannelTx channel-artifacts/`+channel.ChannelName+`.tx -channelID `+channel.ChannelName+`
    if [ $? -ne 0 ]; then echo "failed to generate `+channel.ChannelName+`.tx"; exit 1; fi`
	}
	str+=`
}

export PATH=../../bin:$PATH

generate`

	file,err:=os.Create(dstPath)
	if err!=nil {
		return err
	}
	defer file.Close()

	file.Write([]byte(str))
	return nil
}

func main()  {
	config:=Config{}
	err:=loadConfig("generate.json",&config)
	if err!=nil {
		panic(err)
	}
	fmt.Println("generating crypto-config.yaml")
	err=generateCryptoConfig("crypto-config.yaml",config)
	if err!=nil {
		panic(err)
	}
	fmt.Println("generating configtx.yaml")
	err=generateConfigtx("configtx.yaml",config)
	if err!=nil {
		panic(err)
	}
	fmt.Println("generating _generate.sh")
	err=generateShell("_generate.sh",config)
	if err!=nil {
		panic(err)
	}
}
