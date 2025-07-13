# 🔐 OAuth Process Explained

## 🎯 What is OAuth?

OAuth (Open Authorization) is an industry-standard protocol that allows users to grant third-party applications access to their information without sharing their passwords. Instead of giving your password to every app, OAuth lets you say "Yes, this app can access my basic info" through a secure, standardized process.

## 🌟 Why Use OAuth?

### For Users:
- **Security**: Never share passwords with third-party apps
- **Control**: Grant specific permissions and revoke access anytime
- **Convenience**: One-click login with existing accounts
- **Trust**: Login through providers you already trust (Google, Facebook, Apple)

### For Developers:
- **Security**: Don't store user passwords
- **User Experience**: Faster signup/login process
- **Trust**: Users more likely to sign up with familiar providers
- **Compliance**: Reduced responsibility for user data security

## 🔄 The OAuth Flow (Simplified)

```
1. User clicks "Login with Google" on your app
2. Your app redirects user to Google's login page
3. User logs in to Google (not your app!)
4. Google asks: "Allow Interview Prep App to access your basic info?"
5. User clicks "Yes"
6. Google redirects back to your app with a special code
7. Your app exchanges this code for an access token
8. Your app uses the token to get user's info from Google
9. Your app creates/logs in the user
```

## 🔍 Detailed OAuth 2.0 Flow

### Step 1: Authorization Request
```
User → Your App → OAuth Provider (Google/Facebook/Apple)

Your app redirects user to:
https://accounts.google.com/oauth/authorize?
  client_id=YOUR_CLIENT_ID&
  redirect_uri=YOUR_REDIRECT_URI&
  scope=email profile&
  response_type=code
```

### Step 2: User Authorization
```
User ← OAuth Provider

- User sees familiar login page (Google/Facebook/Apple)
- User enters their credentials
- Provider shows permission screen
- User grants/denies permission
```

### Step 3: Authorization Grant
```
User → Your App ← OAuth Provider

Provider redirects back to your app with:
https://yourapp.com/callback?code=AUTHORIZATION_CODE

Or if user denies:
https://yourapp.com/callback?error=access_denied
```

### Step 4: Access Token Request
```
Your App → OAuth Provider

POST https://oauth2.googleapis.com/token
{
  "client_id": "YOUR_CLIENT_ID",
  "client_secret": "YOUR_CLIENT_SECRET",
  "code": "AUTHORIZATION_CODE",
  "grant_type": "authorization_code",
  "redirect_uri": "YOUR_REDIRECT_URI"
}
```

### Step 5: Access Token Response
```
Your App ← OAuth Provider

{
  "access_token": "ya29.a0AfH6SMC...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "1//04...",
  "scope": "email profile"
}
```

### Step 6: Resource Access
```
Your App → OAuth Provider

GET https://www.googleapis.com/oauth2/v2/userinfo
Authorization: Bearer ya29.a0AfH6SMC...

Response:
{
  "id": "123456789",
  "email": "user@gmail.com",
  "name": "John Doe",
  "picture": "https://..."
}
```

## 🏗️ Implementation in Interview Prep App

### Frontend Flow

#### 1. User Clicks OAuth Button
```typescript
// frontend/src/pages/Login.tsx
const handleGoogleLogin = async () => {
  setIsLoading(true);
  setError('');
  
  const success = await loginWithGoogle();
  if (!success) {
    setError('Google login failed. Please try again.');
  }
  
  setIsLoading(false);
};
```

#### 2. Initialize OAuth Client
```typescript
// frontend/src/contexts/AuthContext.tsx
const loginWithGoogle = async (): Promise<boolean> => {
  try {
    // Check if Google SDK is loaded
    if (!window.google) {
      console.error('Google OAuth not loaded');
      return false;
    }

    return new Promise((resolve) => {
      // Initialize Google OAuth client
      window.google.accounts.oauth2.initTokenClient({
        client_id: process.env.REACT_APP_GOOGLE_CLIENT_ID || '',
        scope: 'email profile',
        callback: async (response: any) => {
          if (response.access_token) {
            const success = await handleOAuthLogin('google', response.access_token);
            resolve(success);
          } else {
            resolve(false);
          }
        },
      }).requestAccessToken();
    });
  } catch (error) {
    console.error('Google login error:', error);
    return false;
  }
};
```

#### 3. Send Token to Backend
```typescript
// frontend/src/contexts/AuthContext.tsx
const handleOAuthLogin = async (provider: string, accessToken: string): Promise<boolean> => {
  try {
    const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';
    const response = await fetch(`${API_BASE_URL}/auth/oauth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        provider,
        access_token: accessToken,
      }),
    });

    if (response.ok) {
      const data = await response.json();
      
      setUser(data.user);
      localStorage.setItem('auth_token', data.token);
      localStorage.setItem('auth_user', JSON.stringify(data.user));
      
      return true;
    }
    return false;
  } catch (error) {
    console.error('OAuth login error:', error);
    return false;
  }
};
```

### Backend Flow

#### 1. Receive OAuth Request
```go
// backend/internal/handlers/auth_handler.go
func (h *AuthHandler) OAuthLogin(c *gin.Context) {
	var req models.OAuthLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	log.Printf("OAuth login attempt with provider: %s", req.Provider)

	// Authenticate user with OAuth
	user, err := h.userService.LoginWithOAuth(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT tokens
	token, err := h.generateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User:  user,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
}
```

#### 2. Validate OAuth Token
```go
// backend/internal/services/user_service.go
func (s *UserService) LoginWithOAuth(req *models.OAuthLoginRequest) (*models.User, error) {
	// Validate OAuth token and get user info
	userInfo, err := s.validateOAuthToken(req)
	if err != nil {
		return nil, fmt.Errorf("invalid OAuth token: %w", err)
	}

	// Try to find existing user by provider ID
	user, err := s.userRepo.GetByProviderID(req.Provider, userInfo.ProviderID)
	if err == nil {
		// User exists, update last login
		err = s.userRepo.UpdateLastLogin(user.ID)
		return user, nil
	}

	// Create new user
	user = &models.User{
		Email:        userInfo.Email,
		Name:         userInfo.Name,
		Avatar:       userInfo.Avatar,
		AuthProvider: req.Provider,
		ProviderID:   userInfo.ProviderID,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
```

#### 3. Validate Google Token
```go
// backend/internal/services/user_service.go
func (s *UserService) validateGoogleToken(token string) (*OAuthUserInfo, error) {
	// Make request to Google's userinfo endpoint
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", token)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to validate Google token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid Google token")
	}

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	err = json.NewDecoder(resp.Body).Decode(&googleUser)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Google user info: %w", err)
	}

	return &OAuthUserInfo{
		ProviderID: googleUser.ID,
		Email:      googleUser.Email,
		Name:       googleUser.Name,
		Avatar:     googleUser.Picture,
	}, nil
}
```

## 🔐 Security Considerations

### 1. Token Validation
- **Always validate tokens server-side**
- Never trust tokens from the frontend
- Check token expiration
- Verify token with the OAuth provider

### 2. Secure Storage
- Store tokens securely (encrypted if possible)
- Use HTTPS in production
- Implement token refresh mechanisms
- Revoke tokens when users logout

### 3. Scope Management
- Request minimal necessary permissions
- Clearly explain what data you'll access
- Allow users to revoke permissions
- Regular security audits

## 🌍 Provider-Specific Details

### Google OAuth
```
Authorization URL: https://accounts.google.com/oauth/authorize
Token URL: https://oauth2.googleapis.com/token
User Info URL: https://www.googleapis.com/oauth2/v2/userinfo
Scopes: email, profile, openid
```

### Facebook OAuth
```
Authorization URL: https://www.facebook.com/v18.0/dialog/oauth
Token URL: https://graph.facebook.com/v18.0/oauth/access_token
User Info URL: https://graph.facebook.com/me?fields=id,email,name,picture
Scopes: email, public_profile
```

### Apple Sign In
```
Authorization URL: https://appleid.apple.com/auth/authorize
Token URL: https://appleid.apple.com/auth/token
User Info: Included in JWT token
Scopes: name, email
```

## 🚨 Common Issues & Solutions

### 1. "Invalid Client ID"
- **Cause**: Wrong client ID in environment variables
- **Solution**: Verify client ID matches OAuth provider settings

### 2. "Redirect URI Mismatch"
- **Cause**: Redirect URI doesn't match provider configuration
- **Solution**: Add correct URIs to OAuth provider settings

### 3. "Scope Not Granted"
- **Cause**: User denied permissions or app not approved
- **Solution**: Request minimal scopes, get app approved for production

### 4. "Token Expired"
- **Cause**: Access token has expired
- **Solution**: Implement token refresh or re-authenticate

### 5. "CORS Errors"
- **Cause**: Cross-origin requests blocked
- **Solution**: Configure CORS properly, use HTTPS

## 🔄 Token Refresh Flow

```
1. Access token expires (usually 1 hour)
2. App detects 401 Unauthorized response
3. App uses refresh token to get new access token
4. Continue using new access token
5. If refresh token expires, user must re-authenticate
```

## 📱 Mobile vs Web OAuth

### Web (Your Implementation)
- Uses JavaScript SDK
- Popup or redirect flow
- Tokens handled in browser

### Mobile
- Uses native SDKs
- Deep linking for redirects
- More secure token storage

## 🎯 Best Practices

### 1. User Experience
- Clear permission explanations
- Graceful error handling
- Loading states during OAuth flow
- Fallback to email/password

### 2. Security
- Validate all tokens server-side
- Use HTTPS everywhere
- Implement CSRF protection
- Regular security updates

### 3. Privacy
- Request minimal permissions
- Clear privacy policy
- Allow data deletion
- Transparent data usage

## 🔧 Testing OAuth

### Development
```bash
# Use localhost URLs
http://localhost:3000

# Test with personal accounts
# Use provider test environments
```

### Production
```bash
# Use HTTPS URLs
https://yourdomain.com

# Get apps approved by providers
# Test with real user accounts
```

## 📚 Additional Resources

- [OAuth 2.0 RFC](https://tools.ietf.org/html/rfc6749)
- [Google OAuth Documentation](https://developers.google.com/identity/protocols/oauth2)
- [Facebook Login Documentation](https://developers.facebook.com/docs/facebook-login/)
- [Apple Sign In Documentation](https://developer.apple.com/sign-in-with-apple/)
- [OAuth Security Best Practices](https://tools.ietf.org/html/draft-ietf-oauth-security-topics)

The OAuth implementation in your Interview Prep App follows industry best practices and provides a secure, user-friendly authentication experience! 🚀 