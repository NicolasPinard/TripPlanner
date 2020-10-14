package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"googlemaps.github.io/maps"
)

func check(e error, errMsg string) {
	if e != nil {
		if strings.Contains(errMsg, "%v") {
			log.Panicf(errMsg, e)
		} else {
			log.Panic(errMsg)
		}
	}
}

func readPlaces(path string) []string {
	file, err := os.Open(path)
	check(err, "Unable to open file: %v")

	var places []string
	reader := bufio.NewReader(file)
	line, isPrefix, err := reader.ReadLine()
	places = append(places, string(line))
	for err == nil && !isPrefix {
		line, isPrefix, err = reader.ReadLine()
		if err == nil {
			places = append(places, string(line))
		}
	}
	file.Close()

	return places
}

func validatePlaces(places []string, mapsClient *maps.Client) {
	var placeSearch *maps.TextSearchRequest
	for _, place := range places {
		placeSearch = &maps.TextSearchRequest{
			Query: place,
		}
		response, err := mapsClient.TextSearch(context.Background(), placeSearch)
		check(err, "The Google Maps Places API returned an error: %v")

		totalResponses := len(response.Results)
		if totalResponses == 0 {
			log.Printf("No results found for place %s\n", place)
		} else if totalResponses == 1 {
			resultedPlace := response.Results[0]
			log.Printf("Found %s.", resultedPlace.Name)
		} else if totalResponses > 1 {
			log.Println("Found multiple matches, please choose from the following and press enter")
			iterNumber := totalResponses / 10
			if (totalResponses % 10) != 0 {
				iterNumber++
			}
			for i := 0; i < iterNumber; i++ {
				remainingResponses := totalResponses - 10*i
				iterUpBound := 10
				if remainingResponses < 10 {
					iterUpBound = remainingResponses
				}
				for j := 0; j < iterUpBound; j++ {
					currentResult := response.Results[i*10+j]
					fmt.Printf("[%d]    %s at %s\n", j, currentResult.Name, currentResult.FormattedAddress)
				}

				fmt.Println()
				reader := bufio.NewReader(os.Stdin)
				input, err := reader.ReadString('\n')
				check(err, "Error while reading input from console: %v")
				input = strings.Replace(input, "\n", "", -1)
				if choice, err := strconv.Atoi(input); err == nil && choice > 0 && choice < iterUpBound {
					log.Printf("Confirmed %s as your desired destination.\n", response.Results[i*10+choice].Name)
				} else {
					log.Printf("Please enter a number from 1 to %d\n", iterUpBound)
					i--
				}
			}
		}
	}
}

func main() {
	var path string = "/tmp/places.txt"
	places := readPlaces(path)
	for idx, place := range places {
		fmt.Println(idx, place)
	}

}
