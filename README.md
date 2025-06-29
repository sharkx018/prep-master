# Interview Prep App

A full-stack application for tracking interview preparation progress across different categories (DSA, LLD, HLD) with subcategories, status tracking, and progress analytics. **Now with JWT Authentication for secure API access.**

## 🔐 Security Features

- ✅ **JWT Authentication**: All API endpoints are protected with JWT tokens
- ✅ **Login System**: Modern login interface with secure credential handling
- ✅ **Token Management**: Automatic token refresh and logout functionality
- ✅ **Protected Routes**: Frontend routes are protected and redirect to login when unauthenticated

## 🏗️ Project Structure

```
interview-prep-app/
├── backend/                 # Go backend application
│   ├── cmd/
│   │   └── server/
│   │       └── main.go      # Application entry point
│   │   ├── internal/
│   │   │   ├── config/          # Configuration management
│   │   │   │   └── config.go
│   │   │   ├── models/          # Data models and types
│   │   │   │   ├── item.go
│   │   │   │   └── stats.go
│   │   │   ├── database/        # Database connection and migrations
│   │   │   │   ├── connection.go
│   │   │   │   └── migrations.go
│   │   │   ├── repositories/    # Data access layer
│   │   │   │   ├── item_repository.go
│   │   │   │   └── stats_repository.go
│   │   │   ├── services/        # Business logic layer
│   │   │   │   ├── item_service.go
│   │   │   │   └── stats_service.go
│   │   │   ├── handlers/        # HTTP request handlers
│   │   │   │   ├── item_handler.go
│   │   │   │   ├── stats_handler.go
│   │   │   │   └── auth_handler.go    # JWT authentication
│   │   │   └── middleware/      # HTTP middleware
│   │   │       └── auth_middleware.go # JWT validation
│   │   ├── pkg/
│   │   │   └── server/          # HTTP server setup
│   │   │       └── server.go
│   │   ├── go.mod
│   ├── internal/
│   │   ├── config/          # Configuration management
│   │   │   └── config.go
│   │   ├── models/          # Data models and types
│   │   │   ├── item.go
│   │   │   └── stats.go
│   │   ├── database/        # Database connection and migrations
│   │   │   ├── connection.go
│   │   │   └── migrations.go
│   │   ├── repositories/    # Data access layer
│   │   │   ├── item_repository.go
│   │   │   └── stats_repository.go
│   │   ├── services/        # Business logic layer
│   │   │   ├── item_service.go
│   │   │   └── stats_service.go
│   │   ├── handlers/        # HTTP request handlers
│   │   │   ├── item_handler.go
│   │   │   ├── stats_handler.go
│   │   │   └── auth_handler.go    # JWT authentication
│   │   └── middleware/      # HTTP middleware
│   │       └── auth_middleware.go # JWT validation
│   ├── pkg/
│   │   └── server/          # HTTP server setup
│   │       └── server.go
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   ├── docker-compose.yml
│   ├── env.example
│   └── Makefile
├── frontend/                # React frontend application
│   ├── public/
│   ├── src/
│   │   ├── components/
│   │   │   ├── Layout.tsx
│   │   │   └── ProtectedRoute.tsx  # Route protection
│   │   ├── contexts/
│   │   │   └── AuthContext.tsx     # Authentication context
│   │   ├── pages/
│   │   │   └── Login.tsx           # Login page
│   │   ├── services/
│   │   │   └── api.ts              # API client with JWT
│   │   └── App.tsx
│   ├── package.json
│   └── README.md
├── README.md
└── SHOWCASE_README.md
```

## 🚀 Features

- ✅ **Clean Architecture**: Separation of concerns with repository, service, and handler layers
- ✅ **Scalable Structure**: Easy to add new features and maintain
- ✅ **JWT Authentication**: Secure API access with JSON Web Tokens
- ✅ **Protected Frontend**: Login system with automatic token management
- ✅ **Category Management**: Support for DSA, LLD, HLD categories
- ✅ **Subcategory Organization**: Organize items within categories by subcategories
- ✅ **Progress Tracking**: Track completion status and statistics at category and subcategory levels
- ✅ **Random Selection**: Get random pending items to study
- ✅ **Completion Cycles**: Track how many times all items have been completed
- ✅ **RESTful API**: Both legacy and versioned API endpoints
- ✅ **Database Migrations**: Automatic schema setup
- ✅ **Environment Configuration**: Easy configuration management

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- Node.js 16+ and npm
- Git

## 🛠️ Installation

### Backend Setup

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd interview-prep-app
   ```

2. **Set up backend environment variables**
   ```bash
   cd backend
   cp env.example .env
   ```
   Edit `.env` with your database and authentication credentials:
   ```
   DATABASE_URL=postgresql://username:password@localhost:5432/interview_prep
   PORT=8080
   NODE_ENV=development
   
   # Authentication (REQUIRED)
   AUTH_USERNAME=admin
   AUTH_PASSWORD=secure123
   JWT_SECRET=your_super_secret_jwt_key_here_make_it_long_and_random
   ```

3. **Install Go dependencies**
   ```bash
   go mod tidy
   ```

4. **Run the backend**
   ```bash
   go run cmd/server/main.go
   ```

### Frontend Setup

1. **Navigate to frontend directory**
   ```bash
   cd ../frontend
   ```

2. **Set up frontend environment variables**
   ```bash
   cp env.example .env
   ```
   Edit `.env` with your API URL:
   ```
   REACT_APP_API_URL=http://localhost:8080
   ```

3. **Install dependencies**
   ```bash
   npm install
   ```

4. **Start the development server**
   ```bash
   npm start
   ```

The frontend will run on http://localhost:3000 and the backend API on http://localhost:8080.

## 🔐 Authentication

### Login Credentials
Use the credentials you set in your backend `.env` file:
- **Username**: Value of `AUTH_USERNAME` (default: `admin`)
- **Password**: Value of `AUTH_PASSWORD` (default: `secure123`)

### JWT Token Management
- Tokens are automatically included in all API requests
- Tokens expire after 24 hours
- Automatic logout on token expiration
- Tokens are stored securely in localStorage

### Security Best Practices
1. **Change default credentials** in production
2. **Use strong JWT secret** (minimum 32 characters)
3. **Use HTTPS** in production
4. **Regularly rotate JWT secrets**

## 🔌 API Endpoints

### Authentication (Public)
- `POST /api/v1/auth/login` - Login with username/password, returns JWT token

### API v1 (Protected - Requires JWT Token)

All API endpoints now require a valid JWT token in the Authorization header:
```bash
Authorization: Bearer <your-jwt-token>
```

#### Items
- `POST /api/v1/items` - Create new item
- `GET /api/v1/items` - List items (with filters)
- `GET /api/v1/items/next` - Get random pending item
- `POST /api/v1/items/skip` - Skip current item and get next
- `GET /api/v1/items/subcategories/:category` - Get common subcategories for a category
- `GET /api/v1/items/:id` - Get specific item
- `PUT /api/v1/items/:id` - Update item
- `PUT /api/v1/items/:id/complete` - Mark item as complete
- `DELETE /api/v1/items/:id` - Delete item
- `POST /api/v1/items/reset` - Reset all items to pending

#### Statistics
- `GET /api/v1/stats` - Get overall statistics
- `GET /api/v1/stats/detailed` - Get detailed stats with category and subcategory breakdown
- `GET /api/v1/stats/category/:category` - Get stats for specific category
- `GET /api/v1/stats/category/:category/subcategory/:subcategory` - Get stats for specific subcategory
- `POST /api/v1/stats/reset-completed-all` - Reset completion counter

### Legacy Endpoints (Protected - Backward Compatible)
All legacy endpoints are also protected and require JWT authentication.

## 📝 API Examples

### Login to Get JWT Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "secure123"
  }'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": "admin"
}
```

### Create an Item (with JWT Token)
```bash
curl -X POST http://localhost:8080/api/v1/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "title": "Two Sum Problem",
    "link": "https://leetcode.com/problems/two-sum/",
    "category": "dsa",
    "subcategory": "arrays"
  }'
```

### Get Items (with JWT Token)
```bash
curl -H "Authorization: Bearer <your-jwt-token>" \
  http://localhost:8080/api/v1/items
```

### Error Response for Missing/Invalid Token
```json
{
  "error": "Authorization header required"
}
```

```json
{
  "error": "Invalid or expired token"
}
```

## 🚀 Deployment

### Environment Variables for Production

#### Backend (.env)
```bash
DATABASE_URL=postgresql://user:pass@host:5432/dbname
PORT=8080
NODE_ENV=production

# Authentication - CHANGE THESE IN PRODUCTION!
AUTH_USERNAME=your_secure_username
AUTH_PASSWORD=your_very_secure_password
JWT_SECRET=your_super_long_random_jwt_secret_key_at_least_32_characters
```

#### Frontend (.env)
```bash
REACT_APP_API_URL=https://your-api-domain.com
```

### Security Considerations for Production

1. **Use strong, unique credentials**
2. **Use a long, random JWT secret (32+ characters)**
3. **Enable HTTPS/TLS**
4. **Set secure CORS policies**
5. **Use environment variables, never hardcode secrets**
6. **Consider shorter JWT expiration times**
7. **Implement refresh token mechanism for longer sessions**

## 🔄 Future Security Enhancements

The current JWT implementation provides a solid foundation and can be extended with:

1. **Refresh Tokens**
   - Longer-lived refresh tokens
   - Automatic token refresh

2. **Role-Based Access Control (RBAC)**
   - Multiple user roles
   - Permission-based access

3. **Rate Limiting**
   - Login attempt limiting
   - API rate limiting per user

4. **Session Management**
   - Active session tracking
   - Remote logout capability

5. **Two-Factor Authentication (2FA)**
   - TOTP support
   - SMS/Email verification

6. **OAuth Integration**
   - Google/GitHub login
   - Social authentication

## 🧪 Testing Authentication

### Manual Testing
1. Start the application
2. Try accessing any API endpoint without token (should fail)
3. Login with correct credentials
4. Use returned token for API calls
5. Try with invalid token (should fail)

### Frontend Testing
1. Open http://localhost:3000
2. Should redirect to login page
3. Login with configured credentials
4. Should access the dashboard
5. Logout should redirect back to login

## 📚 Development Guidelines

1. **Adding New Features**
   - Start with models in `internal/models`
   - Add repository methods if needed
   - Implement business logic in services
   - Create handlers for HTTP endpoints
   - Update routes in `pkg/server`

2. **Database Changes**
   - Add migrations in `internal/database/migrations.go`
   - Update models accordingly
   - Modify repository methods

3. **Error Handling**
   - Return errors from repositories/services
   - Format errors appropriately in handlers
   - Use proper HTTP status codes

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

MIT License - feel free to use this for your interview preparation! 