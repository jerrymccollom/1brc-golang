package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {

	var Mins = make(map[string]float64)
	var Maxs = make(map[string]float64)
	var Sums = make(map[string]float64)
	var Count = make(map[string]float64)

	file, err := os.Open("measurements.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ";")
		if len(parts) == 2 {
			city := parts[0]
			temperature, _ := strconv.ParseFloat(parts[1], 64)
			if temperature < Mins[city] || Mins[city] == 0 {
				Mins[city] = temperature
			}
			if temperature > Maxs[city] {
				Maxs[city] = temperature
			}
			Sums[city] += temperature
			Count[city]++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Collect city names and sort them
	cityNames := make([]string, 0, len(Mins))
	for name := range Mins {
		cityNames = append(cityNames, name)
	}
	sort.Strings(cityNames)

	for _, name := range cityNames {
		average := Sums[name] / Count[name]
		fmt.Printf("%s;%.1f;%.1f;%.1f\n", name, Mins[name], Maxs[name], average)
	}

}
