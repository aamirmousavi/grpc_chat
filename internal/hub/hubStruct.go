package hub

import (
	"fmt"
	"grpc_chat/internal/chat"
	"sync"
)

type hub struct {
	users syncUsers
}

type (
	syncUsers struct {
		sm sync.Map
	}
	user struct {
		id        int
		username  string
		firstName string
		lastName  string
		//Conn *grpc.ClientConn
		stream chat.Chat_StreamServer
	}
	message struct {
		profile *user
		value   *userMessage
	}
	userMessage struct {
		//grpc chat_req
	}
)

func newHub() *hub {
	h := &hub{
		users: syncUsers{},
	}
	return h
}

func (s *syncUsers) get(id int) (*user, bool) {
	usr, ok := s.sm.Load(id)
	if ok {
		return usr.(*user), true
	}
	return nil, false
}

func (h *hub) userAdd(usr *user) error {
	_, load := h.users.sm.LoadOrStore(usr.id, usr)
	if load {
		return fmt.Errorf("Already Registerd")
	}
	println("add user: ", usr.id)
	return nil
}

func (h *hub) userRemove(usr *user) {
	fmt.Printf("user %v removed\n", usr.id)
	h.users.sm.Delete(usr.id)
}
