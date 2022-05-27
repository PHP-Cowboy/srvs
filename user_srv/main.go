package main

import (
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"strings"
)

type Carer interface {
	run()
	start()
}

type BYD struct{}

func (b *BYD) run() {
	fmt.Println("run")
}

func (b BYD) start() {
	fmt.Println("start")
}

func GeneratePwd(pwd string) {

	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(pwd, options)

	newPwd := fmt.Sprintf("pbkdf2-sha512$%s$%s", salt, encodedPwd)

	fmt.Println(salt)
	fmt.Println(encodedPwd)
	fmt.Println(newPwd)

	pwdSlice := strings.Split(newPwd, "$")

	fmt.Println(pwdSlice)

	check := password.Verify(pwd, pwdSlice[1], pwdSlice[2], options)
	fmt.Println(check) // true
}

func main() {
	GeneratePwd("sgagrrgweg")
}
