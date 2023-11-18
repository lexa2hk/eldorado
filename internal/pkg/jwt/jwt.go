package jwt

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/romankravchuk/eldorado/internal/data"
)

func CreateToken(payload *data.TokenPayload, ttl time.Duration, key *rsa.PrivateKey) (*data.TokenDetails, error) {
	now := time.Now().UTC()
	td := &data.TokenDetails{
		Payload:   *payload,
		ExpiresAt: now.Add(ttl).Unix(),
	}

	claims := data.Claims{
		TokenID: payload.ID,
		UserID:  payload.UserID,
		Email:   payload.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        payload.ID,
		},
	}

	var err error

	td.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, &claims).SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	return td, nil
}

func ValidateToken(token string, key *rsa.PublicKey) (*data.TokenPayload, error) {
	parsedToken, err := jwt.ParseWithClaims(
		token,
		&data.Claims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return key, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*data.Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	payload := &data.TokenPayload{
		ID:     claims.TokenID,
		UserID: claims.UserID,
		Email:  claims.Email,
	}

	return payload, nil
}

func ParseKeyPairs(strPem, strPub string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	decodedPem, err := base64.StdEncoding.DecodeString(strPem)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	decodedPub, err := base64.StdEncoding.DecodeString(strPub)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	pem, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPem)
	if err != nil {
		return nil, nil, err
	}

	pub, err := jwt.ParseRSAPublicKeyFromPEM(decodedPub)
	if err != nil {
		return nil, nil, err
	}

	return pem, pub, nil
}
