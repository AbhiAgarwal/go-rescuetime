package rescuetime

import (
	"fmt"
	"os"
	"testing"
)

var RescueTimeApiKey = os.Getenv("RESCUE_TIME_KEY")

func TestMain(m *testing.M) {
	if RescueTimeApiKey == "" {
		fmt.Println("No API key provided in the RESCUE_TIME_KEY environment variable!")
		os.Exit(1)
	} else {
		os.Exit(m.Run())
	}
}

func TestDailySummary(t *testing.T) {
	var rescue RescueTime
	rescue.ApiKey = RescueTimeApiKey
	response, err := rescue.DailySummary()
	if err != nil {
		t.Log(err)
	}
	t.Log(response)
}

func TestGetData(t *testing.T) {
	var rescue RescueTime
	rescue.ApiKey = RescueTimeApiKey
	response, err := rescue.GetData("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response)
}
