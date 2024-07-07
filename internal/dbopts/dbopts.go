package dbopts

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

type Jobs struct {
	Id        int
	Script    string
	Exec_time string
	Every     string
	Status    string
}

/*func Opt(operation, user string, job Jobs, cfgFile string) ([]any, error) {
	password := config.LoadConfig(cfgFile)
	fmt.Println(password)
	//pass_cmd := fmt.Sprintf("ont:%s@/ontime", password)
	// db, err := sql.Open("mysql", "ont:password@/ontime")

	// if err != nil {
	// 	return err, nil
	// }
	// defer db.Close()
	//fmt.Printf("Operation: %s -- User: %s\n", operation, user)

	switch operation {
	case "create":
		fmt.Println("Operation: create")
		err = Create(db, user)
		if err != nil {
			return err, nil
		}
	case "insert":
		//arg[1] = user, arg[2] = script path, arg[3] = next_run
		//fmt.Println("Hello there")
		err := Insert(db, user, job, true)
		if err != nil {
			return err, nil
		}
	case "list":
		//arg[1] = user
		if job.Id == 0 {
			err, jobs := PrintJobs(db, user)
			if err != nil {
				return err, nil
			}
			return nil, jobs
		} else if job.Id > 0 {
			err, jobs := PrintOneJob(db, user, job.Id)
			if err != nil {
				return err, nil
			}
			return nil, jobs
		} else {
			return errors.New("please enter a valid jobID"), nil
		}
	case "stop":
		err := ChangeJobStatus(db, user, "Inactive", job)
		if err != nil {
			return err, nil
		}
	case "start":
		err := ChangeJobStatus(db, user, "Active", job)
		if err != nil {
			return err, nil
		}
	case "remove":
		if err := RemoveJob(db, user, job); err != nil {
			return err, nil
		}

	}

	return nil, nil
}*/

func List(db *sql.DB, job []Jobs, user string) (error, []Jobs) {
	if job.Id == 0 {
		err, jobs := PrintJobs(db, user)
		if err != nil {
			return err, nil
		}
		return nil, jobs
	} else if job.Id > 0 {
		err, jobs := PrintOneJob(db, user, job.Id)
		if err != nil {
			return err, nil
		}
		return nil, jobs
	} else {
		return errors.New("please enter a valid jobID"), nil
	}
}
