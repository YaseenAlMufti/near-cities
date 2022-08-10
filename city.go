package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
	"strings"
)

type Point struct {
	lat float64
	lng float64
}

type CityRow struct {
	city     string
	state    string
	point    Point
	density  float64
	timezone string
}

type Cities []CityRow

func (cr *CityRow) fromString(s string, sep string) []error {
	var errors []error
	rowArr := strings.Split(s, sep)
	if len(rowArr) > 0 {
		cr.city = strings.TrimSpace(rowArr[0])
		cr.state = strings.TrimSpace(rowArr[1])
		latFloat, err := strconv.ParseFloat(strings.TrimSpace(rowArr[2]), 64)
		if err != nil {
			errors = append(errors, err)
			cr.point.lat = 0
		} else {
			cr.point.lat = latFloat
		}

		lngFloat, err := strconv.ParseFloat(strings.TrimSpace(rowArr[3]), 64)
		if err != nil {
			errors = append(errors, err)
			cr.point.lng = 0
		} else {
			cr.point.lng = lngFloat
		}

		densityInt, err := strconv.ParseFloat(strings.TrimSpace(rowArr[4]), 64)
		if err != nil {
			errors = append(errors, err)
			cr.density = 0
		} else {
			cr.density = densityInt
		}
		cr.timezone = strings.TrimSpace(rowArr[5])
	}
	return errors
}

func getDistanceBetweenTwoPoints(
	point1 Point,
	point2 Point,
) float64 {
	radlat1 := float64(math.Pi * point1.lat / 180)
	radlat2 := float64(math.Pi * point2.lat / 180)

	theta := float64(point1.lng - point2.lng)
	radtheta := float64(math.Pi * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515
	return math.Round(dist)
}

func (c Cities) findCityByName(s string) CityRow {
	foundCity := CityRow{}
	for _, city := range c {
		if strings.EqualFold(city.city, s) {
			foundCity = city
			break
		}
	}
	return foundCity
}

// radius is in miles
func (c Cities) getCitiesWithinRadius(city CityRow, radius float64) Cities {
	var citiesThatMeetCriteria Cities

	for _, cityRef := range c {
		distance := getDistanceBetweenTwoPoints(city.point, cityRef.point)
		if distance < radius {
			citiesThatMeetCriteria = append(citiesThatMeetCriteria, cityRef)
		}
	}

	return citiesThatMeetCriteria
}

func getCitiesFromFile(fileName string, sep string) (Cities, []error) {
	var errors []error
	var cities Cities
	bytesCities, err := ioutil.ReadFile(fileName)
	if err != nil {
		errors = append(errors, err)
		return cities, errors
	}
	stringifiedCitiesArr := strings.Split(string(bytesCities), "\n")
	for index, stringifiedCity := range stringifiedCitiesArr {
		// skipping headers line
		if index > 0 {
			var city CityRow
			errs := city.fromString(stringifiedCity, sep)
			if len(errs) > 0 {
				fmt.Println(stringifiedCity)
			}
			errors = append(errors, errs...)
			if len(city.city) > 0 {
				cities = append(cities, city)
			}
		}
	}
	return cities, errors
}

func (c Cities) toJSON() string {
	jsonString := "[\n"
	for index, city := range c {
		commaAtEnd := ","
		if index == len(c)-1 {
			commaAtEnd = ""
		}
		jsonString = jsonString + fmt.Sprintf(
			"\t{ \"city\": \"%s\", \"state\": \"%s\", \"lat\": %v, \"lng\": %v, \"density\": %v, \"timezone\": \"%s\" }%s\n",
			city.city,
			city.state,
			city.point.lat,
			city.point.lng,
			city.density,
			city.timezone,
			commaAtEnd,
		)
	}
	jsonString = jsonString + "]"
	return jsonString
}

func (c Cities) saveToFile(fileName string) error {
	bytesCities := []byte(c.toJSON())
	return ioutil.WriteFile(fileName, bytesCities, fs.ModeAppend)
}

func (c Cities) sortByDensity() {
	sort.SliceStable(c, func(i, j int) bool {
		return c[i].density > c[j].density
	})
}
