package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os/signal"
	"raptor/internal/database"
	"sync"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	loggerCfg := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths: []string{"raptor_server.log"},
	}
	logger, err := loggerCfg.Build()
	if err != nil {
		log.Fatal(err)
	}

	address := "0.0.0.0:8091"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	handleQueries(ctx, listener, logger)
}

func handleQueries(ctx context.Context, listener net.Listener, logger *zap.Logger) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			connection, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}

				logger.Error("failed to accept", zap.Error(err))
				continue
			}

			go func(connection net.Conn) {
				handleConnection(ctx, connection, logger)
			}(connection)
		}
	}()

	<-ctx.Done()
	listener.Close()

	wg.Wait() // don't wait for connections to complete
}

func handleConnection(ctx context.Context, connection net.Conn, logger *zap.Logger) {
	defer func() {
		if v := recover(); v != nil {
			logger.Error("captured panic", zap.Any("panic", v))
		}

		if err := connection.Close(); err != nil {
			logger.Warn("failed to close connection", zap.Error(err))
		}
	}()

	// reuse buffer for requests
	request := make([]byte, 4<<10)

	db := database.New()

	for {
		count, err := connection.Read(request)
		if err != nil && err != io.EOF {
			logger.Warn(
				"failed to read data",
				zap.String("address", connection.RemoteAddr().String()),
				zap.Error(err),
			)
			break
		}

		if response, err := db.HandleQuery(ctx, string(request[:count])); err != nil {
			logger.Warn("failed to handle query", zap.Error(err))
			break
		} else if _, err := connection.Write([]byte(response + "\n")); err != nil {
			logger.Warn(
				"failed to write data",
				zap.String("address", connection.RemoteAddr().String()),
				zap.Error(err),
			)
			break
		} else {
			logger.Warn("msg", zap.Error(fmt.Errorf("response sent: %s", response)))
		}
	}
}
