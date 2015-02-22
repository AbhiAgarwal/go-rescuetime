package rescuetime

import "testing"

func TestDailySummary(t *testing.T) {
	var rescue RescueTime
	rescue.apiKey = ""
	rescue.DailySummary()
}

func TestGetData(t *testing.T) {
	var rescue RescueTime
	rescue.apiKey = ""
	rescue.GetData()
}
