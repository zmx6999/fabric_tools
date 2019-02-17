# Hyperledger fabric+kafka+docker-compose Multi-host Deployment

## Assume that we have 4 hosts A,B,C,D and each host's IP and roles are as following:
```
A 139.180.138.179 zookeeper0 kafka0 orderer0.house.com

B 45.77.251.25 zookeeper1 kafka1 orderer1.house.com

C 45.77.250.9 zookeeper2 kafka2 peer0.orgauth.house.com peer1.orgauth.house.com peer0.orgcert.house.com peer1.orgcert.house.com ca_OrgAuth ca_OrgCert cli

D 139.180.146.33 kafka3 peer0.orgcredit.house.com peer1.orgcredit.house.com ca_OrgCredit
```

## 1.Prepare

(a)Download fabric_tools from https://github.com/zmx6999/fabric_tools

(b)Copy init.sh in downloaded fabric_tools to each host of A,B,C,D and execute
```
./init.sh 190216
```
It will install docker,docker-compose and go,and download fabric.git,fabric-samples and docker images related to fabric,and create a directory /root/fabric/scripts/fabric-samples/190216/network on each host.

## 2.Genarate crypto-config files,genesis block and channels' configuration files

(a)Copy generate.json,generate.go,generate.sh in downloaded fabric_tools to the directory /root/fabric/scripts/fabric-samples/190216/network of host C and enter the directory.

(b)Edit generate.json as following:
```
{
  "domain": "house.com",
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
      "org_name": "OrgAuth",
      "peer_count": 2,
      "user_count": 1,
      "anchor_peers": [
        "peer0"
      ]
    },
    {
      "org_name": "OrgCert",
      "peer_count": 2,
      "user_count": 1,
      "anchor_peers": [
        "peer0"
      ]
    },
    {
      "org_name": "OrgCredit",
      "peer_count": 2,
      "user_count": 1,
      "anchor_peers": [
        "peer0"
      ]
    }
  ],
  "genesis_profile": "ThreeOrgsOrdererGenesis",
  "channels": [
    {
      "channel_name": "auth",
      "orgs": [
        "OrgAuth"
      ]
    },
    {
      "channel_name": "cert",
      "orgs": [
        "OrgCert"
      ]
    },
    {
      "channel_name": "credit",
      "orgs": [
        "OrgCredit"
      ]
    }
  ]
}
```

(c)Execute
```
chmod +x generate.sh
./generate.sh
```
It will generate crypto-config files,genesis.block,and auth.tx,cert.tx,credit.tx which are channel config files.

## 3.Generate docker-compose configuration files and start services including zookeeper,kafka,orderer,peer,ca,cli

(a)Copy docker_compose_cfggen.json,docker_compose_cfggen.go,docker_compose_cfggen.sh in fabric_tools to the directory /root/fabric/scripts/fabric-samples/190216/network of each host of A,B,C,D and enter the directory.

(b)Edit docker_compose_cfggen.json on host A as following:
```
{
  "domain": "house.com",
  "zookeepers": [
    {
      "host_name": "zookeeper0",
      "ports": [
        "2181",
        "2888",
        "3888"
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
      "ports": [
        "7050"
      ]
    }
  ],
  "hosts": [
    "zookeeper0:139.180.138.179",
    "kafka0:139.180.138.179",
    "orderer0.house.com:139.180.138.179",
    "zookeeper1:45.77.251.25",
    "kafka1:45.77.251.25",
    "orderer1.house.com:45.77.251.25",
    "zookeeper2:45.77.250.9",
    "kafka2:45.77.250.9",
    "peer0.orgauth.house.com:45.77.250.9",
    "peer1.orgauth.house.com:45.77.250.9",
    "peer0.orgcert.house.com:45.77.250.9",
    "peer1.orgcert.house.com:45.77.250.9",
    "kafka3:139.180.146.33",
    "peer0.orgcredit.house.com:139.180.146.33",
    "peer1.orgcredit.house.com:139.180.146.33"
  ]
}
```
Edit docker_compose_cfggen.json on host B as following:
```
{
  "domain": "house.com",
  "zookeepers": [
    {
      "host_name": "zookeeper1",
      "ports": [
        "2181",
        "2888",
        "3888"
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
      "ports": [
        "8050"
      ]
    }
  ],
  "hosts": [
    "zookeeper0:139.180.138.179",
    "kafka0:139.180.138.179",
    "orderer0.house.com:139.180.138.179",
    "zookeeper1:45.77.251.25",
    "kafka1:45.77.251.25",
    "orderer1.house.com:45.77.251.25",
    "zookeeper2:45.77.250.9",
    "kafka2:45.77.250.9",
    "peer0.orgauth.house.com:45.77.250.9",
    "peer1.orgauth.house.com:45.77.250.9",
    "peer0.orgcert.house.com:45.77.250.9",
    "peer1.orgcert.house.com:45.77.250.9",
    "kafka3:139.180.146.33",
    "peer0.orgcredit.house.com:139.180.146.33",
    "peer1.orgcredit.house.com:139.180.146.33"
  ]
}
```
Edit docker_compose_cfggen.json on host C as following:
```
{
  "domain": "house.com",
  "cas": [
    {
      "peer_org_name": "OrgAuth",
      "ports": [
        "7054"
      ],
      "admin_name": "admin",
      "admin_password": "adminpw"
    },
    {
      "peer_org_name": "OrgCert",
      "ports": [
        "8054"
      ],
      "admin_name": "admin",
      "admin_password": "adminpw"
    }
  ],
  "zookeepers": [
    {
      "host_name": "zookeeper2",
      "ports": [
        "2181",
        "2888",
        "3888"
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
      "org_name": "OrgAuth",
      "ports": [
        "7051",
        "7052",
        "7053"
      ]
    },
    {
      "peer_name": "peer1",
      "org_name": "OrgAuth",
      "ports": [
        "8051",
        "8052",
        "8053"
      ]
    },
    {
      "peer_name": "peer0",
      "org_name": "OrgCert",
      "ports": [
        "9051",
        "9052",
        "9053"
      ]
    },
    {
      "peer_name": "peer1",
      "org_name": "OrgCert",
      "ports": [
        "10051",
        "10052",
        "10053"
      ]
    }
  ],
  "clis": [
    {
      "cli_name": "cli",
      "core_peer_name": "peer0",
      "core_peer_org": "OrgAuth",
      "depends": [
        "peer0.orgauth.house.com",
        "peer1.orgauth.house.com",
        "peer0.orgcert.house.com",
        "peer1.orgcert.house.com"
      ]
    }
  ],
  "hosts": [
    "zookeeper0:139.180.138.179",
    "kafka0:139.180.138.179",
    "orderer0.house.com:139.180.138.179",
    "zookeeper1:45.77.251.25",
    "kafka1:45.77.251.25",
    "orderer1.house.com:45.77.251.25",
    "zookeeper2:45.77.250.9",
    "kafka2:45.77.250.9",
    "peer0.orgauth.house.com:45.77.250.9",
    "peer1.orgauth.house.com:45.77.250.9",
    "peer0.orgcert.house.com:45.77.250.9",
    "peer1.orgcert.house.com:45.77.250.9",
    "kafka3:139.180.146.33",
    "peer0.orgcredit.house.com:139.180.146.33",
    "peer1.orgcredit.house.com:139.180.146.33"
  ]
}
```
Edit docker_compose_cfggen.json on host D as following:
```
{
  "domain": "house.com",
  "cas": [
    {
      "peer_org_name": "OrgCredit",
      "ports": [
        "9054"
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
      "org_name": "OrgCredit",
      "ports": [
        "11051",
        "11052",
        "11053"
      ]
    },
    {
      "peer_name": "peer1",
      "org_name": "OrgCredit",
      "ports": [
        "12051",
        "12052",
        "12053"
      ]
    }
  ],
  "hosts": [
    "zookeeper0:139.180.138.179",
    "kafka0:139.180.138.179",
    "orderer0.house.com:139.180.138.179",
    "zookeeper1:45.77.251.25",
    "kafka1:45.77.251.25",
    "orderer1.house.com:45.77.251.25",
    "zookeeper2:45.77.250.9",
    "kafka2:45.77.250.9",
    "peer0.orgauth.house.com:45.77.250.9",
    "peer1.orgauth.house.com:45.77.250.9",
    "peer0.orgcert.house.com:45.77.250.9",
    "peer1.orgcert.house.com:45.77.250.9",
    "kafka3:139.180.146.33",
    "peer0.orgcredit.house.com:139.180.146.33",
    "peer1.orgcredit.house.com:139.180.146.33"
  ]
}
```

(c)Copy crypto-config files and genesis.block from host C to other hosts.

Execute on host A
```
cd /root/fabric/scripts/fabric-samples/190216/network
mkdir channel-artifacts
mkdir -p crypto-config/ordererOrganizations/house.com/orderers
```
Copy crypto-config files and genesis.block from host C to host A,executing in the directory /root/fabric/scripts/fabric-samples/190216/network of host C
```
scp channel-artifacts/genesis.block root@139.180.138.179:/root/fabric/scripts/fabric-samples/190216/network/channel-artifacts
scp -r crypto-config/ordererOrganizations/house.com/orderers/orderer0.house.com root@139.180.138.179:/root/fabric/scripts/fabric-samples/190216/network/crypto-config/ordererOrganizations/house.com/orderers
```
Execute on host B
```
cd /root/fabric/scripts/fabric-samples/190216/network
mkdir channel-artifacts
mkdir -p crypto-config/ordererOrganizations/house.com/orderers
```
Copy crypto-config files and genesis.block from host C to host B,executing in the directory /root/fabric/scripts/fabric-samples/190216/network of host C
```
scp channel-artifacts/genesis.block root@45.77.251.25:/root/fabric/scripts/fabric-samples/190216/network/channel-artifacts
scp -r crypto-config/ordererOrganizations/house.com/orderers/orderer1.house.com root@45.77.251.25:/root/fabric/scripts/fabric-samples/190216/network/crypto-config/ordererOrganizations/house.com/orderers
```
Execute on host D
```
cd /root/fabric/scripts/fabric-samples/190216/network
mkdir -p crypto-config/peerOrganizations/orgcredit.house.com
```
Copy crypto-config files from host C to host D,executing in the directory /root/fabric/scripts/fabric-samples/190216/network of host C
```
scp -r crypto-config/peerOrganizations/orgcredit.house.com/peers root@139.180.146.33:/root/fabric/scripts/fabric-samples/190216/network/crypto-config/peerOrganizations/orgcredit.house.com
scp -r crypto-config/peerOrganizations/orgcredit.house.com/ca root@139.180.146.33:/root/fabric/scripts/fabric-samples/190216/network/crypto-config/peerOrganizations/orgcredit.house.com
```

(d)Execute on each host of A,B,C,D
```
cd /root/fabric/scripts/fabric-samples/190216/network
chmod +x docker_compose_cfggen.sh
./docker_compose_cfggen.sh house # house is COMPOSE_PROJECT_NAME
```
It will generate docker-compose configuration files including zookeeper.yaml,kafka.yaml,docker-compose.yaml.

(e)Start zookeeper on each host of A,B,C,executing in the directory /root/fabric/scripts/fabric-samples/190216/network of each host
```
docker-compose -f zookeeper.yaml up -d
```

(f)Start kafka on each host of A,B,C,D,executing in the directory /root/fabric/scripts/fabric-samples/190216/network of each host
```
docker-compose -f kafka.yaml up -d
```

(g)Execute in the directory /root/fabric/scripts/fabric-samples/190216/network of each host of A,B,C,D
```
docker-compose -f docker-compose.yaml up -d
```
It will start service including orderer0.house.com on host A,orderer1.house.com on host B,peer0.orgauth.house.com,peer1.orgauth.house.com,peer0.orgcert.house.com,peer1.orgcert.house.com,ca_OrgAuth,ca_OrgCert,cli on host C,peer0.orgcredit.house.com,peer1.orgcredit.house.com,ca_OrgCredit on host D.

## 4.Create channels,and make each peer node join different channels

(a)Copy channel.go,channel.sh in fabric_tools to the directory /root/fabric/scripts/fabric-samples/190216/network of host C and enter the directory.

(b)Create and edit channel.json as following:
```
{
  "domain": "house.com",
  "channels": [
    {
      "channel_name": "auth",
      "orgs": [
        {
          "org_name": "OrgAuth",
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
        }
      ]
    },
    {
      "channel_name": "cert",
      "orgs": [
        {
          "org_name": "OrgCert",
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
        }
      ]
    },
    {
      "channel_name": "credit",
      "orgs": [
        {
          "org_name": "OrgCredit",
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
  ],
  "orderer": {
    "orderer_name": "orderer0",
    "port": "7050"
  },
  "cli_name": "cli"
}
```

(c)Execute
```
chmod +x channel.sh
./channel.sh channel.json
```
It will create a channel named auth which includes peer nodes in OrgAuth,a channel named cert which includes peer nodes in OrgCert and a channel named credit which includes peer nodes in OrgCredit.

## 5.Edit chaincodes

(a)Create a directory named chaincode in the directory /root/fabric/scripts/fabric-samples/190216 of host C and enter the created chaincode directory.Then created 3 directories named auth(stores the chaincode of renter authority),cert(stores the certificate of landlords),credit(stores the chaincode of renters' credit) respectively.

(b)Create and edit auth.go in the created directory auth as following:
```
package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"errors"
	"fmt"
)

func main()  {
	shim.Start(new(AuthChaincode))
}

type AuthChaincode struct {

}

func (this *AuthChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (this *AuthChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn,args:=stub.GetFunctionAndParameters()
	if fn=="check" {
		return this.check(stub,args)
	} else if fn=="add" {
		return this.add(stub,args)
	}
	return shim.Error("Method doesn't exist")
}

func (this *AuthChaincode) check(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}
	id:=args[0]
	name:=args[1]
	data,err:=stub.GetState(id)
	if err!=nil {
		return shim.Error(err.Error())
	}
	if data==nil {
		return shim.Success([]byte("false"))
	}
	if string(data)==name {
		return shim.Success([]byte("true"))
	} else {
		return shim.Success([]byte("false"))
	}
}

func (this *AuthChaincode) add(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}
	id:=args[0]
	name:=args[1]
	err=stub.PutState(id,[]byte(name))
	if err!=nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func checkArgsNum(args []string,n int) error {
	if len(args)!=n {
		return errors.New(fmt.Sprintf("%d parameter(s) required",n))
	}
	return nil
}
```
Create and edit cert.go in the created directory cert as following:
```
package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"errors"
	"fmt"
)

func main()  {
	shim.Start(new(CertChaincode))
}

type CertChaincode struct {

}

func (this *CertChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (this *CertChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn,args:=stub.GetFunctionAndParameters()
	if fn=="check" {
		return this.check(stub,args)
	} else if fn=="add" {
		return this.add(stub,args)
	}
	return shim.Error("Method doesn't exist")
}

func (this *CertChaincode) check(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}
	id:=args[0]
	name:=args[1]
	data,err:=stub.GetState(id)
	if err!=nil {
		return shim.Error(err.Error())
	}
	if data==nil {
		return shim.Success([]byte("false"))
	}
	if string(data)==name {
		return shim.Success([]byte("true"))
	} else {
		return shim.Success([]byte("false"))
	}
}

func (this *CertChaincode) add(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}
	id:=args[0]
	name:=args[1]
	err=stub.PutState(id,[]byte(name))
	if err!=nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func checkArgsNum(args []string,n int) error {
	if len(args)!=n {
		return errors.New(fmt.Sprintf("%d parameter(s) required",n))
	}
	return nil
}
```
Create and edit credit.go in the created directory credit as following:
```
package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"errors"
	"fmt"
)

func main()  {
	shim.Start(new(CreditChaincode))
}

type CreditChaincode struct {

}

func (this *CreditChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (this *CreditChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn,args:=stub.GetFunctionAndParameters()
	if fn=="check" {
		return this.check(stub,args)
	} else if fn=="add" {
		return this.add(stub,args)
	}
	return shim.Error("Method doesn't exist")
}

func (this *CreditChaincode) check(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}
	id:=args[0]
	data,err:=stub.GetState(id)
	if err!=nil {
		return shim.Error(err.Error())
	}
	if data==nil {
		return shim.Success([]byte("false"))
	}
	return shim.Success(data)
}

func (this *CreditChaincode) add(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}
	id:=args[0]
	credit:=args[1]
	err=stub.PutState(id,[]byte(credit))
	if err!=nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func checkArgsNum(args []string,n int) error {
	if len(args)!=n {
		return errors.New(fmt.Sprintf("%d parameter(s) required",n))
	}
	return nil
}
```

## 6.Install,instantiate and invoke chaincodes

(a)Copy chaincode.go,chaincode.sh in fabric_tools to the directory /root/fabric/scripts/fabric-samples/190216/network of host C and enter the directory.

(b)Create and edit auth.json as following:
```
{
  "domain": "house.com",
  "channels": [
    {
      "channel_name": "auth",
      "orgs": [
        {
          "org_name": "OrgAuth",
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
        }
      ]
    }
  ],
  "chaincode_name": "auth",
  "chaincode_version": "1.0",
  "orderer": {
    "orderer_name": "orderer0",
    "port": "7050"
  },
  "endorse": "AND ('OrgAuthMSP.member')",
  "cli_name": "cli"
}
```

(c)Execute
```
chmod +x chaincode.sh
./chaincode.sh -i auth.json
```
It will install and instantiate the chaincode auth.

(d)Create and edit cert.json as following:
```
{
  "domain": "house.com",
  "channels": [
    {
      "channel_name": "cert",
      "orgs": [
        {
          "org_name": "OrgCert",
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
        }
      ]
    }
  ],
  "chaincode_name": "cert",
  "chaincode_version": "1.0",
  "orderer": {
    "orderer_name": "orderer0",
    "port": "7050"
  },
  "endorse": "AND ('OrgCertMSP.member')",
  "cli_name": "cli"
}
```

(e)Execute
```
./chaincode.sh -i cert.json
```
It will install and instantiate the chaincode cert.

(f)Create and edit credit.json as following:
```
{
  "domain": "house.com",
  "channels": [
    {
      "channel_name": "credit",
      "orgs": [
        {
          "org_name": "OrgCredit",
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
  ],
  "chaincode_name": "credit",
  "chaincode_version": "1.0",
  "orderer": {
    "orderer_name": "orderer0",
    "port": "7050"
  },
  "endorse": "AND ('OrgCreditMSP.member')",
  "cli_name": "cli"
}
```

(g)Execute
```
./chaincode.sh -i credit.json
```
It will install and instantiate the chaincode credit.

(h)Invoke the chaincode auth.

Create and edit auth_test.sh as following:
```
#!/bin/bash

docker exec cli peer chaincode invoke -n auth -C auth -c '{"args":["add","1003","Jenny"]}' -o orderer1.house.com:8050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/house.com/orderers/orderer1.house.com/msp/tlscacerts/tlsca.house.com-cert.pem --peerAddresses peer0.orgauth.house.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgauth.house.com/peers/peer0.orgauth.house.com/tls/ca.crt

sleep 10
docker exec -e "CORE_PEER_ADDRESS=peer1.orgauth.house.com:8051" -e "CORE_PEER_LOCALMSPID=OrgAuthMSP" -e "CORE_PEER_TLS_ENABLED=true" -e "CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgauth.house.com/peers/peer1.orgauth.house.com/tls/server.crt" -e "CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgauth.house.com/peers/peer1.orgauth.house.com/tls/server.key" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgauth.house.com/peers/peer1.orgauth.house.com/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgauth.house.com/users/Admin@orgauth.house.com/msp" cli peer chaincode query -n auth -C auth -c '{"args":["check","1003","Jenny"]}'
```
Then execute
```
chmod +x auth_test.sh
./auth_test.sh
```

(i)Invoke the chaincode cert.

Create and edit cert_test.sh as following:
```
#!/bin/bash

docker exec -e "CORE_PEER_ADDRESS=peer0.orgcert.house.com:9051" -e "CORE_PEER_LOCALMSPID=OrgCertMSP" -e "CORE_PEER_TLS_ENABLED=true" -e "CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcert.house.com/peers/peer0.orgcert.house.com/tls/server.crt" -e "CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcert.house.com/peers/peer0.orgcert.house.com/tls/server.key" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcert.house.com/peers/peer0.orgcert.house.com/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcert.house.com/users/Admin@orgcert.house.com/msp" cli peer chaincode invoke -n cert -C cert -c '{"args":["add","1004","LaMeMei"]}' -o orderer1.house.com:8050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/house.com/orderers/orderer1.house.com/msp/tlscacerts/tlsca.house.com-cert.pem --peerAddresses peer0.orgcert.house.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcert.house.com/peers/peer0.orgcert.house.com/tls/ca.crt

sleep 10
docker exec -e "CORE_PEER_ADDRESS=peer1.orgcert.house.com:10051" -e "CORE_PEER_LOCALMSPID=OrgCertMSP" -e "CORE_PEER_TLS_ENABLED=true" -e "CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcert.house.com/peers/peer1.orgcert.house.com/tls/server.crt" -e "CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcert.house.com/peers/peer1.orgcert.house.com/tls/server.key" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcert.house.com/peers/peer1.orgcert.house.com/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcert.house.com/users/Admin@orgcert.house.com/msp" cli peer chaincode query -n cert -C cert -c '{"args":["check","1004","LaMeMei"]}'
```
Then execute
```
chmod +x cert_test.sh
./cert_test.sh
```

(j)Invoke the chaincode credit.

Create and edit credit_test.sh as following:
```
#!/bin/bash

docker exec -e "CORE_PEER_ADDRESS=peer0.orgcredit.house.com:11051" -e "CORE_PEER_LOCALMSPID=OrgCreditMSP" -e "CORE_PEER_TLS_ENABLED=true" -e "CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcredit.house.com/peers/peer0.orgcredit.house.com/tls/server.crt" -e "CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcredit.house.com/peers/peer0.orgcredit.house.com/tls/server.key" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcredit.house.com/peers/peer0.orgcredit.house.com/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcredit.house.com/users/Admin@orgcredit.house.com/msp" cli peer chaincode invoke -n credit -C credit -c '{"args":["add","1003","true"]}' -o orderer1.house.com:8050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/house.com/orderers/orderer1.house.com/msp/tlscacerts/tlsca.house.com-cert.pem --peerAddresses peer0.orgcredit.house.com:11051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcredit.house.com/peers/peer0.orgcredit.house.com/tls/ca.crt

sleep 10
docker exec -e "CORE_PEER_ADDRESS=peer1.orgcredit.house.com:12051" -e "CORE_PEER_LOCALMSPID=OrgCreditMSP" -e "CORE_PEER_TLS_ENABLED=true" -e "CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcredit.house.com/peers/peer1.orgcredit.house.com/tls/server.crt" -e "CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcredit.house.com/peers/peer1.orgcredit.house.com/tls/server.key" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcredit.house.com/peers/peer1.orgcredit.house.com/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgcredit.house.com/users/Admin@orgcredit.house.com/msp" cli peer chaincode query -n credit -C credit -c '{"args":["check","1003"]}'
```
Then execute
```
chmod +x credit_test.sh
./credit_test.sh
```
