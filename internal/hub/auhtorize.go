package hub

import (
	"fmt"
	"grpc_chat/internal/chat"
	"grpc_chat/internal/pg"
)

func authorization(stream chat.Chat_StreamServer) (*user, error) {
	var usr *pg.User
	for {
		req, err := stream.Recv()
		if err != nil {
			return nil, err
		}

		autho := req.GetReq_Authorization()
		if autho == nil {
			return nil, fmt.Errorf("first message must be authorization")
		}

		usr, err = pg.GetUser(autho.GetName())
		if err != nil {
			return nil, err
		}

		break
	}
	return &user{
		id:        usr.Id,
		username:  usr.Username,
		firstName: usr.FirstName,
		lastName:  usr.LastName,
		stream:    stream,
	}, nil
}
