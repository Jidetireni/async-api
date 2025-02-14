package store

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserStore struct {
	db *sqlx.DB
}

type User struct {
	Id                 uuid.UUID `db:"id"`
	Email              string    `db:"email"`
	HashedPasswdBase64 string    `db:"password_hash"`
	CreatedAt          time.Time `db:"created_at"`
}

func (u *User) Validate(password string) error {
	byts, err := base64.StdEncoding.DecodeString(u.HashedPasswdBase64)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(byts, []byte(password))
	if err != nil {
		return fmt.Errorf("invalid password %w", err)
	}

	return nil

}

func (s *UserStore) CreateUser(ctx context.Context, email, password string) (*User, error) {
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING *;`

	var user User
	byts, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	HashedPasswdBase64 := base64.StdEncoding.EncodeToString(byts)

	err = s.db.GetContext(ctx, &user, query, email, HashedPasswdBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return &user, nil

}

func (s *UserStore) ByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT * FROM users WHERE email = $1`

	var user User
	err := s.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to return user: %w", err)
	}
	return &user, nil
}

func (s *UserStore) ByID(ctx context.Context, userId uuid.UUID) (*User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	var user User
	err := s.db.GetContext(ctx, &user, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to return user: %w", err)
	}
	return &user, nil
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: sqlx.NewDb(db, "postgres"),
	}
}
