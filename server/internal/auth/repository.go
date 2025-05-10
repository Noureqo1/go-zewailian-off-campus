package auth

import (
	"context"
	"database/sql"
)

type Repository interface {
	UpsertUser(ctx context.Context, user *User) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByGoogleID(ctx context.Context, googleID string) (*User, error)
	CreateSession(ctx context.Context, session *Session) (*Session, error)
	GetSessionByToken(ctx context.Context, token string) (*Session, error)
	DeleteSession(ctx context.Context, token string) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) UpsertUser(ctx context.Context, user *User) (*User, error) {
	query := `
		INSERT INTO users (id, email, name, picture, password_hash, google_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (email) DO UPDATE
		SET name = $3, picture = $4, password_hash = COALESCE($5, users.password_hash), google_id = COALESCE($6, users.google_id)
		RETURNING id, email, name, picture, password_hash, google_id, created_at
	`

	var result User
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.Name,
		user.Picture,
		user.PasswordHash,
		user.GoogleID,
		user.CreatedAt,
	).Scan(
		&result.ID,
		&result.Email,
		&result.Name,
		&result.Picture,
		&result.PasswordHash,
		&result.GoogleID,
		&result.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *PostgresRepository) GetUserByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, email, name, picture, password_hash, google_id, created_at
		FROM users
		WHERE id = $1
	`

	var user User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Picture,
		&user.PasswordHash,
		&user.GoogleID,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, name, picture, password_hash, google_id, created_at
		FROM users
		WHERE email = $1
	`

	var user User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Picture,
		&user.PasswordHash,
		&user.GoogleID,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresRepository) GetUserByGoogleID(ctx context.Context, googleID string) (*User, error) {
	query := `
		SELECT id, email, name, picture, password_hash, google_id, created_at
		FROM users
		WHERE google_id = $1
	`

	var user User
	err := r.db.QueryRowContext(ctx, query, googleID).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Picture,
		&user.PasswordHash,
		&user.GoogleID,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresRepository) CreateSession(ctx context.Context, session *Session) (*Session, error) {
	query := `
		INSERT INTO sessions (id, user_id, token, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, token, expires_at
	`

	var result Session
	err := r.db.QueryRowContext(
		ctx,
		query,
		session.ID,
		session.UserID,
		session.Token,
		session.ExpiresAt,
	).Scan(
		&result.ID,
		&result.UserID,
		&result.Token,
		&result.ExpiresAt,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *PostgresRepository) GetSessionByToken(ctx context.Context, token string) (*Session, error) {
	query := `
		SELECT id, user_id, token, expires_at
		FROM sessions
		WHERE token = $1
	`

	var session Session
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&session.ID,
		&session.UserID,
		&session.Token,
		&session.ExpiresAt,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *PostgresRepository) DeleteSession(ctx context.Context, token string) error {
	query := `DELETE FROM sessions WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}
