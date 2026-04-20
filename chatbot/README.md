# Chatbot UI

A minimal React chat interface for demoing the Smart Home Cat Demo system end-to-end. It acts as a simple command interpreter that calls the API Gateway to exercise the backend microservices.

## Supported Commands

| Command | Action |
|---------|--------|
| `list cats` | `GET /api/cats` — lists all cat profiles |
| `list devices` | `GET /api/devices` — lists all registered devices |
| `add cat` | `POST /api/cats` — creates a sample cat profile |
| `add device` | `POST /api/devices` — registers a sample device |

## Running Locally

```bash
npm install
npm run dev
```

The dev server starts on [http://localhost:3000](http://localhost:3000).

Set the API Gateway URL via environment variable:

```bash
VITE_API_URL=http://localhost:8080 npm run dev
```

## Running with Docker

```bash
docker build -t chatbot-ui .
docker run -p 3000:3000 chatbot-ui
```

Or as part of the full stack:

```bash
docker-compose up chatbot-ui
```

## Tech Stack

- Vite + React 18 + TypeScript
- No external UI libraries — plain CSS
- Nginx for production serving
