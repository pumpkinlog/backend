# 🎃 Pumpkinlog Backend

`pumpkinlog` is a presence-based tax evaluation system. It evaluates users travel histories against complex, jurisdiction-specific tax residency rules to determine  residency status. This repository holds the backend code.

This is a portfolio project built to demonstrate backend system design. It is under active development.

## ✨ Features

- Includes a flexible rule engine supporting the vast majority of tax residency rules.
- Stateless and horizontally scalable.
- Full suite of unit and integration tests.

## 🌎 How It Works

Pumpkinlog models complex residency logic using in `Rules` using child `Nodes`. The general app structure follows:

- `Region` is an isolated tax jurisdiction, be it a `country`, `state` or `zone`.
    - `Rule` is a child of a `Region`. It is the high level structure that contains an inital `Node`.
    - `Condition` is a child of a `Region`. It defines a question, and the user-response can be used as a `Rule` dependency. Answers are stored as an `Answer`.
- `Node` is a child of a `Rule`. Nodes can be the following types:
    - `Strategy` nodes evaluate a users presence in a `Region` and generate a residency profile.
    - `Condition` nodes depend on a region condition and subsequent user answer.
    - `And` and `Any` nodes can be used to combine two or more child nodes to create complex branching logic.

With this design, `pumpkinlog` can effictively model any day-based tax residency criteria for any tax jurisdiction globally.

## 🏗️ Architecture Overview

`pumpkinlog` is structured using a clean architecture architecture to keep concerns well separated:

```
├── cmd/                    # App entrypoint
├── internal/
│ ├── api/                  # HTTP API handlers
│ ├── app/                  # Core business logic
│ ├── cmd/                  # Command definitions
│ ├── cmdutil/              # CLI helper utilities
│ ├── domain/               # Core types and interfaces
│ ├── engine/               # Evaluation engine logic
│ ├── engine/strategies/    # Tax rule strategies
│ ├── repository/           # PostgreSQL data access layer
│ ├── service/              # Business logic
│ ├── seed/                 # App data seeder 
│ ├── worker/               # RabbitMQ workers
│ └── test/mocks/           # Mocks for testing
├── migrations/             # Postgres schema migrations
├── docker-compose.yml
├── Dockerfile
```

The app is designed to be stateless and horizontally scalable.

## 🚀 Getting Started

Start a local development instance using Docker:

```
docker compose up -d --build
```

Run database migrations and seed initial data:
```
make migrate_up
make seed
```

Start the backend API:
```
make run
```

- **API** ->                                ```http://localhost:4000```
- **API OpenAPI Documentation** ->          ```http://localhost:4000/docs```
- **Go Runtime Info & Exported Metrics** -> ```http://localhost:6060/debug/vars```

## 🛠️ Further Improvements

- Implement a cache repository layer using Redis.
- Build out a notification worker to alert users of residency thresholds.