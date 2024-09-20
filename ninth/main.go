package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"os"
	"sort"
)

type City struct {
	Name  string
	Hash  int
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

const (
	Name = iota
	Sign
	Temp
	Fraction
	EndLine
)

const Chunks = 1280

func main() {

	reader, err := os.ReadFile("measurements.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	readerLen := len(reader)
	partSize := readerLen / Chunks

	cities := make(map[string]*City)
	cityMaps := make([]map[string]*City, Chunks)
	start := 0
	wg := errgroup.Group{}
	wg.SetLimit(Chunks)

	for i := 0; i < Chunks; i++ {
		end := start + partSize
		if i == Chunks-1 || end >= readerLen {
			end = readerLen
		} else if reader[end] != '\n' {
			for end > start && reader[end] != '\n' {
				end--
			}
		}
		func(a int, b int) {
			wg.Go(func() error {
				cityMaps[i] = processPart(reader, readerLen, a, b+1)
				return nil
			})
		}(start, end)
		start = end + 1
	}
	err = wg.Wait()
	if err != nil {
		fmt.Println("Error processing file:", err)
		return
	}

	for _, parts := range cityMaps {
		for key, val := range parts {
			c, ok := cities[key]
			if !ok {
				cities[key] = val
			} else {
				if val.Min < c.Min {
					c.Min = val.Min
				}
				if val.Max > c.Max {
					c.Max = val.Max
				}
				c.Sum += val.Sum
				c.Count += val.Count
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
		fmt.Printf("%s;%.1f;%.1f;%.1f\n", city.Name, city.Min, city.Max, city.Sum/float64(city.Count))
	}

}

func processPart(reader []byte, readerLen int, first int, last int) map[string]*City {
	cities := make(map[string]*City)
	state := Name
	name := ""
	sign := 1.0
	temp := 0.0
	start := first

	for i := first; i < last && i < readerLen; i++ {
		c := &reader[i]
		switch state {
		case Name:
			if *c == ';' {
				name = string(reader[start:i])
				state = Sign
			}
			break
		case Sign:
			if *c == '-' {
				sign = -1.0
			} else {
				state = Temp
				temp = float64(*c - '0')
			}
			break
		case Temp:
			if *c == '.' {
				state = Fraction
			} else {
				temp = temp*10.0 + float64(*c-'0')
			}
			break
		case Fraction:
			state = EndLine
			temp = temp + float64(*c-'0')/10.0
			break
		default:
			if *c == '\n' {
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
				// reset
				state = Name
				start = i + 1
				name = ""
				sign = 1.0
				temp = 0.0
			}
			break
		}
	}

	return cities
}
