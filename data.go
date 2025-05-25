package wcollect

import (
	"time"
)

type featureCollection struct {
	Members []member `xml:"member"`
}

type member struct {
	Time  time.Time `xml:"BsWfsElement>Time"`
	Name  string    `xml:"BsWfsElement>ParameterName"`
	Value string    `xml:"BsWfsElement>ParameterValue"`
}

type exception struct {
	Message []string `xml:"Exception>ExceptionText"`
}

type ClickhouseData interface {
	table() string
}

type ObservationData struct {
	Location        string
	ObservationTime time.Time
	CreationTime    time.Time

	Pressure               float32 // Pressure (msl), unit: hPa
	PrecipitationAmount    float32 // Precipitation amount (1h), unit: mm
	PrecipitationIntensity float32 // Precipitation intensity, unit: mm/h
	RelativeHumidity       float32 // Relative humidity, unit: %
	SnowDepth              float32 // Snow depth, unit: cm
	AirTemperature         float32 // Air Temperature, unit: degC
	Dewpoint               float32 // Dew-point temperature, unit: degC
	Visibility             float32 // Horizontal Visibility, unit: m
	WindDirection          float32 // Wind direction (10 min avg), unit: deg
	GustSpeed              float32 // Gust speed (10 min avg), unit: m/s
	WindSpeed              float32 // Wind speed (10 min avg), unit: m/s
	SmartSymbol            int32   // Smart Symbol Code
}

func (ObservationData) table() string { return "observations" }

type ForecastData struct {
	Location        string
	ObservationTime time.Time
	CreationTime    time.Time

	Pressure            float32 // Pressure (msl), unit: hPa
	PrecipitationAmount float32 // Precipitation amount (1h), unit: mm
	RelativeHumidity    float32 // Relative humidity, unit: %
	AirTemperature      float32 // Air Temperature, unit: degC
	Dewpoint            float32 // Dew-point temperature, unit: degC
	WindDirection       float32 // Wind direction (10 min avg), unit: deg
	GustSpeed           float32 // Gust speed (10 min avg), unit: m/s
	WindSpeed           float32 // Wind speed (10 min avg), unit: m/s
	SmartSymbol         int32   // Smart Symbol Code
}

func (ForecastData) table() string { return "forecasts" }
