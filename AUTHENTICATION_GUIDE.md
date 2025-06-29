# 🔐 Authentication Guide

This guide explains how to set up and use the JWT authentication system in the Interview Prep App.

## 🚀 Quick Start

### Option 1: Automated Setup (Recommended)
```bash
# Run the setup script
./setup.sh
```

### Option 2: Manual Setup

#### Backend Setup
1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Create `.env` file:
   ```bash
   cp env.example .env
   ```

3. Edit `.env` with your credentials:
   ```env
   DATABASE_URL=postgresql://username:password@localhost:5432/interview_prep
   PORT=8080
   NODE_ENV=development
   
   # Authentication
   AUTH_USERNAME=admin
   AUTH_PASSWORD=secure123
   JWT_SECRET=your_super_secret_jwt_key_here_make_it_long_and_random
   ```

4. Install dependencies and build:
   ```bash
   go mod tidy
   go build -o server ./cmd/server
   ```

#### Frontend Setup
1. Navigate to the frontend directory:
   ```bash
   cd ../frontend
   ```

2. Create `.env` file:
   ```bash
   cp env.example .env
   ```

3. Edit `.env`:
   ```env
   REACT_APP_API_URL=http://localhost:8080
   ```

4. Install dependencies:
   ```bash
   npm install
   ```

## 🔑 Authentication Flow

### 1. Login Process
1. User enters username and password on the login page
2. Frontend sends POST request to `/api/v1/auth/login`
3. Backend validates credentials against environment variables
4. If valid, backend generates JWT token and returns it
5. Frontend stores token in localStorage
6. Frontend includes token in all subsequent API requests

### 2. Token Validation
- All API endpoints (except login) require a valid JWT token
- Token must be included in the Authorization header: `Bearer <token>`
- Tokens expire after 24 hours
- Invalid/expired tokens result in 401 Unauthorized response

### 3. Logout Process
- User clicks logout button
- Frontend removes token from localStorage
- User is redirected to login page

## 🛡️ Security Features

### JWT Token Structure
```json
{
  "username": "admin",
  "exp": 1234567890,
  "iat": 1234567890
}
```

### Token Security
- Tokens are signed with HMAC-SHA256
- Secret key is stored in environment variables
- Tokens include expiration time (24 hours)
- Automatic logout on token expiration

### API Protection
- All endpoints require authentication except:
  - `/health` - Health check
  - `/api/v1/auth/login` - Login endpoint
- CORS is configured to allow frontend access
- Error messages don't expose sensitive information

## 🔧 Configuration

### Environment Variables

#### Backend (.env)
```env
# Database
DATABASE_URL=postgresql://user:pass@host:5432/dbname
PORT=8080
NODE_ENV=development

# Authentication
AUTH_USERNAME=your_username
AUTH_PASSWORD=your_password
JWT_SECRET=your_jwt_secret_key_minimum_32_characters
```

#### Frontend (.env)
```env
REACT_APP_API_URL=http://localhost:8080
```

### Security Best Practices

1. **Strong Credentials**
   - Use complex usernames and passwords
   - Avoid default credentials in production

2. **JWT Secret**
   - Use a long, random string (minimum 32 characters)
   - Include letters, numbers, and special characters
   - Never share or commit to version control

3. **Environment Variables**
   - Never hardcode credentials in source code
   - Use different credentials for different environments
   - Keep .env files in .gitignore

4. **HTTPS**
   - Always use HTTPS in production
   - Tokens are vulnerable to interception over HTTP

## 🧪 Testing Authentication

### Backend API Testing

#### 1. Test Login Endpoint
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "secure123"
  }'
```

Expected response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": "admin"
}
```

#### 2. Test Protected Endpoint
```bash
# Without token (should fail)
curl http://localhost:8080/api/v1/items

# With token (should succeed)
curl -H "Authorization: Bearer <your-token>" \
  http://localhost:8080/api/v1/items
```

#### 3. Test Invalid Token
```bash
curl -H "Authorization: Bearer invalid-token" \
  http://localhost:8080/api/v1/items
```

Expected response:
```json
{
  "error": "Invalid or expired token"
}
```

### Frontend Testing

1. **Login Flow**
   - Open http://localhost:3000
   - Should redirect to login page
   - Enter correct credentials
   - Should redirect to dashboard

2. **Token Persistence**
   - Login successfully
   - Refresh the page
   - Should remain logged in

3. **Logout Flow**
   - Click logout button
   - Should redirect to login page
   - Token should be removed from localStorage

4. **Expired Token**
   - Wait for token to expire (24 hours)
   - Or manually set an expired token in localStorage
   - Should automatically logout

## 🚨 Troubleshooting

### Common Issues

#### 1. "Authorization header required"
- **Cause**: Missing or malformed Authorization header
- **Solution**: Ensure token is included as `Bearer <token>`

#### 2. "Invalid or expired token"
- **Cause**: Token is invalid, expired, or JWT secret mismatch
- **Solutions**:
  - Login again to get a new token
  - Check JWT_SECRET in backend .env
  - Verify token hasn't expired

#### 3. "Invalid credentials"
- **Cause**: Wrong username or password
- **Solutions**:
  - Check AUTH_USERNAME and AUTH_PASSWORD in backend .env
  - Ensure frontend is sending correct credentials

#### 4. CORS errors
- **Cause**: Frontend and backend on different domains
- **Solutions**:
  - Check REACT_APP_API_URL in frontend .env
  - Ensure backend CORS settings allow frontend origin

#### 5. "Network Error" on login
- **Cause**: Backend not running or wrong API URL
- **Solutions**:
  - Start backend server
  - Check REACT_APP_API_URL matches backend address

### Debug Steps

1. **Check Backend Logs**
   ```bash
   cd backend
   go run cmd/server/main.go
   ```

2. **Check Frontend Console**
   - Open browser developer tools
   - Look for network errors or console messages

3. **Verify Environment Variables**
   ```bash
   # Backend
   cat backend/.env
   
   # Frontend
   cat frontend/.env
   ```

4. **Test API Directly**
   ```bash
   # Health check
   curl http://localhost:8080/health
   
   # Login
   curl -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username": "admin", "password": "secure123"}'
   ```

## 🔄 Production Deployment

### Security Checklist

- [ ] Change default username and password
- [ ] Use strong, unique JWT secret
- [ ] Enable HTTPS/TLS
- [ ] Set secure CORS policies
- [ ] Use environment variables for all secrets
- [ ] Consider shorter token expiration times
- [ ] Set up proper logging and monitoring
- [ ] Regular security audits

### Environment Variables for Production

```env
# Backend
DATABASE_URL=postgresql://produser:strongpass@prod-db:5432/interview_prep
PORT=8080
NODE_ENV=production
AUTH_USERNAME=secure_admin_username
AUTH_PASSWORD=very_secure_password_with_special_chars_123!
JWT_SECRET=super_long_random_jwt_secret_for_production_environment_2024_abcdef1234567890

# Frontend
REACT_APP_API_URL=https://your-api-domain.com
```

## 📚 Additional Resources

- [JWT.io](https://jwt.io/) - JWT token debugger
- [OWASP JWT Security](https://cheatsheetseries.owasp.org/cheatsheets/JSON_Web_Token_for_Java_Cheat_Sheet.html)
- [Go JWT Library Documentation](https://pkg.go.dev/github.com/golang-jwt/jwt/v4)

## 🆘 Support

If you encounter issues:

1. Check this guide first
2. Review the troubleshooting section
3. Check the main README.md
4. Look at the example environment files
5. Test with the provided curl commands

Remember: **Never commit .env files or share JWT secrets publicly!** 