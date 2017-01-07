//go:generate ./gen-pass.sh && go-bindata -pkg handler -prefix view/ -ignore \.*\.less -o handler/assets.go view/...

package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"simas/model"
	"strings"
)

var (
	startConfig = flag.Bool("config", false, "Menjalankan proses konfigurasi aplikasi")
	portNumber  = flag.Int("p", 8081, "Port yang digunakan oleh aplikasi")
)

func main() {
	// Parse flags
	flag.Parse()

	// Check if user want to configure
	if *startConfig {
		startConfiguration()
		return
	}

	// Load configuration file
	configFile, err := ioutil.ReadFile("./config")
	if err != nil {
		log.Fatalln("Lakukan konfigurasi terlebih dahulu")
		os.Exit(1)
	}

	// Decrypt configuration file
	decrypted, err := decrypt([]byte(configPassword), configFile)
	checkError(err)

	// Decode configuration
	config := model.Configuration{}
	buffer := bytes.NewBuffer(decrypted)
	err = gob.NewDecoder(buffer).Decode(&config)
	checkError(err)

	// Create file directory if needed
	err = os.MkdirAll(config.FileDirectory, os.ModePerm)
	checkError(err)

	// Create backend
	backEnd := NewBackEnd(*portNumber, config)
	defer backEnd.Close()

	// Serve app
	backEnd.ServeApp()
}

func startConfiguration() {
	// Accept configuration input from user
	config := model.Configuration{}

	fmt.Print("01/11", "\t", "Nama domain yang digunakan (contoh www.simas.com) :", "\n\t")
	fmt.Scanln(&config.AppDomain)

	fmt.Print("\n", "02/11", "\t", "Nama user database (contoh root) :", "\n\t")
	fmt.Scanln(&config.DatabaseUser)

	fmt.Print("\n", "03/11", "\t", "Password user database :", "\n\t")
	fmt.Scanln(&config.DatabasePassword)

	fmt.Print("\n", "04/11", "\t", "Nama database :", "\n\t")
	fmt.Scanln(&config.DatabaseName)

	fmt.Print("\n", "05/11", "\t", "User key Zenziva untuk SMS gateway :", "\n\t")
	fmt.Scanln(&config.ZenzivaUserKey)

	fmt.Print("\n", "06/11", "\t", "Pass key Zenziva untuk SMS gateway :", "\n\t")
	fmt.Scanln(&config.ZenzivaPassKey)

	fmt.Print("\n", "07/11", "\t", "Alamat email yang digunakan untuk email gateway (contoh m.radhi.f@gmail.com):", "\n\t")
	fmt.Scanln(&config.EmailAddress)

	fmt.Print("\n", "08/11", "\t", "Password alamat email yang digunakan :", "\n\t")
	fmt.Scanln(&config.EmailPassword)

	fmt.Print("\n", "09/11", "\t", "Server email yang digunakan (contoh smtp.gmail.com) :", "\n\t")
	fmt.Scanln(&config.EmailServer)

	fmt.Print("\n", "10/11", "\t", "Port server email yang digunakan (contoh 587 untuk Gmail) :", "\n\t")
	fmt.Scanln(&config.EmailServerPort)

	fmt.Print("\n", "11/11", "\t", "Direktori untuk menyimpan file surat yang diupload (contoh /home/imageDir) :", "\n\t")
	fmt.Scanln(&config.FileDirectory)

	// Remove trailing path from file directory
	fileDir := strings.TrimSpace(config.FileDirectory)
	if fileDir[len(fileDir)-1:] == "/" {
		fileDir = fileDir[:len(fileDir)-1]
	}
	config.FileDirectory = fileDir

	// Encrypt configuration
	buffer := bytes.Buffer{}
	err := gob.NewEncoder(&buffer).Encode(&config)
	checkError(err)

	encrypted, err := encrypt([]byte(configPassword), buffer.Bytes())
	checkError(err)

	// Save config to file
	configFile, _ := os.Create("./config")
	defer configFile.Close()

	_, err = configFile.Write(encrypted)
	checkError(err)
}

func encrypt(key, value []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(value))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], value)

	return ciphertext, nil
}

func decrypt(key, value []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(value) < aes.BlockSize {
		return nil, errors.New("Cipher text is too short")
	}

	iv := value[:aes.BlockSize]
	value = value[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(value, value)

	return value, nil
}

func checkError(err error) {
	if err != nil && err != sql.ErrNoRows {
		log.Fatalln("Error:", err)
	}
}
