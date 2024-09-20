package main

import (
	"fmt"
	"os"
	"sort"
)

type City struct {
	Name  string
	Hash  int
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

	HashSize = 1000000
)

func main() {

	reader, err := os.ReadFile("measurements.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	var cities [HashSize]*City
	var cityList []*City

	max := len(reader)
	state := Name
	name := ""
	sign := 1.0
	temp := 0.0
	hash := 0
	start := 0
	tick := 0

	for i := 0; i < max; i++ {
		c := reader[i]
		switch state {
		case Name:
			if c == ';' {
				if hash < 0 {
					hash = -hash
				}
				hash = hash % HashSize
				name = string(reader[start:i])
				state = Sign
			} else {
				hash = 31*hash + tick + int(c)
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
			break
		case EndLine:
			if c == '\n' {
				temp *= sign
				// Add city
				city := cities[hash]
				if city == nil {
					city = &City{Name: name, Min: temp, Max: temp, Sum: temp, Count: 1, Hash: hash}
					cities[hash] = city
					cityList = append(cityList, city)
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
				// reset
				state = Name
				start = i + 1
				name = ""
				sign = 1.0
				temp = 0.0
				hash = 0
				tick = 0
			}
			break
		}
	}

	// Sort by names
	sort.Slice(cityList, func(i, j int) bool {
		return cityList[i].Name < cityList[j].Name
	})

	for _, city := range cityList {
		fmt.Printf("%s;%.1f;%.1f;%.1f\n", city.Name, city.Min, city.Max, city.Sum/city.Count)
	}
}
