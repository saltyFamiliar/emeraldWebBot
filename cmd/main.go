package main

import (
	"fmt"
	"github.com/saltyFamiliar/emeraldWebBot/tools"
	"github.com/saltyFamiliar/tgramAPIBotLib/api"
	"github.com/saltyFamiliar/tgramAPIBotLib/bot"
	"strconv"
	"strings"
	"time"
)

func main() {
	tGramBot := bot.NewTgramBot(api.GetAPIKey("token.txt"))
	jobCh := make(chan *api.Message, 10)

	go func() {
		for {
			for _, update := range tGramBot.GetUpdates() {
				tGramBot.Offset = int(update.UpdateId) + 1
				if update.Message == nil {
					continue
				}
				jobCh <- update.Message
				ackMsg := fmt.Sprintf("Received request: %s", update.Message.Text)
				go tGramBot.SendMsg(ackMsg, update.Message.Chat.Id)
			}
			time.Sleep(4 * time.Second)
		}
	}()

	for job := range jobCh {
		go func(reqMsg *api.Message) {
			msgText := reqMsg.Text
			parts := strings.Split(msgText, " ")
			if len(parts) == 3 {
				startPort, err := strconv.Atoi(parts[1])
				if err != nil {
					tGramBot.SendMsg("Error parsing first port", reqMsg.Chat.Id)
					return
				}
				endPort, err := strconv.Atoi(parts[2])
				if err != nil {
					tGramBot.SendMsg("Error parsing last port", reqMsg.Chat.Id)
					return
				}

				open, closed := tools.ScanPorts(parts[0], startPort, endPort, 4)
				msg := fmt.Sprintf("Open: %v Closed: %v ", open, closed)
				go tGramBot.SendMsg(msg, reqMsg.Chat.Id)
			}
		}(job)
	}
}
