// Package main contains HTTP request handlers for the Gratitude Jar application.
// These handlers process incoming HTTP requests and generate appropriate responses.
package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/darynforman/gratitude-jar1/internal/config"
	"github.com/darynforman/gratitude-jar1/internal/data"
	"github.com/darynforman/gratitude-jar1/internal/session"
	"github.com/darynforman/gratitude-jar1/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

// home handles requests to the root path ("/").
// It displays the welcome page of the Gratitude Jar application.
func home(w http.ResponseWriter, r *http.Request) {
	// Only handle the root path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Create a PageData struct with a title for the homepage
	data := PageData{
		Title: "Welcome to Gratitude Jar",
	}
	// Render the home template with the given data
	render(w, r, "home.tmpl", data)
}

// getGratitudeModel returns a new GratitudeModel instance with the current database connection
func getGratitudeModel() *data.GratitudeModel {
	return data.NewModels(config.DB).Gratitudes
}

// getUserModel returns a new UserModel instance with the current database connection
func getUserModel() *data.UserModel {
	return data.NewModels(config.DB).Users
}

// viewNotes handles requests to view all gratitude notes.
// It supports both full page loads and HTMX partial updates.
func viewNotes(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling view notes request")

	// Get user info from session
	userID := session.Manager.GetInt(r, "userID")

	// Get notes from database with context
	notes, err := getGratitudeModel().GetAll(r.Context(), userID)
	if err != nil {
		log.Printf("Error fetching notes: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get user role from session
	role := session.Manager.GetString(r, "role")

	data := PageData{
		Title:           "My Gratitude Notes",
		Notes:           notes,
		IsAuthenticated: userID > 0,
		UserRole:        role,
	}
	log.Printf("Created PageData with %d notes", len(data.Notes))

	// Check if the request is from HTMX (for partial updates)
	if r.Header.Get("HX-Request") == "true" {
		log.Printf("HTMX request detected, rendering partial template")
		// Load and parse only the notes-list section from view-notes.tmpl
		tmpl := template.Must(template.ParseFiles("ui/html/view-notes.tmpl"))
		tmpl.ExecuteTemplate(w, "notes-list", data)
		return
	}

	log.Printf("Rendering view notes template")
	// Otherwise, render the view notes template
	render(w, r, "notes.tmpl", data)
}

// gratitude handles requests to the gratitude page where users can add new notes.
// It supports both full page loads and HTMX partial updates.
func gratitude(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling gratitude page request")

	data := PageData{
		Title:  "Add Gratitude Note",
		Emojis: []string{"✨", "🌟", "💫", "🙏", "❤️", "🌈", "🌞", "🌺", "🎉", "💝", "🌱", "⭐"},
	}

	// Check if the request is from HTMX (for partial updates)
	if r.Header.Get("HX-Request") == "true" {
		log.Printf("HTMX request detected, rendering partial template")
		// Load and parse only the notes-list section from gratitude.tmpl
		tmpl := template.Must(template.ParseFiles("ui/html/gratitude.tmpl"))
		tmpl.ExecuteTemplate(w, "notes-list", data)
		return
	}

	log.Printf("Rendering gratitude template")
	// Otherwise, render the gratitude template
	render(w, r, "gratitude.tmpl", data)
}

// createGratitude handles form submissions for creating new gratitude notes.
// It processes POST requests and supports both regular form submissions and HTMX requests.
func createGratitude(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling create gratitude note request")

	if r.Method != http.MethodPost {
		log.Printf("Invalid method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get form values
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	category := r.PostForm.Get("category")
	emoji := r.PostForm.Get("emoji")

	// Validate the form data
	v := validator.ValidateGratitudeNote(title, content, category, emoji)
	if !v.ValidData() {
		// If validation fails, return the errors
		if r.Header.Get("HX-Request") == "true" {
			// For HTMX requests, return the errors as JSON
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(v.Errors)
			return
		}
		// For regular requests, render the form with errors
		data := PageData{
			Title:  "Add Gratitude Note",
			Errors: v.Errors,
		}
		render(w, r, "add-note.tmpl", data)
		return
	}

	// Get user ID from session
	userID := session.Manager.GetInt(r, "userID")
	if userID == 0 {
		log.Printf("No user ID found in session")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Create a new gratitude note from form data
	note := &data.GratitudeNote{
		Title:     title,
		Content:   content,
		Category:  category,
		Emoji:     emoji,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert the note into the database with context
	err = getGratitudeModel().Insert(r.Context(), note)
	if err != nil {
		log.Printf("Error inserting note into database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If this is an HTMX request, return the thank you message
	if r.Header.Get("HX-Request") == "true" {
		// Render the thank you message template
		tmpl := template.Must(template.ParseFiles("ui/html/partials/thank-you-message.tmpl"))
		if err := tmpl.ExecuteTemplate(w, "thank-you-message", nil); err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	// For regular form submissions, redirect to the notes page
	http.Redirect(w, r, "/notes", http.StatusSeeOther)
}

// updateGratitude handles both updating and deleting gratitude notes.
// It processes PUT and DELETE requests and supports HTMX partial updates.
func updateGratitude(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if override := r.Header.Get("X-HTTP-Method-Override"); override != "" {
		method = override
	}
	log.Printf("Handling update/delete request with method: %s", method)
	log.Printf("Request URL: %s", r.URL.Path)

	// Get user info from session
	userID := session.Manager.GetInt(r, "userID")
	if userID == 0 {
		log.Printf("No user ID found in session")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract ID from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		log.Printf("Invalid URL path: %s", r.URL.Path)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID: %v", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	log.Printf("Parsed ID: %d", id)

	// Handle DELETE request
	if method == http.MethodDelete {
		log.Printf("Processing DELETE request for note ID: %d", id)
		err = getGratitudeModel().Delete(r.Context(), id, userID)
		if err != nil {
			log.Printf("Error deleting note: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// Return empty response for HTMX to remove the element
		w.WriteHeader(http.StatusOK)
		return
	}

	// Handle PUT request (or POST with override)
	if method != http.MethodPut {
		log.Printf("Invalid method: %s", method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("Processing PUT request for note ID: %d", id)

	// Parse form data
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Debug: Log all incoming form data
	log.Printf("Raw form data: %v", r.PostForm)

	// Get form values
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	category := r.PostForm.Get("category")
	emoji := r.PostForm.Get("emoji")

	// Log form values for debugging
	log.Printf("Form values - Title: %s, Content: %s, Category: %s, Emoji: %s", title, content, category, emoji)

	// Validate the form data
	v := validator.ValidateGratitudeNote(title, content, category, emoji)
	if !v.ValidData() {
		log.Printf("Validation errors: %v", v.Errors)
		// If validation fails, return the errors
		if r.Header.Get("HX-Request") == "true" {
			// For HTMX requests, return the errors as JSON
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(v.Errors)
			return
		}
		// For regular requests, render the form with errors
		data := PageData{
			Title:  "Edit Gratitude Note",
			Errors: v.Errors,
		}
		render(w, r, "edit-note.tmpl", data)
		return
	}

	// Create updated note
	note := &data.GratitudeNote{
		ID:        id,
		Title:     title,
		Content:   content,
		Category:  category,
		Emoji:     emoji,
		UserID:    userID,
		UpdatedAt: time.Now(),
	}
	log.Printf("Created note object: %+v", note)

	// Update note in database with context
	err = getGratitudeModel().Update(r.Context(), note)
	if err != nil {
		log.Printf("Error updating note in database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully updated note in database")

	// Fetch the updated note with context
	updatedNote, err := getGratitudeModel().Get(r.Context(), id)
	if err != nil {
		log.Printf("Error fetching updated note: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Printf("Fetched updated note: %+v", updatedNote)

	// For HTMX requests, return the updated note HTML
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		tmpl := template.Must(template.ParseFiles("ui/html/partials/note-card.tmpl"))
		if err := tmpl.ExecuteTemplate(w, "note-card", updatedNote); err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	// For regular requests, redirect to notes page
	http.Redirect(w, r, "/notes", http.StatusSeeOther)
}

// getNoteForEdit handles requests to get a note for editing.
// It retrieves a specific note by ID and renders it in the edit form.
func getNoteForEdit(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling get note for edit request")

	// Extract ID from URL
	idStr := r.URL.Path[len("/gratitude/edit/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID: %v", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Get note from database with context
	note, err := getGratitudeModel().Get(r.Context(), id)
	if err != nil {
		log.Printf("Error fetching note: %v", err)
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	// Prepare template data
	data := PageData{
		Title:  "Edit Gratitude Note",
		Note:   note,
		Emojis: []string{"✨", "🌟", "💫", "🙏", "❤️", "🌈", "🌞", "🌺", "🎉", "💝", "🌱", "⭐"},
	}

	// Render edit form
	w.Header().Set("Content-Type", "text/html")
	render(w, r, "edit-form.tmpl", data)
}

// registerHandler handles user registration (GET shows form, POST processes registration).
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := PageData{
			Title: "Register",
			Form:  map[string]string{}, // Explicitly set empty form data
		}
		render(w, r, "register.tmpl", data)
		return
	}
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		// Validate all registration fields
		v := validator.ValidateRegistration(username, email, password, confirmPassword)

		// Check for existing username/email with context
		userModel := getUserModel()
		if user, _ := userModel.GetByUsername(r.Context(), username); user != nil {
			v.AddError("username", "Username already taken")
		}
		if user, _ := userModel.GetByEmail(r.Context(), email); user != nil {
			v.AddError("email", "Email already registered")
		}

		if !v.ValidData() {
			data := PageData{
				Title:  "Register",
				Errors: v.Errors,
				Form: map[string]string{
					"username": username,
					"email":    email,
				},
			}
			render(w, r, "register.tmpl", data)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Insert new user with context
		err = userModel.Insert(r.Context(), username, email, string(hash), "user")
		if err != nil {
			data := PageData{
				Title:  "Register",
				Errors: map[string]string{"generic": "Registration failed: " + err.Error()},
				Form: map[string]string{
					"username": username,
					"email":    email,
				},
			}
			render(w, r, "register.tmpl", data)
			return
		}

		// Set flash message for successful registration
		session.Manager.Put(r, "flash", "Registration successful! Please log in.")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// loginHandler handles user login (GET shows form, POST processes login).
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// If this is an HTMX request, return just the navigation
		if r.Header.Get("HX-Request") == "true" {
			// Get user info from session
			userID := session.Manager.GetInt(r, "userID")
			role := session.Manager.GetString(r, "role")

			data := PageData{
				Title:           "Navigation",
				IsAuthenticated: userID > 0,
				UserRole:        role,
			}
			render(w, r, "partials/nav.tmpl", data)
			return
		}

		// For regular requests, show the login form with empty form data
		data := PageData{
			Title: "Login",
			Form:  map[string]string{}, // Explicitly set empty form data
		}
		render(w, r, "login.tmpl", data)
		return
	}
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Validate input fields
		v := validator.NewValidator()
		v.Check(validator.NotBlank(username), "username", "Username is required")

		// Validate password format
		if validator.NotBlank(password) {
			passwordValidator := validator.ValidatePassword(password)
			if !passwordValidator.ValidData() {
				v.AddError("password", "Invalid password format")
			}
		} else {
			v.AddError("password", "Password is required")
		}

		// If there are validation errors
		if !v.ValidData() {
			errorMessage := "Please check your username and password format"
			if r.Header.Get("HX-Request") == "true" {
				// For HTMX requests, render just the error message template
				tmpl := template.Must(template.ParseFiles("ui/html/login.tmpl"))
				tmpl.ExecuteTemplate(w, "error-message", errorMessage)
				return
			}
			// For regular requests, render the full page with error
			data := PageData{
				Title:  "Login",
				Errors: v.Errors,
			}
			render(w, r, "login.tmpl", data)
			return
		}

		// Attempt authentication
		userModel := getUserModel()
		user, err := userModel.GetByUsername(r.Context(), username)

		// Check for authentication errors
		var errorMessage string
		if err != nil || user == nil {
			errorMessage = "Invalid username or password"
		} else if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
			errorMessage = "Invalid username or password"
		}

		// If there's an authentication error
		if errorMessage != "" {
			if r.Header.Get("HX-Request") == "true" {
				// For HTMX requests, render just the error message template
				tmpl := template.Must(template.ParseFiles("ui/html/login.tmpl"))
				tmpl.ExecuteTemplate(w, "error-message", errorMessage)
				return
			}
			// For regular requests, render the full page with error
			data := PageData{
				Title:  "Login",
				Errors: map[string]string{"generic": errorMessage},
			}
			render(w, r, "login.tmpl", data)
			return
		}

		// Set session values
		session.Manager.Put(r, "userID", user.ID)
		session.Manager.Put(r, "role", user.Role)
		session.Manager.Put(r, "flash", "Successfully logged in!")

		// Check if this is an HTMX request
		if r.Header.Get("HX-Request") == "true" {
			// Set HX-Redirect header for successful login
			w.Header().Set("HX-Redirect", "/")
			// Also trigger navigation update
			w.Header().Set("HX-Trigger", "loginSuccess")
			return
		}

		// For regular requests, redirect to home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// logoutHandler logs out the user by destroying the session.
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear session data first
	session.Manager.Put(r, "userID", nil)
	session.Manager.Put(r, "role", nil)

	// Force the session to be saved with cleared values
	session.Manager.Put(r, "_cleared", time.Now().Unix())

	// Redirect to login page
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// contact handles requests to the contact page.
// It displays the contact information and form.
func contact(w http.ResponseWriter, r *http.Request) {
	log.Printf("Contact handler called with path: %s and method: %s", r.URL.Path, r.Method)

	// Only handle exact /contact path
	if r.URL.Path != "/contact" {
		log.Printf("Invalid contact path: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	// Only handle GET requests
	if r.Method != http.MethodGet {
		log.Printf("Invalid method for contact: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := PageData{
		Title: "Contact Us",
	}

	// Render the contact template with the given data
	render(w, r, "contact.tmpl", data)
	log.Printf("Contact page rendered successfully")
}

// about handles requests to the about page.
func about(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/about" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	data := PageData{
		Title: "About",
	}
	render(w, r, "about.tmpl", data)
}
