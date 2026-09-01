package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type PostsStore struct {
	db *sql.DB
}
type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	CreatedAt string    `json:"created_at"`
	Tags      []string  `json:"tags"`
	UpdatedAt string    `json:"updated_at"`
	Version   int64     `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostMetadata struct {
	Post
	CommentsCount int64 `json:"comments_count"`
}

func (s *PostsStore) GetById(ctx context.Context, postId int64) (*Post, error) {
	var post Post
	query := `SELECT * FROM posts WHERE id=$1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query, postId).Scan(
		&post.ID,
		&post.Title,
		&post.UserID,
		&post.Content,
		&post.CreatedAt,
		pq.Array(&post.Tags),
		&post.UpdatedAt,
		&post.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	fmt.Println("ERROR RETURNED")
	return &post, nil
}

func (s *PostsStore) GetUserFeed(ctx context.Context, userId int64, fq PaginationStructQuery) ([]PostMetadata, error) {
	var feed []PostMetadata
	//
	query := `
		SELECT 
			p.id, p.user_id, p.content, p.created_at, p.version, p.tags, 
			u.username,
			COUNT (c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		JOIN followers f ON f.follower_id = p.user_id OR p.user_id=$1
		WHERE 
			f.user_id = $1 AND
			(p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') AND
			(p.tags @> $5::varchar[] OR $5::varchar[] = '{}')
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + fq.Sort + `
		LIMIT $2 OFFSET $3

	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query, userId, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags))

	if err != nil {
		return nil, ErrorFactoryDB(err)
	}
	defer rows.Close()
	for rows.Next() {
		var feedResponse PostMetadata
		err := rows.Scan(
			&feedResponse.ID,
			&feedResponse.UserID,
			&feedResponse.Content,
			&feedResponse.CreatedAt,
			&feedResponse.Version,
			pq.Array(&feedResponse.Tags),
			&feedResponse.User.Username,
			&feedResponse.CommentsCount,
		)
		if err != nil {

			return nil, ErrorFactoryDB(err)
		}
		feed = append(feed, feedResponse)

	}
	return feed, nil
}

func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1,$2,$3,$4) RETURNING id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostsStore) DeleteById(ctx context.Context, postId int64) error {

	query := `DELETE FROM posts WHERE id=$1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	res, err := s.db.ExecContext(ctx, query, postId)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound

		default:
			return err
		}
	}
	println(res)
	return nil

}

func (s *PostsStore) UpdatePost(ctx context.Context, postId int64, post *Post) (*Post, error) {
	var updatedPost Post
	query := `UPDATE posts 
						SET content=$1, title=$2,  version = version + 1 
						WHERE id=$3 AND version = $4
						RETURNING id,content,title,updated_at, version
						`
	// Avoiding usuing old data when updating, keep track of each post based on his current version
	// THis will no match if you are using an old version of the post.
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		postId,
		post.Version,
	).Scan(
		&updatedPost.ID,
		&updatedPost.Content,
		&updatedPost.Title,
		&updatedPost.UpdatedAt,
		&updatedPost.Version,
	)
	if err != nil {
		return nil, err
	}
	return &updatedPost, nil
}
