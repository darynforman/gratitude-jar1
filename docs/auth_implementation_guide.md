# Authentication and Authorization Implementation Guide

This guide provides detailed implementation instructions for enhancing the authentication and authorization system in the Gratitude Jar application.

## 1. Session Security Enhancements

### 1.1 Replace Hardcoded Session Secret

**File**: `internal/session/manager.go`

**Current Code**:
```go
func init() {
    Manager = sessions.New([]byte("your-secret-key"))
    Manager.Lifetime = 24 * time.Hour
    Manager.Secure = false // Set to true in production with HTTPS
}
```

**Implementation**:

1. Update the session manager to use an environment variable:

```go
func init() {
    // Get session secret from environment variable or use a default in development
    secretKey := os.Getenv("SESSION_SECRET")
    if secretKey == "" {
        // Only use this default in development
        secretKey = "dev-session-secret-replace-in-production"
        log.Println("WARNING: Using default session secret. Set SESSION_SECRET environment variable in production.")
    }
    
    Manager = sessions.New([]byte(secretKey))
    Manager.Lifetime = 24 * time.Hour
    
    // Enable secure cookies in production
    secureMode := os.Getenv("SECURE_COOKIES") == "true"
    Manager.Secure = secureMode
    
    // Always set HttpOnly to prevent JavaScript access
    Manager.HttpOnly = true
    
    // Set SameSite attribute to prevent CSRF
    Manager.SameSite = http.SameSiteStrictMode
}
```

2. Update the README to document the new environment variables:

```markdown
## Environment Variables

- `SESSION_SECRET`: Secret key for session encryption (required in production)
- `SECURE_COOKIES`: Set to "true" to enable secure cookies (recommended in production)
```

### 1.2 Implement Session Timeout and Renewal

**File**: Create new file `internal/auth/session_middleware.go`

**Implementation**:

```go
package auth

import (
    "net/http"
    "time"

    "github.com/darynforman/gratitude-jar1/internal/session"
)

// SessionTimeoutMiddleware checks if the session has timed out due to inactivity
func SessionTimeoutMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Skip for unauthenticated users
        userID := session.Manager.GetInt(r, "userID")
        if userID == 0 {
            next.ServeHTTP(w, r)
            return
        }

        // Get last activity time from session
        lastActivityVal := session.Manager.Get(r, "last_activity")
        
        // Check if last_activity exists and is a valid time
        var lastActivity time.Time
        if lastActivityVal != nil {
            if lastTime, ok := lastActivityVal.(time.Time); ok {
                lastActivity = lastTime
            }
        }

        now := time.Now()
        
        // If session has been inactive for too long (30 minutes), log out
        if !lastActivity.IsZero() && now.Sub(lastActivity) > 30*time.Minute {
            session.LogoutUser(w, r)
            http.Redirect(w, r, "/user/login", http.StatusSeeOther)
            return
        }
        
        // Update last activity time
        session.Manager.Put(r, "last_activity", now)
        
        next.ServeHTTP(w, r)
    })
}
```

**File**: Update `cmd/web/routes.go` to add the middleware

```go
// Add the import
import (
    "github.com/darynforman/gratitude-jar1/internal/auth"
)

// In the routes function, add the middleware to the chain
func routes() http.Handler {
    // Existing code...
    
    // Chain middleware in the correct order
    handler := LoggingMiddleware(mux)
    handler = RateLimitMiddleware(handler)
    handler = SecureHeadersMiddleware(handler)
    handler = auth.SessionTimeoutMiddleware(handler) // Add this line
    handler = RecoverPanicMiddleware(handler)
    
    return handler
}
```

## 2. User-Specific Data Access Verification

### 2.1 Create Resource Ownership Middleware

**File**: Create new file `internal/auth/ownership_middleware.go`

**Implementation**:

```go
package auth

import (
    "net/http"
    "strconv"
    "strings"

    "github.com/darynforman/gratitude-jar1/internal/config"
    "github.com/darynforman/gratitude-jar1/internal/data"
    "github.com/darynforman/gratitude-jar1/internal/session"
)

// RequireOwnership ensures the user owns the requested resource
func RequireOwnership(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Get user ID from session
        userID := session.Manager.GetInt(r, "userID")
        if userID == 0 {
            http.Redirect(w, r, "/user/login", http.StatusSeeOther)
            return
        }

        // Extract resource ID from URL
        // Assuming URLs like /gratitude/edit/123 or /notes/123
        parts := strings.Split(r.URL.Path, "/")
        if len(parts) < 3 {
            next(w, r)
            return
        }

        resourceIDStr := parts[len(parts)-1]
        resourceID, err := strconv.Atoi(resourceIDStr)
        if err != nil {
            // If we can't parse the ID, just continue
            next(w, r)
            return
        }

        // Check if this is a gratitude note
        if strings.Contains(r.URL.Path, "/gratitude/") || strings.Contains(r.URL.Path, "/notes/") {
            // Get the note
            gratitudeModel := data.NewGratitudeModel(config.DB)
            note, err := gratitudeModel.Get(resourceID)
            if err != nil || note == nil {
                http.NotFound(w, r)
                return
            }

            // Check ownership
            if note.UserID != userID {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }
        }

        // If we get here, the user owns the resource or it's not a resource that needs checking
        next(w, r)
    })
}
```

### 2.2 Apply Ownership Checks to Routes

**File**: Update `cmd/web/routes.go`

**Implementation**:

```go
// Update route definitions to include ownership checks
mux.Handle("/gratitude/edit/", auth.RequireLogin(auth.RequireOwnership(http.HandlerFunc(getNoteForEdit))))
mux.Handle("/notes/", auth.RequireLogin(auth.RequireOwnership(http.HandlerFunc(updateGratitude))))
```

## 3. CSRF Protection

### 3.1 Create CSRF Package

**File**: Create new file `internal/csrf/csrf.go`

**Implementation**:

```go
package csrf

import (
    "crypto/rand"
    "encoding/base64"
    "errors"
    "fmt"
    "net/http"
    "sync"
    "time"
)

var (
    // tokens stores valid CSRF tokens
    tokens = struct {
        sync.RWMutex
        m map[string]time.Time
    }{m: make(map[string]time.Time)}

    // ErrInvalidToken is returned when a CSRF token is invalid
    ErrInvalidToken = errors.New("invalid CSRF token")
)

// GenerateToken creates a new CSRF token
func GenerateToken() (string, error) {
    // Generate 32 bytes of random data
    b := make([]byte, 32)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }

    // Convert to base64
    token := base64.StdEncoding.EncodeToString(b)

    // Store token with expiration time (1 hour)
    tokens.Lock()
    tokens.m[token] = time.Now().Add(1 * time.Hour)
    tokens.Unlock()

    return token, nil
}

// ValidateToken checks if a token is valid
func ValidateToken(token string) bool {
    tokens.RLock()
    expiry, exists := tokens.m[token]
    tokens.RUnlock()

    // Check if token exists and hasn't expired
    if !exists || time.Now().After(expiry) {
        return false
    }

    return true
}

// CleanupExpiredTokens removes expired tokens
func CleanupExpiredTokens() {
    tokens.Lock()
    defer tokens.Unlock()

    now := time.Now()
    for token, expiry := range tokens.m {
        if now.After(expiry) {
            delete(tokens.m, token)
        }
    }
}

// StartCleanupRoutine starts a goroutine to periodically clean up expired tokens
func StartCleanupRoutine() {
    go func() {
        ticker := time.NewTicker(15 * time.Minute)
        defer ticker.Stop()

        for range ticker.C {
            CleanupExpiredTokens()
        }
    }()
}

// Middleware creates middleware to check CSRF tokens
func Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Skip for GET, HEAD, OPTIONS, TRACE
        if r.Method == http.MethodGet || r.Method == http.MethodHead ||
           r.Method == http.MethodOptions || r.Method == http.MethodTrace {
            next.ServeHTTP(w, r)
            return
        }

        // Check CSRF token
        token := r.Header.Get("X-CSRF-Token")
        if token == "" {
            token = r.FormValue("csrf_token")
        }

        if !ValidateToken(token) {
            http.Error(w, "Invalid CSRF Token", http.StatusForbidden)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

### 3.2 Initialize CSRF Protection

**File**: Update `cmd/web/main.go`

**Implementation**:

```go
// Add import
import (
    "github.com/darynforman/gratitude-jar1/internal/csrf"
)

// In main function, start the CSRF cleanup routine
func main() {
    // Existing code...
    
    // Start CSRF token cleanup routine
    csrf.StartCleanupRoutine()
    
    // Start the server
    startServer(app)
}
```

### 3.3 Add CSRF Token to Templates

**File**: Update `cmd/web/render.go`

**Implementation**:

```go
// Add import
import (
    "github.com/darynforman/gratitude-jar1/internal/csrf"
)

// Update render function to include CSRF token
func render(w http.ResponseWriter, r *http.Request, name string, data PageData) {
    // Get the template from the cache
    tmpl, err := getTemplate(name)
    if err != nil {
        log.Printf("Template %s not found in cache: %v", name, err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Generate CSRF token for forms
    csrfToken, err := csrf.GenerateToken()
    if err != nil {
        log.Printf("Error generating CSRF token: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Add session data to the template data
    userID := session.Manager.GetInt(r, "userID")
    role := session.Manager.GetString(r, "role")
    flash := session.Manager.PopString(r, "flash")

    // Create a new data struct that includes session data
    templateData := struct {
        PageData
        IsAuthenticated bool
        UserID          int
        UserRole        string
        Flash           string
        CurrentYear     int
        CSRFToken       string
    }{
        PageData:        data,
        IsAuthenticated: userID > 0,
        UserID:          userID,
        UserRole:        role,
        Flash:           flash,
        CurrentYear:     time.Now().Year(),
        CSRFToken:       csrfToken,
    }

    // Execute the template
    err = tmpl.ExecuteTemplate(w, "base", templateData)
    if err != nil {
        log.Printf("Error executing template: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}
```

### 3.4 Update Form Templates

**File**: Update all form templates to include CSRF token

Example for a form template:

```html
<form method="POST" action="/some-endpoint">
    <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
    <!-- Rest of the form -->
</form>
```

### 3.5 Apply CSRF Middleware

**File**: Update `cmd/web/routes.go`

**Implementation**:

```go
// Add import
import (
    "github.com/darynforman/gratitude-jar1/internal/csrf"
)

// In the routes function, add the middleware to the chain
func routes() http.Handler {
    // Existing code...
    
    // Chain middleware in the correct order
    handler := LoggingMiddleware(mux)
    handler = RateLimitMiddleware(handler)
    handler = csrf.Middleware(handler) // Add this line
    handler = SecureHeadersMiddleware(handler)
    handler = auth.SessionTimeoutMiddleware(handler)
    handler = RecoverPanicMiddleware(handler)
    
    return handler
}
```

## 4. Rate Limiting for Authentication

### 4.1 Enhance Rate Limiting for Auth Endpoints

**File**: Update `cmd/web/middleware_ratelimit.go`

**Implementation**:

```go
// Add imports if needed
import (
    "strings"
    "sync"
    "time"
)

// Create a stricter limiter for authentication endpoints
var (
    authLimiters = struct {
        sync.RWMutex
        m map[string]*rate.Limiter
    }{m: make(map[string]*rate.Limiter)}
)

// GetAuthLimiter returns a rate limiter for authentication endpoints
func GetAuthLimiter(ip string) *rate.Limiter {
    authLimiters.RLock()
    limiter, exists := authLimiters.m[ip]
    authLimiters.RUnlock()

    if !exists {
        // Create a new limiter: 5 requests per minute with burst of 5
        limiter = rate.NewLimiter(rate.Limit(5)/60, 5)
        
        authLimiters.Lock()
        authLimiters.m[ip] = limiter
        authLimiters.Unlock()
    }

    return limiter
}

// Update RateLimitMiddleware to apply stricter limits to auth endpoints
func RateLimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Skip rate limiting for static files
        if strings.HasPrefix(r.URL.Path, "/static/") {
            next.ServeHTTP(w, r)
            return
        }

        // Get client IP
        ip := getClientIP(r)
        
        // Check if this is an auth endpoint
        isAuthEndpoint := r.URL.Path == "/user/login" || r.URL.Path == "/register"
        
        var allow bool
        if isAuthEndpoint {
            // Use stricter limiter for auth endpoints
            limiter := GetAuthLimiter(ip)
            allow = limiter.Allow()
        } else {
            // Use regular limiter for other endpoints
            limiter := globalLimiter.GetLimiter(ip)
            allow = limiter.Allow()
        }
        
        if !allow {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

## 5. Security Headers Enhancement

### 5.1 Add Content Security Policy

**File**: Update `cmd/web/middleware.go`

**Implementation**:

```go
// Update SecureHeadersMiddleware to add Content Security Policy
func SecureHeadersMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Add security headers to enhance protection against common attacks
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "deny")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        
        // Add Content Security Policy
        // Adjust this policy based on your application's needs
        w.Header().Set("Content-Security-Policy", 
            "default-src 'self'; "+
            "script-src 'self' 'unsafe-inline'; "+
            "style-src 'self' 'unsafe-inline'; "+
            "img-src 'self' data:; "+
            "connect-src 'self'; "+
            "font-src 'self'; "+
            "object-src 'none'; "+
            "frame-src 'none'; "+
            "base-uri 'self'; "+
            "form-action 'self'")

        // Pass request to the next handler
        next.ServeHTTP(w, r)
    })
}
```

## Next Steps

After implementing these core security enhancements, consider:

1. Writing tests for each new component
2. Implementing the optional features (account recovery, account lockout)
3. Conducting security testing
4. Updating deployment documentation

Remember to update the task tracker as you complete each task.
