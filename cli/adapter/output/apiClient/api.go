package apiClient

import (
	"io/ioutil"
	"log"
	"net/http"
)

var baseURL = "http://localhost:9001/lic-api"
var client = &http.Client{}

func GetCharacterInfo(session string, characterId string) string {
	req, _ := http.NewRequest("GET", baseURL+"/characters/"+characterId, nil)
	req.Header.Set("Authorization", "Bearer "+session)
	charResp, _ := client.Do(req)

	charBody, err := ioutil.ReadAll(charResp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(charBody)
}

func GetCharacterSkills(session string, characterId string) string {
	req, _ := http.NewRequest("GET", baseURL+"/characters/"+characterId+"/skills", nil)
	req.Header.Set("Authorization", "Bearer "+session)
	charResp, _ := client.Do(req)

	charBody, err := ioutil.ReadAll(charResp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(charBody)
}
