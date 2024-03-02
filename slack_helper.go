package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func SendMessageToSlack(message string) error {
	postUrl := "hook_url_here"

	body := []byte(fmt.Sprintf(`{"text":"asdasd"}`))

	post, err := http.Post(postUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer post.Body.Close()
	return nil
}
