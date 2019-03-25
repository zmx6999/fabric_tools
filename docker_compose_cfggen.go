package main

import (
	"os"
	"strings"
	"errors"
	"encoding/json"
	"fmt"
	"github.com/sparrc/go-ping"
	"time"
	"math/rand"
)

type Ca struct {
	PeerOrgName string `json:"peer_org_name"`
	Ports []string `json:"ports"`
	AdminName string `json:"admin_name"`
	AdminPassword string `json:"admin_password"`
}

type Zookeeper struct {
	HostName string `json:"host_name"`
	Ports []string `json:"ports"`
	ZooMyID string `json:"zoo_my_id"`
	ZooServers string `json:"zoo_servers"`
	DataBackupDir string `json:"data_backup_dir"`
	DataLogBackupDir string `json:"data_log_backup_dir"`
}

type Kafka struct {
	HostName string `json:"host_name"`
	BrokerID string `json:"broker_id"`
	Zookeepers []string `json:"zookeepers"`
	Ports []string `json:"ports"`
	BackupDir string `json:"backup_dir"`
}

type Orderer struct {
	OrdererName string `json:"orderer_name"`
	KafkaBrokers []string `json:"kafka_brokers"`
	Ports []string `json:"ports"`
	BackupDir string `json:"backup_dir"`
}

type Peer struct {
	PeerName string `json:"peer_name"`
	OrgName string `json:"org_name"`
	Ports []string `json:"ports"`
	Couchdb Couchdb `json:"couchdb"`
	BackupDir string `json:"backup_dir"`
}

type Cli struct {
	CliName string `json:"cli_name"`
	CorePeerName string `json:"core_peer_name"`
	CorePeerOrg string `json:"core_peer_org"`
	Depends []string `json:"depends"`
}

type DockerComposeConfig struct {
	Domain string `json:"domain"`
	Cas []Ca `json:"cas"`
	Zookeepers []Zookeeper `json:"zookeepers"`
	Kafkas []Kafka `json:"kafkas"`
	Orderers []Orderer `json:"orderers"`
	Peers []Peer `json:"peers"`
	Hosts []string `json:"hosts"`
	Clis []Cli `json:"clis"`
}

type Couchdb struct {
	CouchdbName string `json:"couchdb_name"`
	Ports []string `json:"ports"`
	BackupDir string `json:"backup_dir"`
}

func generateNonceStr(length int) string {
	rand.Seed(time.Now().UnixNano())
	x:="abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	nx:=len(x)
	xbyte:=[]byte(x)
	str:=""
	for i:=0; i<length; i++ {
		str+=string(xbyte[rand.Intn(nx)])
	}
	return str
}

func loadDockerComposeConfig(configPath string,config *DockerComposeConfig) error {
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

func generateDockerComposeConfig(dstPath string,config DockerComposeConfig) error {
	str:=`
version: '2'

services:`

	caStr,err:=caConfigStr(config)
	if err!=nil {
		return err
	}
	str+=caStr

	ordererStr:=ordererConfigStr(config)
	str+=ordererStr

	peerStr:=peerConfigStr(config)
	str+=peerStr

	cliStr:=cliConfigStr(config)
	str+=cliStr

	file,err:=os.Create(dstPath)
	if err!=nil {
		return err
	}
	defer file.Close()

	file.Write([]byte(str))
	return nil
}

func getCaServerTLSKeyFile(peerOrgName string,domain string) (string,error) {
	path:="crypto-config/peerOrganizations/"+strings.ToLower(peerOrgName)+"."+domain+"/ca"
	info,err:=os.Stat(path)
	if err!=nil {
		return "", err
	}

	if info.IsDir() {
		dir,err:=os.Open(path)
		if err!=nil {
			return "", err
		}
		defer dir.Close()

		files,err:=dir.Readdir(0)
		if err!=nil {
			return "", err
		}

		for _,file:=range files{
			fileName:=file.Name()
			if strings.HasSuffix(fileName,"_sk") {
				return fileName,nil
			}
		}
	}
	return "",errors.New("CA SERVER TLS KEYFILE NOT FOUND")
}

func caConfigStr(config DockerComposeConfig) (string,error) {
	str:=``
	for _,ca:=range config.Cas{
		key,err:=getCaServerTLSKeyFile(ca.PeerOrgName,config.Domain)
		if err!=nil {
			return "", err
		}
		_str:=`
  ca_`+ca.PeerOrgName+`:
    image: hyperledger/fabric-ca
    environment:
      - FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server
      - FABRIC_CA_SERVER_CA_NAME=ca_`+ca.PeerOrgName+`
      - FABRIC_CA_SERVER_TLS_ENABLED=true
      - FABRIC_CA_SERVER_TLS_CERTFILE=/etc/hyperledger/fabric-ca-server-config/ca.`+strings.ToLower(ca.PeerOrgName)+`.`+config.Domain+`-cert.pem
      - FABRIC_CA_SERVER_TLS_KEYFILE=/etc/hyperledger/fabric-ca-server-config/`+key+`
    ports:`
		for _,port:=range ca.Ports{
			_str+=`
      - "`+port+`:7054"`
		}
		_str+=`
    command: sh -c 'fabric-ca-server start -b `+ca.AdminName+`:`+ca.AdminPassword+` -d'
    volumes:
      - ./crypto-config/peerOrganizations/`+strings.ToLower(ca.PeerOrgName)+`.`+config.Domain+`/ca/:/etc/hyperledger/fabric-ca-server-config
    container_name: ca_`+ca.PeerOrgName+`
`
		str+=_str
	}
	return str,nil
}

func ordererConfigStr(config DockerComposeConfig) string {
	str:=``
	for _,orderer:=range config.Orderers{
		kafkaBrokers:=[]string{}
		for _,broker:=range orderer.KafkaBrokers{
			kafkaBrokers=append(kafkaBrokers,broker)
		}
		_str:=`
  `+orderer.OrdererName+`.`+config.Domain+`:
    container_name: `+orderer.OrdererName+`.`+config.Domain+`
    image: hyperledger/fabric-orderer
    environment:
      - ORDERER_GENERAL_LOGLEVEL=INFO
      - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
      - ORDERER_GENERAL_GENESISMETHOD=file
      - ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block
      - ORDERER_GENERAL_LOCALMSPID=OrdererMSP
      - ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp
      - ORDERER_GENERAL_TLS_ENABLED=true
      - ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key
      - ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt
      - ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]
      - ORDERER_KAFKA_RETRY_LONGINTERVAL=10s
      - ORDERER_KAFKA_RETRY_LONGTOTAL=100s
      - ORDERER_KAFKA_RETRY_SHORTINTERVAL=1s
      - ORDERER_KAFKA_RETRY_SHORTTOTAL=30s
      - ORDERER_KAFKA_VERBOSE=true
      - ORDERER_KAFKA_BROKERS=[`+strings.Join(kafkaBrokers,",")+`]
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: orderer
    volumes:
      - ./channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block
      - ./crypto-config/ordererOrganizations/`+config.Domain+`/orderers/`+orderer.OrdererName+`.`+config.Domain+`/msp:/var/hyperledger/orderer/msp
      - ./crypto-config/ordererOrganizations/`+config.Domain+`/orderers/`+orderer.OrdererName+`.`+config.Domain+`/tls/:/var/hyperledger/orderer/tls
      - `+orderer.BackupDir+`:/var/hyperledger/production
    ports:`
		for _,port:=range orderer.Ports{
			_str+=`
      - "`+port+`:7050"`
		}
		_str+=`
    extra_hosts:`
		for _,host:=range config.Hosts{
			_str+=`
      - "`+host+`"`
		}
		_str+=`
`
		str+=_str
	}
	return str
}

func peerConfigStr(config DockerComposeConfig) string {
	str:=``
	innerPorts:=[]string{"7051","7052","7053"}
	for _,peer:=range config.Peers{
		_str:=`
  `+peer.PeerName+`.`+strings.ToLower(peer.OrgName)+`.`+config.Domain+`:
    container_name: `+peer.PeerName+`.`+strings.ToLower(peer.OrgName)+`.`+config.Domain+`
    image: hyperledger/fabric-peer
    environment:
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=${COMPOSE_PROJECT_NAME}_default
      - CORE_LOGGING_LEVEL=INFO
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_GOSSIP_USELEADERELECTION=true
      - CORE_PEER_GOSSIP_ORGLEADER=false
      - CORE_PEER_PROFILE_ENABLED=true
      - CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt
      - CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key
      - CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt
      - CORE_PEER_ID=`+peer.PeerName+`.`+strings.ToLower(peer.OrgName)+`.`+config.Domain+`
      - CORE_PEER_ADDRESS=`+peer.PeerName+`.`+strings.ToLower(peer.OrgName)+`.`+config.Domain+`:7051
      - CORE_PEER_CHAINCODEADDRESS=`+peer.PeerName+`.`+strings.ToLower(peer.OrgName)+`.`+config.Domain+`:7052
      - CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:7052
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT=`+peer.PeerName+`.`+strings.ToLower(peer.OrgName)+`.`+config.Domain+`:7051
      - CORE_PEER_GOSSIP_BOOTSTRAP=peer0.`+strings.ToLower(peer.OrgName)+`.`+config.Domain+`:7051
      - CORE_PEER_LOCALMSPID=`+peer.OrgName+`MSP`
		var couchdbUsername string
		var couchdbPassword string
		if peer.Couchdb.CouchdbName!="" {
			couchdbUsername=generateNonceStr(16)
			couchdbPassword=generateNonceStr(16)
			_str+=`
      - CORE_LEDGER_STATE_STATEDATABASE=CouchDB
      - CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=`+peer.Couchdb.CouchdbName+`.`+strings.ToLower(peer.OrgName)+`.`+config.Domain+`:5984
      - CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=`+couchdbUsername+`
      - CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=`+couchdbPassword
		}
		_str+=`
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
    command: peer node start
    volumes:
      - /var/run/:/host/var/run/
      - ./crypto-config/peerOrganizations/`+strings.ToLower(peer.OrgName)+`.`+config.Domain+`/peers/`+peer.PeerName+`.`+strings.ToLower(peer.OrgName)+`.`+config.Domain+`/msp:/etc/hyperledger/fabric/msp
      - ./crypto-config/peerOrganizations/`+strings.ToLower(peer.OrgName)+`.`+config.Domain+`/peers/`+peer.PeerName+`.`+strings.ToLower(peer.OrgName)+`.`+config.Domain+`/tls:/etc/hyperledger/fabric/tls
      - `+peer.BackupDir+`:/var/hyperledger/production
    ports:`
		for k,port:=range peer.Ports{
			_str+=`
      - "`+port+`:`+innerPorts[k]+`"`
		}
		_str+=`
    extra_hosts:`
		for _,host:=range config.Hosts{
			_str+=`
      - "`+host+`"`
		}
		_str+=`
`
		if peer.Couchdb.CouchdbName!="" {
			_str+=`
  `+peer.Couchdb.CouchdbName+`.`+strings.ToLower(peer.OrgName)+`.`+config.Domain+`:
    container_name: `+peer.Couchdb.CouchdbName+`.`+strings.ToLower(peer.OrgName)+`.`+config.Domain+`
    image: hyperledger/fabric-couchdb
    environment:
      - COUCHDB_USER=`+couchdbUsername+`
      - COUCHDB_PASSWORD=`+couchdbPassword+`
    ports:`
		for _,port:=range peer.Couchdb.Ports{
			_str+=`
      - "127.0.0.1:`+port+`:5984"`
		}
		_str+=`
    volumes:
      - `+peer.Couchdb.BackupDir+`:/opt/couchdb/data
`
		}
		str+=_str
	}
	return str
}

func cliConfigStr(config DockerComposeConfig) string {
	str:=``
	for _,cli:=range config.Clis{
		_str:=`
  `+cli.CliName+`:
    container_name: `+cli.CliName+`
    image: hyperledger/fabric-tools:$IMAGE_TAG
    tty: true
    stdin_open: true
    environment:
      - GOPATH=/opt/gopath
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_LOGGING_LEVEL=INFO
      - CORE_PEER_ID=`+cli.CliName+`
      - CORE_PEER_ADDRESS=`+cli.CorePeerName+`.`+strings.ToLower(cli.CorePeerOrg)+`.`+config.Domain+`:7051
      - CORE_PEER_LOCALMSPID=`+cli.CorePeerOrg+`MSP
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/`+strings.ToLower(cli.CorePeerOrg)+`.`+config.Domain+`/peers/`+cli.CorePeerName+`.`+strings.ToLower(cli.CorePeerOrg)+`.`+config.Domain+`/tls/server.crt
      - CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/`+strings.ToLower(cli.CorePeerOrg)+`.`+config.Domain+`/peers/`+cli.CorePeerName+`.`+strings.ToLower(cli.CorePeerOrg)+`.`+config.Domain+`/tls/server.key
      - CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/`+strings.ToLower(cli.CorePeerOrg)+`.`+config.Domain+`/peers/`+cli.CorePeerName+`.`+strings.ToLower(cli.CorePeerOrg)+`.`+config.Domain+`/tls/ca.crt
      - CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/`+strings.ToLower(cli.CorePeerOrg)+`.`+config.Domain+`/users/Admin@`+strings.ToLower(cli.CorePeerOrg)+`.`+config.Domain+`/msp
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
    command: /bin/bash
    volumes:
      - /var/run/:/host/var/run/
      - ./../chaincode/:/opt/gopath/src/github.com/chaincode
      - ./crypto-config:/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/
      - ./channel-artifacts:/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts
    depends_on:`
		for _,depend:=range cli.Depends{
			_str+=`
      - "`+depend+`"`
		}
		_str+=`
    extra_hosts:`
		for _,host:=range config.Hosts{
			_str+=`
      - "`+host+`"`
		}
		_str+=`
`
		str+=_str
	}
	return str
}

func generateZookeeper(dstPath string,config DockerComposeConfig) error {
	str:=`
version: '2'

services:`
	// innerPorts:=[]string{"2181","2888","3888"}
	for _,zookeeper:=range config.Zookeepers{
		_str:=`
  `+zookeeper.HostName+`:
    container_name: `+zookeeper.HostName+`
    hostname: `+zookeeper.HostName+`
    image: hyperledger/fabric-zookeeper
    restart: always
    ports:`
		/*
		for k,port:=range zookeeper.Ports{
			_str+=`
      - "`+port+`:`+innerPorts[k]+`"`
		}
		 */
		for _,port:=range zookeeper.Ports{
			_str+=`
      - "`+port+`:`+port+`"`
		}
		/*
		_str+=`
    expose:`
		for _,port:=range zookeeper.Ports{
			_str+=`
      - "`+port+`"`
		}
		 */
		_str+=`
    environment:`
		if len(zookeeper.Ports)>0 {
			_str+=`
      - ZOO_PORT=`+zookeeper.Ports[0]
		}
		_str+=`
      - ZOO_MY_ID=`+zookeeper.ZooMyID+`
      - ZOO_SERVERS=`+zookeeper.ZooServers+`
    extra_hosts:`
		for _,host:=range config.Hosts{
			_str+=`
      - "`+host+`"`
		}
		_str+=`
    volumes:
      - `+zookeeper.DataBackupDir+`:/data
      - `+zookeeper.DataLogBackupDir+`:/datalog
`
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

func generateKafka(dstPath string,config DockerComposeConfig) error {
	str:=`
version: '2'

services:`

	for _,kafka:=range config.Kafkas{
		_str:=`
  `+kafka.HostName+`:
    container_name: `+kafka.HostName+`
    hostname: `+kafka.HostName+`
    image: hyperledger/fabric-kafka
    restart: always
    environment:`
		if len(kafka.Ports) >0 {
			_str+=`
      - KAFKA_PORT=`+kafka.Ports[0]
		}
		_str+=`
      - KAFKA_MESSAGE_MAX_BYTES=103809024 # 99 * 1024 * 1024 B
      - KAFKA_REPLICA_FETCH_MAX_BYTES=103809024 # 99 * 1024 * 1024 B
      - KAFKA_UNCLEAN_LEADER_ELECTION_ENABLE=false
      - KAFKA_BROKER_ID=`+kafka.BrokerID+`
      - KAFKA_MIN_INSYNC_REPLICAS=2
      - KAFKA_DEFAULT_REPLICATION_FACTOR=3
      - KAFKA_ZOOKEEPER_CONNECT=`+strings.Join(kafka.Zookeepers,",")+`
    ports:`
		for _,port:=range kafka.Ports{
			_str+=`
      - "`+port+`:`+port+`"`
		}
		_str+=`
    extra_hosts:`
		for _,host:=range config.Hosts{
			_str+=`
      - "`+host+`"`
		}
		_str+=`
    volumes:
      - `+kafka.BackupDir+`:/tmp/kafka-logs
`
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

func updateHosts(dstPath string,config DockerComposeConfig) error {
	ch:=make(chan map[string]interface{})
	hosts:=config.Hosts
	for _,host:=range hosts{
		go func(host string) {
			r:=strings.Split(host,":")
			url:=r[0]
			ip:=r[1]
			fmt.Println("ping "+url)
			pinger,err:=ping.NewPinger(url)
			if err!=nil {
				ch<- map[string]interface{}{
					"url":url,
					"ip":ip,
					"total":3,
					"rev":0,
				}
				return
			}

			pinger.Timeout=time.Second*20
			pinger.Count=3
			pinger.Run()
			total:=pinger.PacketsSent
			rev:=pinger.PacketsRecv
			ch<- map[string]interface{}{
				"url":url,
				"ip":ip,
				"total":total,
				"rev":rev,
			}
		}(host)
	}
	str:=""
	for i:=0; i<len(hosts); i++ {
		m:=<-ch
		total:=float64(m["total"].(int))
		rev:=float64(m["rev"].(int))
		if rev/total<1.0/3 {
			str+="\n"+m["ip"].(string)+" "+m["url"].(string)
		}
	}

	file,err:=os.OpenFile(dstPath,os.O_APPEND|os.O_WRONLY,0644)
	if err!=nil {
		return err
	}
	defer file.Close()

	_,err=file.Write([]byte(str))
	if err!=nil {
		return err
	}
	return nil
}

func main()  {
	config:=DockerComposeConfig{}
	err:=loadDockerComposeConfig("docker_compose_cfggen.json",&config)
	if err!=nil {
		panic(err)
	}
	fmt.Println("generating docker-compose.yaml")
	err=generateDockerComposeConfig("docker-compose.yaml",config)
	if err!=nil {
		panic(err)
	}
	fmt.Println("generating zookeeper.yaml")
	err=generateZookeeper("zookeeper.yaml",config)
	if err!=nil {
		panic(err)
	}
	fmt.Println("generating kafka.yaml")
	err=generateKafka("kafka.yaml",config)
	if err!=nil {
		panic(err)
	}
	err=updateHosts("/etc/hosts",config)
	if err!=nil {
		panic(err)
	}
}
