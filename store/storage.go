package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrInternalServer      = errors.New("something went wrong")
	ErrNotFound            = errors.New("record not found")
	QueryTimeoutDuration   = time.Second * 5
	ErrDuplicateEmail      = errors.New("duplicated Email")
	ErrDupliucatedUsername = errors.New("duplicated Username")
)

func ErrorFactoryDB(err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrNotFound
	case err.Error() == `pq: duplicated key value violates unique constraint "users_email_key"`:
		return ErrDuplicateEmail
	case err.Error() == `pq: duplicated key value vioates unique constraint "users_username_key"`:
		return ErrDupliucatedUsername
	default:
		return ErrInternalServer
	}

}

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, int64) (*Post, error)
		DeleteById(context.Context, int64) error
		UpdatePost(context.Context, int64, *Post) (*Post, error)
		GetUserFeed(context.Context, int64, PaginationStructQuery) ([]PostMetadata, error)
	}
	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		Delete(context.Context, int64) error
		GetUserByUsername(context.Context, string) (*QueryUserByUsername, error)
		GetUserById(context.Context, int64) (*User, error)
		Deactivate(context.Context, int64) error
		Activate(context.Context, string) error
		CreateAndInvite(context.Context, *User, string, time.Duration) error
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetCommentsByUserId(context.Context, int64) ([]Comment, error)
	}
	Followers interface {
		Follow(ctx context.Context, user_id, follower_id int64) error
		Unfollow(ctx context.Context, user_id, follower_id int64) error
	}
	Roles interface {
		GetRoleByName(ctx context.Context, roleName string) (*Role, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostsStore{db},
		Users:     &UsersStorage{db},
		Comments:  &CommentStorage{db},
		Followers: &FollowerStorage{db},
		Roles:     &RolesStorage{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
