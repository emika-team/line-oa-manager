package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"

	pb "github.com/emika-team/grpc-proto/line-oa/go"
	grpc_server "github.com/emika-team/line-oa-manager/cmd/grpc"
	channel "github.com/emika-team/line-oa-manager/pkg/http/channel"
	message "github.com/emika-team/line-oa-manager/pkg/http/message"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterLineOAMessageServer(s, grpc_server.NewGRPCHandler())
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})
	e.POST("/webhook", message.ReceiveMessage)
	e.GET("/channels", channel.GetChannel)
	e.GET("/channels/:channelId/messages/:messageId", message.GetContent)
	e.POST("/channels/:channelId/messages", message.SendMessage)
	e.Logger.Fatal(e.Start(":1323"))
}
