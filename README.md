# Chirpy-Go

A Twitter-like social media REST API backend built with Go and PostgreSQL. Chirpy allows users to post short messages ("chirps"), manage their accounts, and upgrade to premium membership. ⭐

## What is Chirpy?

Chirpy-Go is a robust RESTful API service that provides the backend functionality for a Twitter-style microblogging platform. It features:

- **User Management**: Registration, authentication, and profile updates with secure password hashing (Argon2id)
- **JWT Authentication**: Token-based authentication with access and refresh tokens
- **Chirps (Posts)**: Create, retrieve, and delete short messages (max 140 characters)
- **Content Moderation**: Automatic profanity filtering
- **Premium Upgrades**: Integration with webhook system for user upgrades to "Chirpy Red"
- **Metrics Dashboard**: Admin interface to track application usage

## Why Chirpy?

Chirpy-Go demonstrates modern Go web development best practices:

- Clean architecture with separation of concerns
- Type-safe database queries using [sqlc](https://sqlc.dev/)
- Secure authentication with JWT and Argon2id password hashing
- RESTful API design following standard HTTP conventions
- PostgreSQL for reliable data persistence
- Proper error handling and logging

Whether you're learning Go web development or need a foundation for a social media API, Chirpy provides a solid, production-ready starting point.

## ⚙️ Installation

### Prerequisites

- Go 1.25.0 or higher
- PostgreSQL database
- [sqlc](https://sqlc.dev/) (for database code generation, if modifying queries)

### Clone the Repository

```bash
git clone https://github.com/gooneraki/chirpy-go.git
cd chirpy-go
```

### Install Dependencies

```bash
go mod download
```

## 🔧 Configuration

Create a `.env` file in the root directory with the following environment variables:

```env
DB_URL=postgres://username:password@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
JWT_SECRET=your-secret-key-here
POLKA_KEY=your-polka-api-key-here
```

### Environment Variables

- `DB_URL`: PostgreSQL connection string
- `PLATFORM`: Deployment platform (e.g., "dev", "prod")
- `JWT_SECRET`: Secret key for signing JWT tokens
- `POLKA_KEY`: API key for Polka webhook authentication

## 🗄️ Database Setup

### Run Migrations

Apply the database schema using your preferred migration tool, or manually execute the SQL files in order:

```bash
psql -d chirpy -f sql/schema/001_users.sql
psql -d chirpy -f sql/schema/002_chirps.sql
psql -d chirpy -f sql/schema/003_user_password.sql
psql -d chirpy -f sql/schema/004_refresh_tokens.sql
psql -d chirpy -f sql/schema/005_user_is_red.sql
```

### Generate Database Code (Optional)

If you modify the SQL queries, regenerate the Go code:

```bash
sqlc generate
```

## 🚀 Running the Application

Start the server:

```bash
go run .
```

The server will start on port `8080` by default. You should see:

```
Serving on port: 8080
```

## 📡 API Endpoints

### Health Check

- `GET /api/healthz` - Check if the API is running

### Users

- `POST /api/users` - Create a new user account
  ```json
  {
    "email": "user@example.com",
    "password": "securepassword"
  }
  ```

- `PUT /api/users` - Update user information (requires authentication)
  ```json
  {
    "email": "newemail@example.com",
    "password": "newpassword"
  }
  ```

### Authentication

- `POST /api/login` - Login and receive access/refresh tokens
  ```json
  {
    "email": "user@example.com",
    "password": "securepassword"
  }
  ```

- `POST /api/refresh` - Refresh access token using refresh token
- `POST /api/revoke` - Revoke a refresh token

### Chirps (Posts)

All chirp endpoints except `GET` require authentication via Bearer token.

- `GET /api/chirps` - Retrieve all chirps
  - Query params: `author_id` (filter by author), `sort` (asc/desc)
  
- `GET /api/chirps/{chirpID}` - Get a specific chirp by ID

- `POST /api/chirps` - Create a new chirp (requires authentication)
  ```json
  {
    "body": "This is my first chirp!"
  }
  ```

- `DELETE /api/chirps/{chirpID}` - Delete a chirp (requires authentication, author only)

### Webhooks

- `POST /api/polka/webhooks` - Webhook endpoint for premium upgrades (requires API key)

### Admin

- `GET /admin/metrics` - View application metrics (page visit counter)
- `POST /admin/reset` - Reset application state (development only)

### Static Files

- `/app/*` - Serves static files from the root directory

## 📁 Project Structure

```
chirpy-go/
├── main.go                      # Application entry point and server setup
├── chirps.go                    # Chirp creation handlers
├── users.go                     # User creation handler
├── handler_login.go             # Login authentication
├── handler_refresh.go           # Token refresh logic
├── handler_revoke.go            # Token revocation
├── handler_chirps_get.go        # Chirp retrieval
├── handler_chirps_get_by_id.go  # Single chirp retrieval
├── handler_chirps_delete.go     # Chirp deletion
├── handler_user_update.go       # User profile updates
├── handler_webhooks.go          # Webhook processing
├── metrics.go                   # Metrics middleware and handlers
├── readiness.go                 # Health check handler
├── reset.go                     # Reset handler (dev)
├── json.go                      # JSON response helpers
├── internal/
│   ├── auth/
│   │   ├── auth.go              # Authentication utilities (JWT, password hashing)
│   │   └── auth_test.go         # Auth tests
│   └── database/
│       ├── db.go                # Database connection
│       ├── models.go            # Database models
│       └── *.sql.go             # Generated sqlc query code
├── sql/
│   ├── schema/                  # Database schema migrations
│   └── queries/                 # SQL queries for sqlc
└── sqlc.yaml                    # sqlc configuration
```

## 🧪 Testing

Run all tests:

```bash
go test ./... -v
```

Run tests with coverage:

```bash
go test ./... -cover
```

## 🛡️ Security Features

- **Password Hashing**: Uses Argon2id for secure password storage
- **JWT Authentication**: Stateless authentication with signed tokens
- **Content Filtering**: Automatic profanity detection and replacement
- **API Key Authentication**: Webhook endpoints protected with API keys
- **Input Validation**: Enforces message length limits and validates user input

## 🔨 Development

### Code Generation

This project uses [sqlc](https://sqlc.dev/) for type-safe SQL queries. After modifying SQL queries:

```bash
sqlc generate
```

### Environment

For development, set `PLATFORM=dev` to enable development-only features like the reset endpoint.

## 📄 License

This project is available as open source for educational purposes.

## 👏 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. Ensure that:

- Your code passes all existing tests
- You add tests for new functionality
- Code follows Go conventions and formatting (`go fmt`)
- Commit messages are clear and descriptive

## 🤝 Acknowledgments

Built with:
- [Go](https://golang.org/) - Programming language
- [PostgreSQL](https://www.postgresql.org/) - Database
- [sqlc](https://sqlc.dev/) - SQL compiler
- [golang-jwt](https://github.com/golang-jwt/jwt) - JWT implementation
- [argon2id](https://github.com/alexedwards/argon2id) - Password hashing
- [godotenv](https://github.com/joho/godotenv) - Environment variable loading
