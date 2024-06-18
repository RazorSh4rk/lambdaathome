package setup

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

var PassCode string

func Setup() {
	if _, err := os.Stat("passfile"); os.IsNotExist(err) {
		f, err := os.OpenFile("passfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		passCode := uuid.New().String()
		f.Write([]byte(passCode))
		PassCode = passCode
		fmt.Printf("Set up first time passcode, copy this somewhere to access the server later: %s\n", passCode)
	} else {
		f, err := os.ReadFile("passfile")
		if err != nil {
			log.Fatal(err)
		}
		PassCode = string(f)
	}
}
