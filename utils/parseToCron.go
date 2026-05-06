package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var weekdays = map[string]int{
	"вс": 0, "воскресенье": 0,
	"пн": 1, "понедельник": 1,
	"вт": 2, "вторник": 2,
	"ср": 3, "среда": 3,
	"чт": 4, "четверг": 4,
	"пт": 5, "пятница": 5,
	"сб": 6, "суббота": 6,
}

func parseToCronInternal(input string) ([]string, error) {
	input = normalize(input)

	// -------------------------
	// каждые N ...
	// -------------------------
	if strings.HasPrefix(input, "каждые") {
		return parseInterval(input)
	}

	// -------------------------
	// будни / выходные
	// -------------------------
	if strings.Contains(input, "будни") {
		return buildCronForDays(input, "1-5")
	}

	if strings.Contains(input, "выход") {
		return buildCronForDays(input, "0,6")
	}

	// -------------------------
	// каждый день
	// -------------------------
	if strings.Contains(input, "каждый день") {
		return buildCronForDays(input, "*")
	}

	// -------------------------
	// дни недели (множественные)
	// -------------------------
	var days []string
	for name, num := range weekdays {
		if strings.Contains(input, name) {
			days = append(days, strconv.Itoa(num))
		}
	}

	if len(days) > 0 {
		return buildCron(input, "*", "*", strings.Join(days, ","))
	}

	// -------------------------
	// каждый месяц X числа
	// -------------------------
	if strings.Contains(input, "числа") {
		return parseMonthDays(input)
	}

	return nil, fmt.Errorf("непонятный формат", input)
}

// -------------------------
// ВСПОМОГАТЕЛЬНЫЕ
// -------------------------

func normalize(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " в ", " ")
	s = strings.ReplaceAll(s, ",", " ")
	return s
}

// -------------------------

func extractTimes(input string) ([]time.Time, error) {
	parts := strings.Fields(input)

	var times []time.Time

	for _, p := range parts {
		if strings.Contains(p, ":") {
			t, err := time.Parse("15:04", p)
			if err == nil {
				times = append(times, t)
			}
		}
	}

	if len(times) == 0 {
		return nil, fmt.Errorf("не найдено время", input)
	}

	return times, nil
}

// -------------------------

func buildCronForDays(input, dayExpr string) ([]string, error) {
	return buildCron(input, "*", "*", dayExpr)
}

// -------------------------

func buildCron(input, dom, mon, dow string) ([]string, error) {
	times, err := extractTimes(input)
	if err != nil {
		return nil, err
	}

	var res []string

	for _, t := range times {
		res = append(res,
			fmt.Sprintf("%d %d %s %s %s",
				t.Minute(), t.Hour(), dom, mon, dow))
	}

	return res, nil
}

// -------------------------

func parseInterval(input string) ([]string, error) {
	parts := strings.Fields(input)

	if len(parts) < 3 {
		return nil, fmt.Errorf("непонятный формат", input)
	}

	n, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	unit := parts[2]

	switch unit {
	case "минута", "минуты", "минут":
		return []string{fmt.Sprintf("@every %dm", n)}, nil

	case "час", "часа", "часов":
		return []string{fmt.Sprintf("@every %dh", n)}, nil

	case "день", "дня", "дней":
		return []string{fmt.Sprintf("@every %dh", n*24)}, nil
	}

	return nil, fmt.Errorf("непонятный интервал", input)
}

// -------------------------

func parseMonthDays(input string) ([]string, error) {
	parts := strings.Fields(input)

	var days []string

	for i, p := range parts {
		if p == "числа" && i > 0 {
			days = append(days, parts[i-1])
		}
		if p == "и" && i > 0 && i < len(parts)-1 {
			days = append(days, parts[i-1])
			days = append(days, parts[i+1])
		}
	}

	if len(days) == 0 {
		return nil, fmt.Errorf("не найдены дни месяца", input)
	}

	return buildCron(input, strings.Join(days, ","), "*", "*")
}

func ParseToCron(input string) (string, error) {
	exprs, err := parseToCronInternal(input)
	if err != nil {
		return "", err
	}

	return strings.Join(exprs, ";"), nil
}
