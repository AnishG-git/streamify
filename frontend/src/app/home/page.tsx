"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { TypographySmall } from "@/components/small";

const Home = () => {
  const router = useRouter();
  const [joinCode, setJoinCode] = useState<string>("");
  const [name, setName] = useState<string>("");

  const handleCreateRoom = async () => {
    const response = await fetch("http://localhost:8080/room/generate", {
      method: "GET",
    });

    if (!response.ok) {
      throw new Error("Failed to generate room code");
    }

    interface RoomCode {
      code: string;
    }

    const roomCode: RoomCode = await response.json();
  
    if (roomCode.code.length !== 5) {
      alert("Failed to generate room");
      return;
    }

    try {
      router.push(`/room/${roomCode.code}`);
    } catch {
      alert("Failed to generate room");
    }
  };

  const handleJoinRoom = () => {
    if (joinCode.length === 5) {
      router.push(`/room/${joinCode}`);
    } else {
      alert("Please enter a valid 5-character room code.");
    }
  };

  return (
    <div className="flex flex-col items-center justify-center h-screen space-y-6 bg-gray-100">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-4xl font-bold text-center">Streamify</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <Input
            type="text"
            placeholder="Enter Your Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
          <div className="flex space-x-4">
            <Input
              type="text"
              placeholder="Enter Room Code"
              value={joinCode}
              onChange={(e) => setJoinCode(e.target.value.toUpperCase())}
            />
            <Button onClick={handleJoinRoom}>
              Join Room
            </Button>
          </div>
          <div className="flex justify-center items-center">
            <TypographySmall className="text-center">
              or
            </TypographySmall>
          </div>
          <div className="flex justify-center">
            <Button onClick={handleCreateRoom}>
              Create Room
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default Home;
