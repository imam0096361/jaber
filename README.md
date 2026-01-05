# News Portal 

## Features

- ✅ Article Management (Create, Read, Update)
- ✅ Category System
- ✅ Image Upload
- ✅ User Authentication (Optional)
- ✅ Input Validation
- ✅ SQLite Database
- ✅ RESTful API

## Project Structure

```
src/
├── main.go              # Application entry point
├── config/             # Configuration management
│   └── app.go
├── router/             # Route definitions
│   └── routes.go
├── controller/         # API request handlers
│   ├── article_controller.go
│   └── auth_controller.go
├── service/           # Business logic
│   ├── article_service.go
│   └── user_service.go
├── model/             # Data structures
│   ├── article.go
│   └── user.go
├── middleware/        # HTTP middleware
│   └── jwt.go
├── database/          # Database setup & migrations
│   ├── sqlite.go
│   ├── seeder.go
│   └── migrations/
├── utils/            # Utility functions
│   ├── bcrypt.go
│   ├── logger.go
│   └── token.go
└── validation/       # Input validation
    └── validator.go
```

## Installation

```bash
# Clone repository
cd /home/jaber/Desktop/news-portal

# Install dependencies
go mod download
go mod tidy

# Setup environment
cp .env.example .env
```

## Running the Server

```bash
# Development
go run src/main.go

# Production
go build -o news-portal src/main.go
./news-portal
```

Server will run on `http://localhost:9999`

## API Endpoints

### Public Endpoints

#### Get All Articles
```bash
GET /api/articles
GET /api/articles?category=জাতীয়
GET /api/articles?search=keyword
```

#### Get Featured Articles
```bash
GET /api/featured
```

#### Get Single Article
```bash
GET /api/article?id=1
```

#### Get Categories
```bash
GET /api/categories
```

#### Create Article
```bash
POST /api/add-article
Content-Type: application/json

{
  "title": "Article Title",
  "content": "Article content...",
  "category": "Category Name",
  "author": "Author Name",
  "featured": false
}
```

#### Upload Image
```bash
POST /api/upload-image
Content-Type: multipart/form-data

[file upload]
```

### Authentication Endpoints (Optional)

#### Register
```bash
POST /api/auth/register
{
  "name": "User Name",
  "email": "email@example.com",
  "password": "password123"
}
```

#### Login
```bash
POST /api/auth/login
{
  "email": "email@example.com",
  "password": "password123"
}
```

#### Get Profile
```bash
GET /api/auth/profile
Authorization: Bearer TOKEN
```

## Configuration

`.env` file থেকে পড়া হয়:

```env
APP_PORT=9999
APP_ENV=development
DB_PATH=./news.db
LOG_LEVEL=info
MAX_UPLOAD_SIZE=10485760
```

## Database

SQLite database automatically created at `./news.db`

### Tables

- **users** - User accounts
- **articles** - News articles

### Migrations

See `src/database/migrations/` for SQL schema

## Development

### Add New Feature

1. Create model in `src/model/`
2. Create service in `src/service/`
3. Create controller in `src/controller/`
4. Add routes in `src/router/routes.go`
5. Add validation in `src/validation/`

### Running Tests

```bash
go test ./...
```

## Troubleshooting

### Port already in use
```bash
lsof -ti:9999 | xargs kill -9
```

### Database lock
```bash
rm news.db
```

## License

MIT License
## API Documentation
