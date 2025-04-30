# Authentication and Authorization Implementation Plan

## Overview

This document outlines the plan for implementing and enhancing authentication and authorization in the Gratitude Jar application. The goal is to ensure that users can only access their own data while maintaining a simple and secure authentication system.

## Current Implementation Assessment

The application already has several authentication and authorization components:

- User registration with password hashing (bcrypt)
- Login functionality with session management
- Protected routes using middleware
- Input validation for registration and login forms
- Basic security headers
- Database models for users and gratitude notes with user_id relationships

## Implementation Plan

### 1. Session Security Enhancements

**Status**: Needs improvement

**Current Implementation**:
```go
// In internal/session/manager.go
func init() {
    Manager = sessions.New([]byte("your-secret-key"))
    Manager.Lifetime = 24 * time.Hour
    Manager.Secure = false // Set to true in production with HTTPS
}
```

**Proposed Changes**:
- Replace hardcoded session secret with environment variable
- Enable Secure flag in production environments
- Add SameSite and HttpOnly flags
- Implement session expiration and renewal

**Implementation Steps**:
1. Modify `internal/session/manager.go` to load secret from environment
2. Add configuration for cookie security settings
3. Add session timeout and renewal logic

### 2. User-Specific Data Access Verification

**Status**: Partially implemented

**Current Implementation**:
- Some data access methods filter by user_id
- Need to verify all data access points enforce user-specific access

**Proposed Changes**:
- Audit all data access methods to ensure they filter by user_id
- Add middleware to verify resource ownership for all protected routes

**Implementation Steps**:
1. Review all model methods in `internal/data/` that retrieve data
2. Ensure all methods include user_id filtering
3. Create a resource ownership middleware

### 3. CSRF Protection

**Status**: Not implemented

**Current Implementation**:
- No CSRF protection currently exists

**Proposed Changes**:
- Add CSRF token generation and validation
- Include CSRF tokens in all forms
- Add middleware to verify CSRF tokens for state-changing operations

**Implementation Steps**:
1. Create CSRF package in `internal/csrf/`
2. Modify templates to include CSRF tokens
3. Add middleware to verify tokens

### 4. Rate Limiting for Authentication

**Status**: Basic rate limiting exists

**Current Implementation**:
- General rate limiting middleware exists
- No specific protection for authentication endpoints

**Proposed Changes**:
- Enhance rate limiting for authentication endpoints
- Implement IP-based and username-based rate limiting for login attempts

**Implementation Steps**:
1. Modify `cmd/web/middleware_ratelimit.go` to add stricter limits for auth endpoints
2. Add tracking for failed login attempts

### 5. Security Headers Enhancement

**Status**: Basic headers implemented

**Current Implementation**:
- Some security headers are set in `SecureHeadersMiddleware`

**Proposed Changes**:
- Add Content Security Policy
- Review and update existing security headers

**Implementation Steps**:
1. Update `SecureHeadersMiddleware` in `cmd/web/middleware.go`
2. Test headers with security scanning tools

### 6. Account Recovery (Optional)

**Status**: Not implemented

**Current Implementation**:
- No password reset functionality

**Proposed Changes**:
- Add simple password reset functionality
- Implement secure token generation and verification

**Implementation Steps**:
1. Create password reset token generation and storage
2. Add handlers for password reset request and confirmation
3. Create email templates for reset instructions

### 7. Account Lockout (Optional)

**Status**: Not implemented

**Current Implementation**:
- No account lockout mechanism

**Proposed Changes**:
- Track failed login attempts
- Temporarily lock accounts after multiple failed attempts

**Implementation Steps**:
1. Add failed login tracking to user model
2. Implement account lockout logic in login handler
3. Add unlock mechanism (time-based or manual)

## Testing Plan

1. **Unit Tests**:
   - Test password hashing and verification
   - Test session management functions
   - Test CSRF token generation and validation

2. **Integration Tests**:
   - Test authentication flow (register, login, logout)
   - Test protected route access
   - Test data access with different user sessions

3. **Security Testing**:
   - Test for common vulnerabilities (OWASP Top 10)
   - Test rate limiting effectiveness
   - Test CSRF protection

## Deployment Considerations

1. **Environment Variables**:
   - SESSION_SECRET: Strong random string for session encryption
   - CSRF_SECRET: Secret for CSRF token generation
   - SECURE_COOKIES: Boolean to enable secure cookies in production

2. **Database Updates**:
   - Add columns for tracking failed login attempts
   - Add table for password reset tokens (if implementing account recovery)

3. **Monitoring**:
   - Log authentication attempts and failures
   - Monitor for unusual activity patterns
