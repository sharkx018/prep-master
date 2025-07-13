# 🔐 OAuth Authentication Setup Guide

This guide will help you set up OAuth authentication with Google, Facebook, and Apple for the Interview Prep App.

## 🚀 Overview

The app supports four authentication methods:
- **Email/Password**: Traditional email and password authentication
- **Google OAuth**: Sign in with Google account
- **Facebook OAuth**: Sign in with Facebook account
- **Apple Sign In**: Sign in with Apple ID

## 📋 Prerequisites

- Node.js 16+ and npm
- Go 1.21+
- PostgreSQL database
- Developer accounts for OAuth providers

## 🔧 Backend Setup

### 1. Environment Variables

Add these to your `backend/.env` file:

```env
# Database
DATABASE_URL=postgresql://username:password@localhost:5432/interview_prep
PORT=8080
NODE_ENV=development

# JWT Secret (generate a secure random string)
JWT_SECRET=your_super_secure_jwt_secret_key_here

# OAuth Configuration (optional - for additional validation)
GOOGLE_CLIENT_SECRET=your-google-client-secret
FACEBOOK_APP_SECRET=your-facebook-app-secret
APPLE_PRIVATE_KEY_ID=your-apple-private-key-id
APPLE_TEAM_ID=your-apple-team-id
```

### 2. Database Migration

The app will automatically create the necessary tables:
- `users` - User profiles and authentication info
- `refresh_tokens` - Secure refresh token storage
- `user_progress` - User progress tracking
- `user_stats` - User statistics

## 🌐 Frontend Setup

### 1. Environment Variables

Add these to your `frontend/.env` file:

```env
# API Configuration
REACT_APP_API_URL=http://localhost:8080

# OAuth Configuration
REACT_APP_GOOGLE_CLIENT_ID=your-google-client-id
REACT_APP_FACEBOOK_APP_ID=your-facebook-app-id
REACT_APP_APPLE_CLIENT_ID=your-apple-client-id
```

## 🔑 OAuth Provider Setup

### Google OAuth Setup

1. **Go to Google Cloud Console**
   - Visit [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project or select existing one

2. **Enable Google+ API**
   - Navigate to "APIs & Services" > "Library"
   - Search for "Google+ API" and enable it

3. **Create OAuth Credentials**
   - Go to "APIs & Services" > "Credentials"
   - Click "Create Credentials" > "OAuth client ID"
   - Choose "Web application"
   - Add authorized origins:
     - `http://localhost:3000` (development)
     - `https://yourdomain.com` (production)
   - Add authorized redirect URIs:
     - `http://localhost:3000` (development)
     - `https://yourdomain.com` (production)

4. **Copy Client ID**
   - Copy the Client ID to `REACT_APP_GOOGLE_CLIENT_ID`

### Facebook OAuth Setup

1. **Go to Facebook Developers**
   - Visit [Facebook Developers](https://developers.facebook.com/)
   - Create a new app or select existing one

2. **Add Facebook Login Product**
   - In your app dashboard, click "Add Product"
   - Find "Facebook Login" and click "Set Up"

3. **Configure Facebook Login**
   - Go to Facebook Login > Settings
   - Add Valid OAuth Redirect URIs:
     - `http://localhost:3000` (development)
     - `https://yourdomain.com` (production)

4. **Copy App ID**
   - From App Dashboard, copy App ID to `REACT_APP_FACEBOOK_APP_ID`

5. **App Review (for production)**
   - Submit your app for review to use with real users
   - For development, you can add test users

### Apple Sign In Setup

1. **Apple Developer Account**
   - You need a paid Apple Developer account ($99/year)
   - Visit [Apple Developer](https://developer.apple.com/)

2. **Create App ID**
   - Go to Certificates, Identifiers & Profiles
   - Create a new App ID
   - Enable "Sign In with Apple" capability

3. **Create Service ID**
   - Create a new Services ID
   - Enable "Sign In with Apple"
   - Configure domains and redirect URLs:
     - `localhost:3000` (development)
     - `yourdomain.com` (production)

4. **Create Private Key**
   - Create a new Key with "Sign In with Apple" enabled
   - Download the private key file

5. **Configure Environment**
   - Copy Service ID to `REACT_APP_APPLE_CLIENT_ID`
   - Add private key info to backend environment

## 🚀 Running the Application

### 1. Start Backend

```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

### 2. Start Frontend

```bash
cd frontend
npm install
npm start
```

### 3. Access the App

- Open http://localhost:3000
- You'll see the login page with all authentication options

## 🔒 Security Considerations

### Production Setup

1. **HTTPS Required**
   - OAuth providers require HTTPS in production
   - Use SSL certificates for your domain

2. **Environment Variables**
   - Never commit OAuth secrets to version control
   - Use environment variables or secret management services

3. **CORS Configuration**
   - Configure CORS properly for your production domain
   - Update backend CORS settings

4. **Token Security**
   - JWT tokens expire after 24 hours
   - Refresh tokens expire after 7 days
   - Tokens are stored securely in localStorage

### Security Best Practices

1. **Validate Tokens**
   - All OAuth tokens are validated server-side
   - Invalid tokens are rejected

2. **User Data Protection**
   - Passwords are hashed with bcrypt
   - Sensitive data is never logged

3. **Rate Limiting**
   - Consider implementing rate limiting for login attempts
   - Monitor for suspicious activity

## 🧪 Testing OAuth Integration

### Testing Google OAuth

1. **Development Testing**
   - Use your personal Google account
   - Test with localhost:3000

2. **Production Testing**
   - Add test users in Google Cloud Console
   - Test with production domain

### Testing Facebook OAuth

1. **Development Testing**
   - Add test users in Facebook App Dashboard
   - Test with localhost:3000

2. **Production Testing**
   - Submit app for review
   - Test with approved domain

### Testing Apple Sign In

1. **Development Testing**
   - Use Apple ID for testing
   - Test with localhost:3000

2. **Production Testing**
   - Verify domain ownership
   - Test with production domain

## 🚨 Troubleshooting

### Common Issues

1. **"OAuth provider not loaded"**
   - Check if SDK scripts are loaded in index.html
   - Verify internet connection

2. **"Invalid client ID"**
   - Verify client ID in environment variables
   - Check OAuth provider configuration

3. **"Unauthorized domain"**
   - Add your domain to OAuth provider settings
   - Check redirect URI configuration

4. **"Token validation failed"**
   - Check if token is expired
   - Verify OAuth provider API access

### Debug Steps

1. **Check Browser Console**
   - Look for JavaScript errors
   - Verify OAuth SDK loading

2. **Check Network Tab**
   - Monitor API requests
   - Check for CORS errors

3. **Check Backend Logs**
   - Monitor server logs for errors
   - Check database connections

## 📚 Additional Resources

- [Google OAuth Documentation](https://developers.google.com/identity/protocols/oauth2)
- [Facebook Login Documentation](https://developers.facebook.com/docs/facebook-login/)
- [Apple Sign In Documentation](https://developer.apple.com/sign-in-with-apple/)
- [JWT Best Practices](https://auth0.com/blog/a-look-at-the-latest-draft-for-jwt-bcp/)

## 🔄 Migration from Basic Auth

If you're migrating from the basic username/password system:

1. **Backup Data**
   - Export existing user data
   - Backup database

2. **Run Migrations**
   - The new system will create user tables
   - Old auth system remains as fallback

3. **Update Environment**
   - Add OAuth configuration
   - Update frontend environment

4. **Test Thoroughly**
   - Test all authentication methods
   - Verify user data integrity

## 🎯 Next Steps

1. **Set up OAuth providers** following the guides above
2. **Configure environment variables** for all providers
3. **Test authentication flows** in development
4. **Deploy to production** with proper HTTPS setup
5. **Monitor and maintain** OAuth integrations

The authentication system is now ready for production use with enterprise-grade security and user experience! 