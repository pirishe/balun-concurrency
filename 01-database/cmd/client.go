package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const bufferSize = 4 << 10

func main() {
	loggerCfg := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths: []string{"raptor_client.log"},
	}
	logger, err := loggerCfg.Build()
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(os.Stdin)
	connection, err := net.Dial("tcp", "0.0.0.0:8091")
	if err != nil {
		log.Fatal(err)
	}

	for {
		fmt.Print("[raptor] > ")
		request, err := reader.ReadString('\n')
		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal("connection was closed", zap.Error(err))
		} else if err != nil {
			logger.Error("failed to read query", zap.Error(err))
		}

		response, err := Send(connection, []byte(request))
		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal("connection was closed", zap.Error(err))
		} else if err != nil {
			logger.Error("failed to send query", zap.Error(err))
		}

		fmt.Println("response: " + string(response))
	}
}

func Send(connection net.Conn, request []byte) ([]byte, error) {
	if _, err := connection.Write(request); err != nil {
		return nil, err
	}

	response := make([]byte, bufferSize)
	count, err := connection.Read(response)
	if err != nil && err != io.EOF {
		return nil, err
	} else if count == bufferSize {
		return nil, errors.New("small buffer size")
	}

	return response[:count], nil
}
