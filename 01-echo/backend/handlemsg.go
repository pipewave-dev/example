package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/pipewave-dev/go-pkg/core/delivery"
	voAuth "github.com/pipewave-dev/go-pkg/core/domain/value-object/auth"
)

const (
	MsgTypeEchoReq = "ECHO_REQ"
	MsgTypeEchoRes = "ECHO_RES"
)

type handleMsg struct {
	i delivery.ModuleDelivery // Inject ModuleDelivery to call service layer and other utilities
}

// Text decoder for echo message, you can implement your own decoder based on your message format
func textDecoder(data []byte) string {
	return string(data)
}

func textEncoder(msg string) []byte {
	return []byte(msg)
}

func (h *handleMsg) HandleMessage(ctx context.Context, auth voAuth.WebsocketAuth, inputType string, data []byte) (outputType string, res []byte, err error) {
	if h.i == nil {
		panic("Module is not ready, do not serve server at this moment")
	}

	switch inputType {
	case MsgTypeEchoReq:
		msg := textDecoder(data)
		slog.Info("Received echo request", "message", msg)
		// Acknowledge the sender
		resMsg := textEncoder(
			fmt.Sprintf("Got [%s] at %s", msg, time.Now().Format(time.TimeOnly)),
		)
		return MsgTypeEchoRes,
			resMsg, nil

	default:
		slog.Warn("Received unsupported message type", "type", inputType)
		return "", nil, fmt.Errorf("unsupported message type: %s", inputType)
	}
}
