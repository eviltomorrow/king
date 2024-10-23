package service

import "testing"

func TestLoadHolidayFromFile(t *testing.T) {
	LoadHolidayFromFile("../../data/2024_放假通知.txt")
}
