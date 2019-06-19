package types

import (
	"time"
)

//Alert combination of prometheus alert and webhook (webhook with one alert)
type Alert struct {
	address      string
	status       string
	annotations  map[string]string
	labels       map[string]string
	startsAt     time.Time
	endsAt       time.Time
	generatorURL string

	//Prometheus webhook data
	receiver          string
	groupLabels       map[string]string
	commonLabels      map[string]string
	commonAnnotations map[string]string
	externalURL       string
	//Get from Url
	values map[string][]string
}
