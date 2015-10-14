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
)

func main() {
	// ch, err := slack.NewChannel("monkeytacos", "api-test")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // DEBUG
	// fmt.Println(ch)

	mittens, err := slack.NewUser("monkeytacos", "sgtmittens")
	if err != nil {
		log.Fatal(err)
	}

	// DEBUG
	fmt.Println(mittens)

	// err = ch.UpdateMembers()

	// DEBUG
	// fmt.Println("ERR AFTER UPDATEMEMBERS()?:", err)

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

func randomMember(ch slack.Channel, mittensId string, nextMember chan string) {
	var err error
	for {
		err = ch.UpdateMembers()
		if err != nil {
			log.Println(err)
		}

		i := rand.Intn(len(ch.Members))

		// prevent mittens from picking self
		for ch.Members[i] == mittensId {
			i = rand.Intn(len(ch.Members))
		}

		nextMember <- ch.Members[i]
	}
}

func randPushUps(min, max int) int {
	return rand.Intn(max-min+1) + min
}
