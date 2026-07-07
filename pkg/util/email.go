package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateVerificationCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	//确保生成的是6位数,不足向前补足0
	return fmt.Sprintf("%06d", n.Int64()), nil
}
