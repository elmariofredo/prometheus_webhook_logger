package logger

import (
	"strings"

	logrus "github.com/Sirupsen/logrus"
	template "github.com/sysincz/prometheus_webhook_logger/template"
	types "github.com/sysincz/prometheus_webhook_logger/types"
)

func callLog(alert types.Alert) {
	//a := RunTemplate("{{ printf \"%#v\" . }}", alert)
	// fmt.Printf("%+v\n", a)
	log.WithFields(logrus.Fields{"alertData": &alert}).Info()

}

//RunTemplate translate template string to string + trimSpace
func RunTemplate(text string, data interface{}) string {
	tmpl := template.Init()

	value, err := tmpl.Execute(text, data)
	if err != nil {
		log.Errorf("Error loading templates from %s: %s", text, err)
		return ""
	}
	value = strings.TrimSpace(value)
	return value
}
