package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type City struct {
	Name    string
	Min     float64
	Max     float64
	Sum     float64
	Count   float64
	Average float64
}

func main() {

	var Cities = make(map[string]*City, 1000)

	file, err := os.Open("measurements.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		semi := strings.Index(line, ";")
		dot := strings.LastIndex(line, ".")
		name := line[0:semi]
		city, ok := Cities[name]
		if !ok {
			city = &City{Name: name, Min: 0.0, Max: 0.0, Sum: 0.0, Count: 0.0, Average: 0.0}
			Cities[name] = city
		}
		fractionS := line[dot+1:]
		sign := 1.0
		if line[semi+1] == '-' {
			sign = -1.0
			semi++
		}
		temperatureS := line[semi+1 : dot]
		temperature := float64(0.0)
		switch dot - semi - 1 {
		case 1:
			temperature = float64(temperatureS[0]-'0') + float64(fractionS[0]-'0')/10.0
			break
		case 2:
			temperature = float64(temperatureS[1]-'0') + float64(temperatureS[0]-'0')*10.0 + float64(fractionS[0]-'0')/10.0
			break
		default:
			temperature = float64(temperatureS[2]-'0') + float64(temperatureS[1]-'0')*10.0 + float64(temperatureS[0]-'0')*100.0 + float64(fractionS[0]-'0')/10.0
		}
		temperature *= sign
		if temperature < city.Min || city.Min == 0 {
			city.Min = temperature
		}
		if temperature > city.Max {
			city.Max = temperature
		}
		city.Sum += temperature
		city.Count++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	// Collect city names and sort them
	var cityNames []string
	for name := range Cities {
		cityNames = append(cityNames, name)
	}
	sort.Strings(cityNames)

	for _, name := range cityNames {
		city := Cities[name]
		city.Average = city.Sum / city.Count
		fmt.Printf("%s;%.1f;%.1f;%.1f\n", city.Name, city.Min, city.Max, city.Average)
	}
}
