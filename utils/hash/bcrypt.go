package hash

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// 密码加密及验证密码的方法

func BcryptMake(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func BcryptMakeCheck(pwd []byte, hashedPwd string) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, pwd)
	if err != nil {
		return false
	}
	return true
}
