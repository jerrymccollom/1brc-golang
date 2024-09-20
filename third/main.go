package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
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

	state := Name
	name := ""
	c := rune(0)
	sign := 1.0
	temp := 0.0

	scanner := bufio.NewReader(file)
	for {
		if c, _, err = scanner.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}
		switch state {
		case Name:
			if c == ';' {
				state = Sign
			} else {
				name += string(c)
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
			temp = temp + float64(c-'0')/10.0
			temp *= sign
			state = EndLine
			break
		case EndLine:
			if c == '\n' {
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
				// Reset
				name = ""
				state = Name
				c = rune(0)
				sign = 1.0
			}
		}
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
