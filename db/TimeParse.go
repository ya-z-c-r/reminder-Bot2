package db

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseHumanTime(input string) (time.Time, error) {
	now := time.Now()

	input = strings.ToLower(strings.TrimSpace(input))
	input = strings.ReplaceAll(input, " в ", " ")
	parts := strings.Fields(input)

	// --- 0. Полная дата (ISO)
	if t, err := time.Parse("2006-01-02 15:04", input); err == nil {
		return t, nil
	}

	// --- 0.1 Европейский формат
	if t, err := time.Parse("02.01.2006 15:04", input); err == nil {
		return t, nil
	}

	// --- 1. "15:04"
	if t, err := time.Parse("15:04", input); err == nil {
		return time.Date(
			now.Year(), now.Month(), now.Day(),
			t.Hour(), t.Minute(), 0, 0, time.Local,
		), nil
	}

	// --- 2. сегодня
	if len(parts) >= 2 && parts[0] == "сегодня" {
		t, err := time.Parse("15:04", parts[len(parts)-1])
		if err == nil {
			return time.Date(
				now.Year(), now.Month(), now.Day(),
				t.Hour(), t.Minute(), 0, 0, time.Local,
			), nil
		}
	}

	// --- 3. завтра
	if len(parts) >= 2 && parts[0] == "завтра" {
		t, err := time.Parse("15:04", parts[len(parts)-1])
		if err == nil {
			day := now.AddDate(0, 0, 1)
			return time.Date(
				day.Year(), day.Month(), day.Day(),
				t.Hour(), t.Minute(), 0, 0, time.Local,
			), nil
		}
	}

	// --- 4. послезавтра
	if len(parts) >= 2 && parts[0] == "послезавтра" {
		t, err := time.Parse("15:04", parts[len(parts)-1])
		if err == nil {
			day := now.AddDate(0, 0, 2)
			return time.Date(
				day.Year(), day.Month(), day.Day(),
				t.Hour(), t.Minute(), 0, 0, time.Local,
			), nil
		}
	}

	// --- 5. через N
	if len(parts) >= 3 && parts[0] == "через" {
		n, err := strconv.Atoi(parts[1])
		if err == nil {
			switch parts[2] {
			case "минута", "минуты", "минут":
				return now.Add(time.Duration(n) * time.Minute), nil
			case "час", "часа", "часов":
				return now.Add(time.Duration(n) * time.Hour), nil
			case "день", "дня", "дней":
				return now.AddDate(0, 0, n), nil
			}
		}
	}

	return time.Time{}, fmt.Errorf("непонятный формат")
}
