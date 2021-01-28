package metar

import (
	"io/ioutil"
	"net/http"
	"time"
)

// CloudAmount enum type
type CloudAmount int

// CloudAmount enum
const (
	CloudsFew CloudAmount = iota
	CloudsScattered
	CloudsBroken
	CloudsOvercast
)

// SkyCoverLevel describes a single sky cover level as reported by METARs
type SkyCoverLevel struct {
	Amount   CloudAmount
	Altitude int
}

// AirportData is a metar data for a single airport
type AirportData struct {
	RawText                   string
	StationID                 string
	ObservationTime           time.Time
	Latitude                  float64
	Longitude                 float64
	TempC                     float64
	DewPointC                 float64
	WindDirectionDegrees      int64
	WindSpeedKts              int64
	WindGustKts               int64
	VisibilityStatuteMi       float64
	AltimeterHg               float64
	QNH                       int64
	SeaLevelPressureHg        string
	Corrected                 string
	Auto                      string
	AutoStation               string
	MaintenanceIndicatorOn    string
	NoSignal                  string
	LightningSensorOff        string
	FreezingRainSensorOff     string
	PresentWeatherSensorOff   string
	WxString                  string
	FlightCat                 string
	ThreeHrPressureTendencyMb string
	MaxTempC                  string
	MinTempC                  string
	MaxTemp24hrC              string
	MinTemp24hrC              string
	PrecipIn                  string
	Precip3hrIn               string
	Precip6hrIn               string
	Precip24hrIn              string
	SnowIn                    string
	VertVisFt                 string
	MetarType                 string
	ElevationM                string
	SkyCover                  []SkyCoverLevel
}

const (
	dataURL          = "https://aviationweather.gov/adds/dataserver_current/current/metars.cache.csv"
	dataFetchTimeout = 1 * time.Second

	maxRetries = 5
)

func (m *Metar) fetch() {

	retries := maxRetries

	for retries > 0 {
		c := http.Client{Timeout: dataFetchTimeout}
		resp, err := c.Get(dataURL)
		if err != nil {
			retries--
			continue
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			retries--
			continue
		}

		aipData, err := parse(data)
		if err != nil {
			retries--
			continue
		}

		m.lock.Lock()
		for _, airport := range aipData {
			m.data[airport.StationID] = airport
		}
		m.lock.Unlock()
		return
	}
}
