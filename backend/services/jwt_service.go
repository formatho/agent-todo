package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	secretKey []byte
}

func NewJWTService() *JWTService {
	return &JWTService{
		secretKey: []byte(os.Getenv("JWT_SECRET")),
	}
}

type Claims struct {
	UserID         string `json:"user_id"`
	Email          string `json:"email"`
	OrganisationID string `json:"organisation_id,omitempty"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for a user
func (s *JWTService) GenerateToken(userID, email string) (string, error) {
	return s.GenerateTokenWithOrg(userID, email, "")
}

// GenerateTokenWithOrg generates a JWT token for a user with organisation context
func (s *JWTService) GenerateTokenWithOrg(userID, email, organisationID string) (string, error) {
	claims := &Claims{
		UserID:         userID,
		Email:          email,
		OrganisationID: organisationID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "agent-todo",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken validates a JWT token and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GenerateAPIKey generates a new API key for agents
func GenerateAPIKey() string {
	return "sk_agent_" + uuid.New().String()
}
