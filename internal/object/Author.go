package object

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Author struct {
	Timestamp time.Time
	Name      string
	Email     string
}

func (a *Author) Type() string {
	return "author"
}

func (a *Author) String() string {
	_, offset := a.Timestamp.Zone()

	abs := offset
	if abs < 0 {
		abs = -abs
	}

	hours := abs / 3600
	minutes := (abs % 3600) / 60
	sign := "+"
	if offset < 0 {
		sign = "-"
	}
	tz := fmt.Sprintf("%s%02d%02d", sign, hours, minutes)

	return fmt.Sprintf("%s <%s> %d %s", a.Name, a.Email, a.Timestamp.Unix(), tz)
}

// Parses a string in Author.String() result format
// Expects the value part only, without keyword prefix (e.g. "John Doe <john@example.com> 1714000000 +0200")
func ParseAuthor(str string) (*Author, error) {
	var result Author

	ltIdx := strings.IndexByte(str, '<')
	gtIdx := strings.IndexByte(str, '>')
	if ltIdx == -1 || gtIdx == -1 || gtIdx < ltIdx {
		return nil, errors.New("Error parsing author string - missing email brackets")
	}

	result.Name = strings.TrimSpace(str[:ltIdx])

	result.Email = str[ltIdx+1 : gtIdx]

	rest := strings.TrimSpace(str[gtIdx+1:])
	parts := strings.Split(rest, " ")
	if len(parts) != 2 {
		return nil, errors.New("Error parsing author string - timestamp/timezone")
	}

	timestamp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, errors.New("Error parsing author string - invalid timestamp")
	}

	tz := parts[1]
	sign := 1
	if tz[0] == '-' {
		sign = -1
	}
	hours, err := strconv.Atoi(tz[1:3])
	if err != nil {
		return nil, errors.New("Error parsing author string - invalid timezone hours")
	}
	minutes, err := strconv.Atoi(tz[3:5])
	if err != nil {
		return nil, errors.New("Error parsing author string - invalid timezone minutes")
	}
	offset := sign * (hours*3600 + minutes*60)

	loc := time.FixedZone(tz, offset)
	result.Timestamp = time.Unix(timestamp, 0).In(loc)

	return &result, nil
}
