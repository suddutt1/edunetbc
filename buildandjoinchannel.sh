
#!/bin/bash -e




	
	echo "Building channel for degreerecordchannel" 
	
	. setpeer.sh EDUNET peer0
	export CHANNEL_NAME="degreerecordchannel"
	peer channel create -o orderer0.edunet.net:7050 -c $CHANNEL_NAME -f ./degreerecordchannel.tx --tls true --cafile $ORDERER_CA -t 1000s
	
		
        
            . setpeer.sh EDUNET peer0
            export CHANNEL_NAME="degreerecordchannel"
			peer channel join -b $CHANNEL_NAME.block
		
	
		
        
            . setpeer.sh IITJ peer0
            export CHANNEL_NAME="degreerecordchannel"
			peer channel join -b $CHANNEL_NAME.block
		
	
		
        
            . setpeer.sh IITKJP peer0
            export CHANNEL_NAME="degreerecordchannel"
			peer channel join -b $CHANNEL_NAME.block
		
	
	
	
	
		. setpeer.sh EDUNET peer0
		export CHANNEL_NAME="degreerecordchannel"
		peer channel update -o  orderer0.edunet.net:7050 -c $CHANNEL_NAME -f ./degreerecordchannelEDUNETMSPAnchor.tx --tls --cafile $ORDERER_CA 
	

	
	
		. setpeer.sh IITJ peer0
		export CHANNEL_NAME="degreerecordchannel"
		peer channel update -o  orderer0.edunet.net:7050 -c $CHANNEL_NAME -f ./degreerecordchannelIITJMSPAnchor.tx --tls --cafile $ORDERER_CA 
	

	
	
		. setpeer.sh IITKJP peer0
		export CHANNEL_NAME="degreerecordchannel"
		peer channel update -o  orderer0.edunet.net:7050 -c $CHANNEL_NAME -f ./degreerecordchannelIITKJPMSPAnchor.tx --tls --cafile $ORDERER_CA 
	




