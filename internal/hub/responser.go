package hub

import (
	"fmt"
	"grpc_chat/internal/chat"
	"grpc_chat/internal/pg"
	"time"
)

func (h *hub) proccessResponse(usr *user, req *chat.StreamRequest) {

	switch event := req.Event.(type) {

	case *chat.StreamRequest_Req_SendMessage:
		fmt.Println("send msg recived")
		h.response_sendMessage(usr, event.Req_SendMessage)

	case *chat.StreamRequest_ReqCreateRoom:
		fmt.Println("create room recived")
		h.response_createRoom(usr, event.ReqCreateRoom)

	}

}

func (h *hub) response_sendMessage(usr *user, req *chat.StreamRequest_SendMessage) {
	room, err := pg.GetRoom(int(req.GetRoomId()))
	if err != nil {
		usr.sendError("SendMessage", err)
		return
	}

	message := &pg.Message{
		RoomId:    int(req.GetRoomId()),
		SenderId:  usr.id,
		Text:      req.GetText(),
		CreatedAt: time.Now(),
	}

	if err := pg.SendMessage(message); err != nil {
		usr.sendError("SendMessage", err)
		return
	}

	for _, roomatesId := range room.Users {
		roomate, ok := h.users.get(int(roomatesId))
		if !ok || roomate.id == usr.id {
			continue
		}
		roomate.send_Message(message)
	}
}

func (h *hub) response_createRoom(usr *user, req *chat.StreamRequest_CreateRoom) {
	room := &pg.Room{
		Users: req.GetUsers(),
	}
	if err := pg.CreateRoom(room); err != nil {
		usr.sendError("CreateRoom", err)
		return
	}

}
