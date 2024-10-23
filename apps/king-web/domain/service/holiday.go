package service

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Holiday struct {
	Type string
	Note string
}

var (
	regTitle = regexp.MustCompile(`(\d{4}).*节假日安排的通知`)
)

func LoadHolidayFromFile(path string) (string, map[string]Holiday, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return "", nil, err
	}
	defer file.Close()

	var (
		scanner = bufio.NewScanner(file)
		lines   = make([]string, 0, 32)
	)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text != "" {
			lines = append(lines, text)
		}
	}

	if len(lines) == 0 {
		return "", nil, fmt.Errorf("no content from holiday text")
	}

	// title
	year, err := func() (string, error) {
		data := regTitle.FindStringSubmatch(lines[0])
		if len(data) != 2 {
			return "", fmt.Errorf("not found year, nest data: %s", data)
		}
		return data[1], nil
	}()
	if err != nil {
		return "", nil, fmt.Errorf("parse year failure, nest error: %v", err)
	}

	for _, line := range lines[1:] {
		isList(line)
	}

	return year, nil, nil
}

func parseLineHoliday(text string) (map[string]Holiday, error) {
	return nil, nil
}

func isList(text string) bool {
	for _, s := range text {
		// for _, no := range noList {
		// 	for _, v := range no.value {
		// 		if r == v {

		// 		}
		// 	}
		// }
	}

	return true
}

type No struct {
	value []string
}

var noList = map[string]No{
	"1": {value: []string{"一", "1"}},
}
