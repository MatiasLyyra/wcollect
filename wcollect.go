package wcollect

import (
	"encoding/xml"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func FetchObservations(location string, start, end time.Time) ([]ObservationData, error) {
	const parameters = "p_sea,r_1h,rh,ri_10min,snow_aws,t2m,td,vis,wd_10min,wg_10min,ws_10min,SmartSymbol"
	fc, err := querySimpleEndpoint(location, "fmi::observations::weather::simple", parameters, start, end)
	if err != nil {
		return nil, err
	}
	return parseObservations(fc, location), nil
}
func FetchForecast(location string, start, end time.Time) ([]ForecastData, error) {
	const parameters = "Pressure,Precipitation1h,Humidity,Temperature,DewPoint,WindDirection,HourlyMaximumGust,WindSpeedMS,SmartSymbol"
	fc, err := querySimpleEndpoint(location, "fmi::forecast::edited::weather::scandinavia::point::simple", parameters, start, end)
	if err != nil {
		return nil, err
	}
	return parseForecasts(fc, location), nil
}

func querySimpleEndpoint(location, endpoint, parameters string, start, end time.Time) (*featureCollection, error) {
	var c http.Client
	var fc featureCollection

	req, err := http.NewRequest(http.MethodGet, "https://opendata.fmi.fi/wfs/fin", nil)
	if err != nil {
		return nil, err
	}
	values := make(url.Values)
	values.Add("service", "WFS")
	values.Add("version", "2.0.0")
	values.Add("request", "GetFeature")
	values.Add("storedquery_id", endpoint)
	values.Add("parameters", parameters)
	values.Add("place", location)
	values.Add("starttime", start.UTC().Format(time.RFC3339))
	values.Add("endtime", end.UTC().Format(time.RFC3339))
	req.URL.RawQuery = values.Encode()

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	dec := xml.NewDecoder(res.Body)
	if res.StatusCode != http.StatusOK {
		var e exception
		if err := dec.Decode(&e); err != nil || len(e.Message) == 0 {
			return nil, fmt.Errorf("api error: unknown error (%v)", res.Status)
		}
		return nil, fmt.Errorf("api error: %v (%v)", e.Message[0], res.Status)
	}
	if err := dec.Decode(&fc); err != nil {
		return nil, err
	}
	return &fc, nil
}

func parseObservations(fc *featureCollection, place string) []ObservationData {
	datapoinsts := make(map[time.Time]ObservationData)
	now := time.Now().UTC()
	for _, mem := range fc.Members {
		point := datapoinsts[mem.Time]
		if point.Location == "" {
			point.Location = place
		}
		if point.ObservationTime.IsZero() {
			point.ObservationTime = mem.Time
		}
		if point.CreationTime.IsZero() {
			point.CreationTime = now
		}
		dataValue, _ := strconv.ParseFloat(mem.Value, 64)
		if math.IsNaN(dataValue) {
			dataValue = 0
		}
		v := float32(dataValue)
		switch mem.Name {
		case "p_sea":
			point.Pressure = v
		case "r_1h":
			point.PrecipitationAmount = v
		case "rh":
			point.RelativeHumidity = v
		case "ri_10min":
			point.PrecipitationIntensity = v
		case "snow_aws":
			if v > 0 {
				point.SnowDepth = v
			}
		case "t2m":
			point.AirTemperature = v
		case "td":
			point.Dewpoint = v
		case "vis":
			point.Visibility = v
		case "wd_10min":
			point.WindDirection = v
		case "wg_10min":
			point.GustSpeed = v
		case "ws_10min":
			point.WindSpeed = v
		case "SmartSymbol":
			point.SmartSymbol = int32(v)
		}
		datapoinsts[mem.Time] = point
	}
	data := make([]ObservationData, 0, len(datapoinsts))
	for _, v := range datapoinsts {
		data = append(data, v)
	}
	return data
}

func parseForecasts(fc *featureCollection, place string) []ForecastData {
	datapoinsts := make(map[time.Time]ForecastData)
	now := time.Now().UTC()
	for _, mem := range fc.Members {
		point := datapoinsts[mem.Time]
		if point.Location == "" {
			point.Location = place
		}
		if point.ObservationTime.IsZero() {
			point.ObservationTime = mem.Time
		}
		if point.CreationTime.IsZero() {
			point.CreationTime = now
		}
		dataValue, _ := strconv.ParseFloat(mem.Value, 64)
		if math.IsNaN(dataValue) {
			dataValue = 0
		}
		v := float32(dataValue)
		switch mem.Name {
		case "Pressure":
			point.Pressure = v
		case "Precipitation1h":
			point.PrecipitationAmount = v
		case "Humidity":
			point.RelativeHumidity = v
		case "Temperature":
			point.AirTemperature = v
		case "DewPoint":
			point.Dewpoint = v
		case "WindDirection":
			point.WindDirection = v
		case "HourlyMaximumGust":
			point.GustSpeed = v
		case "WindSpeedMS":
			point.WindSpeed = v
		case "SmartSymbol":
			point.SmartSymbol = int32(v)
		}
		datapoinsts[mem.Time] = point
	}
	data := make([]ForecastData, 0, len(datapoinsts))
	for _, v := range datapoinsts {
		data = append(data, v)
	}
	return data
}
