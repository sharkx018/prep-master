# Multi-User Authentication Guide

This system now supports multiple users using comma-separated environment variables. You can share the system with your friends by adding their credentials to the environment variables.

## Setup

### Option 1: Multi-User Authentication (Recommended)

Set the following environment variables with comma-separated values:

```bash
# Comma-separated usernames
AUTH_USERS=admin,john,jane,bob

# Comma-separated passwords (must match the order of usernames)
AUTH_PASSWORDS=admin123,john_pass,jane_pass,bob_pass
```

**Important:** 
- Each username at index `i` corresponds to the password at index `i`
- Make sure you have the same number of usernames and passwords
- Spaces around commas are automatically trimmed

### Option 2: Single User Authentication (Legacy)

If you prefer to use single user authentication or want backward compatibility:

```bash
AUTH_USERNAME=admin
AUTH_PASSWORD=password
```

## How it Works

1. **Both Methods Work Together**: The system now supports BOTH single-user and multi-user authentication simultaneously
2. **Multi-User Authentication**: If `AUTH_USERS` and `AUTH_PASSWORDS` are set, all users in the comma-separated lists can log in
3. **Single-User Authentication**: If `AUTH_USERNAME` and `AUTH_PASSWORD` are set, that user can also log in
4. **Combined Authentication**: You can set both types - all users from both methods will be able to log in
5. **Validation**: During login, the system:
   - First checks if the credentials match any user in the multi-user list (`AUTH_USERS`/`AUTH_PASSWORDS`)
   - Then checks if the credentials match the single-user credentials (`AUTH_USERNAME`/`AUTH_PASSWORD`)
   - If either check passes, the login succeeds

### Perfect Backward Compatibility

✅ **Your existing setup will continue to work without any changes!**

- If you only have `AUTH_USERNAME` and `AUTH_PASSWORD` set, the system works exactly as before
- You can add multi-user credentials (`AUTH_USERS`/`AUTH_PASSWORDS`) and your existing single-user credentials will still work
- Both authentication methods work together - no conflicts, no migration needed

## Example Usage

### Environment File (.env)

```bash
# Database Configuration
DATABASE_URL=postgresql://username:password@localhost:5432/interview_prep
PORT=3000
NODE_ENV=development

# Multi-User Authentication
AUTH_USERS=admin,alice,bob,charlie
AUTH_PASSWORDS=secure123,alice_pass,bob_secret,charlie_key

JWT_SECRET=your_jwt_secret_key_here
```

### Login Examples

With the above configuration, these users can log in:
- Username: `admin`, Password: `secure123`
- Username: `alice`, Password: `alice_pass`
- Username: `bob`, Password: `bob_secret`
- Username: `charlie`, Password: `charlie_key`

## Security Notes

1. **Use Strong Passwords**: Make sure each user has a strong, unique password
2. **JWT Secret**: Use a strong JWT secret key
3. **Environment Security**: Keep your environment files secure and never commit them to version control
4. **User Management**: Currently, users are managed via environment variables. For production use, consider implementing a proper user management system with a database

## Testing

The system includes comprehensive tests to ensure the multi-user authentication works correctly:

```bash
cd backend
go test -v ./internal/config/
```

## Migration from Single User

If you're currently using single-user authentication, the system will continue to work without any changes. To migrate to multi-user:

1. Set the `AUTH_USERS` and `AUTH_PASSWORDS` environment variables
2. Include your current user in the new lists
3. The system will automatically use multi-user authentication

### Scenario 1: Keep Using Single User (No Changes Required)
```bash
# Your existing setup - continues to work exactly as before
AUTH_USERNAME=admin
AUTH_PASSWORD=mypassword
```

### Scenario 2: Add Multi-User While Keeping Single-User
```bash
# Your existing credentials (still work!)
AUTH_USERNAME=admin
AUTH_PASSWORD=mypassword

# Add multi-user credentials (both work together!)
AUTH_USERS=friend1,friend2,friend3
AUTH_PASSWORDS=friend1_pass,friend2_pass,friend3_pass

# Now these users can all log in:
# - admin with mypassword (from single-user)
# - friend1 with friend1_pass (from multi-user)
# - friend2 with friend2_pass (from multi-user)
# - friend3 with friend3_pass (from multi-user)
```

### Scenario 3: Multi-User Only
```bash
# Only set multi-user variables
AUTH_USERS=admin,friend1,friend2
AUTH_PASSWORDS=mypassword,friend1_pass,friend2_pass

# AUTH_USERNAME and AUTH_PASSWORD not needed
```

### Scenario 4: Same User, Multiple Passwords
```bash
# If the same username appears in both methods, both passwords work!
AUTH_USERNAME=admin
AUTH_PASSWORD=oldpassword

AUTH_USERS=admin,friend1
AUTH_PASSWORDS=newpassword,friend1_pass

# Admin can log in with EITHER "oldpassword" OR "newpassword"
# Friend1 can log in with "friend1_pass"
``` 