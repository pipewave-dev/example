package handlemsg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/pipewave-dev/go-pkg/core/delivery"
	voAuth "github.com/pipewave-dev/go-pkg/core/domain/value-object/auth"
	wsSv "github.com/pipewave-dev/go-pkg/core/service/websocket"
)

var catFactClient = &http.Client{Timeout: 3 * time.Second}

const catFactAPIURL = "https://catfact.ninja/fact"

const (
	// C2S
	MsgTypeEcho          = "ECHO"
	MsgTypeFireAndForget = "FIRE_AND_FORGET"
	MsgTypeErrorTest     = "ERROR_TEST"
	MsgTypeTestReqRes    = "REQ_RES"
	MsgTypePing          = "PING"

	// S2C
	MsgTypeEchoResponse = "ECHO_RESPONSE"
	MsgTypeMemInfo      = "MEM_INFO"
	MsgTypeCatFact      = "CAT_FACT"
	MsgTypePong         = "PONG"
	MsgTypeWelcome      = "WELCOME"

	//
)

// Handler handles all WebSocket message types for the echo-demo.
type Handler struct {
	pw       delivery.ModuleDelivery
	mu       sync.RWMutex
	allUsers map[string]struct{} // all connected users
}

// New creates a Handler, registers connection lifecycle hooks for user tracking.
func New(pw delivery.ModuleDelivery) *Handler {
	h := &Handler{
		pw:       pw,
		allUsers: make(map[string]struct{}),
	}

	// Track connected users for broadcasting
	pw.Services().OnNewRegister().Register(
		wsSv.OnNewWsKeyName("echoDemo"),
		func(conn wsSv.WebsocketConn) error {
			userID := conn.Auth().UserID
			if userID == "" {
				return nil
			}
			aErr1 := pw.Services().SendToSession(
				context.Background(),
				conn.Auth().UserID,
				conn.Auth().InstanceID,
				MsgTypeWelcome,
				fmt.Appendf(nil, "Welcome, user %s!", userID),
			)
			if aErr1 != nil {
				slog.Error("[OnNewRegister] SendToSession failed", "userID", userID, "instanceID", conn.Auth().InstanceID, "error", aErr1)
			}
			acked, aErr2 := pw.Services().SendToSessionWithAck(
				context.Background(),
				conn.Auth().UserID,
				conn.Auth().InstanceID,
				MsgTypeWelcome,
				fmt.Appendf(nil, "Welcome <has ack>, user %s!", userID),
				time.Second*5,
			)
			if aErr2 != nil {
				slog.Error("[OnNewRegister] SendToSessionWithAck failed", "userID", userID, "instanceID", conn.Auth().InstanceID, "error", aErr2)
			}
			if !acked {
				slog.Warn("[OnNewRegister] client did not ack welcome message", "userID", userID, "instanceID", conn.Auth().InstanceID)
			}
			h.mu.Lock()
			h.allUsers[userID] = struct{}{}
			h.mu.Unlock()
			return nil
		},
	)

	pw.Services().OnCloseRegister().RegisterAll(func(auth voAuth.WebsocketAuth) {
		if auth.UserID == "" {
			return
		}
		h.mu.Lock()
		delete(h.allUsers, auth.UserID)
		h.mu.Unlock()
	})

	return h
}

// HandleMessage handles incoming client messages and returns the response.
// The SDK automatically sets ResponseToId = request.Id on the outgoing frame.
func (h *Handler) HandleMessage(
	ctx context.Context,
	auth voAuth.WebsocketAuth,
	inputType string,
	data []byte,
) (outputType string, res []byte, err error) {
	switch inputType {

	case MsgTypeEcho:
		// Echo back the exact data with type ECHO_RESPONSE
		return MsgTypeEchoResponse, data, nil

	case MsgTypeFireAndForget:
		// No response; just log
		slog.Info("[FIRE_AND_FORGET]", "userID", auth.UserID, "payload", string(data))
		return "", nil, nil

	case MsgTypeErrorTest:
		// Respond with an error frame (same MsgType, same Id via SDK auto-set)
		return "", nil, errors.New("intentional error from server: " + string(data))

	case MsgTypeTestReqRes:
		msg := string(data)
		waitSec := min(len(msg), 6)
		time.Sleep(time.Duration(waitSec) * time.Second)
		return MsgTypeTestReqRes, []byte("Your request has been processed."), nil

	case MsgTypePing:
		// Respond with PONG, echo the data
		return MsgTypePong, data, nil

	default:
		slog.Warn("[handleMsg] unknown message type", "type", inputType)
		return "", nil, nil
	}
}

// StartBroadcasters starts the MEM_INFO (every 2s) and CAT_FACT (every 3s) server-push goroutines.
func (h *Handler) StartBroadcasters() {
	go h.memInfoBroadcaster()
	go h.catFactBroadcaster()
}

func (h *Handler) memInfoBroadcaster() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		h.broadcast(MsgTypeMemInfo, []byte(infoMemUsage()))
	}
}

func (h *Handler) catFactBroadcaster() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		payload := fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), fetchCatFact())
		h.broadcast(MsgTypeCatFact, []byte(payload))
	}
}

func fetchCatFact() string {
	type catFactResponse struct {
		Fact string `json:"fact"`
	}

	const fallbackFact = "Cats can rotate their ears 180 degrees."

	req, err := http.NewRequest(http.MethodGet, catFactAPIURL, nil)
	if err != nil {
		slog.Warn("[catFact] request build failed", "error", err)
		return fallbackFact
	}

	resp, err := catFactClient.Do(req)
	if err != nil {
		slog.Warn("[catFact] request failed", "error", err)
		return fallbackFact
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyPreview, _ := io.ReadAll(io.LimitReader(resp.Body, 256))
		slog.Warn("[catFact] non-200 response", "status", resp.StatusCode, "body", string(bodyPreview))
		return fallbackFact
	}

	var body catFactResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 4*1024)).Decode(&body); err != nil {
		slog.Warn("[catFact] decode failed", "error", err)
		return fallbackFact
	}

	fact := strings.TrimSpace(body.Fact)
	if fact == "" {
		slog.Warn("[catFact] empty fact response")
		return fallbackFact
	}

	return fact
}

// broadcast sends a message to all currently connected users.
func (h *Handler) broadcast(msgType string, data []byte) {
	h.mu.RLock()
	users := make([]string, 0, len(h.allUsers))
	for uid := range h.allUsers {
		users = append(users, uid)
	}
	h.mu.RUnlock()

	ctx := context.Background()
	for _, uid := range users {
		if aErr := h.pw.Services().SendToUser(ctx, uid, msgType, data); aErr != nil {
			slog.Warn("[broadcast] SendToUser failed", "userID", uid, "msgType", msgType, "error", aErr)
		}
	}
}

func infoMemUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each field, see: https://pkg.go.dev
	return fmt.Sprintf("Alloc = %v MiB \tTotalAlloc = %v MiB \tSys = %v MiB \tNumGC = %v", m.Alloc/1024/1024, m.TotalAlloc/1024/1024, m.Sys/1024/1024, m.NumGC)
}
