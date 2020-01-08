#!/bin/bash
if [[ ! -z "$1" ]]; then  
	. setpeer.sh EDUNET peer0 
export CHANNEL_NAME="degreerecordchannel"
	peer chaincode install -n degreerecordmgmt -v $1 -l golang -p  github.com/degreerecordmgmt
	. setpeer.sh IITJ peer0 
export CHANNEL_NAME="degreerecordchannel"
	peer chaincode install -n degreerecordmgmt -v $1 -l golang -p  github.com/degreerecordmgmt
	. setpeer.sh IITKJP peer0 
export CHANNEL_NAME="degreerecordchannel"
	peer chaincode install -n degreerecordmgmt -v $1 -l golang -p  github.com/degreerecordmgmt
	peer chaincode upgrade -o orderer0.edunet.net:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C degreerecordchannel -n degreerecordmgmt -v $1 -c '{"Args":["init",""]}' -P " OR( 'EDUNETMSP.member','IITJMSP.member','IITKJPMSP.member' ) " 
else
	echo ". degreerecordmgmt_updchain.sh  <Version Number>" 
fi
