package blobs

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"sync/atomic"
	"time"
)

func init() {
	go resetTick()
}

var tickAcc uint32

func resetTick() {
	for _ = range time.Tick(1 * time.Second) {
		atomic.StoreUint32(&tickAcc, 0)
	}
}

func tick() []byte {
	t := atomic.AddUint32(&tickAcc, 1)

	b := [4]byte{}
	b[0] = byte(t & 0xff)
	b[1] = byte((t >> 8) & 0xff)
	b[2] = byte((t >> 16) & 0xff)
	b[3] = byte((t >> 24) & 0xff)
	return b[:]
}

var clock = func() int64 {
	return time.Now().UTC().UnixNano()
}

func int64bits(i int64) []byte {
	b := [8]byte{}
	b[0] = byte(i & 0xff)
	b[1] = byte((i >> 8) & 0xff)
	b[2] = byte((i >> 16) & 0xff)
	b[3] = byte((i >> 24) & 0xff)
	b[4] = byte((i >> 32) & 0xff)
	b[5] = byte((i >> 40) & 0xff)
	b[6] = byte((i >> 48) & 0xff)
	b[7] = byte((i >> 56) & 0xff)
	return b[:]
}

func entropy() []byte {
	var buf [2]byte
	rand.Reader.Read(buf[:])
	return buf[:]
}

type Id [14]byte

func (id Id) String() string {
	return base64.URLEncoding.EncodeToString(id[:])
}

func (id Id) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *Id) UnmarshalJSON(buf []byte) error {
	var str string

	if err := json.Unmarshal(buf, &str); err != nil {
		return err
	}

	if buf, err := base64.URLEncoding.DecodeString(str); err != nil {
		return err
	} else {
		copy(id[:], buf)
	}

	return nil
}

func (id Id) IsEmpty() bool {
	for _, b := range id {
		if b != 0x00 {
			return false
		}
	}

	return true
}

func (left Id) Equal(right Id) bool {
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}

	return true
}

func NewId() Id {
	id := Id{}
	copy(id[:], int64bits(clock()))
	copy(id[8:], entropy())
	copy(id[10:], tick())
	return id
}

func ParseId(str string) (Id, error) {
	var id Id

	if buf, err := base64.URLEncoding.DecodeString(str); err != nil {
		return Id{}, err
	} else {
		copy(id[:], buf)
	}

	return id, nil
}
