package main

import (
	"context"
	"fmt"
	"github.com/saltyFamiliar/emeraldWebBot/internal/commands"
	"github.com/saltyFamiliar/tgramAPIBotLib/api"
	"github.com/saltyFamiliar/tgramAPIBotLib/bot"
	"strings"
	"time"
)

func main() {
	tGramBot := bot.NewTgramBot(api.GetAPIKey("token.txt"))
	jobCh := make(chan *api.Message, 10)
	updatesCh := make(chan []api.Update, 10)

	// producer
	// get updates, send them through updates channel
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			updates, err := tGramBot.GetUpdates(ctx)
			if err != nil {
				fmt.Printf("Unable to get updates: %v", err)
			}
			cancel()
			updatesCh <- updates
			time.Sleep(4 * time.Second)
		}
	}()

	// consumer, producer
	// send update msgs through job channel, update bot offset, send ack chat msg
	go func() {
		for updates := range updatesCh {
			for _, update := range updates {
				tGramBot.Offset = int(update.UpdateId) + 1
				if update.Message == nil {
					continue
				}
				jobCh <- update.Message
				go func(msg *api.Message) {
					ackMsg := fmt.Sprintf("Received request: %s", msg.Text)
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
					if err := tGramBot.SendMsg(ctx, ackMsg, msg.Chat.Id); err != nil {
						fmt.Printf("Unable to send ack message: %v", err)
					}
					cancel()
				}(update.Message)
			}
		}
	}()

	// consumer
	// parses requests and sends response message
	for job := range jobCh {
		go func(reqMsg *api.Message) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			msgText := reqMsg.Text
			parts := strings.Split(msgText, " ")
			cmd, err := commands.NewCommand(parts[0], parts[1:])
			if err != nil {
				if msgErr := tGramBot.SendMsg(ctx, err.Error(), reqMsg.Chat.Id); msgErr != nil {
					fmt.Println(msgErr)
				}
				return
			}

			msg, err := cmd.Execute()
			if err != nil {
				if msgErr := tGramBot.SendMsg(ctx, err.Error(), reqMsg.Chat.Id); msgErr != nil {
					fmt.Println(msgErr)
				}
				return
			}

			if err := tGramBot.SendMsg(ctx, msg, reqMsg.Chat.Id); err != nil {
				fmt.Println(err)
			}
		}(job)
	}
}
