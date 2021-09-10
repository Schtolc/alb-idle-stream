package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/Schtolc/alb-idle-stream/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	echo.UnimplementedEchoServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetData(stream echo.Echo_GetDataServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		err = stream.Send(&echo.Data{
			Payload: fmt.Sprintf("echo: %s", in.Payload),
		})
		if err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 9999))
	if err != nil {
		panic(err)
	}

	var opts []grpc.ServerOption
	// server gRPC keepalive ping
	opts = append(opts, grpc.KeepaliveParams(keepalive.ServerParameters{
		Time: 15 * time.Second,
	}))

	grpcServer := grpc.NewServer(opts...)
	echo.RegisterEchoServer(grpcServer, NewServer())
	grpcServer.Serve(lis)
}
