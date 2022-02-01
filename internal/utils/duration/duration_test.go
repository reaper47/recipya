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
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name    string
		args    args
		want    *Duration
		wantErr bool
	}{
		{
			name: "period-only",
			args: args{d: "P4Y"},
			want: &Duration{
				Years: 4,
			},
			wantErr: false,
		},
		{
			name: "time-only-decimal",
			args: args{d: "T2.5S"},
			want: &Duration{
				Seconds: 2.5,
			},
			wantErr: false,
		},
		{
			name: "full",
			args: args{d: "P3Y6M4DT12H30M5.5S"},
			want: &Duration{
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
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
	tests := []struct {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration := &Duration{
				Years:   tt.fields.Years,
				Months:  tt.fields.Months,
				Weeks:   tt.fields.Weeks,
				Days:    tt.fields.Days,
				Hours:   tt.fields.Hours,
				Minutes: tt.fields.Minutes,
				Seconds: tt.fields.Seconds,
			}
			if got := duration.ToTimeDuration(); got != tt.want {
				t.Errorf("ToTimeDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
