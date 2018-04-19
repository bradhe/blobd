package crypt

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
)

type Key [32]byte

func NewRandomKey() (*Key, error) {
	var k Key

	if n, err := rand.Reader.Read(k[:]); err != nil {
		return nil, err
	} else if n != 32 {
		return nil, ErrKeyGenerationFailed
	}

	return &k, nil
}

func (k *Key) Bytes() []byte {
	return (*k)[:]
}

func (k *Key) MarshalJSON() ([]byte, error) {
	return json.Marshal(base64.StdEncoding.EncodeToString(k.Bytes()))
}

func (k *Key) UnmarshalJSON(buf []byte) error {
	var str string

	if err := json.Unmarshal(buf, &str); err != nil {
		return err
	}

	if buf, err := base64.StdEncoding.DecodeString(str); err != nil {
		return err
	} else {
		copy(k[:], buf)
	}

	return nil
}

func (k *Key) Destroy() {
	for i := range k {
		k[i] = 0x00
	}
}
