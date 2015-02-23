package rescuetime

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	simplejson "github.com/bitly/go-simplejson"
)

const (
	dataURL         string = "https://www.rescuetime.com/anapi/data"
	dailySummaryURL string = "https://www.rescuetime.com/anapi/daily_summary_feed"
)

type RescueTime struct {
	ApiKey string
}

type RescueTimeDailySummary struct {
	AllDistractingDurationFormatted             string  `json:"all_distracting_duration_formatted"`
	AllDistractingHours                         float64 `json:"all_distracting_hours"`
	AllDistractingPercentage                    float64 `json:"all_distracting_percentage"`
	AllProductiveDurationFormatted              string  `json:"all_productive_duration_formatted"`
	AllProductiveHours                          float64 `json:"all_productive_hours"`
	AllProductivePercentage                     float64 `json:"all_productive_percentage"`
	BusinessDurationFormatted                   string  `json:"business_duration_formatted"`
	BusinessHours                               float64 `json:"business_hours"`
	BusinessPercentage                          float64 `json:"business_percentage"`
	CommunicationAndSchedulingDurationFormatted string  `json:"communication_and_scheduling_duration_formatted"`
	CommunicationAndSchedulingHours             float64 `json:"communication_and_scheduling_hours"`
	CommunicationAndSchedulingPercentage        float64 `json:"communication_and_scheduling_percentage"`
	Date                                        string  `json:"date"`
	DesignAndCompositionDurationFormatted       string  `json:"design_and_composition_duration_formatted"`
	DesignAndCompositionHours                   float64 `json:"design_and_composition_hours"`
	DesignAndCompositionPercentage              float64 `json:"design_and_composition_percentage"`
	DistractingDurationFormatted                string  `json:"distracting_duration_formatted"`
	DistractingHours                            float64 `json:"distracting_hours"`
	DistractingPercentage                       float64 `json:"distracting_percentage"`
	EntertainmentDurationFormatted              string  `json:"entertainment_duration_formatted"`
	EntertainmentHours                          float64 `json:"entertainment_hours"`
	EntertainmentPercentage                     float64 `json:"entertainment_percentage"`
	ID                                          float64 `json:"id"`
	NeutralDurationFormatted                    string  `json:"neutral_duration_formatted"`
	NeutralHours                                float64 `json:"neutral_hours"`
	NeutralPercentage                           float64 `json:"neutral_percentage"`
	NewsDurationFormatted                       string  `json:"news_duration_formatted"`
	NewsHours                                   float64 `json:"news_hours"`
	NewsPercentage                              float64 `json:"news_percentage"`
	ProductiveDurationFormatted                 string  `json:"productive_duration_formatted"`
	ProductiveHours                             float64 `json:"productive_hours"`
	ProductivePercentage                        float64 `json:"productive_percentage"`
	ProductivityPulse                           float64 `json:"productivity_pulse"`
	ReferenceAndLearningDurationFormatted       string  `json:"reference_and_learning_duration_formatted"`
	ReferenceAndLearningHours                   float64 `json:"reference_and_learning_hours"`
	ReferenceAndLearningPercentage              float64 `json:"reference_and_learning_percentage"`
	ShoppingDurationFormatted                   string  `json:"shopping_duration_formatted"`
	ShoppingHours                               float64 `json:"shopping_hours"`
	ShoppingPercentage                          float64 `json:"shopping_percentage"`
	SocialNetworkingDurationFormatted           string  `json:"social_networking_duration_formatted"`
	SocialNetworkingHours                       float64 `json:"social_networking_hours"`
	SocialNetworkingPercentage                  float64 `json:"social_networking_percentage"`
	SoftwareDevelopmentDurationFormatted        string  `json:"software_development_duration_formatted"`
	SoftwareDevelopmentHours                    float64 `json:"software_development_hours"`
	SoftwareDevelopmentPercentage               float64 `json:"software_development_percentage"`
	TotalDurationFormatted                      string  `json:"total_duration_formatted"`
	TotalHours                                  float64 `json:"total_hours"`
	UncategorizedDurationFormatted              string  `json:"uncategorized_duration_formatted"`
	UncategorizedHours                          float64 `json:"uncategorized_hours"`
	UncategorizedPercentage                     float64 `json:"uncategorized_percentage"`
	UtilitiesDurationFormatted                  string  `json:"utilities_duration_formatted"`
	UtilitiesHours                              float64 `json:"utilities_hours"`
	UtilitiesPercentage                         float64 `json:"utilities_percentage"`
	VeryDistractingDurationFormatted            string  `json:"very_distracting_duration_formatted"`
	VeryDistractingHours                        float64 `json:"very_distracting_hours"`
	VeryDistractingPercentage                   float64 `json:"very_distracting_percentage"`
	VeryProductiveDurationFormatted             string  `json:"very_productive_duration_formatted"`
	VeryProductiveHours                         float64 `json:"very_productive_hours"`
	VeryProductivePercentage                    float64 `json:"very_productive_percentage"`
}

type RescueTimeData struct {
	Notes      string               `json:"notes"`
	RowHeaders []string             `json:"row_headers"`
	Rows       []MiniRescueTimeData `json:"rows"`
}

type MiniRescueTimeData struct {
	Rank           int    `json:"rank"`
	TimeSpent      int    `json:"time_spent"`
	NumberOfPeople int    `json:"number_of_people"`
	Activity       string `json:"activity"`
	Category       string `json:"category"`
	Productivity   int    `json:"productivity"`
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
}

func (r *RescueTime) CheckAPIKey() {
	if r.apiKey == "" {
		fmt.Println("Please provide API key")
		os.Exit(1)
	}
}

func (r *RescueTime) GetResponse(URL string) []byte {
	r.CheckAPIKey()
	response, err := http.Get(URL + "?key=" + r.apiKey + "&format=json")
	checkError(err)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	checkError(err)
	return contents
}

func (r *RescueTime) GetData() RescueTimeData {
	contents := r.GetResponse(dataURL)
	currentJSON := simplejson.New()
	currentJSON.UnmarshalJSON(contents)

	data := RescueTimeData{}

	var notes string
	notes = fmt.Sprintf("%s", (*currentJSON.Get("notes")).MustString())
	data.Notes = notes

	var rowHeaders []string
	rowHeaders, _ = (*currentJSON.Get("row_headers")).StringArray()
	data.RowHeaders = rowHeaders

	var toAppend []MiniRescueTimeData
	for i := 0; i < 36; i++ {
		rows := (*currentJSON.Get("rows")).GetIndex(i)
		current := MiniRescueTimeData{}
		current.Rank, _ = (*rows).GetIndex(0).Int()
		current.TimeSpent, _ = (*rows).GetIndex(1).Int()
		current.NumberOfPeople, _ = (*rows).GetIndex(2).Int()
		current.Activity, _ = (*rows).GetIndex(3).String()
		current.Category, _ = (*rows).GetIndex(4).String()
		current.Productivity, _ = (*rows).GetIndex(5).Int()
		toAppend = append(toAppend, current)
	}
	data.Rows = toAppend
	return data
}

func (r *RescueTime) DailySummary() []RescueTimeDailySummary {
	contents := r.GetResponse(dailySummaryURL)
	keys := make([]RescueTimeDailySummary, 0)
	err := json.Unmarshal(contents, &keys)
	checkError(err)
	return keys
}
