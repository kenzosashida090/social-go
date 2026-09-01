package store

import (
	"context"
	"database/sql"
)

type Followers struct {
	UserId     int64  `json:"user_id"`
	FollowerId int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowerStorage struct {
	db *sql.DB
}

func (s *FollowerStorage) Follow(ctx context.Context, followerId, userId int64) error {
	query := `
		INSERT INTO followers (user_id, follower_id) VALUES ($1,$2)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userId, followerId)

	if err != nil {
		return ErrorFactoryDB(err)
	}

	return nil

}

func (s *FollowerStorage) Unfollow(ctx context.Context, userId, followerID int64) error {

	query := `
			DELETE FROM followers
			WHERE  user_id = $1 AND follower_id =$2
		`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userId, followerID)

	if err != nil {
		return ErrorFactoryDB(err)
	}
	return nil

}
