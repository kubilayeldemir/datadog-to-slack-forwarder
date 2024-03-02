package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"strings"
)

var base64EncodedBasicAuthenticationSecret = ""

func main() {
	base64EncodedBasicAuthenticationSecret = base64.StdEncoding.EncodeToString([]byte("DatadogLogging:DatadogLoggingPw"))
	http.HandleFunc("/", forwardLogsToSlack)
	port := "8080"
	log.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func forwardLogsToSlack(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}

	isAuthenticated := IsRequestAuthenticated(r.Header.Get("Authorization"))
	if !isAuthenticated {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var bodyByte []byte
	if r.Header.Get("Content-Encoding") == "gzip" {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		bodyByte, err = DecompressGzip(data)
		if err != nil {
			http.Error(w, "Error decompressing request body", http.StatusInternalServerError)
			return
		}
	} else {
		var err error
		bodyByte, err = io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
	}

	_ = SendMessageToSlack(string(bodyByte))
}

func DecompressGzip(data []byte) ([]byte, error) {
	gr, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	return io.ReadAll(gr)
}

func IsRequestAuthenticated(authorizationHeader string) bool {
	if authorizationHeader == "" {
		return false
	}
	authSplit := strings.Split(authorizationHeader, " ")
	if len(authSplit) != 2 || authSplit[0] == "" || authSplit[1] == "" {
		return false
	}
	authMethod := authSplit[0]
	authCredentials := authSplit[1]

	switch authMethod {
	case "Basic":
		if authCredentials == base64EncodedBasicAuthenticationSecret {
			return true
		}
	default:
		return false
	}
	return false
}
