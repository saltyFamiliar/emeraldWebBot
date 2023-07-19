package main

import (
	"fmt"
	"github.com/saltyFamiliar/emeraldWebBot/internal/commands"
	"github.com/saltyFamiliar/tgramAPIBotLib/api"
	"github.com/saltyFamiliar/tgramAPIBotLib/pkg/bot"
	"log"
)

func main() {
	apiKey, err := api.GetAPIKey("token.txt")
	if err != nil {
		log.Fatalln(err)
	}
	tGramBot := bot.NewTgramBot(apiKey)

	scanRoutine := bot.NewRoutine(bot.Action{
		Raw: commands.ScanPorts,
		Wrapper: func(args ...interface{}) (string, error) {
			open, closed := commands.ScanPorts(args[0].(string), args[1].(int), args[2].(int), args[3].(int))
			return fmt.Sprintf("Open: %v, Closed: %v", open, closed), nil
		},
	})

	if err := tGramBot.RegisterRoutine("scan", scanRoutine); err != nil {
		fmt.Println(err)
	}

	tGramBot.Run()
}
