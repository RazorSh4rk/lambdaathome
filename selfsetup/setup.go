package setup

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

var PassCode string

func passfilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(home, ".passfile")
}

func Setup() {
	path := passfilePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		passCode := uuid.New().String()
		f.Write([]byte(passCode))
		PassCode = passCode
		fmt.Printf("Set up first time passcode, copy this somewhere to access the server later: %s\n", passCode)
	} else {
		f, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		PassCode = string(f)
	}
}
