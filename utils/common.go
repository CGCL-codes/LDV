package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	mrand "math/rand"
	"os"

	"bufio"
	"github.com/libp2p/go-libp2p-crypto"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func CheckFileExisted(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func ReadStringsFromFile(path string) ([]string, error) {
	if !CheckFileExisted(path) {
		return nil, errors.New(fmt.Sprintf("File %s does not exsit", path))
	} else {
		var ss []string
		file, err := os.Open(path)
		if err!=nil {
			return nil, err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			ss = append(ss, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		return ss, nil
	}
}

// get the private key from a file
// or generate a new private key
func GetPrivKey(createPK bool, path string, seed int) (crypto.PrivKey, error) {
	if createPK {
		log.Println("Create new private key")
		// create new private key
		var r io.Reader
		if seed == 0 {
			r = rand.Reader
		} else {
			r = mrand.New(mrand.NewSource(int64(seed)))
		}
		pk, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
		if err != nil {
			return nil, err
		}

		// write private key to the file for later use
		if CheckFileExisted(path) {
			os.Rename(path, path+".bak")
		}
		f, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		mpk, err := crypto.MarshalPrivateKey(pk)
		if err != nil {
			return nil, err
		}
		_, err = f.Write(mpk)
		if err != nil {
			return nil, err
		}
		return pk, nil
	} else {
		// read private key from the file
		log.Printf("Read private key from it!\n", path)
		data, err := ioutil.ReadFile(path)
		if err!=nil {
			return nil, err
		}
		pk, err := crypto.UnmarshalPrivateKey(data)
		if err!=nil {
			return nil, err
		}
		return pk, nil
	}
}