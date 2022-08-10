package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var cityToFind string
	var radius int64
	var saveToFile bool
	if len(os.Args) == 1 {
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Println("Enter City Name: ")

		scanner.Scan()

		cityToFind = scanner.Text()
		if len(cityToFind) == 0 {
			fmt.Println("City Argument Error: City name must be provided.")
			os.Exit(1)
		}

		fmt.Println("Enter Search Radius (miles): ")

		var radiusIn string
		scanner.Scan()
		radiusIn = scanner.Text()
		parsedRadius, err := strconv.ParseInt(radiusIn, 0, 64)
		if err != nil {
			fmt.Println("Radius Argument Error:", err)
			os.Exit(1)
		}
		radius = parsedRadius

		fmt.Println("Save result to file? (y/n): ")

		scanner.Scan()
		saveToFile = strings.ToLower(scanner.Text()) == "y"
	} else if len(os.Args) != 4 {
		fmt.Println("Error: Incorrect number of arguments. Please provide city-name radius saveToFile(y/n)")
		os.Exit(1)
	} else if len(os.Args) == 4 {
		cityToFind = os.Args[1]
		if len(cityToFind) == 0 {
			fmt.Println("City Argument Error: City name must be provided.")
			os.Exit(1)
		}
		var err error
		radius, err = strconv.ParseInt(os.Args[2], 0, 64)
		if err != nil {
			fmt.Println("Radius Argument Error:", err)
			os.Exit(1)
		}
		saveToFile = strings.ToLower(os.Args[3]) == "y"
	}

	fileName := "uscities-forgo.csv"
	sep := "\t"
	cities, errs := getCitiesFromFile(fileName, sep)
	if len(errs) > 0 {
		fmt.Println(errs)
		os.Exit(1)
	}

	refCity := cities.findCityByName(cityToFind)
	foundCities := cities.getCitiesWithinRadius(refCity, float64(radius))
	foundCities.sortByDensity()
	fmt.Println("")
	if saveToFile {
		err := foundCities.saveToFile("result.json")
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	} else {
		for _, c := range foundCities {
			fmt.Println(c)
		}
		fmt.Println(len(foundCities))
	}
}
