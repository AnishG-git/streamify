"use client";
import { randomUUID } from "crypto";
import { useParams } from "next/navigation";
import { useEffect, useRef } from 'react';

// Define WebSocket message types
interface BaseMessage {
  type: string;
}

interface ErrorMessage extends BaseMessage {
  type: "error";
  error: string;
}

interface ICECandidateMessage extends BaseMessage {
  type: "ice-candidate";
  candidate: RTCIceCandidate;
}

interface OfferMessage extends BaseMessage {
  type: "offer";
  offer: RTCSessionDescriptionInit;
}

interface AnswerMessage extends BaseMessage {
  type: "answer";
  answer: RTCSessionDescriptionInit;
}

type WebSocketMessage = ErrorMessage | ICECandidateMessage | OfferMessage | AnswerMessage;

const Room = () => {
  const { roomCode } = useParams();
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
      if (!socketRef.current) {
        socketRef.current = new WebSocket(`ws://localhost:8080/room/connect/${roomCode}`);
      }
      
      socketRef.current.onopen = () => {
        const state = socketRef.current?.readyState;
        console.log(`WebSocket opened. Current state: ${state}`);
        socketRef.current?.send(JSON.stringify({ type: "join" }));
      };

      socketRef.current.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data);
          switch (message.type) {
            case "error":
              console.error("Received error message:", message.error);
              socketRef.current?.close();
              window.location.href = "http://localhost:3000/home";
              break;
            case "ice-candidate":
              handleICECandidate(message.candidate);
              break;
            case "offer":
              handleOffer(message.offer);
              break;
            case "answer":
              handleAnswer(message.answer);
              break;
            default:
              console.warn("Unknown message type:", message);
          }
        } catch (err) {
          console.error("Failed to parse WebSocket message:", err);
        }
      };

      socketRef.current.onclose = () => {
        console.log(`Disconnected from room ${roomCode}`);
      };

    return () => {
      if (socketRef.current) {
        console.log(`Closing ${socketRef.current.url} with state ${socketRef.current.readyState}`);
        if (socketRef.current.readyState === WebSocket.OPEN || socketRef.current.readyState === WebSocket.CONNECTING) {
          socketRef.current.close();
        }        
        socketRef.current = null;
      }
    };
  }, [roomCode]);

  const handleICECandidate = (candidate: RTCIceCandidate) => {
    console.log("Received ICE Candidate:", candidate);
  };

  const handleOffer = (offer: RTCSessionDescriptionInit) => {
    console.log("Received Offer:", offer);
  };

  const handleAnswer = (answer: RTCSessionDescriptionInit) => {
    console.log("Received Answer:", answer);
  };

  return (
    <div className="flex flex-col items-center justify-center h-screen bg-gray-200">
      <h1 className="text-3xl font-bold">Room Code: {roomCode}</h1>
      <p className="mt-4 text-lg">This is where the screen sharing will happen!</p>
    </div>
  );
};

export default Room;
