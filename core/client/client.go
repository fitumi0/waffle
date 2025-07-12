package client

import (
	"context"

	"github.com/rs/zerolog/log"

	pb "github.com/fitumi0/waffle/gen/messenger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	conn      *grpc.ClientConn
	stream    pb.MessengerService_MessageStreamClient
	cancel    context.CancelFunc
	userID    string
	listeners []MessageListener
}

type MessageListener interface {
	OnMessage(msg *pb.Message)
	// OnPresence(presence *pb.PresenceUpdate)
	OnAck(ack *pb.Ack)
}

func NewClient(ctx context.Context, serverAddr, authToken string) (*Client, error) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	streamCtx, cancel := context.WithCancel(ctx)
	stream, err := pb.NewMessengerServiceClient(conn).MessageStream(streamCtx)
	if err != nil {
		cancel()
		return nil, err
	}

	c := &Client{
		conn:   conn,
		stream: stream,
		cancel: cancel,
		userID: authToken,
	}

	go c.listen()
	return c, nil
}

func (c *Client) listen() {
	for {
		in, err := c.stream.Recv()
		if err != nil {
			log.Printf("recv error: %v", err)
			return
		}

		switch x := in.Event.(type) {
		case *pb.ServerToClient_Message:
			for _, l := range c.listeners {
				l.OnMessage(x.Message)
			}
		}
	}
}

func (c *Client) SendMessage(chatID string, text string) error {
	msg := &pb.ClientToServer{
		Event: &pb.ClientToServer_Message{
			Message: &pb.Message{
				ChatId:    chatID,
				UserId:    c.userID,
				Timestamp: timestamppb.Now(),
				Attachments: []*pb.Attachment{
					{
						Type: pb.AttachmentType_TEXT,
						Id:   text,
					},
				},
			},
		},
	}
	return c.stream.Send(msg)
}
