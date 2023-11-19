package main

import (
	"encoding/csv"
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
)

type Distributor struct {
	Name              string
	IncludeCountry    []string
	IncludeState      []string
	IncludeCity       []string
	ExcludeCountry    []string
	ExcludeState      []string
	ExcludeCity       []string
	ParentDistributor []Distributor
}

// possible combinations for location input
//country
//state
//city
//state-country
//city-state
//city-state-country
//city-country

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
		if index == 0 || len(record) < 6 {
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
		fmt.Println("2. Add New Permission")
		fmt.Println("3. Check for distribution")
		fmt.Println("4. Exit")

		fmt.Println("")

		var choice string
		fmt.Print("Enter your choice: ")
		fmt.Scan(&choice)

		switch choice {
		case "1":
			addDistributor()
		case "2":
		case "3":
			checkDistribution()
		case "4":
			os.Exit(1)
		default:
			fmt.Println("Invalid choice!")
		}

	}

	// // it is a country
	// if _, ok := locations[locationName]; ok {
	// 	locations[locationName][locationName]
	// } else {
	// 	// it can be a state or a city
	// 	for _, states := range locations {
	// 		if _, ok := states[locationName]; ok {
	// 		} else {
	// 			// it is a city
	// 			for _, cities := range states {
	// 				fmt.Println(cities)
	// 			}
	// 		}
	// 	}
	// }

}

func checkDistribution() {

	var (
		name, location string
	)

	for {
		fmt.Print("Enter Distributor's name: ")
		fmt.Scanln(&name)

		if _, ok := distributors[name]; !ok {
			fmt.Printf("No distributor exist with name %s\n", name)
			continue
		}

	HERE:
		fmt.Print("Enter location: ")
		fmt.Scan(&location)

		// validate location entered by user
		if !IsLocationValid(location) {
			fmt.Printf("Invalid location %s!\n", location)
			goto HERE
		}
		break
	}

	if checkIfAuthorized(name, location) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}

}

// INCLUDE-COUNTRY = INDIA
// INCLUDE-STATE = NULL
// INCLUDE-CITY = NULL

// EXCLUDE-COUNTRY = NULL
// EXCLUDE-STATE = MH
// EXCLUDE-CITY = NULL

// INCLUDE-COUNTRY = NULL
// INCLUDE-STATE = GJ
// INCLUDE-CITY = NULL

// EXCLUDE-COUNTRY = NULL
// EXCLUDE-STATE = NULL
// EXCLUDE-CITY = AHMEDABAD

func checkIfAuthorized(distributorName, location string) bool {

	distributor := distributors[distributorName]

	isExcluded := Isauthorized(location, distributor.ExcludeCity, distributor.ExcludeState, distributor.ExcludeCountry)

	if isExcluded {
		return !isExcluded
	} else {
		return Isauthorized(location, distributor.IncludeCity, distributor.IncludeState, distributor.IncludeCountry)
	}

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

func addDistributor() {

	var (
		name, includeLocationStr, excludeLocationStr string
		locs                                         []string
		valid                                        bool
	)

	fmt.Print("Enter Distributor's name: ")
	fmt.Scanln(&name)

	fmt.Println("Location should be added in format CITY-STATE-COUNTRY")

	for {

		fmt.Print("Add Include location: ")
		fmt.Scanln(&includeLocationStr)

		fmt.Print("Add Exclude location: ")
		fmt.Scanln(&excludeLocationStr)

		includeLocations := parseLocationInput(includeLocationStr)
		excludeLocations := parseLocationInput(excludeLocationStr)

		locs = append(locs, includeLocations...)
		locs = append(locs, excludeLocations...)

		// validate locations entered by user
		for _, location := range locs {
			valid = IsLocationValid(location)
			if !valid {
				fmt.Printf("Invalid location %s!\n", location)
			}
		}

		if !valid {
			continue
		}

		// create distributor object and add include and exclude locations in it

		distributor := Distributor{
			Name: name,
		}

		for _, includeLocation := range includeLocations {
			if countries[includeLocation] {
				distributor.IncludeCountry = append(distributor.IncludeCountry, includeLocation)
			}
			if states[includeLocation] {
				distributor.IncludeState = append(distributor.IncludeState, includeLocation)
			}
			if cities[includeLocation] {
				distributor.IncludeCity = append(distributor.IncludeCity, includeLocation)
			}
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

		distributors[name] = distributor

		break

	}

}

func parseLocationInput(locationStr string) []string {

	locs := strings.Split(locationStr, "-")

	var parsedLocations []string
	for _, loc := range locs {
		parsedLocations = append(parsedLocations, strings.ToLower(loc))
	}

	return parsedLocations

}
