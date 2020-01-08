#!/bin/bash
. setpeer.sh EDUNET peer0 
export CHANNEL_NAME="degreerecordchannel"
peer chaincode install -n degreerecordmgmt -v 1.0 -l golang -p  github.com/degreerecordmgmt
. setpeer.sh IITJ peer0 
export CHANNEL_NAME="degreerecordchannel"
peer chaincode install -n degreerecordmgmt -v 1.0 -l golang -p  github.com/degreerecordmgmt
. setpeer.sh IITKJP peer0 
export CHANNEL_NAME="degreerecordchannel"
peer chaincode install -n degreerecordmgmt -v 1.0 -l golang -p  github.com/degreerecordmgmt
peer chaincode instantiate -o orderer0.edunet.net:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C degreerecordchannel -n degreerecordmgmt -v 1.0 -c '{"Args":["init",""]}' -P " OR( 'EDUNETMSP.member','IITJMSP.member','IITKJPMSP.member' ) " 
