
#!/bin/bash
fabric-ca-client enroll  -u https://admin:adminpw@ca.iitj.ac.in:7054 --tls.certfiles /etc/hyperledger/fabric-ca-server-config/ca.iitj.ac.in-cert.pem 
fabric-ca-client affiliation add iitj  -u https://admin:adminpw@ca.iitj.ac.in:7054 --tls.certfiles /etc/hyperledger/fabric-ca-server-config/ca.iitj.ac.in-cert.pem 
