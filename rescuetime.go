package rescuetime

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	simplejson "github.com/bitly/go-simplejson"
)

const (
	analyticDataURL string = "https://www.rescuetime.com/anapi/data"
	dailySummaryURL string = "https://www.rescuetime.com/anapi/daily_summary_feed"
)

// RescueTime contains the user's API key
type RescueTime struct {
	APIKey string
}

// DailySummary is a users summary for a single day
type DailySummary struct {
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

// AnalyticData describes an Analytic Data API result
type AnalyticData struct {
	Notes      string   `json:"notes"`
	RowHeaders []string `json:"row_headers"`
	Rows       []Row    `json:"rows"`
}

// Row is a single row in an Analytic Data API result
type Row struct {
	Date           *time.Time `json:"date,omitempty"`
	Rank           *int       `json:"rank,omitempty"`
	TimeSpent      *int       `json:"timeSpentSeconds,omitempty"`
	NumberOfPeople *int       `json:"numberOfPeople,omitempty"`
	Person         *string    `json:"person,omitempty"`
	Activity       *string    `json:"activity,omitempty"`
	Category       *string    `json:"category,omitempty"`
	Productivity   *int       `json:"productivity,omitempty"`
}

func (r *RescueTime) buildURL(baseURL string, arguments ...[]string) (string, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	query := parsedURL.Query()
	for _, argPair := range arguments {
		query.Add(argPair[0], argPair[1])
	}
	query.Set("key", r.APIKey)
	query.Set("format", "json")
	parsedURL.RawQuery = query.Encode()
	return parsedURL.String(), nil
}

var camelingRegex = regexp.MustCompile("[0-9A-Za-z]+")

func camelCase(src string) string {
	byteSrc := []byte(src)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		if idx > 0 {
			chunks[idx] = bytes.Title(val)
		}
	}
	return string(bytes.Join(chunks, nil))
}

func (r *RescueTime) getResponse(getURL string) ([]byte, error) {
	if r.APIKey == "" {
		return nil, errors.New("Please provide API key")
	}
	response, err := http.Get(getURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

// GetAnalyticData makes a request to the Analytic Data API with the provided (if any) arguments.
// If a timezone is given, all dates will be located in the given timezone.
func (r *RescueTime) GetAnalyticData(timezone string, arguments ...[]string) (AnalyticData, error) {
	var rtd AnalyticData

	builtURL, err := r.buildURL(analyticDataURL, arguments...)
	if err != nil {
		return rtd, err
	}

	contents, err := r.getResponse(builtURL)
	if err != nil {
		return rtd, err
	}
	currentJSON, err := simplejson.NewJson(contents)
	if err != nil {
		return rtd, err
	}

	data := AnalyticData{}

	var notes string
	notes = fmt.Sprintf("%s", currentJSON.Get("notes").MustString())
	data.Notes = notes

	var rowHeaders []string
	headersMap := make(map[int]string)
	for i, s := range currentJSON.Get("row_headers").MustStringArray() {
		rowHeaders = append(rowHeaders, s)
		headersMap[i] = camelCase(strings.ToLower(s))
	}
	data.RowHeaders = rowHeaders

	var toAppend []Row
	for _, entry := range currentJSON.Get("rows").MustArray() {
		out := simplejson.New()
		var entryData Row
		for k, v := range entry.([]interface{}) {
			switch headersMap[k] {
			case "date":
				parsed, err := time.Parse("2006-01-02T15:04:05", v.(string))
				if err != nil {
					return rtd, err
				}
				if timezone != "" {
					location, err := time.LoadLocation(timezone)
					if err != nil {
						return rtd, err
					}
					parsed = parsed.In(location)
				}
				v = parsed.Format(time.RFC3339)
			}
			out.Set(headersMap[k], v)
		}
		encoded, err := out.Encode()
		if err != nil {
			return rtd, err
		}
		err = json.Unmarshal(encoded, &entryData)
		if err != nil {
			return rtd, err
		}
		toAppend = append(toAppend, entryData)
	}
	data.Rows = toAppend
	return data, nil
}

// GetDailySummary returns the daily summary for the user.
func (r *RescueTime) GetDailySummary(arguments ...[]string) ([]DailySummary, error) {
	var summaries []DailySummary
	builtURL, err := r.buildURL(dailySummaryURL, arguments...)
	if err != nil {
		return summaries, err
	}
	contents, err := r.getResponse(builtURL)
	if err != nil {
		return summaries, err
	}
	var keys []DailySummary
	err = json.Unmarshal(contents, &keys)
	if err != nil {
		return summaries, err
	}
	return keys, nil
}
