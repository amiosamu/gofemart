package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/amiosamu/gofemart/internal/domain"
	"github.com/golang-jwt/jwt/v4"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetUser(ctx context.Context, login, password string) (domain.User, error)
}

type Users struct {
	repo   UserRepository
	hasher PasswordHasher

	hmacSecret []byte
	tokenTTL   time.Duration
}

func NewUsers(repo UserRepository, hasher PasswordHasher, secret []byte, ttl time.Duration) *Users {
	return &Users{
		repo:       repo,
		hasher:     hasher,
		hmacSecret: secret,
		tokenTTL:   ttl,
	}
}

func (u *Users) SignUp(ctx context.Context, usr domain.SighUpAndInInput) error {
	password, err := u.hasher.Hash(usr.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Login:        usr.Login,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	return u.repo.Create(ctx, user)
}

func (u *Users) SignIn(ctx context.Context, usr domain.SighUpAndInInput) (string, error) {
	password, err := u.hasher.Hash(usr.Password)
	if err != nil {
		return "", err
	}

	user, err := u.repo.GetUser(ctx, usr.Login, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrUserNotFound
		}
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.Itoa(user.ID),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(u.tokenTTL)),
	})

	return token.SignedString(u.hmacSecret)
}

func (u *Users) ParseToken(ctx context.Context, token string) (int64, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return u.hmacSecret, nil
	})
	if err != nil {
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}
	return int64(id), nil
}
