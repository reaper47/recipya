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

package duration_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/reaper47/recipya/internal/utils/duration"
)

func TestParse(t *testing.T) {
	type args struct {
		d string
	}
	testcases := []struct {
		name    string
		args    args
		want    *duration.Duration
		wantErr bool
	}{
		{
			name: "period-only",
			args: args{d: "P4Y"},
			want: &duration.Duration{
				Years: 4,
			},
			wantErr: false,
		},
		{
			name: "time-only-decimal",
			args: args{d: "T2.5S"},
			want: &duration.Duration{
				Seconds: 2.5,
			},
			wantErr: false,
		},
		{
			name: "full",
			args: args{d: "P3Y6M4DT12H30M5.5S"},
			want: &duration.Duration{
				Years:   3,
				Months:  6,
				Days:    4,
				Hours:   12,
				Minutes: 30,
				Seconds: 5.5,
			},
			wantErr: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := duration.Parse(tc.args.d)
			if (err != nil) != tc.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Parse() got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDuration_ToTimeDuration(t *testing.T) {
	type fields struct {
		Years   float64
		Months  float64
		Weeks   float64
		Days    float64
		Hours   float64
		Minutes float64
		Seconds float64
	}
	testcases := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			name: "seconds",
			fields: fields{
				Seconds: 33.3,
			},
			want: time.Second*33 + time.Millisecond*300,
		},
		{
			name: "hours, minutes, and seconds",
			fields: fields{
				Hours:   2,
				Minutes: 33,
				Seconds: 17,
			},
			want: time.Hour*2 + time.Minute*33 + time.Second*17,
		},
		{
			name: "days",
			fields: fields{
				Days: 2,
			},
			want: time.Hour * 24 * 2,
		},
		{
			name: "weeks",
			fields: fields{
				Weeks: 1,
			},
			want: time.Hour * 24 * 7,
		},
		{
			name: "fractional weeks",
			fields: fields{
				Weeks: 12.5,
			},
			want: time.Hour*24*7*12 + time.Hour*84,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			d := &duration.Duration{
				Years:   tc.fields.Years,
				Months:  tc.fields.Months,
				Weeks:   tc.fields.Weeks,
				Days:    tc.fields.Days,
				Hours:   tc.fields.Hours,
				Minutes: tc.fields.Minutes,
				Seconds: tc.fields.Seconds,
			}
			got := d.ToTimeDuration()
			if got != tc.want {
				t.Errorf("ToTimeDuration() = %v, want %v", got, tc.want)
			}
		})
	}
}
