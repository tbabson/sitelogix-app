package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sitelogix/backend/internal/user"
	"github.com/sitelogix/backend/pkg/config"
	"github.com/sitelogix/backend/pkg/middleware"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

type Service struct {
	repo     *Repository
	userRepo *user.Repository
	cfg      config.JWTConfig
}

func NewService(repo *Repository, userRepo *user.Repository, cfg config.JWTConfig) *Service {
	return &Service{repo: repo, userRepo: userRepo, cfg: cfg}
}

func (s *Service) Register(req RegisterRequest) (*user.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return s.userRepo.Create(user.CreateUserParams{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         string(req.Role),
	})
}

func (s *Service) Login(req LoginRequest) (*TokenPair, error) {
	u, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.generateTokenPair(u.ID, u.Name, u.Email, u.Role)
}

func (s *Service) Refresh(tokenStr string) (*TokenPair, error) {
	hash := hashToken(tokenStr)
	rt, err := s.repo.FindRefreshToken(hash)
	if err != nil {
		return nil, ErrInvalidToken
	}

	u, err := s.userRepo.FindByID(rt.UserID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if err := s.repo.DeleteRefreshToken(hash); err != nil {
		return nil, err
	}

	return s.generateTokenPair(u.ID, u.Name, u.Email, u.Role)
}

func (s *Service) Logout(tokenStr string) error {
	return s.repo.DeleteRefreshToken(hashToken(tokenStr))
}

func (s *Service) generateTokenPair(userID, name, email, role string) (*TokenPair, error) {
	now := time.Now()
	claims := &middleware.Claims{
		UserID: userID,
		Name:   name,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.AccessTokenTTL)),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return nil, err
	}

	rawRefresh, err := generateRandomToken()
	if err != nil {
		return nil, err
	}

	expiresAt := now.Add(s.cfg.RefreshTokenTTL)
	if err := s.repo.SaveRefreshToken(userID, hashToken(rawRefresh), expiresAt); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		ExpiresIn:    int64(s.cfg.AccessTokenTTL.Seconds()),
	}, nil
}

func generateRandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
