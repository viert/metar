package metar

import (
	"bufio"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func parseDate(date string) (time.Time, error) {
	return time.Parse(time.RFC3339, date)
}

func baroPressure(hgString string) (float64, int64, error) {
	hg, err := strconv.ParseFloat(hgString, 64)
	if err != nil {
		return 0, 0, err
	}
	qnh := int64(math.Round(33.86 * hg))
	return hg, qnh, nil
}

func parseSkyCover(tokens []string) ([]SkyCoverLevel, error) {
	var amount CloudAmount
	cover := make([]SkyCoverLevel, 0)
	for i, token := range tokens {
		if i%2 == 0 {
			switch token {
			case "FEW":
				amount = CloudsFew
			case "SCT":
				amount = CloudsScattered
			case "BKN":
				amount = CloudsBroken
			case "OVC":
				amount = CloudsOvercast
			default:
				return nil, fmt.Errorf("invalid sky cover string \"%s\"", token)
			}
		} else {
			alt, err := strconv.ParseInt(token, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid cloud level string \"%s\"", token)
			}
			level := SkyCoverLevel{
				Amount:   amount,
				Altitude: int(alt),
			}
			cover = append(cover, level)
		}
	}
	return cover, nil
}

func parse(data []byte) ([]AirportData, error) {
	var fieldsOrder []string

	aipData := make([]AirportData, 0)

	rd := strings.NewReader(string(data))
	sc := bufio.NewScanner(rd)
	for sc.Scan() {
		line := sc.Text()
		tokens := strings.Split(line, ",")
		if len(tokens) < 2 {
			continue
		}

		if fieldsOrder == nil {
			fieldsOrder = tokens
			continue
		}

		var airport AirportData
		skyCoverStrs := make([]string, 0)
		for i, field := range fieldsOrder {
			value := tokens[i]
			switch field {
			case "raw_text":
				airport.RawText = value
			case "station_id":
				airport.StationID = value
			case "observation_time":
				airport.ObservationTime, _ = parseDate(value)
			case "latitude":
				airport.Latitude, _ = strconv.ParseFloat(value, 64)
			case "longitude":
				airport.Longitude, _ = strconv.ParseFloat(value, 64)
			case "temp_c":
				airport.TempC, _ = strconv.ParseFloat(value, 64)
			case "dewpoint_c":
				airport.DewPointC, _ = strconv.ParseFloat(value, 64)
			case "wind_dir_degrees":
				airport.WindDirectionDegrees, _ = strconv.ParseInt(value, 10, 64)
			case "wind_speed_kt":
				airport.WindSpeedKts, _ = strconv.ParseInt(value, 10, 64)
			case "wind_gust_kt":
				airport.WindGustKts, _ = strconv.ParseInt(value, 10, 64)
			case "visibility_statute_mi":
				airport.VisibilityStatuteMi, _ = strconv.ParseFloat(value, 64)
			case "altim_in_hg":
				airport.AltimeterHg, airport.QNH, _ = baroPressure(value)
			case "sea_level_pressure_hg":
				airport.SeaLevelPressureHg = value
			case "corrected":
				airport.Corrected = value
			case "auto":
				airport.Auto = value
			case "auto_station":
				airport.AutoStation = value
			case "maintenance_indicator_on":
				airport.MaintenanceIndicatorOn = value
			case "no_signal":
				airport.NoSignal = value
			case "lightning_sensor_off":
				airport.LightningSensorOff = value
			case "freezing_rain_sensor_off":
				airport.FreezingRainSensorOff = value
			case "present_weather_sensor_off":
				airport.PresentWeatherSensorOff = value
			case "wx_string":
				airport.WxString = value
			case "flight_cat":
				airport.FlightCat = value
			case "three_hr_pressure_tendency_mb":
				airport.ThreeHrPressureTendencyMb = value
			case "max_temp_c":
				airport.MaxTempC = value
			case "min_temp_c":
				airport.MinTempC = value
			case "max_temp_2_4hr_c":
				airport.MaxTemp24hrC = value
			case "min_temp_2_4hr_c":
				airport.MinTemp24hrC = value
			case "precip_in":
				airport.PrecipIn = value
			case "precip_3hr_in":
				airport.Precip3hrIn = value
			case "precip_6hr_in":
				airport.Precip6hrIn = value
			case "precip_2_4hr_in":
				airport.Precip24hrIn = value
			case "snow_in":
				airport.SnowIn = value
			case "vert_vis_ft":
				airport.VertVisFt = value
			case "metar_type":
				airport.MetarType = value
			case "elevation_m":
				airport.ElevationM = value
			case "cloud_base_ft_agl":
				fallthrough
			case "sky_cover":
				if value != "" {
					skyCoverStrs = append(skyCoverStrs, value)
				}
			}
		}
		airport.SkyCover, _ = parseSkyCover(skyCoverStrs)
		aipData = append(aipData, airport)
	}
	return aipData, nil
}
