package utils

import (
	"encoding/json"
	"math/rand"
	"os"
)

func GetRandomMsg(messages []string) string {
	if len(messages) == 0 {
		return "Array of messages is empty"
	}
	return messages[rand.Intn(len(messages))]
}

func LoadTextMessagges(path string) ([]string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var messages []string
	err = json.Unmarshal(file, &messages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}