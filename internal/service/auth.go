package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"avito/internal/dto"
	"avito/internal/model"
	"avito/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	TokenTypeBearer = "Bearer"
	TokenExpireIn   = time.Hour * 24 * 365
	UserCacheKey    = "user_id_%d"
	UserTTL         = time.Minute * 30
)

var ErrAuth = errors.New("wrong email or password")

type Auth interface {
	IssueToken(ctx context.Context, input *dto.AuthInput) (*dto.AuthOutput, error)
	ParseToken(ctx context.Context, tokenString string) (*model.User, error)
}

type auth struct {
	cacheService   Cache
	userRepository repository.User
	secretKey      string
}

func NewAuth(cacheService Cache, userRepository repository.User, secretKey string) Auth {
	return &auth{
		cacheService:   cacheService,
		userRepository: userRepository,
		secretKey:      secretKey,
	}
}

func (s *auth) IssueToken(ctx context.Context, input *dto.AuthInput) (*dto.AuthOutput, error) {
	user, err := s.userRepository.GetByEmail(input.Email)
	if err != nil {
		return nil, err
	}

	val, err := json.Marshal(&user)
	if err != nil {
		return nil, fmt.Errorf("authService.IssueToken, marshall error: %w", err)
	}

	if err = s.cacheService.Set(ctx, s.getCacheKey(user.ID), val, UserTTL); err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, ErrAuth
	}

	expiresAt := time.Now().Add(TokenExpireIn)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: expiresAt},
		Subject:   strconv.FormatInt(user.ID, 10),
	})

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return nil, fmt.Errorf("authService.IssueToken, token sign error: %w", err)
	}

	return &dto.AuthOutput{
		TokenType:   TokenTypeBearer,
		ExpiresAt:   expiresAt.Unix(),
		AccessToken: tokenString,
	}, nil
}

func (s *auth) ParseToken(ctx context.Context, tokenString string) (*model.User, error) {
	var claims jwt.RegisteredClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(_ *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("authService.ParseToken, token parse error: %w", err)
	}

	if !token.Valid || claims.Subject == "" {
		return nil, errors.New("authService.ParseToken, token is not valid or payload is empty")
	}

	id, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("authService.ParseToken, subject parse %s error", claims.Subject)
	}

	user, err := s.getFromCache(ctx, id)
	if err != nil && !errors.Is(err, ErrKeyNotExist) {
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	user, err = s.userRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	val, err := json.Marshal(&user)
	if err != nil {
		return nil, fmt.Errorf("authService.ParseToken, marshall error: %w", err)
	}

	if err = s.cacheService.Set(ctx, s.getCacheKey(id), val, UserTTL); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *auth) getCacheKey(id int64) string {
	return fmt.Sprintf(UserCacheKey, id)
}

func (s *auth) getFromCache(ctx context.Context, id int64) (*model.User, error) {
	val, err := s.cacheService.Get(ctx, s.getCacheKey(id))
	if err != nil {
		return nil, err
	}

	var user model.User

	if err = json.Unmarshal([]byte(val), &user); err != nil {
		return nil, fmt.Errorf("authService.getFromCache, unmarshall error: %w", err)
	}

	return &user, nil
}
