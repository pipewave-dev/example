package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	pipewave "github.com/pipewave-dev/go-pkg"
	ddbRepo "github.com/pipewave-dev/go-pkg/core/repository/impl-postgres"
	configprovider "github.com/pipewave-dev/go-pkg/provider/config-provider"
	queueprovider "github.com/pipewave-dev/go-pkg/provider/queue"
)

func main() {
	fmt.Println("hello")

	pw := pipewave.NewPipewave(pipewave.PipewaveConfig{
		ConfigStore:       getConfig(),
		RepositoryFactory: ddbRepo.NewPostgresRepo,
		QueueFactory:      queueprovider.QueueValkey,
		// SlogIns           -> Give your slog instance (default is slog.Default())
	})
	pw.SetFns(&pipewave.FunctionStore{
		InspectToken:      inspectToken,
		HandleMessage:     &handleMsg{i: pw},
		OnNewConnection:   nil,
		OnCloseConnection: nil,
	})

	// inject metrics handler

	server := &http.Server{
		Addr:    ":8080",
		Handler: http.StripPrefix("/pipewave", pw.Mux()),
	}

	go func() {
		fmt.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan
	fmt.Println("Shutting down server...")
	pw.Shutdown()
}

func getConfig() configprovider.ConfigStore {
	return configprovider.FromGoStruct(
		pipewave.ConfigEnv{
			Env:     "simple-chat",
			PodName: "pod-1",
			Version: "0.0.1",

			WorkerPool: configprovider.WorkerPoolT{
				Buffer:         100,
				UpperThreshold: 80,
				LowerThreshold: 20,
			},
			RateLimiter: configprovider.RateLimiterT{
				UserRate:       100,
				UserBurst:      200,
				AnonymousRate:  10,
				AnonymousBurst: 20,
			},
			TimeLocation:  time.UTC,
			TraceIDHeader: "X-Trace-ID",
			IpHeader:      "X-Real-IP",
			Cors: configprovider.CorsConfig{
				Enabled:        true,
				ExactlyOrigins: []string{},
				RegexOrigins:   []string{`^(https?://)?localhost:(\d+)/?$`},
			},
			Otel: configprovider.OtelT{
				Enabled: false,
			},
			Valkey: configprovider.ValkeyT{
				PrimaryAddress: "localhost:29100",
				ReplicaAddress: "localhost:29100",
				Password:       "veryStrongP@ssw0rd",
				DatabaseIdx:    0,
			},
			Postgres: configprovider.PostgresT{
				CreateTables: true,
				Host:         "localhost",
				Port:         29102,
				DBName:       "postgres",
				User:         "postgres",
				Password:     "postgres",
				SSLMode:      "disable",
				MaxConns:     15,
				MinConns:     1,
			},
		},
	)
}

func inspectToken(ctx context.Context, token string) (username string, IsAnonymous bool, err error) {
	// For demo purpose: token is username, but in real app, you should inspect token to get real userID
	return trimToken(token), false, nil
}

func trimToken(token string) string {
	token = strings.TrimSpace(token)
	token = strings.TrimPrefix(token, "Bearer ")
	return strings.TrimSpace(token)
}
