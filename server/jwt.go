package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pborman/uuid"

	"github.com/bradhe/blobd/blobs"
	"github.com/bradhe/blobd/crypt"
)

const (
	DefaultTokenSigningKey = "this is a secret"
)

var (
	TokenSigningKey string = DefaultTokenSigningKey
)

type TokenType string

var (
	WritableToken = TokenType("writable")
	ReadOnlyToken = TokenType("read-only")
)

type BlobClaims struct {
	// Type of claim that we've generated.
	Type TokenType `json:"aud"`

	// Time the claim was generated
	NBF int64 `json:"nbf"`

	// The UUID for this blob.
	BlobId uuid.UUID `json:"sub"`

	// Key used to encrypt the message.
	Key *crypt.Key `json:"key"`
}

func (bc *BlobClaims) Valid() error {
	if uuid.Equal(bc.BlobId, uuid.NIL) {
		return ErrInvalidUUID
	}

	if bc.Key == nil {
		return ErrMissingDecryptionKey
	}

	// we don't use >= because if the JWT was created and then immediately used,
	// that's valid. specifically in tests.
	if bc.NBF > now() {
		return ErrInvalidJWT
	}

	switch bc.Type {
	case WritableToken, ReadOnlyToken:
		// Do nothing. Is valid.
	default:
		return ErrInvalidJWT
	}

	return nil
}

func GetJWT(str string) string {
	return strings.TrimPrefix(str, "Bearer ")
}

var now = func() int64 {
	return time.Now().Unix()
}

func GenerateJWT(t TokenType, key *crypt.Key, blob *blobs.Blob) string {
	claims := &BlobClaims{
		Type:   t,
		NBF:    now(),
		BlobId: blob.Id,
		Key:    key,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	str, err := token.SignedString([]byte(TokenSigningKey))

	if err != nil {
		panic(err)
	}

	return str
}

func ParseJWT(authorization string) (*jwt.Token, *BlobClaims, error) {
	str := GetJWT(authorization)

	token, err := jwt.ParseWithClaims(str, &BlobClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("server: unexpected signing method `%v`", token.Header["alg"])
		}

		return []byte(TokenSigningKey), nil
	})

	if err != nil {
		return nil, nil, err
	}

	if !token.Valid {
		return nil, nil, ErrInvalidJWT
	}

	if claims, ok := token.Claims.(*BlobClaims); ok {
		return token, claims, nil
	}

	return nil, nil, ErrInvalidJWT
}
