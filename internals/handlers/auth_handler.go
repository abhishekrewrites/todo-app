package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"todo-app/internals/models"

	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *sql.DB
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{
		DB: db,
	}
}

func (h *AuthHandler) Login(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	var req models.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	var user models.User

	err = h.DB.QueryRow(
		`
		SELECT
			id,
			name,
			email,
			password
		FROM users
		WHERE email = $1
		`,
		req.Email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		http.Error(
			w,
			"invalid email or password",
			http.StatusUnauthorized,
		)
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		http.Error(
			w,
			"invalid email or password",
			http.StatusUnauthorized,
		)
		return
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": user.ID,
			"name":    user.Name,
			"email":   user.Email,
		},
	)

	tokenString, err := token.SignedString(
		[]byte(os.Getenv("JWT_SECRET")),
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"message": "Login successful",
			"token":   tokenString,
			"user": map[string]interface{}{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
			},
		},
	)
}

func (h *AuthHandler) Register(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	var req models.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	if req.Name == "" ||
		req.Email == "" ||
		req.Password == "" ||
		req.ConfirmPassword == "" {

		http.Error(
			w,
			"all fields are required",
			http.StatusBadRequest,
		)
		return
	}

	if req.Password != req.ConfirmPassword {
		http.Error(
			w,
			"passwords do not match",
			http.StatusBadRequest,
		)
		return
	}

	var count int

	err = h.DB.QueryRow(
		`
		SELECT COUNT(*)
		FROM users
		WHERE email = $1
		`,
		req.Email,
	).Scan(&count)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	if count > 0 {
		http.Error(
			w,
			"email already exists",
			http.StatusBadRequest,
		)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	var userID int

	err = h.DB.QueryRow(
		`
		INSERT INTO users(
			name,
			email,
			password
		)
		VALUES($1, $2, $3)
		RETURNING id
		`,
		req.Name,
		req.Email,
		string(hashedPassword),
	).Scan(&userID)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"id":      userID,
			"name":    req.Name,
			"email":   req.Email,
			"message": "User registered successfully",
		},
	)
}