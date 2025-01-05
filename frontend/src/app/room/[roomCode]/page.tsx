"use client";
import { useParams } from "next/navigation";

const Room = () => {
  const { roomCode } = useParams();

  return (
    <div className="flex flex-col items-center justify-center h-screen bg-gray-200">
      <h1 className="text-3xl font-bold">Room Code: {roomCode}</h1>
      <p className="mt-4 text-lg">This is where the screen sharing will happen!</p>
    </div>
  );
};

export default Room;
