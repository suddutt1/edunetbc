
#!/bin/bash
fabric-ca-client enroll  -u https://admin:adminpw@ca.iitkjp.ac.in:7054 --tls.certfiles /etc/hyperledger/fabric-ca-server-config/ca.iitkjp.ac.in-cert.pem 
fabric-ca-client affiliation add iitkjp  -u https://admin:adminpw@ca.iitkjp.ac.in:7054 --tls.certfiles /etc/hyperledger/fabric-ca-server-config/ca.iitkjp.ac.in-cert.pem 
