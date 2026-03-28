package platform

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/ssh"
)

type GetTokenParams struct {
	PrivateKey []byte
	Passphrase string
	TTL        int
	Username   string
	IP         string
	UserAgent  string
}

// JwtHandler provides JWT token creation and validation using RSA key pairs.
type JwtHandler interface {
	// GetToken creates a new RS256-signed JWT containing the username, IP, user-agent,
	// and expiration derived from the TTL in params.
	GetToken(params GetTokenParams) (string, error)
	// GetUsernameFromToken validates the signed JWT using the provided RSA public key
	// and returns the embedded username and the token's issued-at time. An error is
	// returned if the token is invalid or expired.
	GetUsernameFromToken(publicKey []byte, signedToken string) (string, time.Time, error)
}

type jwtHandler struct {
}

type claims struct {
	Username  string `json:"username"`
	Iat       int64  `json:"iat"`
	Exp       int64  `json:"exp"`
	UserAgent string `json:"userAgent"`
	IP        string `json:"ip"`
	jwt.RegisteredClaims
}

type InValidToken struct {
}

func (jwtHandler *jwtHandler) GetToken(params GetTokenParams) (string, error) {
	rsaPrivateKey, err := ssh.ParseRawPrivateKeyWithPassphrase(params.PrivateKey, []byte(params.Passphrase))

	if err != nil {
		return "", err
	}

	now := time.Now()
	expiresAt := now.Add(time.Minute * time.Duration(params.TTL))
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, &claims{
		Username:  params.Username,
		Iat:       now.Unix(),
		Exp:       expiresAt.Unix(),
		UserAgent: params.UserAgent,
		IP:        params.IP,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})

	tokenStr, err := t.SignedString(rsaPrivateKey)

	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (jwtHandler *jwtHandler) GetUsernameFromToken(publicKey []byte, signedToken string) (string, time.Time, error) {
	issuedAt := time.Now()
	token, err := jwt.ParseWithClaims(signedToken, &claims{}, func(token *jwt.Token) (interface{}, error) {
		rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
		return rsaPublicKey, err
	})

	if err != nil {
		return "", issuedAt, err
	}

	if c, ok := token.Claims.(*claims); ok && token.Valid {
		now := time.Now().Unix()
		if c.Exp < now {
			return "", time.Now(), fmt.Errorf("token expired")
		}
		issuedAt = time.Unix(c.Iat, 0)
		return c.Username, issuedAt, nil
	}

	return "", issuedAt, fmt.Errorf("invalid token claims")
}

func NewJwtHandler() JwtHandler {
	return &jwtHandler{}
}
