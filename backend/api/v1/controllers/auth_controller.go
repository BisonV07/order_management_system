package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"oms/backend/api/v1/helpers"
	apitypes "oms/backend/api/v1/types"
	"oms/backend/core/auth"
	"oms/backend/core/model"
	"oms/backend/core/types"
)

// AuthController handles authentication-related HTTP requests
type AuthController struct {
	userStore types.UserStore
}

// NewAuthController creates a new AuthController
func NewAuthController(userStore types.UserStore) *AuthController {
	return &AuthController{
		userStore: userStore,
	}
}

// Signup handles POST /api/v1/auth/signup
func (ac *AuthController) Signup(w http.ResponseWriter, r *http.Request) {
	var req apitypes.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "Username and password are required")
		return
	}

	if len(req.Password) < 4 {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "Password must be at least 4 characters")
		return
	}

	// Hash password
	hashedPassword, err := model.HashPassword(req.Password)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, "internal_error", "Failed to hash password")
		return
	}

	// Create user
	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Role:     model.UserRoleUser, // Regular users by default
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = ac.userStore.Create(r.Context(), user)
	if err != nil {
		if err.Error() == "username already exists" {
			helpers.WriteErrorResponse(w, http.StatusConflict, "conflict", "Username already exists")
			return
		}
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, "internal_error", "Failed to create user")
		return
	}

	helpers.WriteJSONResponse(w, http.StatusCreated, apitypes.SignupResponse{
		Message: "User created successfully",
		UserID:  user.ID,
	})
}

// Login handles POST /api/v1/auth/login
func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req apitypes.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid_credentials", "Username and password are required")
		return
	}

	// Get user by username
	user, err := ac.userStore.GetByUsername(r.Context(), req.Username)
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusUnauthorized, "invalid_credentials", "Invalid username or password")
		return
	}

	// Verify password
	if !model.CheckPassword(req.Password, user.Password) {
		helpers.WriteErrorResponse(w, http.StatusUnauthorized, "invalid_credentials", "Invalid username or password")
		return
	}

	// Generate JWT token with role
	token, err := auth.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, "internal_error", "Failed to generate token")
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, apitypes.LoginResponse{
		Token:  token,
		UserID: user.ID,
		Role:   string(user.Role),
	})
}

