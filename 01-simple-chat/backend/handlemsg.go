package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/pipewave-dev/go-pkg/core/delivery"
	voAuth "github.com/pipewave-dev/go-pkg/core/domain/value-object/auth"
	"github.com/samber/lo"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	// C2S
	MsgTypeChatSendMsg = "CHAT_SEND_MSG"
	MsgTypeChatTyping  = "CHAT_TYPING"

	// S2C
	MsgTypeChatIncomingMsg = "CHAT_INCOMING_MSG"
	MsgTypeChatUserTyping  = "CHAT_USER_TYPING"
	MsgTypeChatAck         = "CHAT_ACK"
	MsgTypeChatFail        = "CHAT_FAIL"
	MsgTypeEchoResponse    = "ECHO_RESPONSE"
)

type handleMsg struct {
	i delivery.ModuleDelivery
}

// --- Chat Event Payloads ---

// C2S (Client-to-Server)
type ChatSendMsg struct {
	ToUserID string `json:"to_user_id" msgpack:"to_user_id"`
	Content  string `json:"content" msgpack:"content"`
}

type ChatSendMsgAck struct {
	Ok bool `json:"ok" msgpack:"ok"`
}
type ChatSendMsgFail struct {
	Reason string `json:"reason" msgpack:"reason"`
}

type ChatTyping struct {
	ToUserID string `json:"to_user_id" msgpack:"to_user_id"`
}

// S2C (Server-to-Client)
type ChatIncomingMsg struct {
	FromUserID string `json:"from_user_id" msgpack:"from_user_id"`
	Content    string `json:"content" msgpack:"content"`
	Timestamp  int64  `json:"timestamp" msgpack:"timestamp"`
}

type ChatUserTyping struct {
	FromUserID string `json:"from_user_id" msgpack:"from_user_id"`
}

func (h *handleMsg) HandleMessage(ctx context.Context, auth voAuth.WebsocketAuth, inputType string, data []byte) (outputType string, res []byte, err error) {
	if h.i == nil {
		panic("Module is not ready, do not serve server at this moment")
	}

	fmt.Printf("Received message [%s] from %s\n", inputType, auth.UserID)

	switch inputType {
	case MsgTypeChatSendMsg:
		var msg ChatSendMsg
		if err := msgpack.Unmarshal(data, &msg); err != nil {
			slog.Error("Invalid payload", "error", err)
			return "", nil, nil
		}

		// Forward message to destination user
		isOnline, err := h.i.Services().CheckOnline(ctx, msg.ToUserID)
		if err != nil {
			slog.Error("Failed to check online", "error", err)
			return "", nil, nil
		}
		if !isOnline {
			return MsgTypeChatFail,
				lo.Must(msgpack.Marshal(ChatSendMsgFail{
					Reason: "User is not online",
				})), nil
		}
		err = h.i.Services().SendToUser(ctx,
			msg.ToUserID,
			MsgTypeChatIncomingMsg,
			lo.Must(msgpack.Marshal(ChatIncomingMsg{
				FromUserID: auth.UserID,
				Content:    msg.Content,
				Timestamp:  time.Now().Unix(),
			})))
		if err != nil {
			slog.Error("Failed to send message", "error", err)
			return MsgTypeChatFail,
				lo.Must(msgpack.Marshal(ChatSendMsgFail{
					Reason: "Failed to send message",
				})), nil
		}
		// Acknowledge the sender
		return MsgTypeChatAck,
			lo.Must(msgpack.Marshal(ChatSendMsgAck{
				Ok: true,
			})), nil

	case MsgTypeChatTyping:
		var typing ChatTyping
		if err := msgpack.Unmarshal(data, &typing); err != nil {
			slog.Error("Invalid payload", "error", err)
			return "", nil, nil
		}

		// Forward typing status to destination user
		err = h.i.Services().SendToUser(ctx, typing.ToUserID, MsgTypeChatUserTyping,
			lo.Must(msgpack.Marshal(ChatUserTyping{
				FromUserID: auth.UserID,
			})))
		if err != nil {
			slog.Error("Failed to send message", "error", err)
		}

		return "", nil, nil

	default:
		// Default echo behavior for unknown events
		return MsgTypeEchoResponse,
			fmt.Appendf(nil, "Got [ %s ] at %s", string(data), time.Now().Format(time.TimeOnly)),
			nil
	}
}
