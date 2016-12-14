package crypto

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
)

var (
	errMACTooLong = errors.New("mac: too long")
	errMACExpired = errors.New("mac: expired")
	errMACInvalid = errors.New("mac: the value is not valid")
)

const defaultMaxLen = 4096
const macLen = 32

// MACConfig contains all the options to encode or decode a message along with
// a proof of integrity and authenticity.
//
// Key is the secret used for the HMAC key. It should contain at least 16 bytes
// and should be generated by a PRNG.
//
// Name is an optional message name that won't be contained in the MACed
// messaged itself but will be MACed against.
type MACConfig struct {
	Key    []byte
	Name   string
	MaxAge int64
	MaxLen int
}

func assertMACConfig(c *MACConfig) {
	if c.Key == nil {
		panic("hash key is not set")
	}
	if len(c.Key) < 16 {
		panic("hash key is not long enough")
	}
}

// EncodeAuthMessage associates the given value with a message authentication
// code for integrity and authenticity.
//
// If the value, when base64 encoded with a fixed size header is longer than
// the configured maximum length, it will panic.
//
// Message format (name prefix is in MAC but removed from message):
//
//  <------- MAC input ------->
//         <---------- message ---------->
//  | name |    time |  blob  |     hmac |
//  |      | 8 bytes |  ----  | 32 bytes |
//
func EncodeAuthMessage(c *MACConfig, value []byte) ([]byte, error) {
	assertMACConfig(c)

	maxLength := c.MaxLen
	if maxLength == 0 {
		maxLength = defaultMaxLen
	}

	time := Timestamp()

	// Create message with MAC
	size := len(c.Name) + binary.Size(time) + len(value) + macLen
	buf := bytes.NewBuffer(make([]byte, 0, size))
	buf.Write([]byte(c.Name))
	binary.Write(buf, binary.BigEndian, time)
	buf.Write(value)

	// Append mac
	buf.Write(createMAC(c.Key, buf.Bytes()))

	// Skip name
	buf.Next(len(c.Name))

	// Check length
	if base64.URLEncoding.EncodedLen(buf.Len()) > maxLength {
		panic("the value is too long")
	}

	// Encode to base64
	return base64Encode(buf.Bytes()), nil
}

// DecodeAuthMessage verifies a message authentified with message
// authentication code and returns the message value algon with the issued time
// of the message.
func DecodeAuthMessage(c *MACConfig, enc []byte) ([]byte, error) {
	assertMACConfig(c)

	maxLength := c.MaxLen
	if maxLength == 0 {
		maxLength = defaultMaxLen
	}

	// Check length
	if len(enc) > maxLength {
		return nil, errMACTooLong
	}

	// Decode from base64
	dec, err := base64Decode(enc)
	if err != nil {
		return nil, err
	}

	// Prepend name
	dec = append([]byte(c.Name), dec...)

	// Verify message with MAC
	{
		if len(dec) < macLen {
			return nil, errMACInvalid
		}
		var mac []byte
		mac = dec[len(dec)-macLen:]
		dec = dec[:len(dec)-macLen]
		if !verifyMAC(c.Key, dec, mac) {
			return nil, errMACInvalid
		}
	}

	// Skip name prefix
	buf := bytes.NewBuffer(dec)
	buf.Next(len(c.Name))

	// Read time and verify time ranges
	var time int64
	if err = binary.Read(buf, binary.BigEndian, &time); err != nil {
		return nil, errMACInvalid
	}
	if c.MaxAge != 0 && time < Timestamp()-c.MaxAge {
		return nil, errMACExpired
	}

	// Returns the value
	return buf.Bytes(), nil
}

// createMAC creates a MAC with HMAC-SHA256
func createMAC(key, value []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(value)
	return mac.Sum(nil)
}

// verifyMAC returns true is the MAC is valid
func verifyMAC(key, value []byte, mac []byte) bool {
	expectedMAC := createMAC(key, value)
	return hmac.Equal(mac, expectedMAC)
}

// base64Encode encodes a value using base64.
func base64Encode(value []byte) []byte {
	enc := make([]byte, base64.URLEncoding.EncodedLen(len(value)))
	base64.URLEncoding.Encode(enc, value)
	return enc
}

// base64Decode decodes a value using base64.
func base64Decode(value []byte) ([]byte, error) {
	dec := make([]byte, base64.URLEncoding.DecodedLen(len(value)))
	b, err := base64.URLEncoding.Decode(dec, value)
	if err != nil {
		return nil, err
	}
	return dec[:b], nil
}