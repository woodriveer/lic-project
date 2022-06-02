// Package cmd
/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/happierall/l"
	"github.com/spf13/cobra"
	"github.com/yanBrandao/lic-cli/adapter/output/apiClient"
	"github.com/yanBrandao/lic-cli/adapter/output/banner"
	"github.com/yanBrandao/lic-cli/adapter/output/i18n"
	"io/ioutil"
	_ "io/ioutil"
	"log"
	_ "log"
	"math"
	"math/rand"
	"net/http"
	_ "net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type Characters struct {
	Characters []Character `json:"characters"`
}

type Character struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Level    int    `json:"level"`
	Class    string `json:"class"`
	Location string `json:"location"`
}

type Dungeon struct {
	Name       string
	Progress   []string
	Difficulty string
	StartTime  time.Time
	EndTime    time.Time
}

var clear map[string]func() //create a map for storing clear func
var globalSession string

const region = "pt-BR"

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return
		}
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return
		}
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		startGameFlow(nil, nil)
	},
}

func startGameFlow(session *string, paramName *string) {
	var token = new(Token)
	var name = paramName

	if session == nil {
		name = renderLogin(token)
		session = &token.AccessToken
		globalSession = *session
	}
	if session != nil {
		CallClear()
		l.Printf("%s %s\n", i18n.Localize("welcome", region), l.Colorize(*name, l.Green))

		var characterAnswers = renderCharacterSelect(session)

		if strings.Contains(characterAnswers.CharacterId, i18n.Localize("new_character", region)) {
			renderCharacterCreate()
			startGameFlow(session, name)
		} else {
			l.Printf("Carregando informação de personagem")
			r := <-loading()
			l.Printf("%d", r)

			charBody := apiClient.GetCharacterInfo(*session, characterAnswers.CharacterId)

			var response Character

			jsonError := json.Unmarshal([]byte(charBody), &response)
			if jsonError != nil {
				return
			}

			renderGameActions(response)

		}

	}
}

func loading() <-chan int32 {
	r := make(chan int32)

	go func() {
		defer close(r)
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
		err := s.Color("yellow")
		if err != nil {
			return
		}
		s.Prefix = "Loading items: "       // Prefix text before the spinner
		s.Start()                          // Start the spinner
		time.Sleep(100 * time.Millisecond) // Run for some time to simulate work
		s.Prefix = "Loading location: "
		time.Sleep(100 * time.Millisecond)
		s.Prefix = "Loading coffee: "
		time.Sleep(100 * time.Millisecond)

		s.Stop() // Stop the spinner
		r <- rand.Int31n(100)
	}()

	return r
}

func loadingDungeon() string {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
	err := s.Color("yellow")
	if err != nil {
		return ""
	}
	s.Prefix = "Selecting random dungeon: " // Prefix text before the spinner
	s.Start()                               // Start the spinner
	time.Sleep(100 * time.Millisecond)      // Run for some time to simulate work
	var dungeon = "Dragon's Cave"
	s.Prefix = "Loading location (" + dungeon + "): "
	time.Sleep(100 * time.Millisecond)
	s.Prefix = "Loading monsters (" + dungeon + "): "
	time.Sleep(100 * time.Millisecond)
	s.Prefix = "Loading knifes and blood (" + dungeon + "): "
	time.Sleep(100 * time.Millisecond)

	s.Stop() // Stop the spinner

	println("")
	return dungeon
}

func renderCharacterBar(character Character) {
	renderCharacterBarWithClear(character, true)
}

func renderCharacterBarWithClear(character Character, isToCallClear bool) {
	if isToCallClear {
		CallClear()
	}

	l.Printf("Localização: %s \t | \t Personagem: %s \t | \t Level: %s \t | \t Classe: %s\n",
		l.Colorize(character.Location, l.Blue),
		l.Colorize(character.Name, l.Green),
		l.Colorize(strconv.Itoa(character.Level), l.Yellow),
		l.Colorize(character.Class, l.Gray),
	)
}

func renderGameActions(character Character) {

	renderCharacterBar(character)

	l.Printf("O que deseja fazer?")

	var gameActions []string

	//gameActions = append(gameActions, "Visualizar conquistas")
	gameActions = append(gameActions, "Iniciar uma caverna")
	gameActions = append(gameActions, "Visualizar inventário")
	//gameActions = append(gameActions, "Visualizar registro de batalha")
	//gameActions = append(gameActions, "Entrar em chats abertos")

	var actionsQuestion = []*survey.Question{
		{
			Name: "action",
			Prompt: &survey.Select{
				Message: "Escolha uma ação para iniciar:",
				Options: gameActions,
			},
		},
	}

	actionAnswer := struct {
		Action string
	}{}

	err := survey.Ask(actionsQuestion, &actionAnswer)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	switch actionAnswer.Action {
	case "Iniciar uma caverna":
		renderStartDungeon(character)
	}
}

func renderStartDungeon(character Character) {
	var currentDungeon = loadingDungeon()
	var battleLog []string
	character.Location = currentDungeon
	var bannerDraw = banner.LoadBanner("1")

	var dungeonProgress []string
	dungeonProgress = append(dungeonProgress, "???")
	dungeonProgress = append(dungeonProgress, "???")
	dungeonProgress = append(dungeonProgress, "???")
	dungeonProgress = append(dungeonProgress, "Boss")

	var dungeon = Dungeon{
		Name:       currentDungeon,
		Progress:   dungeonProgress,
		Difficulty: "Normal",
		StartTime:  time.Now(),
	}

	renderCharacterBar(character)
	banner.PrintBanner(bannerDraw, "red")

	renderDungeonIntroduction("1")
	renderDungeonProgress(dungeon)
	renderDungeonActions(dungeon)

	var firstInteraction = true
	for hasPendingProgress(dungeon) {
		renderCharacterBar(character)
		banner.PrintBanner(bannerDraw, "red")
		renderDungeonProgress(dungeon)
		renderBattleLog(&battleLog)
		if !firstInteraction {
			renderDungeonActions(dungeon)
		}
		renderPlayerInteraction(dungeon, character, &battleLog)
		firstInteraction = false

	}

	renderCharacterBar(character)
	banner.PrintBanner(bannerDraw, "red")
	renderDungeonProgress(dungeon)
	renderBattleLog(&battleLog)
	renderDungeonActions(dungeon)

	renderDungeonEnding(dungeon, character)

}

type Monster struct {
	Name            string
	Health          float64
	Attack          float64
	Defense         float64
	ElementalResist Resistances
}

type Resistances struct {
	Fire  float64
	Water float64
	Wind  float64
	Earth float64
	Light float64
	Dark  float64
}

type Skill struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
	Power int    `json:"power"`
	Type  string `json:"type"`
}

type CharacterSkills struct {
	Skills []Skill
}

func renderPlayerInteraction(dungeon Dungeon, character Character, battleLog *[]string) {
	//identify dificulty dungeon level
	var currentInteraction string
	for i := 0; i < len(dungeon.Progress); i++ {
		if dungeon.Progress[i] == "???" || dungeon.Progress[i] == "Boss" {
			currentInteraction = dungeon.Progress[i-1]
			break
		}
	}
	l.Printf("%s apareceu com uma fome descomunal e uma fúria excedente", currentInteraction)
	*battleLog = append(*battleLog, fmt.Sprintf("%s apareceu com uma fome descomunal e uma fúria excedente", currentInteraction))
	var currentMonster = Monster{
		Name:    currentInteraction,
		Health:  10,
		Attack:  2,
		Defense: 2,
		ElementalResist: Resistances{
			Fire:  100,
			Water: 0,
			Wind:  30,
			Earth: 30,
			Light: 50,
			Dark:  50,
		},
	}

	for isMonsterAlive(currentMonster) {
		renderCharacterBar(character)
		renderBattleLog(battleLog)
		renderMonsterInfo(currentMonster)

		charBody := apiClient.GetCharacterSkills(globalSession, strconv.Itoa(character.Id))

		var response CharacterSkills

		jsonError := json.Unmarshal([]byte(charBody), &response)
		if jsonError != nil {
			return
		}

		var skillsAvailable []string

		for i := 0; i < len(response.Skills); i++ {
			skillsAvailable = append(skillsAvailable, response.Skills[i].Name+" L: "+strconv.Itoa(response.Skills[i].Level))
		}

		var answer string

		var skillsQuestion = []*survey.Question{
			{
				Name: "skillSelected",
				Prompt: &survey.Select{
					Message: "O que deseja fazer",
					Options: skillsAvailable,
				},
			},
		}

		err2 := survey.Ask(skillsQuestion, &answer)
		if err2 != nil {
			fmt.Println(err2.Error())
			return
		}

		l.Printf("%s escolheu atacar você...", currentMonster.Name)
		*battleLog = append(*battleLog, fmt.Sprintf("%s escolheu atacar você...", currentMonster.Name))

		l.Printf("Você atacou o %s utilizando %s e infringiu %f em seus pontos de vida.", currentMonster.Name, answer, 5.0)
		*battleLog = append(*battleLog, fmt.Sprintf("Você atacou o %s utilizando %s e infringiu %f em seus pontos de vida.", currentMonster.Name, answer, 5.0))

		currentMonster.Health -= 5.0

		if currentMonster.Health > 0 {
			l.Printf("%s atacou e você perdeu X pontos de vida.", currentMonster.Name)
			*battleLog = append(*battleLog, fmt.Sprintf("%s atacou e você perdeu X pontos de vida.", currentMonster.Name))
		}
	}
}

func renderBattleLog(battleLog *[]string) {
	for i := 0; i < len(*battleLog); i++ {
		l.Printf("-> %s", (*battleLog)[i])
	}
}

func isMonsterAlive(monster Monster) bool {
	var response = false
	if monster.Health > 0 {
		response = true
	}
	return response
}

func renderMonsterInfo(monster Monster) {
	l.Printf(l.Colorize("Monster Information", l.Cyan))
	l.Printf("%s | Health: %s", monster.Name, l.Colorize(fmt.Sprintf("%.2f", monster.Health), l.Red))
	l.Printf(l.Colorize("Physical Status", l.Cyan))
	l.Printf("Attack: %s - Defense: %s",
		l.Colorize(fmt.Sprintf("%.2f", monster.Attack), l.Gray),
		l.Colorize(fmt.Sprintf("%.2f", monster.Defense), l.Gray))
	l.Printf(l.Colorize("Elemental Status", l.Cyan))
	fmt.Printf("Fire Resistance: %s%%\t", l.Colorize(fmt.Sprint(math.Round(monster.ElementalResist.Fire)), l.Yellow))
	fmt.Printf("Water Resistance: %s%%\n", l.Colorize(fmt.Sprintf("%.0f", monster.ElementalResist.Water), l.Blue))
	fmt.Printf("Wind Resistance: %s%%\t", l.Colorize(fmt.Sprintf("%.0f", monster.ElementalResist.Wind), l.LightGreen))
	fmt.Printf("Earth Resistance: %s%%\n", l.Colorize(fmt.Sprintf("%.0f", monster.ElementalResist.Earth), l.Green))
	fmt.Printf("Light Resistance: %s%%\t", l.Colorize(fmt.Sprintf("%.0f", monster.ElementalResist.Light), l.White))
	fmt.Printf("Dark Resistance: %s%%", l.Colorize(fmt.Sprintf("%.0f", monster.ElementalResist.Dark), l.Black))
}

func hasPendingProgress(dungeon Dungeon) bool {
	var reveal = false
	for i := 0; i < len(dungeon.Progress); i++ {
		if strings.EqualFold(dungeon.Progress[i], "???") {
			reveal = true
			break
		}
	}
	return reveal
}

func renderDungeonEnding(dungeon Dungeon, character Character) {
	dungeon.EndTime = time.Now()
	dungeonTotalDuration := time.Time{}.Add(dungeon.EndTime.Sub(dungeon.StartTime))
	l.Printf("Parabéns você concluiu a caverna %s com sucesso em %s! O que acha de tentar uma dificuldade maior?", dungeon.Name, dungeonTotalDuration.Format("15:04:05"))
	fmt.Println("Press the Enter Key to terminate the Dungeon!")
	var input string
	_, _ = fmt.Scanln(&input)

	renderGameActions(character)
}

func renderDungeonActions(dungeon Dungeon) {
	answers := struct {
		DungeonOption string
	}{}

	var questions []string
	questions = append(questions, "Iniciar próximo passo")
	questions = append(questions, "Sair")

	var dungeonQuestion = []*survey.Question{
		{
			Name: "dungeonOption",
			Prompt: &survey.Select{
				Message: "O que deseja fazer",
				Options: questions,
			},
		},
	}

	err := survey.Ask(dungeonQuestion, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if strings.EqualFold(answers.DungeonOption, "Iniciar próximo passo") {
		var reveal = false
		for i := 0; i < len(dungeon.Progress); i++ {
			if strings.EqualFold(dungeon.Progress[i], "???") {
				dungeon.Progress[i] = "Dragão de Komodo"
				reveal = true
				break
			}
		}
		if !reveal {
			l.Printf("Iniciando Boss...")
		}
	}
}

func renderDungeonProgress(dungeon Dungeon) {
	for i := 0; i < len(dungeon.Progress); i++ {
		if i == len(dungeon.Progress)-1 {
			if strings.EqualFold(dungeon.Progress[i], "???") {
				fmt.Printf("%s\n", dungeon.Progress[i])
			} else if strings.EqualFold(dungeon.Progress[i], "Boss") {
				fmt.Printf("%s\n", l.Colorize(dungeon.Progress[i], l.Red))
			} else {
				fmt.Printf("%s\n", l.Colorize(dungeon.Progress[i], l.Green))
			}
		} else if strings.EqualFold(dungeon.Progress[i], "???") {
			fmt.Printf("%s -> ", dungeon.Progress[i])
		} else if strings.EqualFold(dungeon.Progress[i], "Boss") {
			fmt.Printf("%s -> ", l.Colorize(dungeon.Progress[i], l.Red))
		} else {
			fmt.Printf("%s -> ", l.Colorize(dungeon.Progress[i], l.Green))
		}
	}
}

func renderDungeonIntroduction(dungeonId string) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:9001/lic-api/dungeons/intro/"+dungeonId, nil)
	charResp, _ := client.Do(req)

	story, err := ioutil.ReadAll(charResp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	l.Printf(string(story))
}

func renderLogin(token *Token) *string {
	fmt.Println(i18n.Localize("welcome_lic", region))
	fmt.Println(i18n.Localize("please_login", region))
	var body []byte = nil

	answers := struct {
		Name     string
		Password string
	}{}

	var loginQuestions = []*survey.Question{
		{
			Name:      "name",
			Prompt:    &survey.Input{Message: i18n.Localize("username", region) + ":"},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name: "password",
			Prompt: &survey.Password{
				Message: i18n.Localize("password", region) + ":",
			},
		},
	}

	err := survey.Ask(loginQuestions, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	var jsonStr = []byte(`{"username": "` + answers.Name + `", "password": "` + answers.Password + `"}`)

	resp, err := http.Post(
		"http://localhost:9001/lic-auth/token",
		"application/json", bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	if strings.Contains(resp.Status, "202") {

		//Convert the body to type string
		CallClear()
		err := json.Unmarshal(body, token)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return &answers.Name

}

func renderCharacterCreate() {
	var characterGender []string

	characterGender = append(characterGender, i18n.Localize("male", region))
	characterGender = append(characterGender, i18n.Localize("female", region))

	var gameClass []string

	gameClass = append(gameClass, i18n.Localize("warrior", region))
	gameClass = append(gameClass, i18n.Localize("mage", region))
	gameClass = append(gameClass, i18n.Localize("archer", region))
	gameClass = append(gameClass, i18n.Localize("warlock", region))

	var newCharacterQuestions = []*survey.Question{
		{
			Name:      "name",
			Prompt:    &survey.Input{Message: i18n.Localize("new_name", region) + ":"},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name: "gender",
			Prompt: &survey.Select{
				Message: i18n.Localize("new_gender", region) + ":",
				Options: characterGender,
			},
		},
		{
			Name: "class",
			Prompt: &survey.Select{
				Message: i18n.Localize("new_class", region) + ":",
				Options: gameClass,
			},
		},
	}
	newCharacterObject := struct {
		Name   string
		Gender string
		Class  string
	}{}
	errSurveyCharacter := survey.Ask(newCharacterQuestions, &newCharacterObject)
	if errSurveyCharacter != nil {
		fmt.Println(errSurveyCharacter.Error())
		return
	}
	l.Printf("Personagem criado com sucesso %s", newCharacterObject)
}

func renderCharacterSelect(session *string) struct{ CharacterId string } {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:9001/lic-api/characters", nil)
	req.Header.Set("Authorization", "Bearer "+*session)
	charResp, _ := client.Do(req)

	charBody, err := ioutil.ReadAll(charResp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var response Characters

	charErr := json.Unmarshal(charBody, &response)
	if charErr != nil {
		fmt.Println(err.Error())
	}

	var playerCharacters []string

	for i := 0; i < len(response.Characters); i++ {
		playerCharacters = append(playerCharacters, fmt.Sprintf("%s (%d)", response.Characters[i].Name,
			response.Characters[i].Level))
	}
	playerCharacters = append(playerCharacters, i18n.Localize("new_character", region)+" (+)")

	var loadCharacters = []*survey.Question{
		{
			Name: "characterId",
			Prompt: &survey.Select{
				Message: i18n.Localize("character_select", region) + ":",
				Options: playerCharacters,
			},
		},
	}

	characterAnswers := struct {
		CharacterId string
	}{}

	errSurveyCharacter := survey.Ask(loadCharacters, &characterAnswers)
	if errSurveyCharacter != nil {
		fmt.Println(errSurveyCharacter.Error())
		return struct{ CharacterId string }{CharacterId: "empty"}
	}

	for i := 0; i < len(response.Characters); i++ {
		var composeNameAndLevel = fmt.Sprintf("%s (%d)", response.Characters[i].Name,
			response.Characters[i].Level)
		if strings.EqualFold(composeNameAndLevel, characterAnswers.CharacterId) {
			characterAnswers.CharacterId = strconv.Itoa(response.Characters[i].Id)
			break
		}
	}

	return characterAnswers
}

func init() {
	rootCmd.AddCommand(startCmd)
}
