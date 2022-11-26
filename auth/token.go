package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTServicer interface {
	GenerateJWTToken(userId int64, role string) (string, error)
	ValidateJWTToken(signedString string) (*TokenData, error)
}

type TokenData struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
}

type JWTService struct {
	signingSecret string
	issuer        string
	maxAge        time.Duration
}

func CreateJWTService(signingSecret, issuer string, maxAge time.Duration) (*JWTService, error) {
	if signingSecret == "" || issuer == "" {
		return nil, errors.New("missing signing secret or issuer in environment variables")
	}
	return &JWTService{signingSecret: signingSecret, issuer: issuer, maxAge: maxAge}, nil
}

func (j *JWTService) GenerateJWTToken(userId int64, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	// set claims and expiry
	claims["iss"] = j.issuer
	claims["rights-stuff:role"] = role
	claims["iat"] = time.Now().UTC().Unix()
	claims["rights-stuff:user"] = userId
	claims["exp"] = time.Now().Add(j.maxAge).UTC().Unix()

	// Sign Token
	tokenString, err := token.SignedString([]byte(j.signingSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *TokenData) SetClaims(claims jwt.MapClaims) error {
	userID, ok := claims["rights-stuff:user"].(float64)
	if !ok {
		return jwt.NewValidationError("Invalid userID in JWT token", 0)
	}

	t.UserID = int64(userID)

	role, ok := claims["rights-stuff:role"].(string)
	if !ok {
		return jwt.NewValidationError("Invalid userID in JWT token", 0)
	}
	t.Role = role

	return nil
}

func (j *JWTService) ValidateJWTToken(signedString string) (*TokenData, error) {
	// Parse token
	token, err := jwt.Parse(signedString, func(token *jwt.Token) (any, error) {
		// First validate the algorithm signature
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("Invalid JWT token", 0)
		}

		return []byte(j.signingSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Check claims and if token is valid
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, jwt.NewValidationError("Invalid JWT token", 0)
	}

	data := TokenData{}
	data.SetClaims(claims)

	return &data, nil
}
