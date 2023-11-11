package password

import (
	"encoding/hex"
	"golang.org/x/crypto/argon2"
)

const (
	keySize     = 32        // Size of the derived key
	timeCost    = 1         // Number of iterations
	memory      = 64 * 1024 // Memory usage in KiB
	parallelism = 4
)

type ArgonPasswordService struct {
}

func NewArgonPasswordService() *ArgonPasswordService {
	return &ArgonPasswordService{}
}

func (service *ArgonPasswordService) GenerateHashedPassword(password string, salt []byte) string {
	return hex.EncodeToString(argon2.IDKey([]byte(password), salt, uint32(timeCost), uint32(memory), uint8(parallelism), uint32(keySize)))
}
