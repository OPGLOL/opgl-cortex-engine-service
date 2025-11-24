# OPGL Cortex Engine

AI-powered performance analysis microservice for League of Legends players.

## Features

- Player performance statistics calculation
- Improvement area identification
- Personalized recommendations based on performance metrics

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Service health check |
| `/api/v1/analyze` | POST | Analyze player performance |

## Analyze Endpoint

**POST** `/api/v1/analyze`

**Request Body**:
```json
{
  "summoner": {
    "puuid": "string",
    "name": "string",
    "summonerLevel": 123
  },
  "matches": [
    {
      "matchId": "string",
      "participants": [...]
    }
  ]
}
```

**Response**:
```json
{
  "playerStats": {
    "puuid": "string",
    "summonerName": "string",
    "totalMatches": 20,
    "winRate": 55.5,
    "averageKDA": 3.5,
    "averageCS": 150.0,
    "csPerMinute": 6.5,
    "averageVisionScore": 45.0
  },
  "improvementAreas": [
    {
      "category": "CS (Creep Score)",
      "currentValue": 5.5,
      "expectedValue": 6.0,
      "gap": -0.5,
      "priority": "HIGH",
      "recommendation": "Focus on last-hitting minions..."
    }
  ],
  "analyzedAt": "2024-11-23T18:00:00Z"
}
```

## Setup

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Run the service**:
   ```bash
   go run main.go
   ```

Service runs on port **8082** by default.

## Environment Variables

- `PORT` - Service port (default: 8082)

## Testing

Use Bruno collection at `bruno-collections/opgl/opgl-model/` to test endpoints.
