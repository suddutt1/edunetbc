
#!/bin/bash
fabric-ca-client enroll  -u https://admin:adminpw@ca.edunet.net:7054 --tls.certfiles /etc/hyperledger/fabric-ca-server-config/ca.edunet.net-cert.pem 
fabric-ca-client affiliation add edunet  -u https://admin:adminpw@ca.edunet.net:7054 --tls.certfiles /etc/hyperledger/fabric-ca-server-config/ca.edunet.net-cert.pem 
