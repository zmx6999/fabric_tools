package main

import (
	"os"
	"encoding/json"
	"strings"
	"fmt"
)

type ChaincodeOrderer struct {
	OrdererName string `json:"orderer_name"`
	Port string `json:"port"`
}

type ChaincodePeer struct {
	PeerName string `json:"peer_name"`
	Port string `json:"port"`
}

type ChaincodeOrg struct {
	OrgName string `json:"org_name"`
	Peers []ChaincodePeer `json:"peers"`
}

type ChaincodeConfig struct {
	Domain string `json:"domain"`
	ChaincodeName string `json:"chaincode_name"`
	ChaincodeVersion string `json:"chaincode_version"`
	Channels []ChaincodeChannel `json:"channels"`
	Orderer ChaincodeOrderer `json:"orderer"`
	Endorse string `json:"endorse"`
	CliName string `json:"cli_name"`
}

type ChaincodeChannel struct {
	ChannelName string `json:"channel_name"`
	Orgs []ChaincodeOrg `json:"orgs"`
}

func loadChaincodeConfig(configPath string,config *ChaincodeConfig) error {
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

func generateChaincodeSh(dstPath string,config ChaincodeConfig,mode string) error {
	str:=`
function installChaincode() {
    CHAINCODE_NAME=$1
    CHAINCODE_VERSION=$2
    CLI=$3

    _v=(`+"`"+`echo "$4" | sed 's/:/\n/g'`+"`"+`)
	host=${_v[0]}
	port=${_v[1]}
	org=${_v[2]}
	hv=(`+"`"+`echo "${host}" | sed 's/\./\n/'`+"`"+`)
	peer=${hv[0]}
	org_domain=${hv[1]}

	echo "docker exec -e "CORE_PEER_ADDRESS=${peer}.${org_domain}:${port}" -e "CORE_PEER_LOCALMSPID=${org}MSP" -e "CORE_PEER_TLS_ENABLED=true" -e "CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/server.crt" -e "CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/server.key" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/users/Admin@${org_domain}/msp" ${CLI} peer chaincode install -n ${CHAINCODE_NAME} -v ${CHAINCODE_VERSION} -p github.com/chaincode/${CHAINCODE_NAME}"
	docker exec -e "CORE_PEER_ADDRESS=${peer}.${org_domain}:${port}" -e "CORE_PEER_LOCALMSPID=${org}MSP" -e "CORE_PEER_TLS_ENABLED=true" -e "CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/server.crt" -e "CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/server.key" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/users/Admin@${org_domain}/msp" ${CLI} peer chaincode install -n ${CHAINCODE_NAME} -v ${CHAINCODE_VERSION} -p github.com/chaincode/${CHAINCODE_NAME}
	if [ $? -ne 0 ]; then echo "${peer}.${org_domain}:${port} failed to install chaincode"; exit 1; fi
}

function instantiateChaincode() {
	CHAINCODE_NAME=$1
	CHAINCODE_VERSION=$2
	CHANNEL=$3
	ORDERER=$4
	ENDORSE=$5
	CLI=$6

	r=(`+"`"+`echo "${ORDERER}" | sed 's/:/\n/'`+"`"+`)
	ORDERER_HOST=${r[0]}

	_v=(`+"`"+`echo "$7" | sed 's/:/\n/g'`+"`"+`)
	host=${_v[0]}
	port=${_v[1]}
	org=${_v[2]}
	hv=(`+"`"+`echo "${host}" | sed 's/\./\n/'`+"`"+`)
	peer=${hv[0]}
	org_domain=${hv[1]}

	echo docker exec -e "CORE_PEER_ADDRESS=${peer}.${org_domain}:${port}" -e "CORE_PEER_LOCALMSPID=${org}MSP" -e "CORE_PEER_TLS_ENABLED=true" -e "CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/server.crt" -e "CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/server.key" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/users/Admin@${org_domain}/msp" ${CLI} peer chaincode `+mode+` -n ${CHAINCODE_NAME} -v ${CHAINCODE_VERSION} -C ${CHANNEL} -c '{"args":["Init"]}' -o ${ORDERER} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/`+config.Domain+`/orderers/${ORDERER_HOST}/msp/tlscacerts/tlsca.`+config.Domain+`-cert.pem -P "${ENDORSE}"
	docker exec -e "CORE_PEER_ADDRESS=${peer}.${org_domain}:${port}" -e "CORE_PEER_LOCALMSPID=${org}MSP" -e "CORE_PEER_TLS_ENABLED=true" -e "CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/server.crt" -e "CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/server.key" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/peers/${peer}.${org_domain}/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/${org_domain}/users/Admin@${org_domain}/msp" ${CLI} peer chaincode `+mode+` -n ${CHAINCODE_NAME} -v ${CHAINCODE_VERSION} -C ${CHANNEL} -c '{"args":["Init"]}' -o ${ORDERER} --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/`+config.Domain+`/orderers/${ORDERER_HOST}/msp/tlscacerts/tlsca.`+config.Domain+`-cert.pem -P "${ENDORSE}"
	if [ $? -ne 0 ]; then echo "failed to instantiate chaincode"; exit 1; fi
}

function installAndInstantiateChaincode() {
	CHAINCODE_NAME=$1
	CHAINCODE_VERSION=$2
	CHANNEL=$3
	ORDERER=$4
	ENDORSE=$5
	CLI=$6

	echo "cd ../chaincode/${CHAINCODE_NAME}"
	cd ../chaincode/${CHAINCODE_NAME}
	echo "go build ${CHAINCODE_NAME}.go"
	go build ${CHAINCODE_NAME}.go
	if [ $? -ne 0 ]; then exit 1; fi

	i=1
	for v in "$@";do
		if [ ${i} -gt 6 ]; then 
			installChaincode ${CHAINCODE_NAME} ${CHAINCODE_VERSION} ${CLI} ${v}
			if [ ${i} -eq 7 ]; then
				instantiateChaincode ${CHAINCODE_NAME} ${CHAINCODE_VERSION} ${CHANNEL} ${ORDERER} "${ENDORSE}" ${CLI} ${v}
			fi
		fi
		i=`+"`"+`expr ${i} + 1`+"`"+`
	done
}`
	for _,channel:=range config.Channels{
		str+=`

installAndInstantiateChaincode `+config.ChaincodeName+` `+config.ChaincodeVersion+` `+channel.ChannelName+` `+config.Orderer.OrdererName+`.`+config.Domain+`:`+config.Orderer.Port+` "`+config.Endorse+`" `+config.CliName
		for _,org:=range channel.Orgs{
			for _,peer:=range org.Peers{
				str+=` `+peer.PeerName+`.`+strings.ToLower(org.OrgName)+`.`+config.Domain+`:`+peer.Port+`:`+org.OrgName
			}
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
	if len(os.Args)<3 {
		panic(`Invalid arguments.Usage: chaincode.sh <option> CONFIG_PATH
options:
    -i instantiate chaincode
    -u upgrade chaincode
`)
	}

	configPath:=os.Args[2]
	config:=ChaincodeConfig{}
	err:=loadChaincodeConfig(configPath,&config)
	if err!=nil {
		panic(err)
	}

	mode:=os.Args[1]
	switch mode {
	case "-i":
		mode="instantiate"
		break
	case "-u":
		mode="upgrade"
		break
	default:
		panic(`Invalid arguments.Usage: chaincode.sh <option> CONFIG_PATH
options:
    -i instantiate chaincode
    -u upgrade chaincode
`)
		break
	}
	fmt.Println("generating _chaincode.sh")
	err=generateChaincodeSh("_chaincode.sh",config,mode)
	if err!=nil {
		panic(err)
	}
}
