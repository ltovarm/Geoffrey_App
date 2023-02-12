#! /usr/bin/env python3

#----- A simple TCP based server program in Python using send() function -----
import socket

# Create a stream based socket(i.e, a TCP socket)
# operating on IPv4 addressing scheme
serverSocket = socket.socket(socket.AF_INET, socket.SOCK_STREAM);
serverSocket.settimeout(100)

# Bind and listen
serverSocket.bind(("172.29.108.252",9090));
serverSocket.listen();

dataFromClient = ""
# Accept connections
while(dataFromClient != "exit"):

    (clientConnected, clientAddress) = serverSocket.accept();
    print("Accepted a connection request from %s:%s"%(clientAddress[0], clientAddress[1]));
    dataFromClient = clientConnected.recv(1024)
    if (len(dataFromClient) > 0):
        print(dataFromClient.decode());
    else:
        print("Timeout")

    # Send some data back to the client
    clientConnected.send("Hello Client!".encode());