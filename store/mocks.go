package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}

}

type MockUserStore struct{}

func (m MockUserStore) Create(ctx context.Context, tx *sql.Tx, u *User) error {
	return nil
}

func (m MockUserStore) Delete(ctx context.Context, id int64) error {

	return nil
}
func (m MockUserStore) GetUserByUsername(ctx context.Context, username string) (*QueryUserByUsername, error) {
	return nil, nil
}

func (m MockUserStore) GetUserById(ctx context.Context, userId int64) (*User, error) {
	return nil, nil
}
func (m MockUserStore) Deactivate(ctx context.Context, userId int64) error {
	return nil
}
func (m MockUserStore) Activate(ctx context.Context, token string) error {
	return nil
}
func (m MockUserStore) CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error {
	return nil
}
