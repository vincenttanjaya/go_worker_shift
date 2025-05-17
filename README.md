# Payd Mini Project: Daily Worker Roster Management System

## Overview

This project is a simple scheduling and roster management system for managing shifts of daily workers. It supports real-world scheduling constraints and simulates key features used by companies with hourly/daily staff.

- **Backend:** Go (with SQLite, REST, Dockerized)
- **Frontend (Employee & Admin UIs):** Svelte, Dockerized
- **Deployment:** Multi-service via Docker Compose

---
```
go_worker_shift/
├── backend/           # Go API backend
│   ├── cmd/
│   ├── internal/
│   ├── migrations/
│   └── Dockerfile
├── employee-frontend/ # Svelte Employee UI
│   └── Dockerfile
├── admin-frontend/    # Svelte Admin UI
│   └── Dockerfile
├── docker-compose.yml
├── README.md
```

## Features

### Employee Interface

- Login with Worker ID and see personalized greeting
- View available (unassigned) shifts and request to work them (one click)
- See status of requested shifts (pending, approved, rejected)
- Request button disabled if already requested

### Admin Interface

- View all shift requests, filter by status, approve/reject with one click
- Create, edit, assign/unassign, and delete shifts (inline table editing)
- Create, edit, and delete workers (inline table editing)
- View all workers and shifts with assigned workers’ names
- All data auto-refreshes on changes

---

## Business Rules

- Workers cannot request shifts already assigned to someone else
- No overlapping shift requests allowed per worker
- Max 1 shift per day, max 5 shifts per week per worker
- Admin can override or reassign approved shifts
- Conflict checking occurs on both worker request and admin approval
- Shift times are stored and compared in UTC

---

## Tech Stack

- **Backend:** Go, SQLite, Gorilla Mux, REST
- **Frontend:** Svelte (Employee and Admin apps)
- **Deployment:** Docker, Docker Compose

---

## Local Development and Deployment

### 1. **Prerequisites**
- [Docker](https://www.docker.com/products/docker-desktop/) installed

### 2. **Clone the Repository**
```bash
git clone <your-repo-url>
cd go_worker_shift
docker compose up --build
```

```
API: http://localhost:8080
Employee: http://localhost:3000
Admin: http://localhost:4000
```