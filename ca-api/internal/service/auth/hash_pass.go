package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/dhruvpatel-10/signee/ca-api/internal/domain/auth"
	"golang.org/x/crypto/argon2"
)

var DefaultPassConfig = auth.PassConfig{
	Time:    3,
	Memory:  64 * 1024,
	Threads: 4,
	KeyLen:  32,
}

// Password verification function
func verifyPassword(password, encodedHash string) (bool, error) {
	// Parse PHC format: $argon2id$v=19$m=65536,t=1,p=4$salt$hash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, fmt.Errorf("invalid hash format")
	}

	// Parse version
	version, err := strconv.Atoi(strings.TrimPrefix(parts[2], "v="))
	if err != nil {
		return false, err
	}

	// Parse parameters
	params := strings.Split(parts[3], ",")
	if len(params) != 3 {
		return false, fmt.Errorf("invalid parameters format")
	}

	memory, err := strconv.Atoi(strings.TrimPrefix(params[0], "m="))
	if err != nil {
		return false, err
	}

	time, err := strconv.Atoi(strings.TrimPrefix(params[1], "t="))
	if err != nil {
		return false, err
	}

	threads, err := strconv.Atoi(strings.TrimPrefix(params[2], "p="))
	if err != nil {
		return false, err
	}

	// Decode salt and hash
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	storedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	// Verify version matches
	if version != argon2.Version {
		return false, fmt.Errorf("argon2 version mismatch")
	}

	// Hash the input password with extracted parameters
	inputHash := argon2.IDKey([]byte(password), salt, uint32(time), uint32(memory), uint8(threads), uint32(len(storedHash)))

	// Constant-time comparison
	return subtle.ConstantTimeCompare(storedHash, inputHash) == 1, nil
}

// Salting
func generateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	return salt, err
}

// Hashing password with argon2id
func hashPassword(password string) (string, error) {
	salt, err := generateSalt(32)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, DefaultPassConfig.Time, DefaultPassConfig.Memory, DefaultPassConfig.Threads, DefaultPassConfig.KeyLen)

	// Store in PHC string format
	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, DefaultPassConfig.Memory, DefaultPassConfig.Time, DefaultPassConfig.Threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash))

	return encoded, nil
}

// tests moved to hash_pass_test.go
