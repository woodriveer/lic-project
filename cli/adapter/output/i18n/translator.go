package i18n

import (
	"encoding/json"
	"log"
	"os"
)

var translator map[string]string

func Localize(messageId string, region string) string {
	dir, dirErr := os.Getwd()
	if dirErr != nil {
		log.Fatal(dirErr)
	}
	file, err := os.Open(dir + "/i18n/" + region + "/strings.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var data map[string]string
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		log.Fatal(err)
	}

	return data[messageId]
}
