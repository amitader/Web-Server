# 🐦 Web Server (Chirpy)
A fully functional RESTful web server in Go for a microblogging platform called **Chirpy**. It supports user authentication, content posting (chirps), token-based security, admin monitoring, and webhook integration — backed by PostgreSQL and structured via sqlc.

## 📑 Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Environment Variables](#environment-variables)
- [Usage](#usage)
- [API Overview](#api-overview)
- [Database & Migrations](#database--migrations)
- [Project Structure](#project-structure)
- [Dependencies](#dependencies)
- [Troubleshooting](#troubleshooting)
- [License](#license)

## 🚀 Introduction
**Chirpy** is a backend service written in Go for a Twitter-like application. It offers:

- JWT-based authentication
- Secure password handling (bcrypt)
- Role-based operations
- Token revocation and refresh
- Chirp management (create, read, delete)
- Metrics and readiness endpoints
- Webhook support for external service integration

## ✨ Features

- 🐦 **Chirp Management**: Users can create, fetch, and delete chirps.
- 🔐 **Authentication**: JWT access tokens + refresh tokens.
- 🔁 **Token Rotation**: Support for refresh & revoke mechanisms.
- 📊 **Admin Metrics**: File server access tracking via `/admin/metrics`.
- 🛠️ **Reset Endpoint**: Delete all users in dev environment.
- 📡 **Webhook Integration**: E.g., handle `user.upgraded` events.
- 🧪 **Tests Included**: Core `auth` logic is unit tested.

## ⚙️ Installation
1. **Clone the repository:**

```bash
git clone https://github.com/amitader/web-Server.git
cd web-Server
```
2. **Setup PostgreSQL database:**

Create the database and run migrations in `sql/schema`.

3. **Install dependencies:**

```bash
go mod tidy
```

4. **Start the server:**

```bash
go run main.go
```

## 🔧 Environment Variables
You must create a `.env` file or export the following variables:

```env
DB_URL=postgresql://<user>:<password>@localhost:5432/<dbname>?sslmode=disable
PLATFORM=dev
SECRET=your-jwt-secret
POLKA_KEY=your-api-key
```

## 📡 Usage
Once the server is running on port `8080`, you can interact with the following:

### 🧑 Users
- `POST /api/users`: Create user

- `POST /api/login`: Login & receive access/refresh tokens

- `PUT /api/users`: Change email/password (auth required)

### 🐦 Chirps
- `POST /api/chirps`: Create chirp (auth required)

- `GET /api/chirps`: Fetch all chirps

- `GET /api/chirps/{chirpID}`: Fetch single chirp

- `DELETE /api/chirps/{chirpID}`: Delete chirp (auth & ownership required)

### 🔄 Tokens
- `POST /api/refresh`: Get new access token using refresh token

- `POST /api/revoke`: Invalidate refresh token

### 📈 Admin
- `GET /admin/metrics`: View access count to static files

- `POST /admin/reset`: Dev-only reset of user data

### ⚙️ System
- `GET /api/healthz`: Readiness probe

- `POST /api/polka/webhooks`: Accept webhook calls with API Key

## 🗃️ Database & Migrations
This project uses **PostgreSQL** with `sqlc` for type-safe query generation.

**Schema files**: `sql/schema/*.sql`

**Query definitions**: `sql/queries/*.sql`

**Generated Go code**: `internal/database/`

To regenerate SQL bindings:

```bash
sqlc generate
```
## 🗂️ Project Structure
```plaintext
├── main.go
├── chirps.go
├── users.go
├── ...
├── internal/
│   ├── auth/          # Token logic, password hashing
│   └── database/      # sqlc-generated DB access
├── sql/
│   ├── schema/        # DB schema migrations
│   └── queries/       # SQL query templates
└── assets/            # Static assets (e.g., 
logo)
```

## 📦 Dependencies
From go.mod:

[github.com/golang-jwt/jwt](github.com/golang-jwt/jwt)

[github.com/google/uuid](github.com/google/uuid)

[github.com/joho/godotenv](github.com/joho/godotenv)

[github.com/lib/pq](github.com/lib/pq)

[golang.org/x/crypto](golang.org/x/crypto)

## 🧪 Troubleshooting
- ❌ **JWT Invalid**: Check `SECRET` environment variable and ensure token is not expired.

- ❌ **DB connection failed**: Confirm `DB_URL` is correctly set and reachable.

- ❌ **Webhook Unauthorized**: Make sure `Authorization: ApiKey ...` header matches `POLKA_KEY`.