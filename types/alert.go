package types

import (
	"time"
)

//Alert combination of prometheus alert and webhook (webhook with one alert)
type Alert struct {
	Address      string            `json:"address"`
	Status       string            `json:"status"`
	Annotations  map[string]string `json:"annotations"`
	Labels       map[string]string `json:"labels"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`

	//Prometheus webhook data
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	//Get from Url
	URLValues map[string][]string `json:"urlValues"`
}
