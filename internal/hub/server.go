package hub

import (
	"grpc_chat/internal/chat"
	"io"
)

var Hub = newHub()

func (h *hub) Stream(stream chat.Chat_StreamServer) error {
	return h.onConnect(stream)
}

func (h *hub) onConnect(stream chat.Chat_StreamServer) error {
	user, err := authorization(stream)
	if err != nil {
		return err
	}

	if err := h.userAdd(user); err != nil {
		return err
	}

	defer h.userRemove(user)

	return h.readStreamData(user)
}

func (h *hub) readStreamData(usr *user) error {
	for {
		req, err := usr.stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		// if req == nil {
		// 	continue
		// }
		go h.proccessResponse(usr, req)
	}
	return nil
}
