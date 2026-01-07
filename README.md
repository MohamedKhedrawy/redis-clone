# Key-Value Store (Go)

A Redis-like key-value store built from scratch in Go, with a focus on **systems design, correctness, and failure handling** rather than feature completeness.

This project was built as a learning exercise to understand how real backend systems handle durability, crashes, time-based behavior, and clean architecture.

---

## Why this project

After working mostly on application-level projects, I wanted to:

- Go deeper into **systems-level concepts**
- Understand how databases survive crashes
- Learn how to design software around **failure, time, and concurrency**
- Build something interactive and non-trivial from first principles

This project prioritizes **correctness and clarity** over raw performance or feature parity with Redis.

---

## Features

### Networking
- Custom **TCP server**
- Simple text-based protocol
- Per-connection goroutines
- **Read/write deadlines** to protect against slow or misbehaving clients

### Storage Engine
- In-memory key-value store
- Thread-safe access using locks
- Clean separation between storage logic and other system components

### Write-Ahead Logging (WAL)
- Append-only WAL for durability
- **WAL-before-mutation rule** enforced for all write operations
- `fsync`-based persistence to guarantee crash safety
- **Startup replay** to rebuild state after a crash
- Safe handling of partial or corrupted tail records

### Time-To-Live (TTL)
- Optional TTL per key
- **Lazy expiration** on access for correctness
- Background eviction to reclaim memory
- Absolute expiration timestamps persisted in WAL to ensure correctness across restarts

### Reliability & Safety
- Panic isolation to prevent crashes from taking down the server
- Graceful handling of client disconnects
- Clear execution boundaries to avoid unintended coupling

---

## Architecture Overview

The system is intentionally structured with clear boundaries:

TCP Server
↓
Command Parsing & Execution
↓
Write-Ahead Log (for writes)
↓
In-Memory Store


Key architectural principles:

- The storage layer does **not** depend on networking or persistence implementations
- Persistence (WAL) wraps command execution instead of being embedded in the store
- Go interfaces are used to decouple components and enforce correct ordering

---

## Project Structure

cmd/
server/ # Application entry point

internal/
server/ # TCP server and connection handling
command/ # Command parsing and execution
store/ # In-memory key-value store + TTL logic
wal/ # Write-Ahead Log implementation and replay
log/ # (planned) async logging
snapshot/ # (planned) snapshots & WAL compaction

docs/
journal.md # Engineering notes and design decisions


---

## What this project focuses on

- Correctness under crashes
- Explicit ordering of state transitions
- Clear separation of concerns
- Understanding trade-offs (durability vs performance)
- Designing systems that fail safely

---

## What this project does NOT aim to be

- A production-ready database
- A Redis clone with full feature parity
- Highly optimized or benchmark-driven (yet)

These are intentional trade-offs to keep the learning focused.

---

## Future Improvements

- Snapshotting and WAL compaction
- Batched WAL fsync (group commit)
- Memory limits and eviction policies (LRU/LFU)
- Metrics and observability
- Unix socket support

---

## Notes

Design decisions, mistakes, and lessons learned are documented in  
`docs/journal.md` (SOON). 

This project is as much about **how** the system was built as **what** was built.
