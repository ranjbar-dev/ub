package jwt

import (
	"exchange-go/internal/platform"
	"os"
	"time"
)

// Service provides JWT token issuance and validation for user authentication.
type Service interface {
	// IssueToken creates a new signed JWT for the given username, user agent, and IP address.
	IssueToken(username, userAgent, ip string) (string, error)
	// GetUsernameFromToken validates the signed token and returns the username and
	// expiration time extracted from its claims.
	GetUsernameFromToken(signedToken string) (string, time.Time, error)
}

type service struct {
	jwtHandler     platform.JwtHandler
	privateKeyPath string
	publicKeyPath  string
	passphrase     string
	ttl            int
	env            string
}

func (jts *service) IssueToken(username, userAgent, ip string) (string, error) {
	//if jwtTokenService.env == "test" {
	//	//only calling this so our expectations in unit test pass
	//	_, _ = jwtTokenService.jwtHandler.GetToken([]byte(""), jwtTokenService.passphrase, jwtTokenService.ttl, username)
	//	return "token", nil
	//}
	privateKey, err := os.ReadFile(jts.privateKeyPath)

	if err != nil {
		return "", err
	}
	params := platform.GetTokenParams{
		PrivateKey: privateKey,
		Passphrase: jts.passphrase,
		TTL:        jts.ttl,
		Username:   username,
		IP:         ip,
		UserAgent:  userAgent,
	}
	tokenStr, err := jts.jwtHandler.GetToken(params)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (jts *service) GetUsernameFromToken(signedToken string) (string, time.Time, error) {
	//if jwtTokenService.env == "test" {
	//	_, _ = jwtTokenService.jwtHandler.GetUsernameFromToken([]byte(""), signedToken)
	//	return "username", nil
	//}
	publicKey, _ := os.ReadFile(jts.publicKeyPath)
	return jts.jwtHandler.GetUsernameFromToken(publicKey, signedToken)
}

func NewJwtService(c platform.Configs, jwtHandler platform.JwtHandler) Service {
	privateKeyPath := c.GetString("jwt.private_key")
	publicKeyPath := c.GetString("jwt.public_key")
	passphrase := c.GetString("jwt.passphrase")
	ttl := c.GetInt("jwt.ttl")
	env := c.GetEnv()
	return &service{jwtHandler, privateKeyPath, publicKeyPath, passphrase, ttl, env}
}
