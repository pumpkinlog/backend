# 🎃 Pumpkinlog Backend

**Pumpkinlog** is a backend service for an automated tax residency tracking tool aimed at frequent travelers and digital nomads. It evaluates users travel histories against complex, location-specific tax residency rules to determine their residency status.

This is a portfolio project built to demonstrate backend system design, clean architecture in Go It is under active development.

## ✨ Features

- Idiomatically designed with go clean architecture.
- Stateless and horizontally scalable.
- Supports daily-based tax residency evaluation for countries, states, and zones.
- Includes a flexible rule engine supporting the vast majority of tax residency rules.
- Unit and integration tests.
- Condition-based rules enable modeling of nuanced residency logic.

## 🏗️ Architecture Overview

Pumpkinlog is structured using a clean layered architecture to keep concerns well separated:


```
├── cmd/            # App entrypoint
├── internal/
│ ├── api/          # HTTP API handlers
│ ├── app/          # Core business logic
│ ├── cmd/          # Command definitions
│ ├── cmdutil/      # CLI helper utilities
│ ├── domain/       # Core domain types and interfaces
│ ├── repository/   # PostgreSQL data access layer
│ ├── service/      # Business logic
│ ├── seed/         # App seed data 
│ ├── worker/       # RabbitMQ workers
│ └── test/mocks/   # Mocks for testing
├── migrations/     # Postgres schema migrations
├── docker-compose.yml
```

The service is designed to be stateless and environment-agnostic. Infrastructure layers can be swapped or mocked independently.

## 🌎 Real World Example

Pumpkinlog models complex, real-world residency logic through rule and condition declarations.

For example, Jersey considers an individual a tax resident if they meet any of the following:

```
- Is present in Jersey for 183 days in any one tax year;
- Maintains a place of abode in Jersey and stays one night in Jersey in a tax year
- Do not maintain a place of abode but visit for an average of three months per year over four years.
```

When a request is made to the /evaluate endpoint, Pumpkinlog evaluates all applicable rules for the given region, taking into account any user-provided inputs. This process generates a comprehensive tax residency profile specific to that region. To optimize performance, Pumpkinlog uses a worker-based architecture combined with memoization. Expensive evaluate operations are cached, ensuring that repeated evaluations are served quickly without redundant computation.

The tax residency profile schema can be found in `internal/domain/evaluation.go` and rule specific evaluators are implemented in `internal/app/evaluator/`.

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

The API will be accessible at:

```
http://localhost:4000
```

## ✅ Tech Stack

- Go
- PostgreSQL
- RabbitMQ
- Docker / Docker Compose

## 🛠️ Further Improvements

- Add support for additional tax regions.
- Implement a cache repository layer using Redis.
- Build out a notification worker to alert users of residency thresholds.
