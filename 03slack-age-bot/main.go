package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/shomali11/slacker"
)

func printCommandEvents(analyticsChannl <-chan *slacker.CommandEvent) {
	for event := range analyticsChannl {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println(event.Parameters)
	}
}

func main() {
	os.Setenv("SLACK_BOT_TOKEN", "xoxb-4678264186740-5014276170422-EpiROIlGerxzw1SHJ9M6bFgk")
	os.Setenv("SLACK_APP_TOKEN", "xapp-1-A050LPT9GCS-5017867641333-0ace7053bef1d94fe2e683709dacb2a428087b7a55bec50ea9f817e7f40a42ea")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	go printCommandEvents(bot.CommandEvents())

	bot.Command("My yob is <year>", &slacker.CommandDefinition{
		Description: "yob caculator",
		Examples:    []string{"My yob is 2002"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			year := request.Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil {
				fmt.Println("error")
			}

			age := (2023 - yob)
			r := fmt.Sprint(age)
			response.Reply(r)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
