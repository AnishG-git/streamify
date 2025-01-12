"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";
import { TypographyH1 } from "@/components/ui/typography/h1";
import { TypographyP } from "@/components/ui/typography/p";

const Home = () => {
  const router = useRouter();
  const [joinCode, setJoinCode] = useState<string>("");
  const [name, setName] = useState<string>("");

  const handleCreateRoom = async () => {
    if (!name.trim()) {
      alert("Name cannot be empty");
      return;
    }

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
    if (!name.trim()) {
      alert("Name cannot be empty");
      return;
    }

    if (joinCode.length === 5) {
      router.push(`/room/${joinCode}`);
    } else {
      alert("Please enter a valid 5-character room code.");
    }
  };

  return (
    <div className="flex flex-col items-center justify-center h-screen space-y-6 bg-gray-100">
      <TypographyH1>Streamify</TypographyH1>
      <Carousel className="w-full max-w-md">
        <CarouselContent>
          <CarouselItem>
            <div className="p-1">
              <Card className="w-full h-42"> {/* Adjusted height */}
                <CardContent className="flex flex-col items-start justify-center p-6 space-y-1 h-48"> {/* Set fixed height */}
                  <TypographyP className="font-bold">Display Name</TypographyP>
                  <Input
                  type="text"
                  placeholder="Enter Your Name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  />
                </CardContent>
              </Card>
            </div>
          </CarouselItem>
          <CarouselItem>
            <div className="p-1">
              <Card className="w-full h-42"> {/* Adjusted height */}
                <CardContent className="flex flex-col items-center justify-center p-6 space-y-4 h-48"> {/* Set fixed height */}
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
                  <TypographyP className="text-center font-semibold">
                    or
                  </TypographyP>
                  <Button onClick={handleCreateRoom}>
                    Create Room
                  </Button>
                </CardContent>
              </Card>
            </div>
          </CarouselItem>
        </CarouselContent>
        <CarouselPrevious />
        <CarouselNext />
      </Carousel>
    </div>
  );
};

export default Home;
