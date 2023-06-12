/* MIT License

Copyright (c) 2022 Kyle McGough

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.*/

package duration

import (
	"errors"
	"math"
	"strconv"
	"time"
	"unicode"
)

// Duration holds all the smaller units that make up the duration
type Duration struct {
	Years   float64
	Months  float64
	Weeks   float64
	Days    float64
	Hours   float64
	Minutes float64
	Seconds float64
}

const (
	parsingPeriod = iota
	parsingTime
)

var errUnexpectedInput = errors.New("unexpected input")

// Parse attempts to parse the given duration string into a *Duration
// if parsing fails an error is returned instead
func Parse(d string) (*Duration, error) {
	state := parsingPeriod
	duration := &Duration{}
	num := ""
	var err error

	for _, char := range d {
		switch char {
		case 'P':
			state = parsingPeriod
		case 'T':
			state = parsingTime
		case 'Y':
			if state != parsingPeriod {
				return nil, errUnexpectedInput
			}

			duration.Years, err = strconv.ParseFloat(num, 64)
			if err != nil {
				return nil, err
			}
			num = ""
		case 'M':
			if state == parsingPeriod {
				duration.Months, err = strconv.ParseFloat(num, 64)
				if err != nil {
					return nil, err
				}
				num = ""
			} else if state == parsingTime {
				duration.Minutes, err = strconv.ParseFloat(num, 64)
				if err != nil {
					return nil, err
				}
				num = ""
			}
		case 'W':
			if state != parsingPeriod {
				return nil, errUnexpectedInput
			}

			duration.Weeks, err = strconv.ParseFloat(num, 64)
			if err != nil {
				return nil, err
			}
			num = ""
		case 'D':
			if state != parsingPeriod {
				return nil, errUnexpectedInput
			}

			duration.Days, err = strconv.ParseFloat(num, 64)
			if err != nil {
				return nil, err
			}
			num = ""
		case 'H':
			if state != parsingTime {
				return nil, errUnexpectedInput
			}

			duration.Hours, err = strconv.ParseFloat(num, 64)
			if err != nil {
				return nil, err
			}
			num = ""
		case 'S':
			if state != parsingTime {
				return nil, errUnexpectedInput
			}

			duration.Seconds, err = strconv.ParseFloat(num, 64)
			if err != nil {
				return nil, err
			}
			num = ""
		default:
			if unicode.IsNumber(char) || char == '.' {
				num += string(char)
				continue
			}

			return nil, errUnexpectedInput
		}
	}

	return duration, nil
}

// ToTimeDuration converts the *Duration to the standard library's time.Duration
// note that for *Duration's with period values of a month or year that the duration becomes a bit fuzzy
// since obviously those things vary month to month and year to year
// I used the values that Google's search provided me with as I couldn't find anything concrete on what they should be
func (duration *Duration) ToTimeDuration() (d time.Duration) {
	d += time.Duration(math.Round(duration.Years * 3.154e+16))
	d += time.Duration(math.Round(duration.Months * 2.628e+15))
	d += time.Duration(math.Round(duration.Weeks * 6.048e+14))
	d += time.Duration(math.Round(duration.Days * 8.64e+13))
	d += time.Duration(math.Round(duration.Hours * 3.6e+12))
	d += time.Duration(math.Round(duration.Minutes * 6e+10))
	d += time.Duration(math.Round(duration.Seconds * 1e+9))
	return d
}
