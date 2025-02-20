package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/AnishG-git/streamify/internal/connections"
	"github.com/gorilla/websocket"
)

func GenerateRoomLogic(ctx context.Context, logger *log.Logger, manager connections.ConnManager) (string, error) {
	var roomCode string
	for {
		roomCode = generateRoomCode()
		exists, err := manager.IsRoomActive(ctx, roomCode)
		if err != nil {
			logger.Printf("Failed to check room code existence: %v", err)
			return "", err
		}
		if !exists {
			break
		}
	}

	err := manager.CreateRoom(ctx, roomCode)
	if err != nil {
		logger.Printf("Failed to set room code: %v", err)
		return "", err
	}

	return roomCode, nil
}

func ConnectToRoomLogic(ctx context.Context, logger *log.Logger, manager connections.ConnManager, roomCode string, name string, conn *websocket.Conn) (string, error) {
	var errMsg string
	if err := manager.CanUserJoinRoom(ctx, roomCode, name); err != nil {
		errMsg = "user cannot join room at this time"
		err = fmt.Errorf("user cannot join room: %w", err)
		return errMsg, err
	}

	// checks have passed, adding connection to room
	connDetails := manager.SetConnection(conn)

	marshalledConnDetails, err := json.Marshal(connDetails)
	if err != nil {
		errMsg = "Internal Server Error"
		err = fmt.Errorf("Failed to marshal connection object: %w", err)
		return errMsg, err
	}

	err = manager.AddUserToRoom(ctx, roomCode, name, string(marshalledConnDetails))
	if err != nil {
		errMsg = "Failed to add connection to room"
		err = fmt.Errorf("Failed to add connection to room: %w", err)
		return errMsg, err
	}

	logger.Printf("User %s has joined room %s", name, roomCode)
	ctxWithoutCancel := context.WithoutCancel(ctx)
	for {
		var message map[string]interface{}
		err := conn.ReadJSON(&message)
		if err != nil {
			// Handle normal WebSocket closure without logging an error
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				logger.Printf("Unexpected WebSocket close error for user %s in room %s: %v", name, roomCode, err)
			} else {
				logger.Printf("WebSocket closed for user %s in room %s: %v", name, roomCode, err)
			}

			// Remove connection from room
			go manager.RemoveConnectionFromRoom(ctxWithoutCancel, logger, roomCode, name)
			break
		}

		faultyReceiverName, err := manager.BroadcastToRoom(ctx, logger, roomCode, name, message)
		if err != nil {
			logger.Printf("Failed to send message to room %s: %v", roomCode, err)
			if faultyReceiverName != "" {
				go manager.RemoveConnectionFromRoom(ctxWithoutCancel, logger, roomCode, name) // Remove faulty connection
			}
		} else {
			logger.Printf("Message from %s: %v", roomCode, message)
		}
	}
	return "", nil
}
