# Anagram Service (Go)

A **multilingual Anagram Service** written in Go that provides:

* A REST API for retrieving anagrams
* A concurrent client for load testing
* Unicode-aware anagram matching (supports languages like French)
* Configurable **LRU cache**
* Config-driven server behavior
* Unit tests for core logic, server, and client
* Prometheus metrics collection for request latency and counts
* Auth middleware that allows request origin from loopback interface only
* Docker based deployment (optional) 
* Makefile setup with starting client/server, running tests and docker management

---
# Quick Start

```bash
# install dependencies
make deps

# build all the executables (server, client and docker container)
make build

# start the server as service
make run-server 

# start the server in docker (optional)
make docker-run

# test the api - single request
curl "http://localhost:8080/v1/anagrams?word=trace"

# test the api - concurrent client
make run-client

# run the tests
make test
```
---

# Configuration

All the configuration options saved to config.yaml. 

* Add a new language dictionary to ```dictionary.files``` slice

* Modify cache size by updating ``lru_cache.capacity``

* Enable debug logs by setting ``log.debug`` to true

---

# Architecture Overview

```
                +------------------+
                |   config.yaml    |
                +------------------+
                         |
                         v
                +------------------+
                | Config Loader    |
                +------------------+
                         |
         +---------------+---------------+
         |                               |
         v                               v
 +---------------+               +---------------+
 | Dictionary    |               | LRU Cache     |
 | Loader        |               | (Hot Queries) |
 +---------------+               +---------------+
         |                               |
         +---------------+               |
                         |               |
                                         |
                 +---------------+       |
                 | Anagram Finder|       |
                 +---------------+       |
                         |               |
                         v               |
                 +---------------+       |
                 | REST Server   |       |
                 | /v1/anagrams  |--------
                 |               |
                 | /metrics      |
                 +---------------+
                         |
                         v
                 +----------------+
                 | Auth Middleware|
                 +----------------+
                         |
                         v
                 +-------------------+
                 | Metrics Middleware|
                 +-------------------+
                         |
                         v
                 +---------------+
                 | Client Request|
                 +---------------+

```

Flow for LRU Cache

```
Client Request
      |
      v
REST Server 
      |
      v
  LRU Cache
  |      |
 HIT    MISS
  |      |
Return  Anagram Service
         |
         v
       Cache
```


---

# Project Structure

```
anagrams/
│
├── cmd/
│   ├── server/           # REST API server
│   │   └── main.go
│   │   └── main_test.go
│   │
│   └── client/           # Concurrent load test client
│       └── main.go
│
├── pkg/
│   ├── service/          # Core services
│   │   ├── anagrams.go
│   │   └── anagrams_test.go
│   │
│   ├── cache/            # LRU cache implementation
│   │   ├── lru.go
│   │   └── lru_test.go
│   │
│   |── config/           # Config loader
│   |    └── config.go
│   |    └── config_test.go
│   │
│   |── middleware/           # Middlewares 
│   |    └── auth.go
│   |    └── auth_test.go
│   |    └── metrics.go
│   |    └── metrics_test.go
│
├── data/
│   └── english.txt
│   └── french.txt
│
├── config.yaml
├── go.mod
├── Makefile.                # Build various services, run tests, etc
└── README.md
```

---

## REST API

**Endpoint** : ```GET /v1/anagrams?word=<word>```

```bash

# Example request:
curl http://localhost:8080/v1/anagrams?word=trace

# Example response:

{
  "word": "trace",
  "anagrams": ["trace", "crate", "react"]
}
```

**Endpoint** : ```GET /metrics```

Returns the Prometheus HTTP metrics for endpoints

---

# Future Enhancements and Production Considerations

* **Distributed caching** – Introduce Redis or Memcached to share cached anagram results across multiple service instances and improve performance in horizontally scaled deployments.

* **Optimized anagram encoding** – Replace the current sorted-character key generation (O(N log N)) with a character frequency encoding (O(N)) to reduce preprocessing time for large dictionaries (1M+ entries).

* **Horizontal scaling** – Containerize the service with Docker and run multiple instances behind a load balancer for improved scalability and fault tolerance.

* **Regional language support** – Deploy region-specific dictionaries (e.g., eu, na, apac) to support multiple languages efficiently while retaining shared dictionaries like English.

* **Security enhancements** – Extend authentication middleware to support API keys or JWT-based authorization.

* **Rate limiting** – Implement request throttling to protect the service from abuse or potential DDoS attacks.

* **Enhanced observability** – Expand metrics collection and integrate with Prometheus/Grafana dashboards and alerting for improved operational visibility.

