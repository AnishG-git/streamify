"use client";
import { useParams } from "next/navigation";
import { useEffect, useRef, useState } from 'react';

// Define WebSocket message types
interface BaseMessage {
  type: string;
}

interface ErrorMessage extends BaseMessage {
  type: "error";
  error: string;
}

interface JoinMessage extends BaseMessage {
  type: "join";
  name: string;
}

interface LeaveMessage extends BaseMessage {
  type: "leave";
  name: string;
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

type WebSocketMessage = ErrorMessage | JoinMessage | LeaveMessage | ICECandidateMessage | OfferMessage | AnswerMessage;

const Room = () => {
  const { roomCode } = useParams();
  const name: string | null = sessionStorage.getItem("name");
  if (!name) {
    window.location.href = "http://localhost:3000/home";
  }
  const socketRef = useRef<WebSocket | null>(null);
  const [participants, setParticipants] = useState<string[]>([]);

  useEffect(() => {
    if (!socketRef.current) {
      if (name === null) {
        alert("Failed to retrieve name from session storage");
        window.location.href = "http://localhost:3000/home";
        return
      }
      socketRef.current = new WebSocket(`ws://localhost:8080/room/connect/${roomCode}?name=${name}`);
    }

    socketRef.current.onopen = () => {
      const state = socketRef.current?.readyState;
      console.log(`WebSocket opened. Current state: ${state}`);
      socketRef.current?.send(JSON.stringify({ type: "join", name: name }));
    };

    socketRef.current.onmessage = (event) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data);
        switch (message.type) {
          case "error":
            console.log("Received error from server:", message.error);
            closeSocket(1000, message.error);
            window.location.href = "http://localhost:3000/home";
            break;
          case "join":
            handleJoin(message.name);
            break;
          case "leave":
            handleLeave(message.name);
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
      closeSocket();
    };
  }, [roomCode]);

  const handleJoin = (name: string) => {
    console.log(name + " joined the room");
    if (!participants.includes(name)) {
      setParticipants((prevParticipants) => [...prevParticipants, name]);
    }
  };

  const handleLeave = (name: string) => {
    console.log(name + " left the room");
    setParticipants((prevParticipants) => prevParticipants.filter((participant) => participant !== name));
  };

  const handleICECandidate = (candidate: RTCIceCandidate) => {
    console.log("Received ICE Candidate:", candidate);
  };

  const handleOffer = (offer: RTCSessionDescriptionInit) => {
    console.log("Received Offer:", offer);
  };

  const handleAnswer = (answer: RTCSessionDescriptionInit) => {
    console.log("Received Answer:", answer);
  };

  const handleClientLeave = () => {
    closeSocket(1000, "client disconnected");
    window.location.href = "http://localhost:3000/home";
  };

  const closeSocket = (code?: number, reason?: string) => {
    if (socketRef.current) {
      console.log(`Closing ${socketRef.current.url}, current state: ${socketRef.current.readyState}`);
      if (socketRef.current.readyState === WebSocket.OPEN || socketRef.current.readyState === WebSocket.CONNECTING) {
        socketRef.current.close(code, reason);
      }
      socketRef.current = null;
    }
  };

  return (
    <div className="flex flex-col items-center justify-center h-screen bg-gray-200">
      <h1 className="text-3xl font-bold">Room Code: {roomCode}</h1>
      <p className="mt-4 text-lg">This is where the screen sharing will happen!</p>
      <div className="mt-6">
        <h2 className="text-2xl font-semibold">Participants:</h2>
        <ul className="mt-2">
          {participants.map((participant, index) => (
            <li key={index} className="text-lg">{participant}</li>
          ))}
        </ul>
      </div>
      <button 
        onClick={handleClientLeave} 
        className="mt-6 px-4 py-2 bg-red-500 text-white rounded hover:bg-red-700"
      >
        Leave Room
      </button>
    </div>
  );
};

export default Room;
