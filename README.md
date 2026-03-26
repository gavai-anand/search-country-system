# Country Search Service (Golang)

A simple backend service built in Go that fetches and returns country details such as name, population, currency, and capital using an external API.

---

## Features

* Search country by name
* Returns:
  * Country Name
  * Capital
  * Population
  * Currency (INR, USD, etc.)
* Clean architecture (Handler → Service → Client)
* Caching support
* Unit tests with mocks
* Graceful shutdown implemented

---

## 🏗️ Project Structure

```
internal/
  app/
    handlers/
    services/
      interfaces/
  bootstrap/
  dto/
  mocks/
  routes/
    v1/
main.go
```

---

##  Setup & Run

### 1. Clone the repo

```
git clone https://github.com/gavai-anand/search-country-system.git
cd search-country-system
```

---

### 2. Install dependencies

```
go mod tidy
```

---

### 3. Run the server

```
cp .env.development .env
go run main.go
```

Server runs on:

```
http://0.0.0.0:8080
```

---

## 📡 API Usage

### Get Country Details

```
GET /country?name=india
```

### Sample Response

```
{
  "name": "India",
  "capital": "New Delhi",
  "population": 1417492000,
  "currency": "INR"
}
```

---

## 🧪 Run Tests

```
go test ./...
```

With verbose:

```
go test -v ./...
```

--

## Graceful Shutdown

* Handles OS signals (`SIGINT`, `SIGTERM`)
* Ensures ongoing requests complete before shutdown

---

## Timeouts

* Read Timeout: 5 seconds
* Write Timeout: 30 seconds
* Also support custom timeout

---

## Tech Stack

* Golang
* net/http
* Testify (for testing)
* Mockery (for mocks)

---

## Notes
* Follows Go best practices for naming and structure

---

## Author

Anand Gavai
