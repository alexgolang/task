package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
	issuer      string
	tokenExpiry time.Duration
}

func NewJWTService(privateKeyPEM string, issuer string, tokenExpiry time.Duration) (*JWTService, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))

	if block == nil {
		return nil, fmt.Errorf("jwt service: invalid private key PEM")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {

		parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("jwt service: failed to parse private key: %w", err)
		}

		privateKey, ok := parsedKey.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("jwt service: failed to parse private key")
		}

		return &JWTService{
			privateKey:  privateKey,
			publicKey:   &privateKey.PublicKey,
			issuer:      issuer,
			tokenExpiry: tokenExpiry,
		}, nil
	}

	return &JWTService{
		privateKey:  privateKey,
		publicKey:   &privateKey.PublicKey,
		issuer:      issuer,
		tokenExpiry: tokenExpiry,
	}, nil
}

func (s *JWTService) ValidateClientAssertion(clientAssertion string, clientAssertionType string) (jwt.MapClaims, error) {
	parser := new(jwt.Parser)
	token, _, err := parser.ParseUnverified(clientAssertion, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("jwt service: failed to parse client assertion: %w", err)
	}

	x5c, ok := token.Header["x5c"].([]interface{})
	if !ok || len(x5c) == 0 {
		return nil, fmt.Errorf("jwt service: x5c header missing or invalid")
	}

	certDER, err := base64.StdEncoding.DecodeString(x5c[0].(string))
	if err != nil {
		return nil, fmt.Errorf("jwt service: failed to decode certificate: %w", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, fmt.Errorf("jwt service: failed to parse certificate: %w", err)
	}

	pubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("jwt service: certificate does not contain RSA public key")
	}

	parsedToken, err := jwt.Parse(clientAssertion, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("jwt service: invalid signing method")
		}
		return pubKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("jwt service: failed to verify JWT signature: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("jwt service: invalid claims format")
	}

	issuer, _ := claims["iss"].(string)
	if issuer != s.issuer {
		return nil, fmt.Errorf("jwt service: invalid issuer. Expected: %s, got: %s", s.issuer, issuer)
	}

	if time.Now().After(cert.NotAfter) {
		return nil, fmt.Errorf("jwt service: certificate expired")
	}

	if time.Now().Before(cert.NotBefore) {
		return nil, fmt.Errorf("jwt service: certificate not yet valid")
	}

	return claims, nil
}

func (s *JWTService) CreateAccessToken(clientID string) (string, error) {
	claims := jwt.MapClaims{
		"iss": s.issuer,
		"sub": clientID,
		"aud": s.issuer + "/api",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(s.tokenExpiry).Unix(),
		"jti": fmt.Sprintf("%d", time.Now().UnixNano()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

func (s *JWTService) ValidateAccessToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("jwt service: invalid signing method")
		}
		return s.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("jwt service: failed to parse access token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("jwt service: invalid access token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("jwt service: invalid token claims")
	}

	if exp, ok := claims["exp"].(float64); ok && time.Now().Unix() > int64(exp) {
		return nil, fmt.Errorf("jwt service: token expired")
	}

	return claims, nil
}
