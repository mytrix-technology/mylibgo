package datetime

import "time"

func TimeToString(t time.Time) string {
	return t.Format("2006-01-02")
}

func StringToTime(str string) time.Time {
	layout := "2006-01-02 15:04:05"
	t, _ := time.Parse(layout, str)
	return t
}

func TimeToStringPattern(t time.Time, pattern string) string {
	return t.Format(pattern)
}

func monthToInt(month string) int {
	switch month {
	case "January":
		return 1
	case "February":
		return 2
	case "March":
		return 3
	case "April":
		return 4
	case "May":
		return 5
	case "June":
		return 6
	case "July":
		return 7
	case "August":
		return 8
	case "September":
		return 9
	case "October":
		return 10
	case "November":
		return 11
	case "December":
		return 12
	default:
		panic("Unrecognized month")
	}
}

func ToIndonesianDay(day string) string {
	switch day {
	case "Monday":
		return "Senin"
	case "Tuesday":
		return "Selasa"
	case "Wednesday":
		return "Rabu"
	case "Thursday":
		return "Kamis"
	case "Friday":
		return "Jumat"
	case "Saturday":
		return "Sabtu"
	case "Sunday":
		return "Minggu"
	}
	return ""
}