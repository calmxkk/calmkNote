package main

import (
	"google.golang.org/grpc"
	"grpcdemo/server"
	"net"

	pb "grpcdemo/proto/employee"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterEmployeeServiceServer(grpcServer, &server.Employee{Employees: make(map[int32]*pb.Employee)})
	err = grpcServer.Serve(lis)
	if err != nil {
		return
	}
}
