package salt

import "crypto/rand"

const SaltSize = 16

type SaltFactory struct {
}

func NewSaltFactory() *SaltFactory {
	return &SaltFactory{}
}

func (factory *SaltFactory) GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}
