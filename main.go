package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/y0ssar1an/slack-pushups/internal/slack"
)

const (
	minPushUps  = 10
	maxPushUps  = 30
	openingHour = 10
	openingMin  = 0
	openingSec  = 0
	closingHour = 18
)

func main() {
	// fail early if zoneinfo db not present on server
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Fatal(err)
	}

	ch, err := slack.NewChannel("api-test")
	if err != nil {
		log.Fatal(err)
	}

	// start async routine that chooses users randomly for push-ups
	nextMemberID := make(chan string)
	go randomMember(ch, nextMemberID)

	SgtMittens := slack.Bot{"SgtMitt"}

	for {
		now := time.Now().In(loc)
		if closed, timeToOpen := isAfterHours(now, loc); closed {
			time.Sleep(timeToOpen)
			SgtMittens.PostMessage("RISE AND SHINE, CUPCAKES!", ch)
		}

		var user slack.User
		for user.Name == "" {
			user, err = slack.NewUser(<-nextMemberID)
			if err != nil {
				log.Println(err)
				time.Sleep(1 * time.Minute)
			}
		}

		pushUps := randInt(minPushUps, maxPushUps+1) // +1 because upper bound is non-inclusive
		msg := fmt.Sprintf(
			"@%s %d PUSH-UPS RIGHT MEOW!",
			user.Name,
			pushUps,
		)

		if !closingSoon(now, loc) {
			msg += "\nNext lottery for push-ups in 20 minutes"
		}

		err = SgtMittens.PostMessage(msg, ch)
		if err != nil {
			log.Println(err)
			time.Sleep(1 * time.Minute)
			continue
		}
		time.Sleep(20 * time.Minute)
	}
}

func randomMember(ch slack.Channel, nextMemberID chan string) {
	var err error
	for {
		err = ch.UpdateMembers()
		if err != nil {
			log.Println(err)
		}

		i := randInt(0, len(ch.Members))
		nextMemberID <- ch.Members[i]
	}
}

func randInt(min, max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	if err != nil {
		log.Println(err)
	}
	// return min if rand.Int() call fails
	return int(n.Int64()) + min
}

func isAfterHours(now time.Time, loc *time.Location) (bool, time.Duration) {
	var timeToOpen time.Duration

	if isWeekend(now) {
		days := daysToMonday(now.Weekday())
		mondayOpeningTime := time.Date(now.Year(), now.Month(), now.Day()+days, openingHour, 0, 0, 0, loc)
		timeToOpen = mondayOpeningTime.Sub(now)
		return true, timeToOpen
	}

	if now.Hour() < openingHour {
		openToday := time.Date(now.Year(), now.Month(), now.Day(), openingHour, 0, 0, 0, loc)
		timeToOpen = openToday.Sub(now)
		return true, timeToOpen
	}

	if now.Hour() >= closingHour {
		openTomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, openingHour, 0, 0, 0, loc)
		timeToOpen = openTomorrow.Sub(now)
		return true, timeToOpen
	}

	return false, timeToOpen // 0 time
}

func isWeekend(t time.Time) bool {
	day := t.Weekday()
	return (day == time.Friday && t.Hour() >= closingHour) || day == time.Saturday || day == time.Sunday
}

func daysToMonday(day time.Weekday) int {
	const weekdays = 7
	return (weekdays - int(day-time.Monday)) % weekdays
}

// Return true if closing time is in less than 20 minutes.
func closingSoon(now time.Time, loc *time.Location) bool {
	closingTime := time.Date(now.Year(), now.Month(), now.Day(), closingHour, 0, 0, 0, loc)
	return closingTime.Sub(now) <= (20 * time.Minute)
}
