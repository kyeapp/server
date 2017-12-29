#!/bin/bash
go fmt
GOPATH="C:\Users\KevinD\Documents\GitHub\server" go build
if [ $? -eq 0 ]
then
	echo "spinning up server..."
	./server.exe
else 
	echo "build failed."
fi