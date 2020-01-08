
#!/bin/bash
export ORDERER_CA=/opt/ws/crypto-config/ordererOrganizations/edunet.net/msp/tlscacerts/tlsca.edunet.net-cert.pem

if [ $# -lt 2 ];then
	echo "Usage : . setpeer.sh EDUNET|IITJ|IITKJP| <peerid>"
fi
export peerId=$2

if [[ $1 = "EDUNET" ]];then
	echo "Setting to organization EDUNET peer "$peerId
	export CORE_PEER_ADDRESS=$peerId.edunet.net:7051
	export CORE_PEER_LOCALMSPID=EDUNETMSP
	export CORE_PEER_TLS_CERT_FILE=/opt/ws/crypto-config/peerOrganizations/edunet.net/peers/$peerId.edunet.net/tls/server.crt
	export CORE_PEER_TLS_KEY_FILE=/opt/ws/crypto-config/peerOrganizations/edunet.net/peers/$peerId.edunet.net/tls/server.key
	export CORE_PEER_TLS_ROOTCERT_FILE=/opt/ws/crypto-config/peerOrganizations/edunet.net/peers/$peerId.edunet.net/tls/ca.crt
	export CORE_PEER_MSPCONFIGPATH=/opt/ws/crypto-config/peerOrganizations/edunet.net/users/Admin@edunet.net/msp
fi

if [[ $1 = "IITJ" ]];then
	echo "Setting to organization IITJ peer "$peerId
	export CORE_PEER_ADDRESS=$peerId.iitj.ac.in:7051
	export CORE_PEER_LOCALMSPID=IITJMSP
	export CORE_PEER_TLS_CERT_FILE=/opt/ws/crypto-config/peerOrganizations/iitj.ac.in/peers/$peerId.iitj.ac.in/tls/server.crt
	export CORE_PEER_TLS_KEY_FILE=/opt/ws/crypto-config/peerOrganizations/iitj.ac.in/peers/$peerId.iitj.ac.in/tls/server.key
	export CORE_PEER_TLS_ROOTCERT_FILE=/opt/ws/crypto-config/peerOrganizations/iitj.ac.in/peers/$peerId.iitj.ac.in/tls/ca.crt
	export CORE_PEER_MSPCONFIGPATH=/opt/ws/crypto-config/peerOrganizations/iitj.ac.in/users/Admin@iitj.ac.in/msp
fi

if [[ $1 = "IITKJP" ]];then
	echo "Setting to organization IITKJP peer "$peerId
	export CORE_PEER_ADDRESS=$peerId.iitkjp.ac.in:7051
	export CORE_PEER_LOCALMSPID=IITKJPMSP
	export CORE_PEER_TLS_CERT_FILE=/opt/ws/crypto-config/peerOrganizations/iitkjp.ac.in/peers/$peerId.iitkjp.ac.in/tls/server.crt
	export CORE_PEER_TLS_KEY_FILE=/opt/ws/crypto-config/peerOrganizations/iitkjp.ac.in/peers/$peerId.iitkjp.ac.in/tls/server.key
	export CORE_PEER_TLS_ROOTCERT_FILE=/opt/ws/crypto-config/peerOrganizations/iitkjp.ac.in/peers/$peerId.iitkjp.ac.in/tls/ca.crt
	export CORE_PEER_MSPCONFIGPATH=/opt/ws/crypto-config/peerOrganizations/iitkjp.ac.in/users/Admin@iitkjp.ac.in/msp
fi

	