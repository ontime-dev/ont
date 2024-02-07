package run

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseTimeDate() {

}

func ParseEvryFrom(every, from string) error {

	number, err := strconv.Atoi(every[:len(every)-1])
	if err != nil {
		return err
	}

	last_char := every[len(every)-1:]

	next_run, err := setNextRun(last_char, from, number)
	if err != nil {
		return err
	}
	fmt.Printf("Next run:\n%s\n", next_run)

	return nil

}

func setNextRun(last_char, from string, number int) (string, error) {

	crntTime := time.Now()

	switch from {
	case "now", "today":
		crntTime = time.Now()
	case "tomorrow":
		crntTime = crntTime.AddDate(0, 0, 1)
	default:
		var err error
		if strings.Contains(from, "T") {
			dateTime := strings.Split(from, "T")
			from = dateTime[0] + " " + dateTime[1]
		} else if strings.Contains(from, "-") {
			from = from + " " + crntTime.Format("15:04:05")
		} else if strings.Contains(from, ":") {
			from = crntTime.Format("2006-01-02") + " " + from
		}

		crntTime, err = time.Parse("2006-01-02 15:04:05", from)
		if err != nil {
			return "", err
		}
	}

	switch last_char {
	case "h":
		fmt.Println("every hour")
		next_run := crntTime.Add(time.Hour * time.Duration(number)).Format("15:04:05 Feb 07 2006")

		return next_run, nil
	case "m":
		fmt.Println("every minute")
		next_run := crntTime.Add(time.Minute * time.Duration(number)).Format("15:04:01 07 February 2024")
		return next_run, nil
	case "s":
		fmt.Println("every seconds")
		next_run := crntTime.Add(time.Second * time.Duration(number)).Format("15:04:01 07 February 2024")
		return next_run, nil
	case "d":
		fmt.Println("every day")
		next_run := crntTime.AddDate(0, 0, number).Format("15:04:05 Jan 02 2006")
		return next_run, nil
	case "w":
		fmt.Println("every week")
		next_run := crntTime.AddDate(0, 0, number*7).Format("15:04:05 Jan 02 2006")
		return next_run, nil
	case "M":
		next_run := crntTime.AddDate(0, number, 0).Format("15:04:05 Jan 02 2006")
		fmt.Println("every month")
		return next_run, nil
	case "y":
		next_run := crntTime.AddDate(number, 0, 0).Format("15:04:05 Jan 02 2006")
		fmt.Println("every year")
		return next_run, nil
	default:
		return "", errors.New("Please specify a valid option with --every flag.")
	}

}
