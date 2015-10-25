package rescuetime

import (
	"fmt"
	"os"
	"testing"
)

var RescueTimeAPIKey = os.Getenv("RESCUE_TIME_KEY")

func TestMain(m *testing.M) {
	if RescueTimeAPIKey == "" {
		fmt.Println("No API key provided in the RESCUE_TIME_KEY environment variable!")
		os.Exit(1)
	} else {
		os.Exit(m.Run())
	}
}

func TestDailySummary(t *testing.T) {
	var rescue RescueTime
	rescue.APIKey = RescueTimeAPIKey
	response, err := rescue.GetDailySummary()
	if err != nil {
		t.Log(err)
	}
	t.Log(response)
}

func TestGetData(t *testing.T) {
	var rescue RescueTime
	rescue.APIKey = RescueTimeAPIKey
	response, err := rescue.GetAnalyticData("", &AnalyticDataQueryParameters{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response)
}
