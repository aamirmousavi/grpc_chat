package main

import "grpc_chat/internal/chat"

func parseAndSend(st chat.Chat_StreamClient, text string) error {
	sendMessage, err := Parse[chat.StreamRequest_SendMessage](text)
	if err != nil {
		createRoom, err := Parse[chat.StreamRequest_CreateRoom](text)
		if err != nil {
			// autho, err := Parse[chat.StreamRequest_Authorization](text)
			// if err != nil {
			return err
			// }
			// return sendAuth(st, &autho)
		}
		return sendCreateRoom(st, &createRoom)
	}
	return sendSendMessage(st, &sendMessage)
}

func sendAuth(st chat.Chat_StreamClient, req *chat.StreamRequest_Authorization) error {
	return st.Send(&chat.StreamRequest{
		Event: &chat.StreamRequest_Req_Authorization{
			Req_Authorization: req,
		},
	})
}

func sendCreateRoom(st chat.Chat_StreamClient, req *chat.StreamRequest_CreateRoom) error {
	return st.Send(&chat.StreamRequest{
		Event: &chat.StreamRequest_ReqCreateRoom{
			ReqCreateRoom: req,
		},
	})
}

func sendSendMessage(st chat.Chat_StreamClient, req *chat.StreamRequest_SendMessage) error {
	return st.Send(&chat.StreamRequest{
		Event: &chat.StreamRequest_Req_SendMessage{
			Req_SendMessage: req,
		},
	})
}
