##tftpgo

 **simple tftp protocol written in Golang**
 
 
             2 bytes     string    1 byte     string   1 byte
            ------------------------------------------------
           | Opcode |  Filename  |   0  |    Mode    |   0  |
            ------------------------------------------------

 tftp server listens on port 69/udp (although changed here for non privileged use)
 
 ## the basic flow:
 1. Host sends RRQ packet to server.
 1. Server sends data packet to host
 1. Host sends Aknowledgement packet to server
 
 the transfer is concluded when a packet < 516 bytes is sent/received
 upon receiving packet < 516 byes the host will send aknowledgement packet.
 
  ```
 this only supports read requests (RRQ) and will serve a static payload.txt
 ```
 from network programming with go

rfc1350
