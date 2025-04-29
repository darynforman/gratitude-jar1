package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/darynforman/gratitude-jar1/internal/auth"
	"github.com/darynforman/gratitude-jar1/internal/data"
	"github.com/darynforman/gratitude-jar1/internal/session"
)

type loginForm struct {
	Username string
	Password string
}

type registerForm struct {
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
}

// login handles user login (GET shows form, POST processes login).
func (app *application) login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := PageData{
			Title: "Login",
		}
		render(w, r, "login.tmpl", data)
		return
	}

	if r.Method != http.MethodPost {
		app.errorResponse(w, r, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var form loginForm
	err := r.ParseForm()
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid form data")
		return
	}

	form.Username = r.PostForm.Get("username")
	form.Password = r.PostForm.Get("password")

	// Basic validation
	if form.Username == "" || form.Password == "" {
		app.errorResponse(w, r, http.StatusBadRequest, "Username and password are required")
		return
	}

	// Get user from database
	user, err := app.models.Users.GetByUsername(form.Username)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, "Database error")
		return
	}
	if user == nil {
		app.errorResponse(w, r, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Check password
	err = auth.CheckPassword(form.Password, user.PasswordHash)
	if err != nil {
		app.errorResponse(w, r, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	log.Printf("[Login] User authenticated successfully - ID: %d, Role: %s", user.ID, user.Role)

	// Set session values
	session.Manager.Put(r, "userID", user.ID)
	session.Manager.Put(r, "role", user.Role)
	session.Manager.Put(r, "flash", "Successfully logged in!")
	log.Printf("[Login] Session values set - userID: %d, role: %s", user.ID, user.Role)

	// Verify session values were set
	verifyUserID := session.Manager.GetInt(r, "userID")
	verifyRole := session.Manager.GetString(r, "role")
	log.Printf("[Login] Verifying session values - userID: %d, role: %s", verifyUserID, verifyRole)

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// Return the updated navigation bar
		data := PageData{
			Title:           "Navigation",
			IsAuthenticated: true,
			UserRole:        user.Role,
		}
		render(w, r, "partials/nav.tmpl", data)
		return
	}

	// For regular requests, redirect to home page
	log.Printf("[Login] Redirecting to home page")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// logoutHandler logs out the user by destroying the session.
func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Destroy the session
	session.Manager.Destroy(r)

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// Set HX-Redirect header to redirect to home page
		w.Header().Set("HX-Redirect", "/")
		// Return updated navigation
		data := PageData{
			Title:           "Navigation",
			IsAuthenticated: false,
		}
		render(w, r, "partials/nav.tmpl", data)
		return
	}

	// For regular requests, redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorResponse(w, r, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var form registerForm
	err := r.ParseForm()
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid form data")
		return
	}

	form.Username = r.PostForm.Get("username")
	form.Email = r.PostForm.Get("email")
	form.Password = r.PostForm.Get("password")
	form.ConfirmPassword = r.PostForm.Get("confirm_password")

	// Basic validation
	if form.Username == "" || form.Email == "" || form.Password == "" || form.ConfirmPassword == "" {
		app.errorResponse(w, r, http.StatusBadRequest, "All fields are required")
		return
	}

	if form.Password != form.ConfirmPassword {
		app.errorResponse(w, r, http.StatusBadRequest, "Passwords do not match")
		return
	}

	if len(form.Password) < 8 {
		app.errorResponse(w, r, http.StatusBadRequest, "Password must be at least 8 characters")
		return
	}

	// Check if user already exists
	existingUser, err := app.models.Users.GetByEmail(form.Email)
	if err != nil && !errors.Is(err, data.ErrRecordNotFound) {
		app.errorResponse(w, r, http.StatusInternalServerError, "Database error")
		return
	}
	if existingUser != nil {
		app.errorResponse(w, r, http.StatusConflict, "Email already registered")
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(form.Password)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, "Password hashing error")
		return
	}

	// Create user
	err = app.models.Users.Insert(form.Username, form.Email, hashedPassword, "user")
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			app.errorResponse(w, r, http.StatusConflict, "Username already taken")
			return
		}
		app.errorResponse(w, r, http.StatusInternalServerError, "Database error")
		return
	}

	app.jsonResponse(w, http.StatusCreated, map[string]string{
		"message": "User successfully registered",
	})
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	app.jsonResponse(w, status, map[string]string{"error": message})
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
