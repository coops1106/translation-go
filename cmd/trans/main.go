package main

import (
	"github.com/nicksnyder/go-i18n/i18n"
	"fmt"
	"io/ioutil"
	"util/translationPOC/data/output/constants"
	"time"
)

type Inventory struct {
	Material string
	Count    uint
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

	fmt.Println(gb(constants.YourUnreadEmailCount, 0))
	fmt.Println(gb(constants.YourUnreadEmailCount, 1))
	fmt.Println(gb(constants.YourUnreadEmailCount, 2))


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