# nuitee Integration Service

This project provides a simple API service to interact with the Hotelbeds API, transforming hotel data into a standardized response format. It also includes *currency conversion using exchange rates* from an external service.

## Features

- Fetch hotel rates from the Hotelbeds API
- Convert prices to a specified currency using real-time exchange rates
- Expose an API endpoint to retrieve hotel rates and details

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Go 1.16 or later installed
- Docker (optional, for running the service in a container)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/chiranshu11/nuitee.git
   cd liteapi
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

## Configuration

The service requires API keys for the Hotelbeds API and the exchange rate API. These are hardcoded in the source code for simplicity, but we should use environment variables or a configuration file for sensitive information in a production environment.

```go
const (
    apiKey    = "your_hotelbeds_api_key"
    apiSecret = "your_hotelbeds_api_secret"
    apiUrl    = "https://api.test.hotelbeds.com/hotel-api/1.0"
    ratesAPI  = "https://api.exchangerate-api.com/v4/latest/"
)
```

## Running the Service

To run the service locally:

```bash
go run main.go
```

The service will start on port 9001. You can access the API endpoint at `http://localhost:9001/hotels`.

## Usage

### API Endpoint

#### GET /hotels

Retrieve hotel rates and details.

**Query Parameters:**

- `checkin` (string, required): Check-in date (YYYY-MM-DD)
- `checkout` (string, required): Check-out date (YYYY-MM-DD)
- `currency` (string, required): Target currency code (e.g., "USD")
- `guestNationality` (string, required): Guest nationality code (e.g., "US")
- `hotelIds` (string, required): Comma-separated list of hotel IDs
- `occupancies` (string, required): JSON string of occupancies (e.g., `[{"adults":2,"children":0,"rooms":1}]`)

**Example Request:**

```bash
curl "http://localhost:9001/hotels?checkin=2024-08-01&checkout=2024-08-10&currency=USD&guestNationality=US&hotelIds=1234,5678&occupancies=[{\"adults\":2,\"children\":0,\"rooms\":1}]"
```

**Example Response:**

```json
{
  "data": [
    {
      "hotelId": "1234",
      "currency": "USD",
      "price": 150.00
    },
    {
      "hotelId": "5678",
      "currency": "USD",
      "price": 200.00
    }
  ],
  "supplier": {
    "request": "{\"stay\":{\"checkIn\":\"2024-08-01\",\"checkOut\":\"2024-08-10\"},\"hotels\":{\"hotel\":[\"1234\",\"5678\"]},\"occupancies\":[{\"adults\":2,\"children\":0,\"rooms\":1}]}",
    "response": "{\"hotels\":{\"hotels\":[{\"code\":1234,\"name\":\"Hotel One\",\"minRate\":\"150\",\"currency\":\"USD\"},{\"code\":5678,\"name\":\"Hotel Two\",\"minRate\":\"200\",\"currency\":\"USD\"}]}}"
  }
}
```

## Docker Support

To build and run the service using Docker:

1. Build the Docker image:

   ```bash
   docker build -t liteapi .
   ```

2. Run the Docker container:

   ```bash
   docker run -p 9001:9001 liteapi
   ```

