package main

import (
	"context"
	"fmt"
	"github.com/saltyFamiliar/emeraldWebBot/internal/commands"
	"github.com/saltyFamiliar/tgramAPIBotLib/api"
	"github.com/saltyFamiliar/tgramAPIBotLib/bot"
	"strconv"
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
			//TODO: Ask about how parsing function should be implemented
			msgText := reqMsg.Text
			parts := strings.Split(msgText, " ")
			if len(parts) == 3 {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				startPort, err := strconv.Atoi(parts[1])
				if err != nil {
					if msgErr := tGramBot.SendMsg(ctx, "Error parsing first port", reqMsg.Chat.Id); msgErr != nil {
						fmt.Println(msgErr)
					}
					return
				}
				endPort, err := strconv.Atoi(parts[2])
				if err != nil {
					if msgErr := tGramBot.SendMsg(ctx, "Error parsing last port", reqMsg.Chat.Id); msgErr != nil {
						fmt.Println(msgErr)
					}
					return
				}

				open, closed := commands.ScanPorts(parts[0], startPort, endPort, 4)
				msg := fmt.Sprintf("Open: %v Closed: %v ", open, closed)

				go func() {
					if err := tGramBot.SendMsg(ctx, msg, reqMsg.Chat.Id); err != nil {
						fmt.Println(err)
					}
				}()
			}
		}(job)
	}
}
