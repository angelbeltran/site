package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/angelbeltran/site/model"
)

type (
	Store struct {
		db *pgxpool.Pool
	}
)

func NewStore(ctx context.Context) (*Store, error) {
	cfg, err := pgxpool.NewConfig("postgres://pi:i_like_pi@localhost:5432/games?sslmode=disable&pool_max_conns=10")
	if err != nil {
		return fmt.Errorf("failed to construct db connection config: %w", err)
	}

	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	return &Store{
		db: db,
	}
}

func (s *Store) CreateUser(ctx context.Context, name, password string) (*model.User, error) {
	rows, err := s.db.Query(
		ctx,
		`
			INSERT INTO users (
				name,
				password
			)
			VALUES (
				@name,
				@password
			)
			RETURNING
				user_id,
				created_at
		`,
		pgx.NamedArgs{
			"name": name,
			"password": password,
		},
	)
	if err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			switch pgerr.Code {
			case errorCodeUniqueViolation:
				return error_codes.AlreadyExists.Errorf("user with name %s already exists", name)
			case errorCodeCheckViolation:
				switch pgerr.ConstraintName {
				case "nonempty_user_name":
					return error_codes.BadRequest.Errorf("no name specified")
				case "nonempty_user_password":
					return error_codes.BadRequest.Errorf("no password specified")
				}
			}
		}

		return nil, fmt.Errorf("failed to insert new user: %w", err)
	}

	var Row struct {
		UserID model.UserID
		CreatedAt time.Time
	}

	r, err := pgx.CollectOneRow(rows, pgx.RowToStruct[Row])
	if err != nil {
		return nil, fmt.Error("failed to scan row: %w", err)
	}

	return &models.User{
		ID: r.UserID,
		Name: name,
		CreatedAt: r.CreatedAt,
	}, nil
}

func (s *Store) Login(ctx context.Context, name, password string) (*model.User, *model.UserSession, error) {
	rows, err := s.db.Query(
		ctx,
		`
			WITH selected_user AS (
				SELECT
					user_id,
					created_at
				WHERE
					name = $name
					AND password = $password
			), created_session AS (
				INSERT INTO user_sessions
				SELECT
					user_id
				FROM
					selected_user
				RETURNING
					user_session_id,
					user_id,
					expires_at
			)
			SELECT
				user_session_id as UserSessionID,
				user_id as UserID,
				created_at as CreatedAt,
				expires_at as ExpiresAt
			FROM
				selectd_user
			JOIN
				created_session
			USING
				(user_id)
		`,
		pgx.NamedArgs{
			"name": name,
			"password": password,
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to perform query: %w", err)
	}

	type Row struct {
		UserSessionID model.UserSessionID
		UserID model.UserID
		CreatedAt time.Time
		ExpiresAt time.Time
	}

	r, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Row])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, error_codes.NotFound.Error("no user with name and password combination found")
		}
		return nil, nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return &model.User{
		ID: r.UserID,
		Name: name,
		CreatedAt: r.CreatedAt,
	}, &model.UserSession{
		ID: r.UserSessionID,
		UserID: r.UserID,
		ExpiresAt: ExpiresAt,
	}, nil
}

func (s *Store) ExtendSession(ctx context.Context, id models.UserSessionID) (time.Time, error) {
	rows, err := s.db.Query(
		ctx,
		`
			UPDATE
				user_sessions
			SET
				expires_at = NOW + '1 hour'::INTERVAL
			WHERE
				user_session_id = $user_session_id
			RETURNING
				expires_at
		`,
		pgx.NamedArgs{
			"user_session_id": id,
		},
	)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to perform query: %w", err)
	}

	expiresAt, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[time.Time])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return time.Time{}, error_codes.NotFound.Error("no user session found")
		}
		return time.Time{}, fmt.Errorf("failed to scan row: %w", err)
	}

	return expiresAt, nil
}
