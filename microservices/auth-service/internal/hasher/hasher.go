package hasher

import "golang.org/x/crypto/bcrypt"

type Hasher struct {
	cost int
}

func New() *Hasher {
	return &Hasher{cost: bcrypt.DefaultCost}
}

func (h *Hasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (h *Hasher) Compare(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
