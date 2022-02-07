package session

import (
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
