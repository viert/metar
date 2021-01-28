package metar

import (
	"runtime"
	"sync"
	"time"
)

// Metar is a store and fetcher
type Metar struct {
	data map[string]AirportData
	stop chan bool
	lock sync.RWMutex
}

// New creates a new Metar object
func New(period time.Duration) *Metar {
	m := &Metar{
		data: make(map[string]AirportData),
		stop: make(chan bool),
	}

	runtime.SetFinalizer(m, finalize)

	go m.loop(period)
	return m
}

func (m *Metar) loop(period time.Duration) {
	m.fetch()
	t := time.NewTicker(period)
	for {
		select {
		case <-t.C:
			m.fetch()
		case <-m.stop:
			t.Stop()
			return
		}
	}
}

func finalize(m *Metar) {
	m.stop <- true
	close(m.stop)
}

// GetAirportData returns airport data by a given airport ICAO code
func (m *Metar) GetAirportData(stationID string) *AirportData {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if airport, found := m.data[stationID]; found {
		return &airport
	}
	return nil
}

// AirportList returns a list of airport ICAO codes fetched from METAR dataserver
func (m *Metar) AirportList() []string {
	m.lock.RLock()
	defer m.lock.RUnlock()

	result := make([]string, 0)
	for key := range m.data {
		result = append(result, key)
	}
	return result
}
