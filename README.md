# tftpgo
 simple tftp protocol written in Golang
 ________________________________________>
 ```
 this only supports read requests (RRQ) and will serve a static payload.txt
 ```
 tftp server listens on port 69/udp (although changed here for non privileged use)
 
 ##the basic flow:
 1. Host sends RRQ packet to server.
 1. Server sends data packet to host
 1. Host sends Aknowledgement packet to server
 
 the transfer is concluded when a packet < 516 bytes is sent/received
 upon receiving packet < 516 byes the host will send aknowledgement packet.
 
 from network programming with go

rfc1350
