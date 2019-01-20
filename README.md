# Hyperledger fabric+kafka+docker-compose Multi-host Deployment

## Assume that we have 4 hosts A,B,C,D and each host's IP and roles are as following:
```
A 139.180.137.89 zookeeper0 kafka0 orderer0.trace.com

B 139.180.209.250 zookeeper1 kafka1 orderer1.trace.com

C 139.180.137.0 zookeeper2 kafka2 peer0.orgdairy.trace.com peer1.orgdairy.trace.com peer0.orgprocess.trace.com peer1.orgprocess.trace.com ca_OrgDairy ca_OrgProcess cli

D 198.13.46.60 kafka3 peer0.orgsell.trace.com peer0.orgsell.trace.com ca_OrgSell
```

## 1.Prepare

(a)Download fabric_tools from https://github.com/zmx6999/fabric_tools

(b)Copy init.sh in downloaded fabric_tools to each host of A,B,C,D and execute
```
./init.sh 190116
```
It will install docker,docker-compose and go,and download fabric.git,fabric-samples and docker images related to fabric,and create a  directory /root/fabric/scripts/fabric-samples/190116/network on each host.

## 2.Genarate crypto-config files,genesis block and files related to anchor peers

(a)Copy generate.json,generate.go,generate.sh in downloaded fabric_tools to the directory /root/fabric/scripts/fabric-samples/190116/network of host C and enter the directory.

(b)Edit generate.json as following:
```
{
  "domain": "trace.com",
  "orderers": [
    "orderer0",
    "orderer1"
  ],
  "kafkas": [
    "kafka0",
    "kafka1",
    "kafka2",
    "kafka3"
  ],
  "peer_orgs": [
    {
      "org_name": "OrgDairy",
      "peer_count": 2,
      "user_count": 1,
      "anchor_peers": [
        "peer0"
      ]
    },
    {
      "org_name": "OrgProcess",
      "peer_count": 2,
      "user_count": 1,
      "anchor_peers": [
        "peer0"
      ]
    },
    {
      "org_name": "OrgSell",
      "peer_count": 2,
      "user_count": 1,
      "anchor_peers": [
        "peer0"
      ]
    }
  ],
  "genesis_profile": "ThreeOrgsOrdererGenesis",
  "channel_profile": "ThreeOrgsChannel",
  "channel": "trace"
}
```

(c)Execute
```
chmod +x generate.sh
./generate.sh
```
It will generate crypto-config files,genesis.block,channel.tx and DairyOrgMSPanchors.tx,ProcessOrgMSPanchors.tx,SellOrgMSPanchors.tx which are related to anchor peers.

## 3.Generate docker-compose configuration files and start services including zookeeper,kafka,orderer,peer,ca,cli

(a)Copy docker_compose_cfggen.json,docker_compose_cfggen.go,docker_compose_cfggen.sh in fabric_tools to the directory /root/fabric/scripts/fabric-samples/190116/network of each host of A,B,C,D and enter the directory.

(b)Edit docker_compose_cfggen.json on host A as following:
```
{
  "domain": "trace.com",
  "zookeepers": [
    {
      "host_name": "zookeeper0",
      // outer ports
      "ports": [
        "2181", // The first item corresponds to the inner port 2181
        "2888", // The second item corresponds to the inner port 2888
        "3888" // The third item corresponds to the inner port 3888
      ],
      "zoo_my_id": "1",
      "zoo_servers": "server.1=zookeeper0:2888:3888 server.2=zookeeper1:2888:3888 server.3=zookeeper2:2888:3888"
    }
  ],
  "kafkas": [
    {
      "host_name": "kafka0",
      "broker_id": "0",
      "zookeepers": [
        "zookeeper0:2181",
        "zookeeper1:2181",
        "zookeeper2:2181"
      ]
    }
  ],
  "orderers": [
    {
      "orderer_name": "orderer0",
      "kafka_brokers": [
        "kafka0",
        "kafka1",
        "kafka2",
        "kafka3"
      ],
      // outer ports
      "ports": [
        "7050" // corresponds to the inner port 7050
      ]
    }
  ],
  "hosts": [
    "zookeeper0:139.180.137.89",
    "kafka0:139.180.137.89",
    "orderer0.trace.com:139.180.137.89",
    "zookeeper1:139.180.209.250",
    "kafka1:139.180.209.250",
    "orderer1.trace.com:139.180.209.250",
    "zookeeper2:139.180.137.0",
    "kafka2:139.180.137.0",
    "peer0.orgdairy.trace.com:139.180.137.0",
    "peer1.orgdairy.trace.com:139.180.137.0",
    "peer0.orgprocess.trace.com:139.180.137.0",
    "peer1.orgprocess.trace.com:139.180.137.0",
    "kafka3:198.13.46.60",
    "peer0.orgsell.trace.com:198.13.46.60",
    "peer1.orgsell.trace.com:198.13.46.60"
  ]
}
```
Edit docker_compose_cfggen.json on host B as following:
```
{
  "domain": "trace.com",
  "zookeepers": [
    {
      "host_name": "zookeeper1",
      // outer ports
      "ports": [
        "2181", // The first item corresponds to the inner port 2181
        "2888", // The second item corresponds to the inner port 2888
        "3888" // The third item corresponds to the inner port 3888
      ],
      "zoo_my_id": "2",
      "zoo_servers": "server.1=zookeeper0:2888:3888 server.2=zookeeper1:2888:3888 server.3=zookeeper2:2888:3888"
    }
  ],
  "kafkas": [
    {
      "host_name": "kafka1",
      "broker_id": "1",
      "zookeepers": [
        "zookeeper0:2181",
        "zookeeper1:2181",
        "zookeeper2:2181"
      ]
    }
  ],
  "orderers": [
    {
      "orderer_name": "orderer1",
      "kafka_brokers": [
        "kafka0",
        "kafka1",
        "kafka2",
        "kafka3"
      ],
      // outer ports
      "ports": [
        "8050" // corresponds to the inner port 7050
      ]
    }
  ],
  "hosts": [
    "zookeeper0:139.180.137.89",
    "kafka0:139.180.137.89",
    "orderer0.trace.com:139.180.137.89",
    "zookeeper1:139.180.209.250",
    "kafka1:139.180.209.250",
    "orderer1.trace.com:139.180.209.250",
    "zookeeper2:139.180.137.0",
    "kafka2:139.180.137.0",
    "peer0.orgdairy.trace.com:139.180.137.0",
    "peer1.orgdairy.trace.com:139.180.137.0",
    "peer0.orgprocess.trace.com:139.180.137.0",
    "peer1.orgprocess.trace.com:139.180.137.0",
    "kafka3:198.13.46.60",
    "peer0.orgsell.trace.com:198.13.46.60",
    "peer1.orgsell.trace.com:198.13.46.60"
  ]
}
```
Edit docker_compose_cfggen.json on host C as following:
```
{
  "domain": "trace.com",
  "cas": [
    {
      "peer_org_name": "OrgDairy",
      // outer ports
      "ports": [
        "7054" // corresponds to the inner port 7054
      ],
      "admin_name": "admin",
      "admin_password": "adminpw"
    },
    {
      "peer_org_name": "OrgProcess",
      // outer ports
      "ports": [
        "8054" // corresponds to the inner port 7054
      ],
      "admin_name": "admin",
      "admin_password": "adminpw"
    }
  ],
  "zookeepers": [
    {
      "host_name": "zookeeper2",
      // outer ports
      "ports": [
        "2181", // The first item corresponds to the inner port 2181
        "2888", // The second item corresponds to the inner port 2888
        "3888" // The third item corresponds to the inner port 3888
      ],
      "zoo_my_id": "3",
      "zoo_servers": "server.1=zookeeper0:2888:3888 server.2=zookeeper1:2888:3888 server.3=zookeeper2:2888:3888"
    }
  ],
  "kafkas": [
    {
      "host_name": "kafka2",
      "broker_id": "2",
      "zookeepers": [
        "zookeeper0:2181",
        "zookeeper1:2181",
        "zookeeper2:2181"
      ]
    }
  ],
  "peers": [
    {
      "peer_name": "peer0",
      "org_name": "OrgDairy",
      // outer ports
      "ports": [
        "7051", // The first item corresponds to the inner port 7051
        "7052", // The second item corresponds to the inner port 7052
        "7053" // The third item corresponds to the inner port 7053
      ]
    },
    {
      "peer_name": "peer1",
      "org_name": "OrgDairy",
      // outer ports
      "ports": [
        "8051", // The first item corresponds to the inner port 7051
        "8052", // The second item corresponds to the inner port 7052
        "8053" // The third item corresponds to the inner port 7053
      ]
    },
    {
      "peer_name": "peer0",
      "org_name": "OrgProcess",
      // outer ports
      "ports": [
        "9051", // The first item corresponds to the inner port 7051
        "9052", // The second item corresponds to the inner port 7052
        "9053" // The third item corresponds to the inner port 7053
      ]
    },
    {
      "peer_name": "peer1",
      "org_name": "OrgProcess",
      // outer ports
      "ports": [
        "10051", // The first item corresponds to the inner port 7051
        "10052", // The second item corresponds to the inner port 7052
        "10053" // The third item corresponds to the inner port 7053
      ]
    }
  ],
  "clis": [
    {
      "cli_name": "cli",
      "core_peer_name": "peer0",
      "core_peer_org": "OrgDairy",
      "depends": [
        "peer0.orgdairy.trace.com",
        "peer1.orgdairy.trace.com",
        "peer0.orgprocess.trace.com",
        "peer1.orgprocess.trace.com"
      ]
    }
  ],
  "hosts": [
    "zookeeper0:139.180.137.89",
    "kafka0:139.180.137.89",
    "orderer0.trace.com:139.180.137.89",
    "zookeeper1:139.180.209.250",
    "kafka1:139.180.209.250",
    "orderer1.trace.com:139.180.209.250",
    "zookeeper2:139.180.137.0",
    "kafka2:139.180.137.0",
    "peer0.orgdairy.trace.com:139.180.137.0",
    "peer1.orgdairy.trace.com:139.180.137.0",
    "peer0.orgprocess.trace.com:139.180.137.0",
    "peer1.orgprocess.trace.com:139.180.137.0",
    "kafka3:198.13.46.60",
    "peer0.orgsell.trace.com:198.13.46.60",
    "peer1.orgsell.trace.com:198.13.46.60"
  ]
}
```
Edit docker_compose_cfggen.json on host D as following:
```
{
  "domain": "trace.com",
  "cas": [
    {
      "peer_org_name": "OrgSell",
      // outer ports
      "ports": [
        "9054" // corresponds to the inner port 7054
      ],
      "admin_name": "admin",
      "admin_password": "adminpw"
    }
  ],
  "kafkas": [
    {
      "host_name": "kafka3",
      "broker_id": "3",
      "zookeepers": [
        "zookeeper0:2181",
        "zookeeper1:2181",
        "zookeeper2:2181"
      ]
    }
  ],
  "peers": [
    {
      "peer_name": "peer0",
      "org_name": "OrgSell",
      // outer ports
      "ports": [
        "11051", // The first item corresponds to the inner port 7051
        "11052", // The second item corresponds to the inner port 7052
        "11053" // The third item corresponds to the inner port 7053
      ]
    },
    {
      "peer_name": "peer1",
      "org_name": "OrgSell",
      // outer ports
      "ports": [
        "12051", // The first item corresponds to the inner port 7051
        "12052", // The second item corresponds to the inner port 7052
        "12053" // The third item corresponds to the inner port 7053
      ]
    }
  ],
  "hosts": [
    "zookeeper0:139.180.137.89",
    "kafka0:139.180.137.89",
    "orderer0.trace.com:139.180.137.89",
    "zookeeper1:139.180.209.250",
    "kafka1:139.180.209.250",
    "orderer1.trace.com:139.180.209.250",
    "zookeeper2:139.180.137.0",
    "kafka2:139.180.137.0",
    "peer0.orgdairy.trace.com:139.180.137.0",
    "peer1.orgdairy.trace.com:139.180.137.0",
    "peer0.orgprocess.trace.com:139.180.137.0",
    "peer1.orgprocess.trace.com:139.180.137.0",
    "kafka3:198.13.46.60",
    "peer0.orgsell.trace.com:198.13.46.60",
    "peer1.orgsell.trace.com:198.13.46.60"
  ]
}
```

(c)Copy crypto-config files and genesis.block from host C to other hosts.

Execute on host A
```
cd /root/fabric/scripts/fabric-samples/190116/network
mkdir channel-artifacts
mkdir -p crypto-config/ordererOrganizations/trace.com/orderers
```
Copy crypto-config files and genesis.block from host C to host A,executing in the directory /root/fabric/scripts/fabric-samples/190116/network of host C
```
scp channel-artifacts/genesis.block root@139.180.137.89:/root/fabric/scripts/fabric-samples/190116/network/channel-artifacts
scp -r crypto-config/ordererOrganizations/trace.com/orderers/orderer0.trace.com root@139.180.137.89:/root/fabric/scripts/fabric-samples/190116/network/crypto-config/ordererOrganizations/trace.com/orderers
```
Execute on host B
```
cd /root/fabric/scripts/fabric-samples/190116/network
mkdir channel-artifacts
mkdir -p crypto-config/ordererOrganizations/trace.com/orderers
```
Copy crypto-config files and genesis.block from host C to host B,executing in the directory /root/fabric/scripts/fabric-samples/190116/network of host C
```
scp channel-artifacts/genesis.block root@139.180.209.250:/root/fabric/scripts/fabric-samples/190116/network/channel-artifacts
scp -r crypto-config/ordererOrganizations/trace.com/orderers/orderer1.trace.com root@139.180.209.250:/root/fabric/scripts/fabric-samples/190116/network/crypto-config/ordererOrganizations/trace.com/orderers
```
Execute on host D
```
cd /root/fabric/scripts/fabric-samples/190116/network
mkdir -p crypto-config/peerOrganizations/orgsell.trace.com
```
Copy crypto-config files from host C to host D,executing in the directory /root/fabric/scripts/fabric-samples/190116/network of host C
```
scp -r crypto-config/peerOrganizations/orgsell.trace.com/peers root@198.13.46.60:/root/fabric/scripts/fabric-samples/190116/network/crypto-config/peerOrganizations/orgsell.trace.com
scp -r crypto-config/peerOrganizations/orgsell.trace.com/ca root@198.13.46.60:/root/fabric/scripts/fabric-samples/190116/network/crypto-config/peerOrganizations/orgsell.trace.com
```

(d)Execute on each host of A,B,C,D
```
cd /root/fabric/scripts/fabric-samples/190116/network
chmod +x docker_compose_cfggen.sh
./docker_compose_cfggen.sh trace # trace is COMPOSE_PROJECT_NAME
```
It will generate docker-compose configuration files including zookeeper.yaml,kafka.yaml,docker-compose.yaml.

(e)Start zookeeper on each host of A,B,C,executing in the directory /root/fabric/scripts/fabric-samples/190116/network of each host
```
docker-compose -f zookeeper.yaml up -d
```

(f)Start kafka on each host of A,B,C,D,executing in the directory /root/fabric/scripts/fabric-samples/190116/network of each host
```
docker-compose -f kafka.yaml up -d
```

(g)Execute in the directory /root/fabric/scripts/fabric-samples/190116/network of each host of A,B,C,D
```
docker-compose -f docker-compose.yaml up -d
```
It will start service including orderer0.trace.com on host A,orderer1.trace.com on host B,peer0.orgdairy.trace.com,peer1.orgdairy.trace.com,peer0.orgprocess.trace.com,peer1.orgprocess.trace.com,ca_OrgDairy,ca_OrgProcess,cli on host C,peer0.orgsell.trace.com,peer1.orgsell.trace.com,ca_OrgSell on host D.

## 4.Create a channel,and make each peer node join the channel,and update anchor peers

(a)Copy channel.go,channel.sh in fabric_tools to the directory /root/fabric/scripts/fabric-samples/190116/network of host C and enter the directory.

(b)Create and edit trace.json as following:
```
{
  "domain": "trace.com",
  "channel_name": "trace",
  "orderer": {
    "orderer_name": "orderer0",
    "port": "7050"
  },
  "cli_name": "cli",
  "channel_orgs": [
    {
      "org_name": "OrgDairy",
      "peers": [
        {
          "peer_name": "peer0",
          "port": "7051"
        },
        {
          "peer_name": "peer1",
          "port": "8051"
        }
      ],
      "anchor_peers": [
        {
          "peer_name": "peer0",
          "port": "7051"
        }
      ]
    },
    {
      "org_name": "OrgProcess",
      "peers": [
        {
          "peer_name": "peer0",
          "port": "9051"
        },
        {
          "peer_name": "peer1",
          "port": "10051"
        }
      ],
      "anchor_peers": [
        {
          "peer_name": "peer0",
          "port": "9051"
        }
      ]
    },
    {
      "org_name": "OrgSell",
      "peers": [
        {
          "peer_name": "peer0",
          "port": "11051"
        },
        {
          "peer_name": "peer1",
          "port": "12051"
        }
      ],
      "anchor_peers": [
        {
          "peer_name": "peer0",
          "port": "11051"
        }
      ]
    }
  ]
}
```

(c)Execute
```
chmod +x channel.sh
./channel.sh trace.json
```
It will create a channel named trace,make each peer node join the channel and update anchor peers.

## 5.Edit chaincodes

(a)Create a directory named chaincode in the directory /root/fabric/scripts/fabric-samples/190116 of host C and enter the created chaincode directory.Then created 3 directories named dairy(stores the chaincode of dairies),process(stores the chaincode of process factories),sell(stores the chaincode of sell organizations) respectively.

(b)Create and edit dairy.go in the created directory dairy as following:
```
package main
 
import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"errors"
	"fmt"
	"time"
	"encoding/json"
)
 
type DairyChaincode struct {
 
}
 
func (dc *DairyChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}
 
func (dc *DairyChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	_,args:=stub.GetFunctionAndParameters()
	err:=checkArgs(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}
	fn:=args[0]
	if fn=="set" {
		return dc.set(stub,args[1:])
	} else if fn=="get" {
		return dc.get(stub,args[1:])
	} else if fn=="history" {
		return dc.history(stub,args[1:])
	}
	return shim.Error("METHOD NOT FOUND")
}
 
func (dc *DairyChaincode) set(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgs(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}
	err=stub.PutState(args[0],[]byte(args[1]))
	if err!=nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}
 
func (dc *DairyChaincode) get(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgs(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}
	data,err:=stub.GetState(args[0])
	if err!=nil {
		return shim.Error(err.Error())
	}
	return shim.Success(data)
}
 
func (dc *DairyChaincode) history(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgs(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}
	iter,err:=stub.GetHistoryForKey(args[0])
	if err!=nil {
		return shim.Error(err.Error())
	}
	defer iter.Close()
	var list []string
	for iter.HasNext() {
		item,err:=iter.Next()
		if err!=nil {
			return shim.Error(err.Error())
		}
		v:=fmt.Sprintf("%s|%s",time.Unix(item.Timestamp.Seconds,0).Format("2006-01-02 15:04:05"),item.Value)
		list=append(list,v)
	}
	data,err:=json.Marshal(list)
	if err!=nil {
		return shim.Error(err.Error())
	}
	return shim.Success(data)
}
 
func main()  {
	shim.Start(new(DairyChaincode))
}
 
func checkArgs(args []string,n int) error {
	if len(args)<n {
		return errors.New(fmt.Sprintf("%d argument(s) required",n))
	}
	return nil
}
```
Create and edit process.go in the created directory process as following:
```
package main
 
import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"errors"
	"fmt"
	"time"
	"encoding/json"
)
 
type ProcessChaincode struct {
 
}
 
func (pc *ProcessChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}
 
func (pc *ProcessChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	_,args:=stub.GetFunctionAndParameters()
	err:=checkArgs(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}
	fn:=args[0]
	if fn=="set" {
		return pc.set(stub,args[1:])
	} else if fn=="get" {
		return pc.get(stub,args[1:])
	} else if fn=="history" {
		return pc.history(stub,args[1:])
	}
	return shim.Error("METHOD NOT FOUND")
}
 
func (pc *ProcessChaincode) set(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgs(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}
	err=stub.PutState(args[0],[]byte(args[1]))
	if err!=nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}
 
func (pc *ProcessChaincode) get(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgs(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}
	data,err:=stub.GetState(args[0])
	if err!=nil {
		return shim.Error(err.Error())
	}
	return shim.Success(data)
}
 
func (pc *ProcessChaincode) history(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgs(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}
	iter,err:=stub.GetHistoryForKey(args[0])
	if err!=nil {
		return shim.Error(err.Error())
	}
	defer iter.Close()
	var list []string
	for iter.HasNext() {
		item,err:=iter.Next()
		if err!=nil {
			return shim.Error(err.Error())
		}
		v:=fmt.Sprintf("%s|%s",time.Unix(item.Timestamp.Seconds,0).Format("2006-01-02 15:04:05"),item.Value)
		list=append(list,v)
	}
	data,err:=json.Marshal(list)
	if err!=nil {
		return shim.Error(err.Error())
	}
	return shim.Success(data)
}
 
func main()  {
	shim.Start(new(ProcessChaincode))
}
 
func checkArgs(args []string,n int) error {
	if len(args)<n {
		return errors.New(fmt.Sprintf("%d argument(s) required",n))
	}
	return nil
}
```
Create and edit sell.go in the created directory sell as following:
```
package main
 
import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"errors"
	"fmt"
	"time"
	"encoding/json"
	"strings"
)
 
type SellChaincode struct {
 
}
 
func (sc *SellChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}
 
func (sc *SellChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	_,args:=stub.GetFunctionAndParameters()
	err:=checkArgs(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}
	fn:=args[0]
	if fn=="set" {
		return sc.set(stub,args[1:])
	} else if fn=="get" {
		return sc.get(stub,args[1:])
	} else if fn=="history" {
		return sc.history(stub,args[1:])
	}
	return shim.Error("METHOD NOT FOUND")
}
 
func (sc *SellChaincode) set(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgs(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}
	err=stub.PutState(args[0],[]byte(args[1]))
	if err!=nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}
 
func (sc *SellChaincode) get(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgs(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}
	data,err:=stub.GetState(args[0])
	if err!=nil {
		return shim.Error(err.Error())
	}
	return shim.Success(data)
}
 
func (sc *SellChaincode) history(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgs(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}
	iter,err:=stub.GetHistoryForKey(args[0])
	if err!=nil {
		return shim.Error(err.Error())
	}
	defer iter.Close()
	var list []map[string]interface{}
	for iter.HasNext() {
		item,err:=iter.Next()
		if err!=nil {
			return shim.Error(err.Error())
		}
		v:=fmt.Sprintf("%s|%s",time.Unix(item.Timestamp.Seconds,0).Format("2006-01-02 15:04:05"),item.Value)
		m:=map[string]interface{}{
			"info": v,
		}
 
		process:=strings.Split(v,"|")[1]
		response:=stub.InvokeChaincode("process",[][]byte{[]byte("invoke"),[]byte("history"),[]byte(process)},"trace")
		if response.Status!=shim.OK {
			return shim.Error(response.String())
		}
		var _processList []string
		err=json.Unmarshal(response.Payload,&_processList)
		if err!=nil {
			return shim.Error(err.Error())
		}
 
		var processList []map[string]interface{}
		for _,_v:=range _processList{
			dairy:=strings.Split(_v,"|")[1]
			_response:=stub.InvokeChaincode("dairy",[][]byte{[]byte("invoke"),[]byte("history"),[]byte(dairy)},"trace")
			if _response.Status!=shim.OK {
				return shim.Error(_response.String())
			}
			var dairyList []string
			err=json.Unmarshal(_response.Payload,&dairyList)
			_m:=map[string]interface{}{
				"info": _v,
				"trace": dairyList,
			}
			processList=append(processList,_m)
		}
		m["trace"]=processList
		list=append(list,m)
	}
	data,err:=json.Marshal(list)
	if err!=nil {
		return shim.Error(err.Error())
	}
	return shim.Success(data)
}
 
func main()  {
	shim.Start(new(SellChaincode))
}
 
func checkArgs(args []string,n int) error {
	if len(args)<n {
		return errors.New(fmt.Sprintf("%d argument(s) required",n))
	}
	return nil
}
```

## 6.Install,instantiate and invoke chaincodes

(a)Copy chaincode.go,chaincode.sh in fabric_tools to the directory /root/fabric/scripts/fabric-samples/190116/network of host C and enter the directory.

(b)Create and edit dairy.json as following:
```
{
  "domain": "trace.com",
  "channel_name": "trace",
  "chaincode_name": "dairy",
  "chaincode_version": "1.0",
  "orderer": {
    "orderer_name": "orderer0",
    "port": "7050"
  },
  "endorse": "AND ('OrgDairyMSP.member')",
  "cli_name": "cli",
  "chaincode_orgs": [
    {
      "org_name": "OrgDairy",
      "peers": [
        {
          "peer_name": "peer0",
          "port": "7051"
        },
        {
          "peer_name": "peer1",
          "port": "8051"
        }
      ]
    },
    {
      "org_name": "OrgSell",
      "peers": [
        {
          "peer_name": "peer0",
          "port": "11051"
        },
        {
          "peer_name": "peer1",
          "port": "12051"
        }
      ]
    }
  ]
}
```

(c)Execute
```
chmod +x chaincode.sh
./chaincode.sh -i dairy.json
```
It will install and instantiate the chaincode dairy.

(d)Create and edit process.json as following:
```
{
  "domain": "trace.com",
  "channel_name": "trace",
  "chaincode_name": "process",
  "chaincode_version": "1.0",
  "orderer": {
    "orderer_name": "orderer0",
    "port": "7050"
  },
  "endorse": "AND ('OrgProcessMSP.member')",
  "cli_name": "cli",
  "chaincode_orgs": [
    {
      "org_name": "OrgProcess",
      "peers": [
        {
          "peer_name": "peer0",
          "port": "9051"
        },
        {
          "peer_name": "peer1",
          "port": "10051"
        }
      ]
    },
    {
      "org_name": "OrgSell",
      "peers": [
        {
          "peer_name": "peer0",
          "port": "11051"
        },
        {
          "peer_name": "peer1",
          "port": "12051"
        }
      ]
    }
  ]
}
```

(e)Execute
```
./chaincode.sh -i process.json
```
It will install and instantiate the chaincode process.

(f)Create and edit sell.json as following:
```
{
  "domain": "trace.com",
  "channel_name": "trace",
  "chaincode_name": "sell",
  "chaincode_version": "1.0",
  "orderer": {
    "orderer_name": "orderer0",
    "port": "7050"
  },
  "endorse": "AND ('OrgSellMSP.member')",
  "cli_name": "cli",
  "chaincode_orgs": [
    {
      "org_name": "OrgSell",
      "peers": [
        {
          "peer_name": "peer0",
          "port": "11051"
        },
        {
          "peer_name": "peer1",
          "port": "12051"
        }
      ]
    }
  ]
}
```

(g)Execute
```
./chaincode.sh -i sell.json
```
It will install and instantiate the chaincode sell.

(h)Invoke the chaincode dairy.

Create and edit dairy_test.sh as following:
```
#!/bin/bash
 
docker exec cli peer chaincode invoke -n dairy -C trace -c '{"args":["invoke","set","dairy101","info101"]}' --peerAddresses peer0.orgdairy.trace.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgdairy.trace.com/peers/peer0.orgdairy.trace.com/tls/ca.crt -o orderer1.trace.com:8050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/trace.com/orderers/orderer1.trace.com/msp/tlscacerts/tlsca.trace.com-cert.pem
sleep 5
 
docker exec -e CORE_PEER_ADDRESS=peer1.orgdairy.trace.com:8051 -e CORE_PEER_LOCALMSPID=OrgDairyMSP -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgdairy.trace.com/peers/peer1.orgdairy.trace.com/tls/server.crt -e CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgdairy.trace.com/peers/peer1.orgdairy.trace.com/tls/server.key -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgdairy.trace.com/peers/peer1.orgdairy.trace.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgdairy.trace.com/users/Admin@orgdairy.trace.com/msp cli peer chaincode query -n dairy -C trace -c '{"args":["invoke","get","dairy101"]}'
 
docker exec cli peer chaincode invoke -n dairy -C trace -c '{"args":["invoke","set","dairy102","info102"]}' --peerAddresses peer0.orgdairy.trace.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgdairy.trace.com/peers/peer0.orgdairy.trace.com/tls/ca.crt -o orderer1.trace.com:8050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/trace.com/orderers/orderer1.trace.com/msp/tlscacerts/tlsca.trace.com-cert.pem
sleep 5
 
docker exec -e CORE_PEER_ADDRESS=peer1.orgdairy.trace.com:8051 -e CORE_PEER_LOCALMSPID=OrgDairyMSP -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgdairy.trace.com/peers/peer1.orgdairy.trace.com/tls/server.crt -e CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgdairy.trace.com/peers/peer1.orgdairy.trace.com/tls/server.key -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgdairy.trace.com/peers/peer1.orgdairy.trace.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgdairy.trace.com/users/Admin@orgdairy.trace.com/msp cli peer chaincode query -n dairy -C trace -c '{"args":["invoke","get","dairy102"]}'
 
docker exec cli peer chaincode invoke -n dairy -C trace -c '{"args":["invoke","set","dairy101","info103"]}' --peerAddresses peer0.orgdairy.trace.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgdairy.trace.com/peers/peer0.orgdairy.trace.com/tls/ca.crt -o orderer1.trace.com:8050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/trace.com/orderers/orderer1.trace.com/msp/tlscacerts/tlsca.trace.com-cert.pem
sleep 5
 
docker exec -e CORE_PEER_ADDRESS=peer1.orgdairy.trace.com:8051 -e CORE_PEER_LOCALMSPID=OrgDairyMSP -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgdairy.trace.com/peers/peer1.orgdairy.trace.com/tls/server.crt -e CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgdairy.trace.com/peers/peer1.orgdairy.trace.com/tls/server.key -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgdairy.trace.com/peers/peer1.orgdairy.trace.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgdairy.trace.com/users/Admin@orgdairy.trace.com/msp cli peer chaincode query -n dairy -C trace -c '{"args":["invoke","history","dairy101"]}'
```
Then execute
```
chmod +x dairy_test.sh
./dairy_test.sh
```

(i)Invoke the chaincode process.

Create and edit process_test.sh as following:
```
#!/bin/bash
 
docker exec cli peer chaincode invoke -n process -C trace -c '{"args":["invoke","set","process101","dairy101"]}' --peerAddresses peer0.orgprocess.trace.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgprocess.trace.com/peers/peer0.orgprocess.trace.com/tls/ca.crt -o orderer1.trace.com:8050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/trace.com/orderers/orderer1.trace.com/msp/tlscacerts/tlsca.trace.com-cert.pem
sleep 5
 
docker exec -e CORE_PEER_ADDRESS=peer1.orgprocess.trace.com:10051 -e CORE_PEER_LOCALMSPID=OrgProcessMSP -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgprocess.trace.com/peers/peer1.orgprocess.trace.com/tls/server.crt -e CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgprocess.trace.com/peers/peer1.orgprocess.trace.com/tls/server.key -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgprocess.trace.com/peers/peer1.orgprocess.trace.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgprocess.trace.com/users/Admin@orgprocess.trace.com/msp cli peer chaincode query -n process -C trace -c '{"args":["invoke","get","process101"]}'
 
docker exec cli peer chaincode invoke -n process -C trace -c '{"args":["invoke","set","process101","dairy102"]}' --peerAddresses peer0.orgprocess.trace.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgprocess.trace.com/peers/peer0.orgprocess.trace.com/tls/ca.crt -o orderer1.trace.com:8050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/trace.com/orderers/orderer1.trace.com/msp/tlscacerts/tlsca.trace.com-cert.pem
sleep 5
 
docker exec -e CORE_PEER_ADDRESS=peer1.orgprocess.trace.com:10051 -e CORE_PEER_LOCALMSPID=OrgProcessMSP -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgprocess.trace.com/peers/peer1.orgprocess.trace.com/tls/server.crt -e CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgprocess.trace.com/peers/peer1.orgprocess.trace.com/tls/server.key -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgprocess.trace.com/peers/peer1.orgprocess.trace.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgprocess.trace.com/users/Admin@orgprocess.trace.com/msp cli peer chaincode query -n process -C trace -c '{"args":["invoke","history","process101"]}'
```
Then execute
```
chmod +x process_test.sh
./process_test.sh
```

(j)Invoke the chaincode sell.

Create and edit sell_test.sh as following:
```
#!/bin/bash
 
docker exec cli peer chaincode invoke -n sell -C trace -c '{"args":["invoke","set","sell101","process101"]}' --peerAddresses peer0.orgsell.trace.com:11051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgsell.trace.com/peers/peer0.orgsell.trace.com/tls/ca.crt -o orderer1.trace.com:8050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/trace.com/orderers/orderer1.trace.com/msp/tlscacerts/tlsca.trace.com-cert.pem
sleep 5
 
docker exec -e CORE_PEER_ADDRESS=peer1.orgsell.trace.com:12051 -e CORE_PEER_LOCALMSPID=OrgSellMSP -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgsell.trace.com/peers/peer1.orgsell.trace.com/tls/server.crt -e CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgsell.trace.com/peers/peer1.orgsell.trace.com/tls/server.key -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgsell.trace.com/peers/peer1.orgsell.trace.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgsell.trace.com/users/Admin@orgsell.trace.com/msp cli peer chaincode query -n sell -C trace -c '{"args":["invoke","get","sell101"]}'
 
docker exec -e CORE_PEER_ADDRESS=peer1.orgsell.trace.com:12051 -e CORE_PEER_LOCALMSPID=OrgSellMSP -e CORE_PEER_TLS_ENABLED=true -e CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgsell.trace.com/peers/peer1.orgsell.trace.com/tls/server.crt -e CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgsell.trace.com/peers/peer1.orgsell.trace.com/tls/server.key -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgsell.trace.com/peers/peer1.orgsell.trace.com/tls/ca.crt -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgsell.trace.com/users/Admin@orgsell.trace.com/msp cli peer chaincode query -n sell -C trace -c '{"args":["invoke","history","sell101"]}' --connTimeout 60s
```
Then execute
```
chmod +x sell_test.sh
./sell_test.sh
```
