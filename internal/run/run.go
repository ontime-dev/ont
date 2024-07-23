package run

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// func ParseTimeDate() {

// }

func ParseEvryFrom(every, from string) (string, error) {

	number, last_char := GetLastChar(every)

	next_run, err := setNextRun(last_char, from, number)
	if err != nil {
		return "", err
	}

	return next_run, err

}

func ParseFrom(from string) (time.Time, error) {
	crntTime := time.Now()

	switch from {
	case "now", "today":
		crntTime = time.Now()
		return crntTime, nil
	case "tomorrow":
		crntTime = crntTime.AddDate(0, 0, 1)
		return crntTime, nil
	default:
		var err error
		if strings.Contains(from, "T") {
			dateTime := strings.Split(from, "T")
			from = dateTime[0] + " " + dateTime[1]
		} else if strings.Contains(from, "-") {
			from = from + " " + crntTime.Format("15:04:05")
		} else if strings.Contains(from, ":") {
			from = crntTime.Format("02-01-2006") + " " + from
		} else if strings.Contains(from, "+") {
			from_string := strings.Split(from, "+")

			num, last_char := GetLastChar(from_string[1])
			from, err = ParseEvery(crntTime, last_char, num)
			if err != nil {
				return crntTime, err
			}
			crntTime, err = time.Parse("15:04:05 Jan 02 2006", from)
			if err != nil {
				return crntTime, err
			}
			return crntTime, nil
		}
		crntTime, err = time.Parse("02-01-2006 15:04:05", from)
		if err != nil {
			return crntTime, err
		}

		return crntTime, nil
	}
}

func ParseEvery(crntTime time.Time, last_char string, number int) (string, error) {

	switch last_char {
	case "h":
		next_run := crntTime.Add(time.Hour * time.Duration(number)).Format("15:04:05 Jan 02 2006")
		return next_run, nil
	case "m":
		next_run := crntTime.Add(time.Minute * time.Duration(number)).Format("15:04:05 Jan 02 2006")
		return next_run, nil
	case "s":
		next_run := crntTime.Add(time.Second * time.Duration(number)).Format("15:04:05 Jan 02 2006")
		return next_run, nil
	case "d":
		next_run := crntTime.AddDate(0, 0, number).Format("15:04:05 Jan 02 2006")
		return next_run, nil
	case "w":
		next_run := crntTime.AddDate(0, 0, number*7).Format("15:04:05 Jan 02 2006")
		return next_run, nil
	case "M":
		next_run := crntTime.AddDate(0, number, 0).Format("15:04:05 Jan 02 2006")
		return next_run, nil
	case "y":
		next_run := crntTime.AddDate(number, 0, 0).Format("15:04:05 Jan 02 2006")
		fmt.Println("every year")
		return next_run, nil
	default:
		return "", errors.New("please specify a valid option with --every flag")
	}
}

func lastCharValidity(last_char string) error {
	valid_chars := []string{"h", "m", "s", "d", "W", "M", "y"}
	for _, char := range valid_chars {
		if char == last_char {
			return nil
		}
	}
	return errors.New("please Specify a valid option with --every flag")
}

func setNextRun(last_char, from string, number int) (string, error) {
	crntTime, err := ParseFrom(from)
	if err != nil {
		return "", err
	}

	//function to insert the next_run in the database
	next_run := crntTime.Format("15:04:05 Jan 02 2006")

	err = lastCharValidity(last_char)
	if err != nil {
		return "", err
	}

	return next_run, nil

}

func GetLastChar(every string) (number int, last_char string) {
	number, err := strconv.Atoi(every[:len(every)-1])
	if err != nil {
		return 0, ""
	}
	last_char = every[len(every)-1:]

	return number, last_char
}
