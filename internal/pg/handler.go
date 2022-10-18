package pg

import (
	"fmt"
	"time"
)

type Room struct {
	Id        int
	Users     []int64
	CreatedAt time.Time
}
type User struct {
	Id        int
	Username  string
	FirstName string
	LastName  string
}

func GetUser(userName string) (*User, error) {
	var usr User
	if err := Hand.Db.QueryRow("SELECT id, username, first_name, last_name FROM public.users WHERE username = $1", userName).Scan(&usr.Id, &usr.Username, &usr.FirstName, &usr.LastName); err != nil {
		return nil, err
	}
	return &usr, nil
}

func CreateRoom(room *Room) error {

	err := Hand.Db.QueryRow("INSERT INTO public.room DEFAULT VALUES RETURNING id;").Scan(&room.Id)
	if err != nil {
		return err
	}

	query := "INSERT INTO public.user_list (\"user\", room_id) VALUES "
	param := []interface{}{}
	for i := range room.Users {
		query += fmt.Sprintf("($%d, %v)", i+1, room.Id)
		param = append(param, room.Users[i])
		if i+1 < len(room.Users) {
			query += ","
		}
	}
	if _, err := Hand.Db.Exec(query, param...); err != nil {
		return err
	}

	return nil
}

func GetRoom(id int) (*Room, error) {
	var room Room
	if err := Hand.Db.QueryRow("SELECT * FROM public.room WHERE id = $1", id).Scan(&room.Id, &room.CreatedAt); err != nil {
		return nil, err
	}

	rows, err := Hand.Db.Query("SELECT \"user\" FROM public.user_list WHERE room_id = $1", room.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	usrs := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		usrs = append(usrs, id)
	}
	room.Users = usrs
	return &room, nil
}

type Message struct {
	Id        int       `json:"id"`
	RoomId    int       `json:"room_id"`
	SenderId  int       `json:"sender_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func SendMessage(msg *Message) error {
	err := Hand.Db.QueryRow("INSERT INTO public.message (room_id, sender_id, text) VALUES ($1, $2, $3) RETURNING id;", msg.RoomId, msg.SenderId, msg.Text).Scan(&msg.Id)
	if err != nil {
		return err
	}
	return nil
}

func GetMessage(roomId int, limit, offset int) ([]*Message, error) {
	rows, err := Hand.Db.Query("SELECT * FROM public.message WHERE room_id = $1 LIMIT $2 OFFSET $3", roomId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	msgs := make([]*Message, 0)
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.Id, &msg.RoomId, &msg.SenderId, &msg.Text, &msg.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, &msg)
	}
	return msgs, nil
}
