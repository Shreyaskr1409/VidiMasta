package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Shreyaskr1409/VidiMasta/gateway/data"
	"github.com/Shreyaskr1409/VidiMasta/gateway/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	l  *log.Logger
	db *pgxpool.Pool
}

func NewUserHandler(l *log.Logger, db *pgxpool.Pool) *UserHandler {
	return &UserHandler{
		l:  l,
		db: db,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var userReq data.User
	if err := utils.ParseRequest(r, &userReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(userReq.Username) == "" ||
		strings.TrimSpace(userReq.Email) == "" ||
		strings.TrimSpace(userReq.Email) == "" ||
		strings.TrimSpace(userReq.Fullname) == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	var exists bool
	ctx := context.Background()
	defer ctx.Done()
	err := h.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2)", userReq.Username, userReq.Email).Scan(&exists)
	if err != nil {
		h.l.Println("Database error while uploading the data: ", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		h.l.Println("User with the given credentials already exists")
		http.Error(w, "User with the given credentials already exists", http.StatusConflict)
		return
	}
	h.l.Println("Obtained valid user request")

	id := uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		h.l.Println("Unable to hash the password: ", err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	now := time.Now()
	refreshToken, err := data.GenerateRefreshToken(&userReq)
	if err != nil {
		h.l.Println("Failed to generate refresh token: ", err)
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	_, err = h.db.Exec(ctx,
		`INSERT INTO users 
		(id, username, fullname, email, password, refresh_token, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		id,
		userReq.Username,
		userReq.Fullname,
		userReq.Email,
		hashedPassword,
		refreshToken,
		now,
		now,
	)
	if err != nil {
		h.l.Println("Database error creating user: ", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	newUser := data.User{
		Id:           id.String(),
		Username:     userReq.Username,
		Fullname:     userReq.Fullname,
		Email:        userReq.Email,
		RefreshToken: refreshToken,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	h.l.Println("User created successfully: ", newUser)

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(newUser); err != nil {
		h.l.Println("Error encoding response:", err)
		http.Error(w, "Internal Server Error (But the user is created successfully)", http.StatusInternalServerError)
		return
	}
	h.l.Println("Response written successfully")
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world")
}
