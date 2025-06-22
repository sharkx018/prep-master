# Interview Prep App

A full-stack application for tracking interview preparation progress across different categories (DSA, LLD, HLD) with subcategories, status tracking, and progress analytics.

## 🏗️ Project Structure

```
interview-prep-app/
├── backend/                 # Go backend application
│   ├── cmd/
│   │   └── server/
│   │       └── main.go      # Application entry point
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
│   │   └── handlers/        # HTTP request handlers
│   │       ├── item_handler.go
│   │       └── stats_handler.go
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
│   │   ├── pages/
│   │   ├── services/
│   │   └── App.tsx
│   ├── package.json
│   └── README.md
├── README.md
└── SHOWCASE_README.md
```

## 🚀 Features

- ✅ **Clean Architecture**: Separation of concerns with repository, service, and handler layers
- ✅ **Scalable Structure**: Easy to add new features and maintain
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
   Edit `.env` with your database credentials:
   ```
   DATABASE_URL=postgresql://username:password@localhost:5432/interview_prep
   PORT=8080
   NODE_ENV=development
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

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Start the development server**
   ```bash
   npm start
   ```

The frontend will run on http://localhost:3000 and the backend API on http://localhost:8080.

## 🔌 API Endpoints

### API v1 (Recommended)

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

### Legacy Endpoints (Backward Compatible)
- `POST /items` - Create item
- `GET /items` - List items
- `GET /items/next` - Get next item
- `POST /items/skip` - Skip current item
- `PUT /items/:id/complete` - Complete item
- `GET /stats` - Get statistics
- `POST /reset` - Reset all items

## 📝 API Examples

### Create an Item with Subcategory
```bash
curl -X POST http://localhost:8080/api/v1/items \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Two Sum Problem",
    "link": "https://leetcode.com/problems/two-sum/",
    "category": "dsa",
    "subcategory": "arrays"
  }'
```

### Get Common Subcategories for a Category
```bash
curl http://localhost:8080/api/v1/items/subcategories/dsa
```

Response:
```json
{
  "category": "dsa",
  "subcategories": [
    "arrays", "strings", "linked-lists", "trees", "graphs", 
    "dynamic-programming", "sorting", "searching", "hashing", 
    "stack", "queue", "heap", "recursion", "backtracking",
    "greedy", "bit-manipulation", "math", "two-pointers",
    "sliding-window", "divide-conquer", "other"
  ]
}
```

### Get Next Item to Study
```bash
curl http://localhost:8080/api/v1/items/next
```

### Skip Current Item
```bash
curl -X POST http://localhost:8080/api/v1/items/skip
```

### Get Detailed Statistics with Subcategory Breakdown
```bash
curl http://localhost:8080/api/v1/stats/detailed
```

Response includes subcategory statistics:
```json
{
  "overall": {
    "total_items": 100,
    "completed_items": 75,
    "pending_items": 25,
    "progress_percentage": 75.0,
    "completed_all_count": 2
  },
  "categories": [
    {
      "category": "dsa",
      "total_items": 50,
      "completed_items": 40,
      "pending_items": 10,
      "progress_percentage": 80.0,
      "subcategories": [
        {
          "subcategory": "arrays",
          "total_items": 15,
          "completed_items": 12,
          "pending_items": 3,
          "progress_percentage": 80.0
        },
        ...
      ]
    },
    ...
  ]
}
```

### Filter Items by Category and Subcategory
```bash
curl "http://localhost:8080/api/v1/items?category=dsa&subcategory=arrays&status=pending"
```

### Get Subcategory Statistics
```bash
curl http://localhost:8080/api/v1/stats/category/dsa/subcategory/arrays
```

### Get Items
```bash
# Get all items
curl http://localhost:8080/api/v1/items

# Filter by category
curl http://localhost:8080/api/v1/items?category=dsa

# Filter by subcategory
curl http://localhost:8080/api/v1/items?subcategory=arrays

# Filter by status
curl "http://localhost:8080/api/v1/items?status=pending"

# Combine filters
curl "http://localhost:8080/api/v1/items?category=dsa&subcategory=arrays&status=pending"

# Pagination
curl "http://localhost:8080/api/v1/items?limit=10&offset=20"
```

## 📚 Subcategories

### DSA (Data Structures & Algorithms)
- **arrays** - Array manipulation problems
- **strings** - String processing and manipulation
- **linked-lists** - Linked list operations
- **trees** - Binary trees, BST, etc.
- **graphs** - Graph algorithms
- **dynamic-programming** - DP problems
- **sorting** - Sorting algorithms
- **searching** - Search algorithms
- **hashing** - Hash table problems
- **stack** - Stack-based problems
- **queue** - Queue-based problems
- **heap** - Heap/Priority queue problems
- **recursion** - Recursive solutions
- **backtracking** - Backtracking algorithms
- **greedy** - Greedy algorithms
- **bit-manipulation** - Bit operations
- **math** - Mathematical problems
- **two-pointers** - Two pointer technique
- **sliding-window** - Sliding window problems
- **divide-conquer** - Divide and conquer
- **other** - Miscellaneous

### LLD (Low Level Design)
- **design-patterns** - Software design patterns
- **solid-principles** - SOLID principles
- **object-modeling** - Object-oriented modeling
- **class-design** - Class structure design
- **database-design** - Database schema design
- **api-design** - API design
- **system-components** - Component design
- **scalability** - Scalability considerations
- **caching** - Caching strategies
- **other** - Miscellaneous

### HLD (High Level Design)
- **distributed-systems** - Distributed system design
- **microservices** - Microservices architecture
- **load-balancing** - Load balancing strategies
- **caching** - Caching at scale
- **databases** - Database choices and scaling
- **messaging** - Message queues and pub/sub
- **storage** - Storage systems
- **cdn** - Content delivery networks
- **monitoring** - System monitoring
- **security** - Security considerations
- **scalability** - Scaling strategies
- **reliability** - Reliability patterns
- **other** - Miscellaneous

## 🏗️ Architecture Benefits

### 1. **Repository Pattern**
- Isolates database logic
- Easy to mock for testing
- Can switch databases without changing business logic

### 2. **Service Layer**
- Contains all business logic
- Validates inputs
- Orchestrates complex operations

### 3. **Handler Layer**
- Handles HTTP-specific concerns
- Request/response transformation
- Error formatting

### 4. **Configuration Management**
- Centralized configuration
- Environment-based settings
- Type-safe config struct

## 🚀 Deployment

### Using Docker

#### Backend
```bash
cd backend
docker build -t interview-prep-backend .
docker run -p 8080:8080 --env-file .env interview-prep-backend
```

#### Frontend
```bash
cd frontend
docker build -t interview-prep-frontend .
docker run -p 3000:3000 interview-prep-frontend
```

#### Docker Compose (Full Stack)
```bash
cd backend
docker-compose up
```

### Deployment Platforms

#### Railway (Recommended)
1. Push to GitHub
2. Connect Railway to repository
3. Add PostgreSQL database
4. Set environment variables
5. Deploy

#### Render
1. Create new Web Service
2. Connect GitHub repository
3. Add PostgreSQL database
4. Configure environment variables
5. Deploy

#### Fly.io
1. Install flyctl
2. Run `fly launch`
3. Add PostgreSQL: `fly postgres create`
4. Deploy: `fly deploy`

## 🔄 Future Enhancements

The current architecture makes it easy to add:

1. **Authentication & Authorization**
   - Add middleware in `pkg/server`
   - Create auth service in `internal/services`

2. **Caching Layer**
   - Add Redis repository
   - Implement caching in service layer

3. **API Rate Limiting**
   - Add rate limit middleware
   - Configure per-endpoint limits

4. **Webhooks & Notifications**
   - Add notification service
   - Implement webhook handlers

5. **Import/Export Features**
   - Add CSV/JSON import handlers
   - Implement bulk operations

6. **Tags & Advanced Filtering**
   - Extend item model with tags
   - Add tag-based filtering

7. **Study Sessions**
   - Track study time per item
   - Implement spaced repetition

## 🧪 Testing

The architecture supports easy testing:

```go
// Example unit test for service
func TestItemService_CreateItem(t *testing.T) {
    // Mock repository
    mockRepo := &MockItemRepository{}
    service := services.NewItemService(mockRepo, nil)
    
    // Test logic
    item, err := service.CreateItem(&models.CreateItemRequest{
        Title:       "Test Item",
        Link:        "https://example.com",
        Category:    models.CategoryDSA,
        Subcategory: "arrays",
    })
    
    assert.NoError(t, err)
    assert.NotNil(t, item)
}
```

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