package main

import (
	"text/template"
	"os"
	"time"
)

var tmpl *template.Template

func main() {
	initializeTemplate()
	resolve("There's a problem with sensor {{ResolveSensorName \"123456789\"}} at {{LocalTime \"2017-12-12T00:00:00Z\"}}")
}

func resolve(tmplDefn string) {
	tmpl, err := tmpl.Parse(tmplDefn)
	if err != nil { panic(err) }

	args := make(map[string]interface{})
	tmpl.Execute(os.Stdout, args)
}

func initializeTemplate() {
	tmpl = template.New("example")

	funcMap := make(map[string]interface{})
	funcMap["ResolveSensorName"] = resolveSensorName
	funcMap["LocalTime"] = resolveLocalTime

	tmpl.Funcs(funcMap)
}

func resolveLocalTime(lt string) string {
	t, err := time.Parse(time.RFC3339, lt)
	if err != nil { panic(err) }
	loc, _ := time.LoadLocation("America/Winnipeg")
	tt := t.In(loc)
	return tt.Format(time.RFC3339)
}

func resolveSensorName(sensorGuid string) string {
	if sensorGuid == "123456789" {
		return "Deep Blue"
	} else if sensorGuid == "987654321" {
		return "Aleph"
	}
	return "HAL9000"
}