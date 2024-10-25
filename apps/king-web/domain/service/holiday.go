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

func LoadHolidayFromFile(path string) (string, map[string]Holiday, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0o644)
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

	if len(lines) <= 3 {
		return "", nil, fmt.Errorf("no content from holiday text")
	}

	// title
	year, err := func() (string, error) {
		reg, err := regexp.Compile(`(\d{4}).*节假日安排的通知`)
		if err != nil {
			return "", fmt.Errorf("compile title regexp failure, nest error: %v", err)
		}
		data := reg.FindStringSubmatch(lines[0])
		if len(data) != 2 {
			return "", fmt.Errorf("not found year, nest data: %s", data)
		}
		return data[1], nil
	}()
	if err != nil {
		return "", nil, fmt.Errorf("parse year failure, nest error: %v", err)
	}

	// target
	// target, err := func() ([]string, error) {
	// 	// 各省、自治区、直辖市人民政府，国务院各部委、各直属机构：
	// 	pattern := "[\u4e00-\u9fa5]+(、[\u4e00-\u9fa5]+)*"
	// 	reg, err := regexp.Compile(pattern)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("compilre target regexp failure, nest error: %v", err)
	// 	}

	// 	return reg.FindAllString(lines[1], -1), nil
	// }()
	// if err != nil {
	// 	return "", nil, fmt.Errorf("parse target failure, nest error: %v", err)
	// }
	// fmt.Println(year, target)

	pattern := `^[一二三四五六七]、`
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return "", nil, err
	}
	for _, line := range lines[1:] {
		if reg.MatchString(line) {
			line = reg.ReplaceAllString(line, "")
			data, err := parseLineHoliday(line)
			if err != nil {
				return "", nil, err
			}
			_ = data
		}
	}
	return year, nil, nil
}

func parseLineHoliday(line string) (map[string]Holiday, error) {
	festival, err := func() (string, error) {
		pattern := "[\u4e00-\u9fa5]{2,3}："
		reg, err := regexp.Compile(pattern)
		if err != nil {
			return "", err
		}
		festival := reg.FindString(line)
		festival = strings.Trim(festival, "：")

		if festival == "" {
			return "", fmt.Errorf("parse festival failure, nest error: nil")
		}
		return festival, nil
	}()
	if err != nil {
		return nil, err
	}

	fmt.Println(festival)
	return nil, nil
}
