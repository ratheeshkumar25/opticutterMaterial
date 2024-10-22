package server

import (
	"fmt"
	"log"
	"net"

	"github.com/ratheeshkumar25/opt_cut_material_service/pkg/handlers"
	pb "github.com/ratheeshkumar25/opt_cut_material_service/pkg/pb"
	"google.golang.org/grpc"
)

func NewGrpcMaterialServer(port string, handlr *handlers.MaterialHandler) error {
	log.Println("Connecting to gRPC server")
	addr := fmt.Sprintf(":%s", port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("error creating listener on %v", addr)
		return err
	}
	grpc := grpc.NewServer()
	pb.RegisterMaterialServiceServer(grpc, handlr)
	log.Printf("listening on gRPC server %v", port)
	err = grpc.Serve(lis)
	if err != nil {
		log.Println("error connecting to gRPC server")
		return err
	}
	return nil

}
