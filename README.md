# Blockchain network setup

## Copy source files
```sh
cd /home/blockchain
mkdir edunet
cd edunet
mkdir network api ui
cd network
git clone git remote add origin git@github.com:suddutt1/edunetbc.git . # Do not miss the .   

```

## Setup the network
All the commands must be executed from  network directory.  

First time setup. Run the following commands
 1. Download the binaries

```sh
 . ./downloadbin.sh

```

 2. To Generate the cryto config and other configurations
```sh
  . ./generateartifacts.sh
```


 3. Start the netowrk  

```sh
  . setenv.sh
  docker-compose up -d
```

 4. Build and join channel. Make sure that network is running

```sh

   docker exec -it cli bash -e ./buildandjoinchannel.sh

```

 5. Install and intantiate the chain codes
```sh
  docker exec -it cli bash -e  ./degreerecordmgmt_install.sh
```
6. Make the following /etc/hosts entries ( This will require sudo access for the blockchain user)

```sh

127.0.0.1	ca.iitj.ac.in        
127.0.0.1	ca.iitkjp.ac.in     
127.0.0.1	ca.edunet.net       
127.0.0.1	orderer0.edunet.net
127.0.0.1	peer0.iitj.ac.in    
127.0.0.1	peer0.edunet.net    
127.0.0.1	peer0.iitkjp.ac.in  

```
## When chain code is modified
To update the chain code , first update the chain code source files in chaincode/github.com/degreerecordmgmt directory.
Then run the following commands as appropriate

```sh
  docker exec -it cli bash -e  ./degreerecordmgmt_update.sh <version>
```

## To bring up an existing network
```sh

cd /home/blockchain/edunet/network
. setenv.sh
docker-compose up -d

```
## To destory  an existing network
```sh
cd /home/blockchain/edunet/network
. setenv.sh
docker-compose down
./clearVols.sh
./removeImages.sh

```
