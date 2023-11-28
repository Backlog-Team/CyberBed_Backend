package converter

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func StringToTime(timeStr string) (time.Time, error) {
	parts := strings.Split(timeStr, "_")
	if len(parts) != 2 {
		return time.Time{}, errors.New("invalid input format")
	}

	duration, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, errors.New("invalid duration")
	}

	unit := parts[1]
	var durationTime time.Duration

	switch unit {
	case "s":
		durationTime = time.Second
	case "m":
		durationTime = time.Minute
	case "h":
		durationTime = time.Hour
	case "d":
		durationTime = time.Hour * 24
	default:
		return time.Time{}, errors.New("invalid unit")
	}

	futureTime := time.Now().Add(time.Duration(duration) * durationTime)
	return futureTime, nil
}
