package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"logger-service/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// Write the log
	logEntry := data.LogEntry{
		Name: input.GetName(),
		Data: input.GetData(),
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{
			Result: "failed",
		}
		return res, err
	}

	// return response
	res := &logs.LogResponse{
		Result: "logged",
	}
	return res, nil
}

func (app *Config) gRPCListen() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	srv := grpc.NewServer()
	logs.RegisterLogServiceServer(srv, &LogServer{
		Models: app.Models,
	})

	log.Printf("gRPC server listening on port %s", grpcPort)

	if err := srv.Serve(listen); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)

	}
}
