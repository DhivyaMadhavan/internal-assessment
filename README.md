# Full-Stack Engineering Intern Assessment
## API Multiplexer with Custom Circuit Breaker, Real-Time Dashboard & Chaos Testing

---

## Overview

This project implements a resilient API Multiplexer built using **Go**, **React**, **Docker**, and **WebSockets**.

The application acts as a reverse proxy that routes incoming client requests to a **Primary API**. When the Primary API becomes unavailable or exceeds a **200 ms timeout**, a custom-built **Circuit Breaker** automatically redirects requests to a **Secondary API** without dropping the client request.

The frontend provides a real-time monitoring dashboard displaying circuit breaker state transitions, request statistics, routing information, and telemetry streamed over WebSockets.

---

# Architecture

```
                Browser / Postman
                        │
                        ▼
                  Go Backend
                        │
                 helloHandler()
                        │
                 proxyRequest()
                        │
                Circuit Breaker
                 │            │
          Allowed        Not Allowed
             │                │
             ▼                ▼
        Toxiproxy       Secondary API
             │
             ▼
        Primary API
             │
     Success / Failure
             │
             ▼
        Update Metrics
             │
             ▼
        WebSocket
             │
             ▼
     Frontend Dashboard

```

---

# Features

## Backend (Go)

- Custom HTTP Reverse Proxy using Go standard library (`net/http`)
- Custom Circuit Breaker implementation (no external resilience libraries)
- Thread-safe state management using Mutex
- Automatic failover to Secondary API
- 200 ms request timeout using Context
- WebSocket telemetry
- Real-time metrics broadcasting
- Dockerized using Multi-stage builds

---

## Circuit Breaker States

The implementation supports all three circuit breaker states.

### CLOSED

- Requests are routed to the Primary API.
- Failure count is monitored.

---

### OPEN

- Triggered after **3 consecutive failures**.
- Primary API is skipped.
- Requests are immediately routed to the Secondary API.

---

### HALF_OPEN

After a 5-second recovery timeout:

- One trial request is sent to the Primary API.
- If successful:
  - Circuit transitions to **CLOSED**
- If unsuccessful:
  - Circuit transitions back to **OPEN**

---

## Frontend (React)

Real-time monitoring dashboard displaying:

- Total Requests
- Active Route
- Circuit Breaker State
- Request Graph
- Transition Log
- Live WebSocket Updates

---

## Telemetry

The backend continuously streams:

- Request Count
- Active Route
- Circuit Breaker State

using WebSockets.

---



## 🧠 Architectural Choices

### 1. Custom Circuit Breaker State Machine (Golang)
*   **Zero-Dependency Design**: Built entirely from scratch using `sync.Mutex` blocks to enforce total thread safety. This prevents race conditions and ensures that state transitions (`CLOSED` ↔ `OPEN` ↔ `HALF-OPEN`) resolve atomically even under high concurrency.
*   **Context Deadline Enforcement**: A strict `context.WithTimeout(200 * time.Millisecond)` boundary is wrapped around all outbound Primary requests. If the network or upstream drops packets or introduces latency, the runtime instantly cancels the execution loop to avoid cascading component exhaustion.
*   **Time-Based Cooldown Lifecycle**: When consecutive failures cross the threshold of **3**, the breaker transitions to `OPEN` and routes traffic directly to the Secondary API without wasting network bandwidth. Once a time-based cooling window expires, the state machine steps into `HALF-OPEN` on the very next incoming request to issue a single trial probe. A successful result smoothly recovers the circuit back to `CLOSED`.

### 2. High-Frequency Rendering Architecture (React)
*   **Optimized Rendering Pipelines**: The frontend UI leverages high-throughput WebSocket listeners designed to parse real-time server telemetries (RPS, state metrics, active routes) without triggering costly full-component tree re-renders or freezing the browser main execution thread.

### 3. Memory Constraint Engineering
*   **128 MB RAM Envelope**: The Go backend is constrained to a strict `128M` memory layout using Docker kernel `cgroups` limits. To review active runtime stats, execute: `docker stats backend --no-stream`.
*   **Deterministic Garbage Prevention**: Resource leaks are completely avoided by using `defer response.Body.Close()` across all endpoints, ensuring underlying system file descriptors are returned to the platform allocator immediately. Telemetry variables are tracked using pre-allocated structures rather than dynamically appending slices, maintaining a flat memory profile over prolonged stress testing.

### 4. Multi-Stage Docker Builds
*   The Go build matrix uses `CGO_ENABLED=0 GOOS=linux` during compilation. This decouples the resulting application from host glibc binaries, producing an ultra-lean, statically linked native execution binary tailored to execute inside minimal `alpine` distribution envelopes safely.

---

# Chaos Engineering

The project integrates **Toxiproxy** to simulate network failures.

Supported scenarios:

- 500 ms latency
- Packet loss
- Automatic failover
- Recovery verification

---

# Project Structure

```
📁 Project Root
├── .dockerignore
├── go.mod
├── go.sum
├── docker-compose.yml
│
├── backend
│   ├── Dockerfile
│   ├── main.go
│   ├── proxy.go
│   ├── handler.go
│   ├── telemetry_helper.go
│   ├── ws.go
│   ├── circuitbreaker
│   │     └── breaker.go
│   └── telemetry
│         └── metrics.go
│
├── frontend
│   ├── Dockerfile
│   ├── package.json
│   ├── package-lock.json
│   ├── App.js
│   ├── App.css
│   ├── index.js
│   ├── index.css
│   ├── reportWebVitals.js
│   ├── setupTests.js
│   └── src/components
│
├── primary-api
│   ├── Dockerfile
│   └── main.go
│
└── secondary-api
    ├── Dockerfile
    └── main.go
```

---

# Technologies Used

| Technology | Purpose |
|------------|---------|
| Go | Backend |
| React | Frontend Dashboard |
| Gorilla WebSocket | Real-time Telemetry |
| Docker | Containerization |
| Docker Compose | Multi-container orchestration |
| Toxiproxy | Chaos Engineering |
| Context API | Request Timeout |
| Mutex | Thread-safe Circuit Breaker |

---

# Running the Project

## 1. Clone Repository

```bash
git clone <repository-url>
cd <repository-name>
```

---

## 2. Build Containers

```bash
docker compose build --no-cache
```

---

## 3. Start Containers

```bash
docker compose up -d
```

---

## 4. Verify Running Containers

```bash
docker compose ps
```

Expected containers:

- backend
- frontend
- primary-api
- secondary-api
- toxiproxy

---

# Application URLs

Frontend

```
http://localhost:3000
```

Backend

```
http://localhost:8080/hello
```

Primary API

```
http://localhost:8081/hello
```

Secondary API

```
http://localhost:8082/hello
```

---

# Creating the Toxiproxy Proxy

Run once after starting the containers.

PowerShell:

```powershell
Invoke-RestMethod -Method Post `
-Uri "http://localhost:8474/proxies" `
-ContentType "application/json" `
-Headers @{"User-Agent"="PowerShell"} `
-Body '{"name":"primary","listen":"0.0.0.0:8666","upstream":"primary-api:8081","enabled":true}'
```

---

# Inject 500 ms Latency

```powershell
Invoke-RestMethod -Method Post `
-Uri "http://localhost:8474/proxies/primary/toxics" `
-Headers @{"User-Agent"="PowerShell"} `
-ContentType "application/json" `
-Body '{"name":"latency","type":"latency","attributes":{"latency":500}}'
```

---

# Remove Latency

```powershell
Invoke-RestMethod -Method Delete `
-Uri "http://localhost:8474/proxies/primary/toxics/latency" `
-Headers @{"User-Agent"="PowerShell"}
```

---

# Packet Loss (20%)

```powershell
 curl.exe -X POST http://localhost:8474/proxies/primary/toxics -H "Content-Type: application/json" -d '{\"name\":\"packet_drop_chaos\",\"type\":\"limit_data\",\"stream\":\"upstream\",\"toxicity\":0.20,\"attributes\":{\"bytes\":0}}'

```

---

# Remove Packet Loss

```powershell
 Invoke-RestMethod -Method Delete -Uri "http://localhost:8474/proxies/primary/toxics/packet_drop_chaos" -Headers @{"User-Agent"="PowerShell"}

```

---

# Expected Behaviour

### Normal

```
Primary API
Circuit : CLOSED
Route : Primary
```

---

### After Latency Injection

```
Primary Timeout

↓

Circuit OPEN

↓

Traffic routed to Secondary API
```

---

### Recovery

```
OPEN

↓

HALF_OPEN

↓

CLOSED
```

---

# Memory Optimization

- Multi-stage Docker builds
- Minimal runtime images
- Thread-safe Circuit Breaker
- Proper HTTP response body cleanup
- Context cancellation after timeout
- Limited frontend history buffers
- Efficient WebSocket broadcasting
- Backend container configured for memory-constrained execution

---

# Author

**Dhivya**

Engineering Intern Assessment Submission
