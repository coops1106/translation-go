package main

import (
	"github.com/nicksnyder/go-i18n/i18n"
	"fmt"
	"io/ioutil"
	"util/translationPOC/data/output/constants"
	"time"
	"html/template"
	"os"
)

type Inventory struct {
	Material string
	Count    uint
}

type Sensor struct {
	Id string
	ResolveName func(id string) string
}

const translationFilesDir = "./data/output/final"

func main() {
	mustLoadFiles(translationFilesDir)

	gb := mustCreateTranslationFunc("en-gb")
	thai := mustCreateTranslationFunc("th")
	us := mustCreateTranslationFunc("en-us")
	nyk := mustCreateTranslationFunc("nyk")

	fmt.Println(us(constants.ProgramGreeting)) // translated
	fmt.Println(thai(constants.ProgramGreeting)) // translated

	fmt.Println(us(constants.Question)) // translated
	fmt.Println(thai(constants.Question)) // not translated - uses en-gb as default

	args := make(map[string]interface{})
	args["colour"] = "red"
	fmt.Println(us(constants.Answer, args)) // translated, injects arg
	args["colour"] = "blue"
	fmt.Println(thai(constants.Answer, args)) // not translated - used en-gb and injects arg

	fmt.Println(nyk(constants.ProgramGreeting)) // will use the default defined when creating the Tfunc

	t := time.Now()
	args["time"] = t.Format(time.RFC3339)
	fmt.Println(gb(constants.LocalTime, args))
	args["time"] = timeToLocation(t, "America/Winnipeg")
	fmt.Println(us(constants.LocalTime, args))
	args["time"] = timeToLocation(t, "Asia/Bangkok")
	fmt.Println(thai(constants.LocalTime, args))

	//Handling Plurals
	fmt.Println(gb(constants.YourUnreadEmailCount, 0))
	fmt.Println(gb(constants.YourUnreadEmailCount, 1))
	fmt.Println(gb(constants.YourUnreadEmailCount, 2))

	//Condition Demo
	args = make(map[string]interface{})
	user := make(map[string]interface{})
	fmt.Println(gb(constants.IfAndDemo, args)) //Denied
	args["User"] = user
	fmt.Println(gb(constants.IfAndDemo, args)) //Denied
	user["Admin"] = true
	fmt.Println(gb(constants.IfAndDemo, args)) //You are an admin

	//Use functions
	s := Sensor{Id:"abcd"}
	args["Sensor"] = s
	fmt.Println(gb(constants.SensorProblem1, args)) //There's a problem with sensor Uncool sensor

	s = Sensor{Id:"1234"}
	args["Sensor"] = s
	fmt.Println(gb(constants.SensorProblem1, args)) //There's a problem with sensor Cool sensor

	//Use Dynamic functions
	s = Sensor{
		Id: "tttt",
		ResolveName: func(id string) string {
			if id == "asdf" {
				return "Power sensor"
			}
			return "Rubbish sensor"
		},
	}
	args["Sensor"] = s
	fmt.Println(gb(constants.SensorProblem2, args)) //There's a problem with sensor Rubbish sensor

	s.Id = "asdf"
	args["Sensor"] = s
	fmt.Println(gb(constants.SensorProblem2, args)) //There's a problem with sensor Power sensor


	// Defining Template Functions
	funcMap := map[string]interface{}{
		"T": i18n.IdentityTfunc,
		"resolveName": func(in map[string]interface{}) map[string]interface{} { return in }, // no op
	}
	tmpl := template.Must(template.New("Eng").Funcs(funcMap).Parse(`{{T "sensor_problem_3" (resolveName .)}}`))
	tmpl2 := template.Must(template.New("Thai").Funcs(funcMap).Parse(`{{T "sensor_problem_3" (resolveName .)}}`))
	T, _ := i18n.Tfunc("en-gb")
	H, _ := i18n.Tfunc("th")
	tmpl.Funcs(map[string]interface{}{
		"T": T,
		"resolveName": func(in map[string]interface{}) map[string]interface{} {
			out := make(map[string]interface{})
			out["sensorName"] = in["sensorName"]
			return out
		},
	})
	tmpl2.Funcs(map[string]interface{}{
		"T": H,
		"resolveName": func(in map[string]interface{}) map[string]interface{} {
			out := make(map[string]interface{})
			out["sensorName"] = in["sensorName"]
			return out
		},
	})
	args["sensorName"] = "Uber"
	tmpl.Execute(os.Stdout, args)
	fmt.Println("\r")
	tmpl2.Execute(os.Stdout, args)
}

func (s Sensor) GetName() string {
	if s.Id == "1234" {
		return "Cool sensor"
	}
	return "Uncool sensor"
}

func mustCreateTranslationFunc(desiredLang string) i18n.TranslateFunc {
	t, err := i18n.Tfunc(desiredLang, "en-gb")
	if err != nil { panic(err) }
	return t
}

func mustLoadFiles(dirname string) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil { panic(err) }

	for _, f := range files {
		i18n.MustLoadTranslationFile(fmt.Sprintf("%s/%s", dirname, f.Name()))
	}
}

func timeToLocation(t time.Time, l string) string {
	loc, _ := time.LoadLocation(l)
	tt := t.In(loc)
	return tt.Format(time.RFC3339)
}