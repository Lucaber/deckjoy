package ipc

//go:generate protoc --go_out=plugins=grpc:.. --go_opt=paths=source_relative daemon.proto
