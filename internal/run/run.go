package run

import (
	"errors"
	"fmt"
	"strconv"
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

	next_run, err := setNextRun(last_char, number)
	fmt.Printf("Next run:\n%s\n", next_run)

	switch from {
	case "now", "today":
	//	fmt.Printf("Starting from %s \n", currentTime)
	default:

	}

	return nil

}

func setNextRun(last_char string, number int) (string, error) {
	crntTime := time.Now()
	fmt.Println(crntTime)
	/*
		hour := crntTime.Hour()
		min := crntTime.Minute()
		sec := crntTime.Second()
		day := crntTime.Day()
		//weekDay := crntTime.Weekday()
		month := crntTime.Month()
		year := crntTime.Year()
	*/
	switch last_char {
	case "h":
		fmt.Println("every hour")
		next_run := crntTime.Add(time.Hour * time.Duration(number)).Format("15:04:05 Feb 07 2006")
		//next_run := strconv.Itoa(hour) + ":" + strconv.Itoa(min) + ":" + strconv.Itoa(sec) + " " + strconv.Itoa(day) + " " + month.String() + " " + strconv.Itoa(year)

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
