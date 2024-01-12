package crypto

import "golang.org/x/crypto/bcrypt"

// HashAndSalt encrypted passwords
func HashAndSalt(pwdStr string) (string, error) {
	pwd := []byte(pwdStr)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	pwdHash := string(hash)
	return pwdHash, nil
}

// ComparePasswords validated password
func ComparePasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		return false
	}
	return true
}
