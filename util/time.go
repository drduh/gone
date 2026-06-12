package util

import "time"

// IsDaytimeAt returns true during daytime,
// approximately based on the current season.
func IsDaytimeAt(t time.Time) bool {
	hour := t.Hour()
	month := t.Month()

	var sunrise, sunset int

	switch {
	case month >= time.March && month <= time.May:
		sunrise, sunset = 6, 20 // spring
	case month >= time.June && month <= time.August:
		sunrise, sunset = 5, 21 // summer
	case month >= time.September && month <= time.November:
		sunrise, sunset = 7, 19 // fall
	default:
		sunrise, sunset = 8, 17 // winter (dec, jan, feb)
	}

	return hour >= sunrise && hour < sunset
}
