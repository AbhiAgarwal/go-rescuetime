// Package rescuetime provides a Golang library for the RescueTime API
package rescuetime

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
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

// AnalyticDataQueryParameters is used to provide parameters to the Analytic Data API
type AnalyticDataQueryParameters struct {
	Perspective    string `field_name:"perspective"`
	ResolutionTime string `field_name:"resolution_time"`
	RestrictGroup  string `field_name:"restrict_group"`
	RestrictBegin  string `field_name:"restrict_begin"`
	RestrictEnd    string `field_name:"restrict_end"`
	RestrictKind   string `field_name:"restrict_kind"`
	RestrictThing  string `field_name:"restrict_thing"`
	RestrictThingy string `field_name:"restrict_thingy"`
}

// AnalyticData describes an Analytic Data API result
type AnalyticData struct {
	Notes      string                       `json:"notes"`
	RowHeaders []string                     `json:"row_headers"`
	Rows       []row                        `json:"rows"`
	Parameters *AnalyticDataQueryParameters `json:"-,omitempty"`
}

// Row is a single row in an Analytic Data API result
type row struct {
	Date             time.Time
	Rank             int
	TimeSpentSeconds int
	NumberOfPeople   int
	Person           string
	Activity         string
	Category         string
	Productivity     int
}

func structToMap(i interface{}) (values url.Values) {
	values = url.Values{}
	iVal := reflect.ValueOf(i).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		// Convert each type into a string for the url.Values string map
		var v string
		switch f.Interface().(type) {
		case int, int8, int16, int32, int64:
			v = strconv.FormatInt(f.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			v = strconv.FormatUint(f.Uint(), 10)
		case float32:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 32)
		case float64:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 64)
		case []byte:
			v = string(f.Bytes())
		case string:
			v = f.String()
		}
		if v == "" {
			continue
		}
		values.Set(typ.Field(i).Tag.Get("field_name"), v)
	}
	return
}

func (r *RescueTime) buildURL(baseURL string, urlValues url.Values) (string, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	urlValues.Set("key", r.APIKey)
	urlValues.Set("format", "json")
	parsedURL.RawQuery = urlValues.Encode()
	return parsedURL.String(), nil
}

var titlingRegex = regexp.MustCompile("[0-9A-Za-z]+")

func titleCase(src string) string {
	byteSrc := []byte(src)
	chunks := titlingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		chunks[idx] = bytes.Title(val)
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

// GetAnalyticData makes a request to the Analytic Data API with the provided parameters.
// If a timezone is given, all dates will be located in the given timezone, otherwise system's local timezone.
func (r *RescueTime) GetAnalyticData(timezone string, parameters *AnalyticDataQueryParameters) (AnalyticData, error) {
	var rtd AnalyticData

	params := structToMap(parameters)

	builtURL, err := r.buildURL(analyticDataURL, params)
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

	data := AnalyticData{
		Parameters: parameters,
	}

	var notes string
	notes = fmt.Sprintf("%s", currentJSON.Get("notes").MustString())
	data.Notes = notes

	var rowHeaders []string
	headersMap := make(map[int]string)
	headerRegex := regexp.MustCompile("[^A-Za-z0-9]+")
	for i, s := range currentJSON.Get("row_headers").MustStringArray() {
		rowHeaders = append(rowHeaders, s)
		headersMap[i] = headerRegex.ReplaceAllString(titleCase(s), "")
	}
	data.RowHeaders = rowHeaders

	var toAppend []Row
	for _, entry := range currentJSON.Get("rows").MustArray() {
		var aRow row
		for index, column := range entry.([]interface{}) {
			thisHeader := headersMap[index]
			field := reflect.ValueOf(&aRow).Elem().FieldByName(thisHeader)
			switch field.Interface().(type) {
			case int, int8, int16, int32, int64:
				intValue, err := column.(json.Number).Int64()
				if err != nil {
					return data, err
				}
				field.SetInt(intValue)
			case time.Time:
				parsed, err := time.Parse("2006-01-02T15:04:05", column.(string))
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
				field.Set(reflect.ValueOf(parsed))
			default:
				field.Set(reflect.ValueOf(column))
			}
		}
		toAppend = append(toAppend, aRow)
	}
	data.Rows = toAppend
	return data, nil
}

// GetDailySummary returns the daily summary for the user.
func (r *RescueTime) GetDailySummary() ([]DailySummary, error) {
	var summaries []DailySummary
	builtURL, err := r.buildURL(dailySummaryURL, url.Values{})
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
