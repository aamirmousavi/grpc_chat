package hub

import (
	"grpc_chat/internal/chat"
	"grpc_chat/internal/pg"
)

func (usr *user) sendError(event string, err error) error {
	return usr.stream.Send(
		&chat.StreamResponse{
			Event: &chat.StreamResponse_RespErrorMessage{
				RespErrorMessage: &chat.StreamResponse_ErrorMessage{
					Event: event,
					Couse: err.Error(),
				},
			},
		},
	)
}

func (usr *user) send_Message(msg *pg.Message) error {
	return usr.stream.Send(
		&chat.StreamResponse{
			Event: &chat.StreamResponse_RespMessage{
				RespMessage: &chat.StreamResponse_Message{
					RoomId:    int64(msg.RoomId),
					UserId:    int64(msg.SenderId),
					Text:      msg.Text,
					CreatedAt: msg.CreatedAt.UnixMilli(),
				},
			},
		},
	)
}

func (usr *user) send_CreateRoom(msg *pg.Room) error {
	return usr.stream.Send(
		&chat.StreamResponse{
			Event: &chat.StreamResponse_RespCreateRoom{
				RespCreateRoom: &chat.StreamResponse_CreateRoom{
					RoomId:    int64(msg.Id),
					Users:     msg.Users,
					CreatedAt: msg.CreatedAt.UnixMilli(),
				},
			},
		},
	)
}
