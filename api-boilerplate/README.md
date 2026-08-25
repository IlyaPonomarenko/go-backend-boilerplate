# Go API Boilerplate

A small HTTP JSON API built with Go's standard library. It demonstrates routing, request decoding, response encoding, validation, SQLite persistence, bearer authentication, structured logging, and graceful shutdown.

## Requirements

- Go 1.22 or newer

This project uses enhanced `http.ServeMux` route patterns such as `GET /api/posts/{id}`, which require Go 1.22 or newer.

## Run the API

From this directory, configure the environment and run:

Copy `.env.example` to `.env` for reference, or set the variables directly. PowerShell does not load `.env` automatically:

```powershell
$env:API_TOKEN = "replace-with-a-long-random-token"
$env:DATABASE_PATH = ".\data\api.db"
go run ./cmd/api
```

The server starts at:

```text
http://localhost:8080
```

The API requires this header on every `/api/` request:

```text
Authorization: Bearer replace-with-a-long-random-token
```

The `/healthz` endpoint is intentionally unauthenticated so deployment systems can check whether the process is ready.

To compile a binary instead:

```powershell
go build -o api.exe ./cmd/api
.\api.exe
```

## Endpoints

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/healthz` | Health check |
| `GET` | `/api/posts` | List all posts |
| `POST` | `/api/posts` | Create a post |
| `GET` | `/api/posts/{id}` | Get one post |
| `PUT` | `/api/posts/{id}` | Replace a post |
| `DELETE` | `/api/posts/{id}` | Delete a post |

A post has this JSON shape:

```json
{
  "userId": 1,
  "id": 1,
  "title": "Hello",
  "body": "Post content"
}
```

## Try the API

Check that the server is running:

```powershell
Invoke-RestMethod http://localhost:8080/healthz
```

Create a post:

```powershell
$body = '{"userId":1,"title":"Hello","body":"My first API post"}'
$headers = @{ Authorization = "Bearer $env:API_TOKEN" }
Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:8080/api/posts `
  -ContentType application/json `
  -Body $body `
  -Headers $headers
```

List posts:

```powershell
$headers = @{ Authorization = "Bearer $env:API_TOKEN" }
Invoke-RestMethod http://localhost:8080/api/posts -Headers $headers
```

Get one post:

```powershell
$headers = @{ Authorization = "Bearer $env:API_TOKEN" }
Invoke-RestMethod http://localhost:8080/api/posts/1 -Headers $headers
```

Update a post:

```powershell
$body = '{"userId":1,"title":"Updated title","body":"Updated content"}'
$headers = @{ Authorization = "Bearer $env:API_TOKEN" }
Invoke-RestMethod `
  -Method Put `
  -Uri http://localhost:8080/api/posts/1 `
  -ContentType application/json `
  -Body $body `
  -Headers $headers
```

Delete a post:

```powershell
Invoke-RestMethod `
  -Method Delete `
  -Uri http://localhost:8080/api/posts/1 `
  -Headers @{ Authorization = "Bearer $env:API_TOKEN" }
```

## How the code works

1. `cmd/api/main.go` loads configuration, opens SQLite, and wires the repository, service, handlers, routes, and server.
2. The multiplexer maps HTTP methods and paths to handler functions.
3. Handlers decode requests, call the service, and return JSON responses.
4. The service validates business rules before calling the repository.
5. The SQLite repository reads and modifies the persistent database.
6. `decodeJSON` accepts only `application/json` request bodies and decodes them into Go structs.
7. `writeJSON` sets the response content type and serializes a Go value as JSON.
8. Bearer authentication protects `/api/` routes, while `/healthz` remains public.
9. `loggingMiddleware` emits structured JSON logs with method, path, and duration.
10. Shutdown signals stop accepting requests and allow up to 10 seconds for active requests to finish.

## Project structure

```text
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   └── service/
├── go.mod
├── .env.example
└── README.md
```

## Important limitation

The production server uses SQLite at `DATABASE_PATH` and creates the schema automatically. The memory repository remains available as a lightweight test double. For a larger deployment, replace the SQLite repository with PostgreSQL while keeping the HTTP handlers and service mostly unchanged.

## Useful next steps

- Add automated handler tests with `net/http/httptest`.
- Add database migrations and connection health checks.
- Add refreshable credentials or an identity provider instead of one shared bearer token.
- Add rate limiting, CORS policy, TLS termination, and metrics for the deployment environment.
- Add database migrations and connection health checks.
- Add refreshable credentials or an identity provider instead of one shared bearer token.
- Add rate limiting, CORS policy, TLS termination, and metrics for the deployment environment.

To use an external API from your Go backend, follow these steps:

1. **Read the API documentation**
   - Find the base URL.
   - Identify endpoints and HTTP methods.
   - Check required headers, parameters, and JSON formats.
   - Understand authentication and rate limits.

2. **Store credentials securely**

Create `.env` locally:

```text
API_KEY=your-secret-key
```

Never commit `.env`; your `.gitignore` already excludes it.

3. **Create request and response models**

```go
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
```

4. **Create an HTTP client**

```go
client := &http.Client{
    Timeout: 10 * time.Second,
}
```

5. **Build and send the request**

```go
request, err := http.NewRequest(
    http.MethodGet,
    "https://api.example.com/users/1",
    nil,
)
if err != nil {
    return err
}

request.Header.Set("Accept", "application/json")
request.Header.Set("Authorization", "Bearer "+apiKey)

response, err := client.Do(request)
if err != nil {
    return err
}
defer response.Body.Close()
```

6. **Check the status code**

```go
if response.StatusCode < 200 || response.StatusCode >= 300 {
    return fmt.Errorf("API returned status: %s", response.Status)
}
```

7. **Decode the response**

```go
var user User

if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
    return err
}
```

8. **Place the code in the project layers**

For the current project:

```text
internal/
├── model/         # Request and response structs
├── repository/    # External API calls
├── service/       # Business rules
└── handler/       # Your own HTTP endpoints
```

The flow would be:

```text
Client
  ↓
Your handler: GET /api/users/1
  ↓
Your service
  ↓
Your repository
  ↓
External API
```

This lets your frontend call your API without exposing the external API key.

For example, your own endpoint could be:

```text
GET /api/users/1
```

Your Go backend receives that request, calls the external service, checks and decodes its response, then returns a controlled response to the client. This is generally preferable to calling the external API directly from a browser because credentials, error handling, caching, and rate limiting stay on your server.