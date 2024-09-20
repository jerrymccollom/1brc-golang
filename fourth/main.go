package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type City struct {
	Name  string
	Min   float64
	Max   float64
	Sum   float64
	Count float64
}

const (
	Name = iota
	Sign
	Temp
	Fraction
	EndLine
)

func main() {

	cities := make(map[string]*City, 1000)

	file, err := os.Open("measurements.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		state := Name
		name := ""
		sign := 1.0
		temp := 0.0
		line := scanner.Text()
		for i := 0; state != EndLine; i++ {
			c := line[i]
			switch state {
			case Name:
				if c == ';' {
					name = line[:i]
					state = Sign
				}
				break
			case Sign:
				if c == '-' {
					sign = -1.0
				} else {
					state = Temp
					temp = float64(c - '0')
				}
				break
			case Temp:
				if c == '.' {
					state = Fraction
				} else {
					temp = temp*10.0 + float64(c-'0')
				}
				break
			case Fraction:
				state = EndLine
				temp = temp + float64(c-'0')/10.0
				temp *= sign
				// Add city
				city, ok := cities[name]
				if !ok {
					city = &City{Name: name, Min: temp, Max: temp, Sum: temp, Count: 1}
					cities[name] = city
				} else {
					if temp < city.Min {
						city.Min = temp
					}
					if temp > city.Max {
						city.Max = temp
					}
					city.Sum += temp
					city.Count++
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Collect city names and sort them
	var cityNames []string
	for name := range cities {
		cityNames = append(cityNames, name)
	}
	sort.Strings(cityNames)

	for _, name := range cityNames {
		city := cities[name]
		fmt.Printf("%s;%.1f;%.1f;%.1f\n", city.Name, city.Min, city.Max, city.Sum/city.Count)
	}

}
