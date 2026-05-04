package common

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword bcrypt加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HashMachineCode bcrypt加密机器码
func HashMachineCode(machineCode string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(machineCode), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckMachineCode 验证机器码
func CheckMachineCode(machineCode, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(machineCode))
	return err == nil
}
