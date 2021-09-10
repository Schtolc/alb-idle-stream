package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Schtolc/alb-idle-stream/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

func main() {
	var opts []grpc.DialOption
	// client gRPC keepalive ping
	opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time: 15 * time.Second,
	}))
	opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(nil)))

	conn, err := grpc.Dial(os.Args[1], opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := echo.NewEchoClient(conn)
	stream, err := client.GetData(context.Background())
	if err != nil {
		panic(err)
	}

	for {
		fmt.Printf("Sending `hello` to server.\n")
		err = stream.Send(&echo.GetDataRequest{
			Payload: "hello",
		})
		if err != nil {
			panic(err)
		}

		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return
			}
			panic(err)
		}

		fmt.Printf("Received `%s` from server. Waiting before next request...\n", in.Payload)
		// ALB Idle timeout is 60 seconds
		time.Sleep(70 * time.Second)
	}
}
