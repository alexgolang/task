package main

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	keyData, err := os.ReadFile("../secret/client.key")
	if err != nil {
		log.Fatal(err)
	}
	block, _ := pem.Decode(keyData)
	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	certData, err := os.ReadFile("../secret/client.crt")
	if err != nil {
		log.Fatal(err)
	}
	certBlock, _ := pem.Decode(certData)
	certB64 := base64.StdEncoding.EncodeToString(certBlock.Bytes)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "ishare-task-api",
		"sub": "test-client-1",
		"aud": "http://localhost:8080/token",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
		"jti": "test-jwt-123",
	})

	token.Header["x5c"] = []string{certB64}

	tokenString, _ := token.SignedString(privKey)
	fmt.Println("Client assertion JWT:")
	fmt.Println(tokenString)
}
