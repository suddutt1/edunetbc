
#!/bin/bash -e
export PWD=`pwd`

export FABRIC_CFG_PATH=$PWD
export ARCH=$(uname -s)
export CRYPTOGEN=$PWD/bin/cryptogen
export CONFIGTXGEN=$PWD/bin/configtxgen

function generateArtifacts() {
	
	echo " *********** Generating artifacts ************ "
	echo " *********** Deleting old certificates ******* "
	
        rm -rf ./crypto-config
	
        echo " ************ Generating certificates ********* "
	
        $CRYPTOGEN generate --config=$FABRIC_CFG_PATH/crypto-config.yaml
        
        echo " ************ Generating tx files ************ "
	
		$CONFIGTXGEN -profile OrdererGenesis -outputBlock ./genesis.block
		
		$CONFIGTXGEN -profile degreerecordchannel -outputCreateChannelTx ./degreerecordchannel.tx -channelID degreerecordchannel
		
		echo "Generating anchor peers tx files for  EDUNET"
		$CONFIGTXGEN -profile degreerecordchannel -outputAnchorPeersUpdate  ./degreerecordchannelEDUNETMSPAnchor.tx -channelID degreerecordchannel -asOrg EDUNETMSP
		
		echo "Generating anchor peers tx files for  IITJ"
		$CONFIGTXGEN -profile degreerecordchannel -outputAnchorPeersUpdate  ./degreerecordchannelIITJMSPAnchor.tx -channelID degreerecordchannel -asOrg IITJMSP
		
		echo "Generating anchor peers tx files for  IITKJP"
		$CONFIGTXGEN -profile degreerecordchannel -outputAnchorPeersUpdate  ./degreerecordchannelIITKJPMSPAnchor.tx -channelID degreerecordchannel -asOrg IITKJPMSP
		

		

}
function generateDockerComposeFile(){
	OPTS="-i"
	if [ "$ARCH" = "Darwin" ]; then
		OPTS="-it"
	fi
	cp  docker-compose-template.yaml  docker-compose.yaml
	
	
	cd  crypto-config/peerOrganizations/edunet.net/ca
	PRIV_KEY=$(ls *_sk)
	cd ../../../../
	sed $OPTS "s/EDUNET_PRIVATE_KEY/${PRIV_KEY}/g"  docker-compose.yaml
	
	
	cd  crypto-config/peerOrganizations/iitj.ac.in/ca
	PRIV_KEY=$(ls *_sk)
	cd ../../../../
	sed $OPTS "s/IITJ_PRIVATE_KEY/${PRIV_KEY}/g"  docker-compose.yaml
	
	
	cd  crypto-config/peerOrganizations/iitkjp.ac.in/ca
	PRIV_KEY=$(ls *_sk)
	cd ../../../../
	sed $OPTS "s/IITKJP_PRIVATE_KEY/${PRIV_KEY}/g"  docker-compose.yaml
	
}
generateArtifacts 
cd $PWD
generateDockerComposeFile
cd $PWD

