package lvsdk

import "github.com/fasthttp/websocket"

func ReadMutation(conn *websocket.Conn) Mutation {
	mt, bytes, err := conn.ReadMessage()
	if err != nil {
		return Mutation{}
	}
	if mt != websocket.TextMessage {
		PanicF("Invalid msg type %v", mt)
	}
	mut, err := DecodeMutation(bytes)
	if err != nil {
		return Mutation{}
	}
	return mut
}

func WriteMutation(conn *websocket.Conn, mut Mutation) {
	bytes, err := EncodeMutation(mut)
	PanicIfError(err)
	err = conn.WriteMessage(websocket.TextMessage, bytes)
	PanicIfError(err)
}
