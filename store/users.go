package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	text *string
	hash []byte
}

func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.text = &text
	p.hash = hash
	return nil
}
func (p *Password) Compare(text string, hash []byte) bool {

	err := bcrypt.CompareHashAndPassword(hash, []byte(text))
	if err != nil {
		return false
	}
	return true
}

type UsersStorage struct {
	db *sql.DB
}
type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  Password `json:"-"`
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
	RoleId    int64    `json:"role_id"`
	Role      Role     `json:"role"`
}

type QueryUserByUsername struct {
	ID       int64  `json:"id"`
	Password []byte `json:"password"`
}

func (s *UsersStorage) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (username,password, email, role_id)
		VALUES ($1,$2, $3, (SELECT id FROM roles WHERE name=$4)) RETURNING id, created_at
	`
	role := user.Role.Name
	if role == "" {
		role = "user"
	}
	err := tx.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password.hash,
		user.Email,
		role,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		fmt.Println("err===", err.Error())
		return ErrorFactoryDB(err)

	}

	return nil
}

func (s *UsersStorage) GetUserByUsername(ctx context.Context, username string) (*QueryUserByUsername, error) {
	var queryResult QueryUserByUsername
	query := `
		SELECT id,password FROM users WHERE username=$1
	`

	err := s.db.QueryRowContext(ctx, query, username).Scan(&queryResult.ID, &queryResult.Password)
	if err != nil {
		fmt.Println(err.Error())
		return nil, ErrorFactoryDB(err)
	}
	return &queryResult, nil
}

func (s *UsersStorage) GetUserById(ctx context.Context, userId int64) (*User, error) {
	var user User
	query := `
		SELECT users.id, users.email, users.username, users.created_at,users.password, roles.*  FROM users
		JOIN roles ON users.role_id = roles.id
		WHERE users.id=$1 AND is_active = true
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		userId,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.CreatedAt,
		&user.Password.hash,
		&user.Role.Id,
		&user.Role.Name,
		&user.Role.Level,
		&user.Role.Description,
	)
	fmt.Println(userId, "-----", user)
	if err != nil {
		fmt.Println("err===", err.Error())
		return nil, ErrorFactoryDB(err)

	}
	return &user, nil

}
func (s *UsersStorage) Delete(ctx context.Context, userId int64) error {
	query := `
			DELETE FROM users
			WHERE id=$1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, userId)
	if err != nil {
		return ErrorFactoryDB(err)
	}
	return nil
}
func (s *UsersStorage) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, tx, user); err != nil {
			return err

		}

		if err := s.createUserInvitation(ctx, tx, token, invitationExp, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (s *UsersStorage) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, invitationExp time.Duration, userID int64) error {
	query := `
		INSERT INTO user_invitation
		(token, user_id, expiry) VALUES ($1,$2,$3)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(invitationExp))

	if err != nil {
		return err
	}
	return nil
}

func (s *UsersStorage) Activate(ctx context.Context, token string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}
		user.IsActive = true

		err = s.update(ctx, tx, user)
		if err != nil {
			return err
		}

		err = s.deleteInvitation(ctx, tx, user.ID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *UsersStorage) Deactivate(ctx context.Context, userId int64) error {
	query :=
		`
					DELETE FROM user_invitation 
					WHERE user_id = $1
				`
	_, err := s.db.ExecContext(ctx, query, userId)
	if err != nil {
		return ErrorFactoryDB(err)
	}
	return nil
}
func (s *UsersStorage) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query :=
		`UPDATE users SET username =$1, email= $2, is_active= $3 WHERE id = $4`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)

	if err != nil {
		return err
	}

	return nil
}
func (s *UsersStorage) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.is_active
		FROM users u
		JOIN user_invitation ui ON u.id = ui.user_id
		WHERE ui.token = $1 AND ui.expiry > $2
		`
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UsersStorage) deleteInvitation(ctx context.Context, tx *sql.Tx, userId int64) error {
	query := `
		DELETE FROM user_invitation WHERE user_id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, userId)

	if err != nil {
		return err
	}
	return nil
}
