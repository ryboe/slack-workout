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
	closingHour = 18
	weekdays    = 7
)

func main() {
	ch, err := slack.NewChannel("api-test")
	if err != nil {
		log.Fatal(err)
	}

	// start async routine that chooses users randomly for push-ups
	nextMemberID := make(chan string)
	go RandomMember(ch, nextMemberID)

	sgtMittens := slack.Bot{"SgtMittens"}

	// DEBUG
	sgtMittens.PostMessage("this is a test", ch)

	loc := time.FixedZone("PST", -8*60*60)
	for {
		now := time.Now().In(loc)
		if closed, timeToOpen := IsAfterHours(now); closed {
			time.Sleep(timeToOpen)
			sgtMittens.PostMessage("RISE AND SHINE, CUPCAKES!", ch)
		}

		var user slack.User
		for user.Name == "" {
			user, err = slack.NewUser(<-nextMemberID)
			if err != nil {
				log.Println(err)
				time.Sleep(1 * time.Minute)
			}
		}

		pushUps := RandInt(minPushUps, maxPushUps+1) // +1 because upper bound is non-inclusive
		msg := fmt.Sprintf(
			"@%s %d PUSH-UPS RIGHT MEOW!",
			user.Name,
			pushUps,
		)

		if !ClosingSoon(now) {
			msg += "\nNext lottery for push-ups in 20 minutes"
		}

		err = sgtMittens.PostMessage(msg, ch)
		if err != nil {
			log.Println(err)
			time.Sleep(1 * time.Minute)
			continue
		}
		time.Sleep(20 * time.Minute)
	}
}

// RandomMember updates the list of members in the given Slack channel and
// returns a random member ID.
func RandomMember(ch slack.Channel, nextMemberID chan string) {
	var err error
	for {
		err = ch.UpdateMembers()
		if err != nil {
			log.Println(err)
		}

		i := RandInt(0, len(ch.Members))
		nextMemberID <- ch.Members[i]
	}
}

// RandInt uses /dev/urandom to select a random integer in the range [min, max).
func RandInt(min, max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	if err != nil {
		log.Println(err)
	}

	return int(n.Int64()) + min // return min if rand.Int() call fails
}

// IsAfterHours returns true if the given time is after work hours at Omaze.
// If Omaze is closed, IsAfterHours will return the duration until opening time
// on Monday.
func IsAfterHours(now time.Time) (closed bool, timeToOpen time.Duration) {
	loc := now.Location()
	if IsWeekend(now) {
		days := DaysToMonday(now.Weekday())
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

// IsWeekend returns true if the given time is Friday after closing time,
// Saturday, or Sunday.
func IsWeekend(t time.Time) bool {
	day := t.Weekday()
	return (day == time.Friday && t.Hour() >= closingHour) || day == time.Saturday || day == time.Sunday
}

// DaysToMonday returns the number of days from the given weekday to Monday.
func DaysToMonday(day time.Weekday) int {
	return (weekdays - int(day-time.Monday)) % weekdays
}

// ClosingSoon returns true if the given time is within 20 minutes of closing
// time.
func ClosingSoon(now time.Time) bool {
	loc := now.Location()
	closingTime := time.Date(now.Year(), now.Month(), now.Day(), closingHour, 0, 0, 0, loc)
	return closingTime.Sub(now) <= (20 * time.Minute)
}
