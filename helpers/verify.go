package helpers

import (

	// # STDLIBS

	"crypto/sha256"
	"encoding/hex"
	"log"

	// # EXTERNAL LIBS
	"golang.org/x/crypto/bcrypt"
)

func Verify(_hashedPassword string, _password string) bool {

	_byteHash := []byte(_hashedPassword)
	_h := sha256.Sum256([]byte(_password))
	endcodeHex := hex.EncodeToString(_h[:])

	err := bcrypt.CompareHashAndPassword(_byteHash, []byte(endcodeHex))

	if err != nil {

		log.Println("Libs > [Bcrypt::ERROR] > ", err.Error())
		return false

	}

	return true

}
