package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	rule, err := parseRepeatRule(repeat)
	if err != nil {
		return "", err
	}

	nextDate := startDate

	for nextDate.Before(now) || nextDate.Equal(startDate) {
		nextDate = applyRepeatRule(nextDate, rule)
	}

	return nextDate.Format("20060102"), nil
}

func parseRepeatRule(repeat string) (repeatRule, error) {
	var rule repeatRule

	if len(repeat) == 1 && repeat[0] == 'y' {
		rule.yearly = true
		return rule, nil
	}

	if len(repeat) < 2 || repeat[0] != 'd' {
		return rule, errors.New("неправильный формат повторения")
	}

	repeatArr := strings.Split(repeat, " ")
	days, err := strconv.Atoi(repeatArr[1])
	if err != nil {
		return rule, errors.New("неправильный формат числа дней")
	}

	if days <= 0 || days > 400 {
		return rule, errors.New("недопустимое количество дней")
	}

	rule.days = days
	return rule, nil
}

func applyRepeatRule(date time.Time, rule repeatRule) time.Time {
	if rule.yearly {
		return date.AddDate(1, 0, 0)
	}
	return date.AddDate(0, 0, rule.days)
}
