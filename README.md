# Streamify

Streamify is a real-time screen-sharing and collaboration application that allows users to create or join rooms using unique codes. Users can stream their screens, including video and audio, with a focus on high-quality streaming for optimal collaboration. Currently, each room supports up to two participants.

## Features

- **Room Creation**: Generate a unique, 5-character alphanumeric room code.
- **Join Room**: Enter a room code to join an existing session.
- **High-Quality Streaming**: Support for 1080p and potentially 1440p video resolution with frame rates up to 120 fps.
- **Screen Sharing with Audio**: Stream both video and audio to other participants.

## Technologies Used

### Frontend
- **Next.js** with TypeScript: For building a responsive and dynamic user interface.
- **Tailwind CSS**: For styling and responsive design.

### Backend
- **Go**: For handling server-side logic and WebSocket connections.
- **PostgreSQL**: For managing active room codes and participant tracking.

### Other Tools
- **Docker**: For containerized development and deployment.
- **WebSockets**: For real-time communication between participants.

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/username/streamify.git
    cd streamify
    ```
2. Install dependencies:
    - For the backend:
      ```bash
      make build
      ```
    - For the frontend:
      ```bash
      npm install
      ```
3. Run the backend server:
    ```bash
    make run
    ```
4. Start the frontend development server:
    ```bash
    npm run dev
    ```

## Usage

- **Create Room**: Navigate to `/home` and click "Create Room" to generate a new room code.
- **Join Room**: Enter a valid room code and click "Join Room" to enter an existing session.
- **Streaming**: Share your screen and audio with another participant.

## Project Structure

```
streamify/
├── backend/
│   ├── main.go
│   ├── handlers/
│   ├── models/
│   └── utils/
├── frontend/
│   ├── src/
│   │   ├── app/
│   │   │   ├── home/
│   │   │   │   └── page.tsx
│   │   │   └── room/
│   │   │       └── [roomCode]/page.tsx
│   └── tailwind.config.ts
└── docker-compose.yml
```

## Roadmap

- Implement advanced session management for more than two participants.
- Add authentication and user profiles.
- Optimize streaming performance for different network conditions.
- Add mobile device support.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any improvements or new features.

## License

This project is licensed under the MIT License.