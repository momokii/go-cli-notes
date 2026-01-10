package util

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// PasswordHasher handles password hashing using Argon2id
type PasswordHasher struct {
	// Argon2id parameters
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

// NewPasswordHasher creates a new password hasher with secure defaults
func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{
		time:    1,       // Number of iterations
		memory:  64 * 1024, // 64 MB
		threads: 4,       // 4 threads
		keyLen:  32,      // 32-byte key
	}
}

// Hash generates a hash for the given password
func (p *PasswordHasher) Hash(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	// Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}

	// Generate the hash
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		p.time,
		p.memory,
		p.threads,
		p.keyLen,
	)

	// Encode as base64 for storage
	// Format: $argon2id$v=%d$m=%d,t=%d,p=%d$base64(salt+hash)
	combined := append(salt, hash...)
	b64Combined := base64.RawStdEncoding.EncodeToString(combined)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s",
		argon2.Version,
		p.memory,
		p.time,
		p.threads,
		b64Combined,
	)

	return encoded, nil
}

// Verify compares a password with a hash
func (p *PasswordHasher) Verify(password, hash string) (bool, error) {
	if password == "" || hash == "" {
		return false, nil
	}

	// Parse the hash
	// Format: $argon2id$v=%d$m=%d,t=%d,p=%d$base64(salt+hash)
	// Example: $argon2id$v=19$m=65536,t=1,p=4$...
	parts := strings.Split(hash, "$")
	if len(parts) != 5 || parts[0] != "" || parts[1] != "argon2id" {
		return false, fmt.Errorf("invalid hash format")
	}

	// Parse version from parts[2]: v=19
	_, err := fmt.Sscanf(parts[2], "v=%d", new(int))
	if err != nil {
		return false, fmt.Errorf("parse version: %w", err)
	}

	// Parse parameters from parts[3]: m=65536,t=1,p=4
	params := strings.Split(parts[3], ",")
	if len(params) != 3 {
		return false, fmt.Errorf("invalid parameters format: got %d params, expected 3", len(params))
	}

	// Helper function to parse "key=value" format
	parseKeyValue := func(s, key string) (int, error) {
		kv := strings.SplitN(s, "=", 2)
		if len(kv) != 2 || kv[0] != key {
			return 0, fmt.Errorf("invalid key-value format: %s, expected %s=...", s, key)
		}
		var val int
		if _, err := fmt.Sscanf(kv[1], "%d", &val); err != nil {
			return 0, fmt.Errorf("parse value: %w", err)
		}
		return val, nil
	}

	// Parse m=65536 (memory)
	memoryVal, err := parseKeyValue(params[0], "m")
	if err != nil {
		return false, fmt.Errorf("parse memory: %w", err)
	}
	memory := uint32(memoryVal)

	// Parse t=1 (time/iterations)
	timeVal, err := parseKeyValue(params[1], "t")
	if err != nil {
		return false, fmt.Errorf("parse time: %w", err)
	}
	time := uint32(timeVal)

	// Parse p=4 (parallelism/threads)
	threadsVal, err := parseKeyValue(params[2], "p")
	if err != nil {
		return false, fmt.Errorf("parse threads: %w", err)
	}
	threads := uint8(threadsVal)

	// Decode base64 combined salt+hash from parts[4]
	combined, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("decode hash: %w", err)
	}

	if len(combined) < 16 {
		return false, fmt.Errorf("hash too short: got %d bytes", len(combined))
	}

	salt := combined[:16]
	decodedHash := combined[16:]

	// Generate hash for the provided password
	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		time,
		memory,
		threads,
		uint32(len(decodedHash)),
	)

	// Use constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare(decodedHash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}
