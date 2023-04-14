# peershare
The project is a distributed file sharing system using peer-to-peer architecture, developed using the Go programming language and the Pion library to establish WebRTC data channel connections for file transfers between peers. Additionally, a signaling server was implemented to help peers to connect to each other.
The file sharing system consists of a sender peer and a receiver peer, where the sender selects the file to be sent via the CLI, and the server returns a unique share code. The receiver can then use the share code to directly obtain the file from the sender peer, with the transfer process being confirmed by both parties, providing added security.
The use of P2P architecture ensures fast and efficient file transfers, eliminating the need for a central server, and enhancing scalability. The system's security is further bolstered by the use of WebRTC data channels which ensure reliable and secure communication.
Overall, the project successfully demonstrates the feasibility of developing a secure and efficient file sharing system using P2P architecture, Go programming language, and the Pion library. The system can be further improved by implementing additional features, such as encryption and compression, to enhance its performance and security.
## Build and Running
```sh
just build
```
### Run Server
```sh
./bin/peershare --host 127.0.0.1 --port 5050 server
```

### Send File
```sh
./bin/peershare client send --file "[file path]"
```

### Recieve File
```sh
./bin/peershare client receive --out "output directory path" --code "shared code"
```