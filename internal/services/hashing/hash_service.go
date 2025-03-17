package hashing

type IHashService interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}
