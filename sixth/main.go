package main

import (
	"bufio"
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

	cities := make([]*City, HashSize)
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
		hash := 0
		line := scanner.Text()
		for i := 0; state != EndLine; i++ {
			c := line[i]
			switch state {
			case Name:
				if c == ';' {
					if hash < 0 {
						hash = -hash
					}
					hash = hash % HashSize
					name = line[:i]
					state = Sign
				} else {
					hash = 31*hash + i + int(c)
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
				city := cities[hash]
				if city == nil {
					city = &City{Name: name, Min: temp, Max: temp, Sum: temp, Count: 1, Hash: hash}
					cities[hash] = city
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

	// Sort by names
	sort.Slice(cities, func(i, j int) bool {
		if cities[i] == nil || cities[j] == nil {
			return cities[i] != nil
		}
		return cities[i].Name < cities[j].Name
	})

	for _, city := range cities {
		if city == nil {
			break
		}
		fmt.Printf("%s;%.1f;%.1f;%.1f\n", city.Name, city.Min, city.Max, city.Sum/city.Count)
	}

}
