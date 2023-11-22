package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	countries    = make(map[string]bool)
	states       = make(map[string]bool)
	cities       = make(map[string]bool)
	locations    = make(map[string]map[string][]string)
	stateToCity  = make(map[string][]string)
	distributors = make(map[string]Distributor)
	scanner      = bufio.NewScanner(os.Stdin)
)

type Distributor struct {
	Name              string
	IncludeCountry    string
	IncludeState      string
	IncludeCity       string
	ExcludeCountry    []string
	ExcludeState      []string
	ExcludeCity       []string
	ParentDistributor string
}

func loadCsv() {

	file, err := os.Open("cities.csv")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	for index, record := range records {
		if index == 0 {
			continue
		} else {
			city := strings.ToLower(record[3])
			state := strings.ToLower(record[4])
			country := strings.ToLower(record[5])

			cities[city] = true
			states[state] = true
			countries[country] = true

			stateToCity[state] = append(stateToCity[state], city)

			if locations[country] != nil {
				locations[country][state] = append(locations[country][state], city)
			} else {
				locations[country] = make(map[string][]string)
				locations[country][state] = append(locations[country][state], city)
			}
		}
	}

}

func IsLocationValid(location string) bool {

	if !cities[location] && !states[location] && !countries[location] {
		return false
	}

	return true

}

func main() {

	// load data from cities.csv
	loadCsv()

	fmt.Println("-------Content Distribution-------")
	for {
		fmt.Println()

		fmt.Println("1. Add Distributor")
		fmt.Println("2. Add Sub Distributor")
		fmt.Println("3. Check for distribution")
		fmt.Println("4. View distributors")
		fmt.Println("5. Exit")

		fmt.Println("")

		var choice string
		fmt.Print("Enter your choice: ")
		fmt.Scan(&choice)

		switch choice {
		case "1":
			addDistributor(false)
		case "2":
			addDistributor(true)
		case "3":
			checkDistribution()
		case "4":
			printDistributors()
		case "5":
			os.Exit(1)
		default:
			fmt.Println("Invalid choice!")
		}

	}

}

func printDistributors() {

	fmt.Println("########################################")

	for name, distributor := range distributors {
		fmt.Printf("DISTRIBUTOR -> %s\n", name)
		fmt.Printf("\t\tLocations Included -> %s %s %s\n", distributor.IncludeCountry,
			distributor.IncludeState, distributor.IncludeCity)
		fmt.Printf("\t\tLocations Excluded -> %v-%v-%v\n", distributor.ExcludeCountry,
			distributor.ExcludeState, distributor.ExcludeCity)
	}

	fmt.Println("########################################")

}

func checkDistribution() {

	var (
		name, location string
	)

	for {
		fmt.Print("Enter Distributor's name: ")
		if scanner.Scan() {
			name = scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}

		if _, ok := distributors[name]; !ok {
			fmt.Printf("No distributor exist with name %s\n", name)
			continue
		}

	HERE:
		fmt.Print("Enter location: ")
		if scanner.Scan() {
			location = scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}

		location = strings.ToLower(location)

		// validate location entered by user
		if !IsLocationValid(location) {
			fmt.Printf("Invalid location %s!\n", location)
			goto HERE
		}
		break
	}

	if checkIfAuthorized(name, location) {
		fmt.Printf("\nYES, %s is authorized to distribute in %s\n\n", name, location)
	} else {
		fmt.Printf("\nNO, %s is not authorized to distribute in %s\n\n", name, location)
	}

}

func checkIfAuthorized(distributorName, location string) bool {

	distributor := distributors[distributorName]

	var authorized bool

	isExcluded := Isauthorized(location, distributor.ExcludeCity, distributor.ExcludeState, distributor.ExcludeCountry)

	if isExcluded {
		return !isExcluded
	} else {
		authorized = Isauthorized(location, []string{distributor.IncludeCity},
			[]string{distributor.IncludeState}, []string{distributor.IncludeCountry})
		if authorized && distributor.ParentDistributor != "" {
			authorized = checkIfAuthorized(distributor.ParentDistributor, location)
		}
	}

	return authorized

}

func Isauthorized(location string, disCities, disStates, disCountries []string) bool {

	var authorized bool

	if cities[location] {
		for _, city := range disCities {
			if city == location {
				authorized = true
			}
		}
		for _, state := range disStates {
			if cities, ok := stateToCity[state]; ok {
				for _, city := range cities {
					if city == location {
						authorized = true
					}
				}
			}
		}
		for _, country := range disCountries {
			if states, ok := locations[country]; ok {
				for _, cities := range states {
					for _, city := range cities {
						if city == location {
							authorized = true
						}
					}
				}
			}
		}
	}

	if states[location] {
		for _, state := range disStates {
			if state == location {
				authorized = true
			}
		}
		for _, country := range disCountries {
			if states, ok := locations[country]; ok {
				if _, ok := states[location]; ok {
					authorized = true
				}
			}
		}
	}

	if countries[location] {
		for _, country := range disCountries {
			if country == location {
				authorized = true
			}
		}
	}

	return authorized

}

func addDistributor(subDistributor bool) {

	var (
		name, includeLocation, excludeLocationStr string
	)

	fmt.Print("Enter Distributor's name: ")
	if scanner.Scan() {
		name = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Println("To enter multiple Exclude Locations, locations can be hyphen separated")

	for {

		fmt.Print("Add Include location: ")
		if scanner.Scan() {
			includeLocation = scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}

		fmt.Print("Add Exclude location: ")
		if scanner.Scan() {
			excludeLocationStr = scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}

		var parentDistributorName string
		if subDistributor {
		HERE1:
			fmt.Print("Enter Parent Distributor's name: ")
			if scanner.Scan() {
				parentDistributorName = scanner.Text()
			}

			if err := scanner.Err(); err != nil {
				panic(err)
			}

			if _, ok := distributors[parentDistributorName]; !ok {
				fmt.Printf("Distributor named %s does not exist\n", parentDistributorName)
				goto HERE1
			}
		}

		includeLocationValidated, excludeLocations, err := validateLocationInput(includeLocation,
			excludeLocationStr, parentDistributorName)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// create distributor object and add include and exclude locations in it

		distributor := Distributor{
			Name: name,
		}

		if countries[includeLocationValidated] {
			distributor.IncludeCountry = includeLocationValidated
		}
		if states[includeLocationValidated] {
			distributor.IncludeState = includeLocationValidated
		}
		if cities[includeLocationValidated] {
			distributor.IncludeCity = includeLocationValidated
		}

		for _, excludeLocation := range excludeLocations {
			if countries[excludeLocation] {
				distributor.ExcludeCountry = append(distributor.ExcludeCountry, excludeLocation)
			}
			if states[excludeLocation] {
				distributor.ExcludeState = append(distributor.ExcludeState, excludeLocation)
			}
			if cities[excludeLocation] {
				distributor.ExcludeCity = append(distributor.ExcludeCity, excludeLocation)
			}
		}

		if subDistributor {
			distributor.ParentDistributor = parentDistributorName
		}

		distributors[name] = distributor

		break

	}

}

func validateLocationInput(includeLocation, excludeLocationStr, parentDistributorName string) (string, []string, error) {

	var (
		excludeLocs []string
	)

	if excludeLocationStr != "" {
		excludeLocs = strings.Split(excludeLocationStr, "-")
	}

	var locs []string

	locs = append(locs, includeLocation)
	locs = append(locs, excludeLocs...)

	// validate locations entered by user
	for _, location := range locs {
		if !IsLocationValid(location) {
			return "", nil, errors.New(fmt.Sprintf("Invalid location %s!\n", location))
		}
	}

	// validation for redundancy
	if checkForSameLocations(locs) {
		return "", nil, errors.New("no two locations can be same")
	}

	// check if include location of sub distributor comes under all the locations included by parent distributor
	if parentDistributorName != "" {
		inCities, inStates, inCountries := getParentLocationsIncluded(parentDistributorName)

		if !Isauthorized(includeLocation, inCities, inStates, inCountries) {
			return "", nil, errors.New("cannot include locations which are not included by parent distributor")
		}
	}

	var includeCities, includeStates, includeCountries []string

	if cities[includeLocation] {
		includeCities = append(includeCities, includeLocation)
	}
	if states[includeLocation] {
		includeStates = append(includeStates, includeLocation)
	}
	if countries[includeLocation] {
		includeCountries = append(includeCountries, includeLocation)
	}

	// check if exclude locations comes under all the locations included by a distributor
	for _, excludeLoc := range excludeLocs {
		if !Isauthorized(excludeLoc, includeCities, includeStates, includeCountries) {
			return "", nil, errors.New("excluded locations doesn't come under included ones")
		}
	}

	excludeLocs = convertToLowerCase(excludeLocs)

	return strings.ToLower(includeLocation), excludeLocs, nil

}

func getParentLocationsIncluded(distributorName string) ([]string, []string, []string) {

	var (
		includeCities, includeStates, includeCountries []string
	)

	distributor := distributors[distributorName]

	includeCities = append(includeCities, distributor.IncludeCity)
	includeStates = append(includeStates, distributor.IncludeState)
	includeCountries = append(includeCountries, distributor.IncludeCountry)

	if distributor.ParentDistributor != "" {
		inCities, inStates, inCountries := getParentLocationsIncluded(distributor.ParentDistributor)
		includeCities = append(includeCities, inCities...)
		includeStates = append(includeStates, inStates...)
		includeCountries = append(includeCountries, inCountries...)
	}

	return includeCities, includeStates, includeCountries

}

func checkForSameLocations(locs []string) bool {

	for out, location := range locs {
		for in, loc := range locs {
			if loc == location && in != out {
				return true
			}
		}
	}

	return false

}

func convertToLowerCase(locs []string) []string {

	var parsedLocations []string
	for _, loc := range locs {
		parsedLocations = append(parsedLocations, strings.ToLower(loc))
	}

	return parsedLocations

}
