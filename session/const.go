package session

import (
	"bytes"
	"encoding/gob"

	"github.com/google/uuid"
)

// SameSiteType cookie session same site constants
const (
	// SameSiteLax lax same site mode
	SameSiteLax string = "Lax"
	// SameSiteStrict strict same site mode
	SameSiteStrict = "Strict"
	// SameSiteNone none same site mode
	SameSiteNone = "None"
)

// UUIDGenerator Generate id using uuid
func UUIDGenerator() string {
	return uuid.New().String()
}

func GetBytes(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
