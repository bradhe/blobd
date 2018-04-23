package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

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

	// Time this claim will expire
	EXP int64 `json:"exp"`

	// The id for this blob.
	BlobId blobs.Id `json:"sub"`

	// Key used to encrypt the message.
	Key *crypt.Key `json:"key"`

	// Expected media type of the content in storage.
	MediaType string `json:"media_type"`
}

func (bc *BlobClaims) IsReadable() bool {
	switch bc.Type {
	case ReadOnlyToken, WritableToken:
		return true
	default:
		return false
	}
}

func (bc *BlobClaims) Destroy() {
	if bc.Key != nil {
		bc.Key.Destroy()
	}
}

func (bc *BlobClaims) IsWritable() bool {
	switch bc.Type {
	case ReadOnlyToken:
		return false
	case WritableToken:
		return true
	default:
		return false
	}
}

func (bc *BlobClaims) Valid() error {
	if bc.BlobId.IsEmpty() {
		return ErrInvalidBlobId
	}

	if bc.Key == nil {
		return ErrMissingDecryptionKey
	}

	t := now()

	// we don't use >= because if the JWT was created and then immediately used,
	// that's valid. specifically in tests.
	if bc.NBF > t {
		return ErrInvalidJWT
	}

	// operator matters less here as default expiration window is pretty far in
	// the future.
	if bc.EXP < t {
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
	return time.Now().UTC().Unix()
}

func GenerateJWT(t TokenType, key *crypt.Key, blob *blobs.Blob) string {
	claims := &BlobClaims{
		Type:      t,
		NBF:       now(),
		EXP:       blob.Expiration().Unix(),
		BlobId:    blob.Id,
		Key:       key,
		MediaType: blob.MediaType,
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
