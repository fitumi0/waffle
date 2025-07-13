package server

import (
	"errors"
	"net"

	pb "github.com/fitumi0/waffle/gen/gmp"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedMessengerServiceServer
}

func (s *Server) MessageStream(stream pb.MessengerService_MessageStreamServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			return err
		}

		if in.GetMessage() == nil {
			return errors.New("Empty message")
		}

		ack := &pb.ServerToClient{
			Event: &pb.ServerToClient_Ack{
				Ack: &pb.Ack{
					Success: true,
				},
			},
		}

		stream.Send(ack)

		go s.broadcast(in.GetMessage())
	}
}

func (s *Server) broadcast(msg *pb.Message) {

}

func Start(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMessengerServiceServer(grpcServer, &Server{})
	log.Printf("gRPC сервер запущен на %s", port)

	return grpcServer.Serve(lis)
}
