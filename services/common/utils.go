package common

import (
	"fmt"
	"math/rand"
)


func  GenerateVerificationCode() string{
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func GenerateRandomId(folder  string)  string{
	return fmt.Sprintf("%012d", rand.Intn(1000000000000))
}

