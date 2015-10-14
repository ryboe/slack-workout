package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/y0ssar1an/slack-pushups/internal/slack"
)

const (
	minPushUps = 10
	maxPushUps = 30
	// TODO: may not need this.
	slackbotURL = "https://%s.slack.com/services/hooks/slackbot?%s"
)

// TODO: may not need this. sending messages as user, not slackbot
// var botToken = os.Getenv("SLACK_BOT_TOKEN")

func main() {
	// TODO; delete this
	// if botToken == "" {
	// 	log.Fatal("SLACK_BOT_TOKEN not set")
	// }

	ch, err := slack.NewChannel("monkeytacos", "api-test")
	if err != nil {
		log.Fatal(err)
	}

	// DEBUG
	fmt.Println(ch)

	err = ch.UpdateMembers()
	fmt.Println("ERRR:", err)

	// nextMember := make(chan string)
	// go randomMember(ch, nextMember)

	// // DEBUG
	// // for i := 0; i < 1; i++ {
	// // 	fmt.Println(<-nextMember)
	// // }

	// var pushUps int
	// for {
	// 	t := time.UTC()
	// 	if closed, timeToOpen := isAfterHours(t); closed {
	// 		time.Sleep(timeToOpen)

	// 		// TODO: rise and shine message here
	// 	}

	// 	// TODO: get user name from users.info
	// 	var user slack.User
	// 	for user.Name == "" {
	// 		user, err = slack.NewUser(<-nextMember)
	// 		if err != nil {
	// 			log.Println(err)
	// 			time.Sleep(1 * time.Minute)
	// 		}
	// 	}

	// 	// TODO: write slack.Bot

	// 	pushUps = randPushUps(minPushUps, maxPushUps)
	// 	msg := fmt.Sprintf(
	// 		"%d PUSH-UPS RIGHT MEOW! @%s\nNext lottery for push-ups in 20 minutes",
	// 		pushUps,
	// 		user.Name,
	// 	)
	// 	err = bot.Msg(msg)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	time.Sleep(20 * time.Minute)
	// }
}

func randomMember(ch slack.Channel, nextMember chan string) {
	var err error
	for {
		err = ch.UpdateMembers()
		if err != nil {
			log.Println(err)
		}

		i := rand.Intn(len(ch.Members))
		nextMember <- ch.Members[i]
	}
}

func randPushUps(min, max int) int {
	return rand.Intn(max-min+1) + min
}
