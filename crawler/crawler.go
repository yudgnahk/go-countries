package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	_ "embed"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/yudgnahk/go-countries/dtos"
)

//go:embed templates/country_map.tmpl
var countryMapTemplate string

var specialCountries = []dtos.SpecialCountry{
	{
		Name:     "England",
		Code:     "GB-ENG",
		FlagCode: "GB-ENG",
	},
	{
		Name:     "Scotland",
		Code:     "GB-SCT",
		FlagCode: "GB-SCT",
	},
	{
		Name:     "Wales",
		Code:     "GB-WLS",
		FlagCode: "GB-WLS",
	},
	{
		Name:     "England Short",
		Code:     "ENG",
		FlagCode: "GB-ENG",
	},
}

func main() {

	// Fetch country data with specific fields: cca2/cca3/cioc for code maps, name for generating constant names
	url := "https://restcountries.com/v3.1/all?fields=cca2,cca3,cioc,name"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	countries := make([]dtos.CountryNew, 0)
	if err = json.Unmarshal(body, &countries); err != nil {
		fmt.Printf("error unmarshaling JSON: %v\n", err)
		return
	}

	var countriesList = dtos.Data{
		Countries:        countries,
		SpecialCountries: specialCountries,
	}

	err = generateTemplate(countryMapTemplate, countriesList)

	fmt.Println("finish crawler")
}

func generateTemplate(layout string, data interface{}) error {
	tmpl, err := template.New("tmpl").Funcs(template.FuncMap{
		"format": func(s string) string {
			caser := cases.Title(language.English)
			s = caser.String(strings.ToLower(s))

			if colonIndex := strings.Index(s, ","); colonIndex > -1 {
				s = s[:(colonIndex - 1)]
			}

			s = strings.ReplaceAll(s, " ", "")
			s = strings.ReplaceAll(s, "(", "")
			s = strings.ReplaceAll(s, ")", "")
			s = strings.ReplaceAll(s, ".", "")
			s = strings.ReplaceAll(s, "-", "")
			s = strings.ReplaceAll(s, "'", "")

			return s
		},
	}).Parse(layout)
	if err != nil {
		log.Printf("error when parse layout: %v", err)
		return err
	}
	file := "./country_map.go"
	fo, err := os.Create(file)
	if err != nil {
		log.Printf("Error when create file: %v", err)
		return err
	}

	defer func() {
		if err = fo.Close(); err != nil {
			log.Printf("error close file: %v", err)
		}
	}()

	err = tmpl.Execute(fo, data)
	if err != nil {
		log.Printf("error when exec template: %v", err)
		return err
	}
	return nil
}
