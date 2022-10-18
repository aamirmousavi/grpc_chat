package main

import (
	"bufio"
	"context"
	"fmt"
	"grpc_chat/internal/chat"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	username := os.Args[1]

	conn, err := grpc.Dial("127.0.0.1:6565", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := chat.NewChatClient(conn)

	stream, err := client.Stream(context.TODO())
	if err != nil {
		panic(err)
	}

	defer stream.CloseSend()

	stream.Send(&chat.StreamRequest{
		Event: &chat.StreamRequest_Req_Authorization{
			Req_Authorization: &chat.StreamRequest_Authorization{
				Name: username,
			},
		},
	})

	go send(stream)

	if err := recive(stream); err != nil {
		panic(err)
	}
}

func recive(stream chat.Chat_StreamClient) error {
	for {
		res, err := stream.Recv()

		if s, ok := status.FromError(err); ok && s.Code() == codes.Canceled {
			log.Printf("Stream canceled \n")
			return nil
		} else if err == io.EOF {
			log.Printf("Stream closed by server\n")
			return nil
		} else if err != nil {
			return err
		}

		switch event := res.Event.(type) {
		case *chat.StreamResponse_RespMessage:
			log.Printf("msg: %#v\n", *event.RespMessage)
		case *chat.StreamResponse_RespCreateRoom:
			log.Printf("createRoom: %#v\n", *event.RespCreateRoom)
		case *chat.StreamResponse_RespErrorMessage:
			log.Printf("eventError: %#v\n", *event.RespErrorMessage)
		default:
			log.Printf("unexpected event from the serverL %T", event)
			return nil
		}
	}
}

func send(stream chat.Chat_StreamClient) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text()
		/*
			{"roomId":6,"text":"hiii grpc"}
		*/
		if err := parseAndSend(stream, text); err != nil {
			fmt.Printf("parse and send Error: %v\n", err.Error())
		}
	}
}
