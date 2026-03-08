# Anagram Service (Go)

A simple **multilingual Anagram Service** written in Go that provides:

* A REST API for retrieving anagrams
* A concurrent client for load testing
* Unicode-aware anagram matching (supports languages like French)
* Configurable **LRU cache**
* Config-driven server and client behavior
* Unit tests for core logic, server, and client

The system is designed to be **simple, extensible, and interview-friendly** while demonstrating common Go service patterns.

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
         +---------------+---------------+
                         |
                         v
                 +---------------+
                 | Anagram Finder|
                 +---------------+
                         |
                         v
                 +---------------+
                 | REST Server   |
                 | /v1/anagrams  |
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

---

# Project Structure

```
anagrams/
│
├── cmd/
│   ├── server/           # REST API server
│   │   └── main.go
│   │   └── main_teset.go
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

# Features

## Unicode / Multilingual Support

The service supports **Unicode words** and can work with multiple locales such as:

* English
* French


Unicode normalization ensures words like:

```
écart
trace
crate
```

match correctly.

---

## Anagram Lookup

Anagrams are computed using:

```
sorted runes → key → lookup map
```

Example:

```
trace → acert
crate → acert
```

Lookup complexity:

```
O(k log k) key generation
O(1) lookup
```

---

## LRU Cache for Hot Queries

To improve performance for frequently requested words, the service uses an **LRU cache**.

Example request flow:

```
Client Request
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

Benefits:

* reduces repeated computation
* improves latency
* improves throughput under load

---

## REST API

Endpoint:

```
GET /v1/anagrams?word=<word>
```

Example request:

```
curl http://localhost:8080/v1/anagrams?word=trace
```

Example response:

```
{
  "word": "trace",
  "anagrams": ["trace", "crate", "react"]
}
```

---

## Concurrent Client

A standalone client simulates load by sending multiple concurrent requests.

Configurable parameters:

```
concurrent_requests
total_requests
timeout
```

Example:

```
50 concurrent goroutines
1000 requests total
```

---

# Configuration

All runtime behavior is controlled by `config.yaml`.

This allows tuning the service without recompiling.

---

# Possible Future Improvements

* Streaming dictionary loading
* Distributed cache (Redis)
* Metrics (Prometheus)
* Request rate limiting
* Trie-based anagram search
* High-performance cache (Ristretto)

---

# Summary

This project demonstrates a clean Go service architecture:

* Config-driven
* REST API
* Concurrent client
* LRU caching
* Multilingual support
* Testable components

It is suitable for:

* learning Go service patterns
* system design exercises
* coding interview preparation
* lightweight production services
