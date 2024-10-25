package service

import "testing"

func TestLoadHolidayFromFile(t *testing.T) {
	year, data, err := LoadHolidayFromFile("../../conf/data/2024_放假通知.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("year: %v, data: %v", year, data)
}
