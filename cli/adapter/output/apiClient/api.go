package client

import (
	"io/ioutil"
	"log"
	"net/http"
)

func GetCharacterInfo(session string, characterId string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:9001/lic-api/characters/"+characterId, nil)
	req.Header.Set("Authorization", "Bearer "+session)
	charResp, _ := client.Do(req)

	charBody, err := ioutil.ReadAll(charResp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(charBody)
}
