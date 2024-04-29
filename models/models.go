package models

import (
	"time"
)

type repeatRule struct {
	days   int
	yearly bool
}

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (t *Task) VerifyChange() error {
	if t.Date == "" {
		t.Date = time.Now().Format("20060102")
	} else {
		dateParsed, err := time.Parse("20060102", t.Date)
		if err != nil {
			return err
		}

		today := time.Now().Format("20060102")

		if dateParsed.Before(time.Now()) && t.Date != today {
			if t.Repeat == "" {
				t.Date = today
			} else {
				if t.Date == today {
					t.Date = today
				} else {
					nextDate, err := NextDate(time.Now(), t.Date, t.Repeat)
					if err != nil {
						return err
					}
					if nextDate >= today {
						t.Date = nextDate
					} else {
						t.Date = today
					}
				}
			}
		}
	}
	return nil
}
