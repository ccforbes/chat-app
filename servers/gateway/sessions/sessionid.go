package sessions

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
)

//InvalidSessionID represents an empty, invalid session ID
const InvalidSessionID SessionID = ""

//idLength is the length of the ID portion
const idLength = 32

//signedLength is the full length of the signed session ID
//(ID portion plus signature)
const signedLength = idLength + sha256.Size

//SessionID represents a valid, digitally-signed session ID.
//This is a base64 URL encoded string created from a byte slice
//where the first `idLength` bytes are crytographically random
//bytes representing the unique session ID, and the remaining bytes
//are an HMAC hash of those ID bytes (i.e., a digital signature).
//The byte slice layout is like so:
//+-----------------------------------------------------+
//|...32 crypto random bytes...|HMAC hash of those bytes|
//+-----------------------------------------------------+
type SessionID string

//ErrInvalidID is returned when an invalid session id is passed to ValidateID()
var ErrInvalidID = errors.New("Invalid Session ID")

//NewSessionID creates and returns a new digitally-signed session ID,
//using `signingKey` as the HMAC signing key. An error is returned only
//if there was an error generating random bytes for the session ID
func NewSessionID(signingKey string) (SessionID, error) {

	if len(signingKey) == 0 {
		err := errors.New("The key can not be empty")
		return InvalidSessionID, err
	}

	salt := make([]byte, idLength)
	rand.Read(salt)

	key := []byte(signingKey)
	h := hmac.New(sha256.New, key)
	h.Write(salt)
	session := h.Sum(salt)

	//session := append(salt, h.Sum(nil)...)

	sessionID := base64.URLEncoding.EncodeToString(session)

	return SessionID(sessionID), nil
}

//ValidateID validates the string in the `id` parameter
//using the `signingKey` as the HMAC signing key
//and returns an error if invalid, or a SessionID if valid
func ValidateID(id string, signingKey string) (SessionID, error) {

	decode, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		return InvalidSessionID, errors.New("Error decoding ID")
	}

	salt := decode[:idLength]
	key := []byte(signingKey)
	h := hmac.New(sha256.New, key)
	h.Write(salt)
	hash1 := h.Sum(nil)
	hash2 := decode[idLength:]
	if subtle.ConstantTimeCompare(hash1, hash2) != 1 {
		return InvalidSessionID, ErrInvalidID
	}
	return SessionID(id), nil
}

//String returns a string representation of the sessionID
func (sid SessionID) String() string {
	return string(sid)
}
