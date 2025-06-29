#!/bin/bash

echo "🚀 Interview Prep App - JWT Authentication Setup"
echo "================================================="

# Check if required tools are installed
command -v go >/dev/null 2>&1 || { echo "❌ Go is required but not installed. Please install Go 1.21+"; exit 1; }
command -v node >/dev/null 2>&1 || { echo "❌ Node.js is required but not installed. Please install Node.js 16+"; exit 1; }
command -v npm >/dev/null 2>&1 || { echo "❌ npm is required but not installed. Please install npm"; exit 1; }

echo "✅ Prerequisites check passed"

# Setup backend
echo ""
echo "🔧 Setting up backend..."
cd backend

if [ ! -f ".env" ]; then
    echo "📝 Creating backend .env file..."
    cat > .env << EOL
DATABASE_URL=postgresql://username:password@localhost:5432/interview_prep
PORT=8080
NODE_ENV=development

# Authentication - CHANGE THESE IN PRODUCTION!
AUTH_USERNAME=admin
AUTH_PASSWORD=secure123
JWT_SECRET=super_secret_jwt_key_for_interview_prep_app_$(date +%s)_$(openssl rand -hex 16 2>/dev/null || echo $(date +%s))
EOL
    echo "✅ Backend .env created"
else
    echo "⚠️  Backend .env already exists, skipping..."
fi

echo "📦 Installing Go dependencies..."
go mod tidy
if [ $? -eq 0 ]; then
    echo "✅ Go dependencies installed"
else
    echo "❌ Failed to install Go dependencies"
    exit 1
fi

echo "🔨 Building backend..."
go build -o server ./cmd/server
if [ $? -eq 0 ]; then
    echo "✅ Backend built successfully"
else
    echo "❌ Failed to build backend"
    exit 1
fi

# Setup frontend
echo ""
echo "🔧 Setting up frontend..."
cd ../frontend

if [ ! -f ".env" ]; then
    echo "📝 Creating frontend .env file..."
    cat > .env << EOL
REACT_APP_API_URL=http://localhost:8080
EOL
    echo "✅ Frontend .env created"
else
    echo "⚠️  Frontend .env already exists, skipping..."
fi

echo "📦 Installing npm dependencies..."
npm install
if [ $? -eq 0 ]; then
    echo "✅ Frontend dependencies installed"
else
    echo "❌ Failed to install frontend dependencies"
    exit 1
fi

echo "🔨 Building frontend..."
npm run build
if [ $? -eq 0 ]; then
    echo "✅ Frontend built successfully"
else
    echo "❌ Failed to build frontend"
    exit 1
fi

cd ..

echo ""
echo "🎉 Setup completed successfully!"
echo ""
echo "📋 Next steps:"
echo "1. Make sure PostgreSQL is running on localhost:5432"
echo "2. Create a database named 'interview_prep'"
echo "3. Start the backend: cd backend && go run cmd/server/main.go"
echo "4. Start the frontend: cd frontend && npm start"
echo ""
echo "🔐 Login credentials:"
echo "   Username: admin"
echo "   Password: secure123"
echo ""
echo "🌐 Access the app:"
echo "   Frontend: http://localhost:3000"
echo "   Backend API: http://localhost:8080"
echo ""
echo "⚠️  IMPORTANT: Change the default credentials in production!"
echo "   Edit backend/.env and update AUTH_USERNAME, AUTH_PASSWORD, and JWT_SECRET"