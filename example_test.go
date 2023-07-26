/*!
 * Copyright 2013 Raymond Hill
 *
 * Project: github.com/gorhill/example_test.go
 * File: example_test.go
 * Version: 1.0
 * License: GPL v3 see <https://www.gnu.org/licenses/gpl.html>
 *
 */

package cronexpr

/******************************************************************************/

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

/******************************************************************************/

func TestCronTimers(t *testing.T) {
	locName := "America/Los_Angeles"
	loc, err := time.LoadLocation(locName)

	cases := []struct {
		name     string
		pattern  string
		initTime time.Time
		expected []time.Time
	}{
		{
			"ExampleMustParse with leap year",
			"0 0 29 2 *",
			time.Date(2013, time.August, 31, 0, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2016, time.February, 29, 0, 0, 0, 0, loc),
				time.Date(2020, time.February, 29, 0, 0, 0, 0, loc),
				time.Date(2024, time.February, 29, 0, 0, 0, 0, loc),
				time.Date(2028, time.February, 29, 0, 0, 0, 0, loc),
				time.Date(2032, time.February, 29, 0, 0, 0, 0, loc),
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			starting := c.initTime
			nextTimes := MustParse(c.pattern)
			for _, next := range c.expected {
				n := nextTimes.Next(starting)
				assert.NoError(t, err)
				assert.Equalf(t, next, n, "next time of %v", starting)

				starting = next
			}
		})
	}
}

func TestSystemdTimers(t *testing.T) {

	locName := "America/Los_Angeles"
	loc, err := time.LoadLocation(locName)

	cases := []struct {
		name     string
		pattern  string
		initTime time.Time
		expected []time.Time
	}{
		{
			"normal time, w/o seconds",
			"05:40",
			time.Date(2019, time.February, 7, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.February, 7, 5, 40, 0, 0, loc),
				time.Date(2019, time.February, 8, 5, 40, 0, 0, loc),
				time.Date(2019, time.February, 9, 5, 40, 0, 0, loc),
			},
		},
		{
			"normal time w seconds",
			"05:40:00",
			time.Date(2019, time.February, 7, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.February, 7, 5, 40, 0, 0, loc),
				time.Date(2019, time.February, 8, 5, 40, 0, 0, loc),
				time.Date(2019, time.February, 9, 5, 40, 0, 0, loc),
			},
		},
		{
			"normal time with seconds",
			"08:05:40",
			time.Date(2019, time.February, 7, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.February, 7, 8, 5, 40, 0, loc),
				time.Date(2019, time.February, 8, 8, 5, 40, 0, loc),
				time.Date(2019, time.February, 9, 8, 5, 40, 0, loc),
				time.Date(2019, time.February, 10, 8, 5, 40, 0, loc),
				time.Date(2019, time.February, 11, 8, 5, 40, 0, loc),
				time.Date(2019, time.February, 12, 8, 5, 40, 0, loc),
			},
		},
		{
			"Date",
			"2023-03-05",
			time.Date(2019, time.February, 7, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2023, time.March, 5, 0, 0, 0, 0, loc),
				//time.Date(2019, time.February, 8, 2, 0, 0, 0, loc),
				//time.Date(2019, time.February, 9, 2, 0, 0, 0, loc),
			},
		},
		{
			"Date with time in past",
			"2003-03-05 05:40",
			time.Date(2019, time.February, 7, 1, 0, 0, 0, loc),
			[]time.Time{
				// deep back date
				time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			"Date with time",
			"2020-06-05 05:40:00",
			time.Date(2019, time.February, 7, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2020, time.June, 5, 5, 40, 0, 0, loc),
				// one date and deep back date
				time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			"Daily",
			"daily",
			time.Date(2019, time.February, 7, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.February, 8, 0, 0, 0, 0, loc),
				time.Date(2019, time.February, 9, 0, 0, 0, 0, loc),
				time.Date(2019, time.February, 10, 0, 0, 0, 0, loc),
			},
		},
		{
			"Date with month range",
			"2019-02..04-05",
			time.Date(2019, time.January, 7, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.February, 5, 0, 0, 0, 0, loc),
				time.Date(2019, time.March, 5, 0, 0, 0, 0, loc),
				time.Date(2019, time.April, 5, 0, 0, 0, 0, loc),
				// end of Next()
				time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			"Date with dom range",
			"2019-02-05..08",
			time.Date(2019, time.January, 4, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.February, 5, 0, 0, 0, 0, loc),
				time.Date(2019, time.February, 6, 0, 0, 0, 0, loc),
				time.Date(2019, time.February, 7, 0, 0, 0, 0, loc),
				time.Date(2019, time.February, 8, 0, 0, 0, 0, loc),
				// end of Next()
				time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			"Date with year range",
			"2019..2023-02-05",
			time.Date(2019, time.January, 4, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.February, 5, 0, 0, 0, 0, loc),
				time.Date(2020, time.February, 5, 0, 0, 0, 0, loc),
				time.Date(2021, time.February, 5, 0, 0, 0, 0, loc),
				time.Date(2022, time.February, 5, 0, 0, 0, 0, loc),
				time.Date(2023, time.February, 5, 0, 0, 0, 0, loc),
				// end of Next()
				time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			"Date with day range and divide",
			"2023-02-05..15/3",
			time.Date(2019, time.January, 4, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2023, time.February, 5, 0, 0, 0, 0, loc),
				time.Date(2023, time.February, 8, 0, 0, 0, 0, loc),
				time.Date(2023, time.February, 11, 0, 0, 0, 0, loc),
				time.Date(2023, time.February, 14, 0, 0, 0, 0, loc),
				// end of Next()
				time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			"Every minute",
			"*-*-* *:*:00",
			time.Date(2019, time.January, 4, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.January, 4, 1, 1, 0, 0, loc),
				time.Date(2019, time.January, 4, 1, 2, 0, 0, loc),
				time.Date(2019, time.January, 4, 1, 3, 0, 0, loc),
				time.Date(2019, time.January, 4, 1, 4, 0, 0, loc),
			},
		},
		{
			"Every ten minutes",
			"*-*-* *:*/10:00",
			time.Date(2019, time.January, 4, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.January, 4, 1, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 1, 20, 0, 0, loc),
				time.Date(2019, time.January, 4, 1, 30, 0, 0, loc),
				time.Date(2019, time.January, 4, 1, 40, 0, 0, loc),
				time.Date(2019, time.January, 4, 1, 50, 0, 0, loc),
				time.Date(2019, time.January, 4, 1, 60, 0, 0, loc),
				time.Date(2019, time.January, 4, 2, 10, 0, 0, loc),
			},
		},
		{
			"Every ten minutes with zero",
			"*-*-* *:0/10:00",
			time.Date(2019, time.January, 4, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.January, 4, 1, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 1, 20, 0, 0, loc),
				time.Date(2019, time.January, 4, 1, 30, 0, 0, loc),
				time.Date(2019, time.January, 4, 1, 40, 0, 0, loc),
				time.Date(2019, time.January, 4, 1, 50, 0, 0, loc),
				time.Date(2019, time.January, 4, 1, 60, 0, 0, loc),
				time.Date(2019, time.January, 4, 2, 10, 0, 0, loc),
			},
		},
		{
			"Every day at one hour",
			"*-*-* 01:00:00",
			time.Date(2019, time.January, 4, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.January, 5, 1, 0, 0, 0, loc),
				time.Date(2019, time.January, 6, 1, 0, 0, 0, loc),
				time.Date(2019, time.January, 7, 1, 0, 0, 0, loc),
				time.Date(2019, time.January, 8, 1, 0, 0, 0, loc),
				time.Date(2019, time.January, 9, 1, 0, 0, 0, loc),
				time.Date(2019, time.January, 10, 1, 0, 0, 0, loc),
				time.Date(2019, time.January, 11, 1, 0, 0, 0, loc),
			},
		},
		{
			"Complex hour rule",
			"*-*-* 0..2,4..5,7..23:10:00",
			time.Date(2019, time.January, 4, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.January, 4, 1, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 2, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 4, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 5, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 7, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 8, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 9, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 10, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 11, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 12, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 13, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 14, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 15, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 16, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 17, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 18, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 19, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 20, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 21, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 22, 10, 0, 0, loc),
				time.Date(2019, time.January, 4, 23, 10, 0, 0, loc),
				time.Date(2019, time.January, 5, 0, 10, 0, 0, loc),
			},
		},
		{
			"Range days, list hours",
			"*-*-1..5 04,12:00:00",
			time.Date(2019, time.January, 4, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.January, 4, 4, 0, 0, 0, loc),
				time.Date(2019, time.January, 4, 12, 0, 0, 0, loc),
				time.Date(2019, time.January, 5, 4, 0, 0, 0, loc),
				time.Date(2019, time.January, 5, 12, 0, 0, 0, loc),
			},
		},
		{
			"Hour divider",
			"*-*-* 0/3:00:00",
			time.Date(2019, time.January, 4, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.January, 4, 3, 0, 0, 0, loc),
				time.Date(2019, time.January, 4, 6, 0, 0, 0, loc),
				time.Date(2019, time.January, 4, 9, 0, 0, 0, loc),
				time.Date(2019, time.January, 4, 12, 0, 0, 0, loc),
			},
		},
		{
			"Leap years",
			"*-02-29 01:00:00",
			time.Date(2019, time.January, 4, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2020, time.February, 29, 1, 0, 0, 0, loc),
				time.Date(2024, time.February, 29, 1, 0, 0, 0, loc),
				time.Date(2028, time.February, 29, 1, 0, 0, 0, loc),
				time.Date(2032, time.February, 29, 1, 0, 0, 0, loc),
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			starting := c.initTime
			nextTimes := MustParseSystemd(c.pattern)
			for _, next := range c.expected {
				n := nextTimes.Next(starting)
				assert.NoError(t, err)
				assert.Equalf(t, next, n, "next time of %v", starting)

				starting = next
			}
		})
	}
}
