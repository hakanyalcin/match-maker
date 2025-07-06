# Matchmaking Service Case Study

## Overview

You are tasked with implementing a basic matchmaking HTTP service for a mobile game.
The service should queue players in FIFO (First-In-First-Out) order and form
matches when enough players are available.

## Key Requirements

1. Implement matchmaker algorithm for players waiting to be matched:
    - Each match should contain exactly **3** players.
    - Concurrent requests should be handled safely
2. Form matches when there are enough players in the queue.
3. Implement HTTP endpoints to interact with the matchmaking service.
4. Add necessary means for observability and monitoring, such as request/response logging and metrics.
5. Add error handling and validation.
6. Ensure the service has tests for the implemented functionality.

## Required HTTP Endpoints

### 1. POST /join

Adds a player to the matchmaking queue.

**Request Body:**
```json
{"id": "p1"}
```

**Success Response (HTTP 200):**

Example response:
```json
{
  "match_id": "<unique-match-id>",
  "status": "waiting",
  "players": ["p1"],
  "updated_at": 1745237600
}
```

- `match_id` is a unique identifier for the match.
- `status` indicates the current status of the match.
- `players` is a list of player IDs currently in the match.
- `updated_at` is a UNIX timestamp indicating when the match was last updated.
- `ready_at` is a UNIX timestamp indicating when the match was formed.

**Error Response (HTTP 500):**
```json
{"error": "Internal Server Error"}
```

### 2. GET /status/{match_id}

Returns the status of a match.

`status` field can be one of the following:
- `waiting`: Match is waiting for more players.
- `ready`: Match is ready to start.

**Success Response (HTTP 200):**
```json
{
  "match_id": "<unique-match-id>",
  "players": ["p1", "p2", "p3"],
  "status": "ready",
  "updated_at": 1745237600,
  "ready_at": 1745237600
}
```
**Error Response (HTTP 500):**
```json
{"error": "Internal Server Error"}
```

**Match Not Found (HTTP 404):**
```json
{"error": "match not found"}
```

## Bonus Features

- Implement a MySQL backend for storing player and match data.
- Implement a graceful shutdown mechanism to wait for inflight requests of the HTTP server.
