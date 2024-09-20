package main

import (
	"bufio"
	"fmt"
	"golang.org/x/exp/mmap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"os"
	"sort"
)

type City struct {
	Name  string
	Min   int64
	Max   int64
	Sum   int64
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

	fileName := "measurements.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}
	reader, err := mmap.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer func() { _ = reader.Close() }()

	readerLen := reader.Len()
	partSize := readerLen / Chunks

	cities := make(map[int64]*City)
	cityMaps := make([]map[int64]*City, Chunks)
	start := 0
	wg := errgroup.Group{}
	wg.SetLimit(Chunks)

	for i := 0; i < Chunks; i++ {
		end := start + partSize
		if i == Chunks-1 || end >= readerLen {
			end = readerLen
		} else if reader.At(end) != '\n' {
			for end > start && reader.At(end) != '\n' {
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
	var citiesByName []*City
	for _, city := range cities {
		citiesByName = append(citiesByName, city)
	}
	collator := collate.New(language.English)
	sort.Slice(citiesByName, func(i, j int) bool {
		return collator.Compare([]byte(citiesByName[i].Name), []byte(citiesByName[j].Name)) < 0
	})

	// Output
	outfile, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer func() { _ = outfile.Close() }()
	out := bufio.NewWriter(outfile)

	for _, city := range citiesByName {
		minSign := int64(-1)
		if city.Min > 0 {
			minSign = int64(1)
		}
		maxSign := int64(-1)
		if city.Max > 0 {
			maxSign = int64(1)
		}
		_, _ = fmt.Fprintf(out, "%s;%d.%.1d;%.1f;%d.%.1d\n", city.Name, city.Min/int64(10), minSign*city.Min%int64(10), float64(city.Sum)/float64(city.Count)/float64(10), city.Max/int64(10), maxSign*city.Max%int64(10))
	}
	_ = out.Flush()
}

const HashSize = 1000000

func processPart(reader *mmap.ReaderAt, readerLen int, first int, last int) map[int64]*City {
	cities := make(map[int64]*City)
	state := Name
	name := ""
	sign := int64(1)
	temp := int64(0)
	start := first
	hash := int64(0)
	buf := make([]byte, 128)
	idx := 0

	for i := first; i < last && i < readerLen; i++ {
		c := reader.At(i)
		switch state {
		case Name:
			if c == ';' {
				hash = hash % HashSize
				if hash < 0 {
					hash = -hash + HashSize
				}
				name = string(buf[:i-start])
				state = Sign
			} else {
				buf[idx] = c
				idx++
				hash = int64(31)*hash + int64(c)
			}
			break
		case Sign:
			if c == '-' {
				sign = int64(-1)
			} else {
				state = Temp
				temp = int64(c - '0')
			}
			break
		case Temp:
			if c == '.' {
				state = Fraction
			} else {
				temp = temp*10 + int64(c-'0')
			}
			break
		case Fraction:
			state = EndLine
			temp = temp*10 + int64(c-'0')
			break
		default:
			if c == '\n' {
				temp *= sign
				// Add city
				city, ok := cities[hash]
				if !ok {
					city = &City{Name: name, Min: temp, Max: temp, Sum: temp, Count: 1}
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
				// reset
				state = Name
				start = i + 1
				name = ""
				sign = int64(1)
				temp = int64(0)
				hash = int64(0)
				idx = 0
			}
			break
		}
	}

	return cities
}
