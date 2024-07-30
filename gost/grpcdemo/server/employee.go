package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "grpcdemo/proto/employee"
)

type Employee struct {
	pb.UnimplementedEmployeeServiceServer
	Employees map[int32]*pb.Employee
}

func (e *Employee) GetEmployee(ctx context.Context, req *pb.EmployeeRequest) (*pb.Employee, error) {
	employee, ok := e.Employees[req.GetId()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "employee not found")
	}
	return employee, nil
}
