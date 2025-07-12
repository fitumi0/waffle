package server

import (
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

		log.Info().Msgf("Received message: %v", in)
	}
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
