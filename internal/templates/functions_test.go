package templates

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTemplatesFunctions(t *testing.T) {
	t.Run("test dec", func(t *testing.T) {
		f, ok := fm["dec"]
		if !ok {
			t.Fatal("dec function must be present in the FuncMap")
		}

		actual := f.(func(i int) int)(5)

		if actual != 4 {
			t.Fatalf("dec: wanted 4 but got %d", actual)
		}
	})

	t.Run("test durationToInput", func(t *testing.T) {
		name := "durationToInput"
		f, ok := fm[name]
		if !ok {
			t.Fatalf("%s function must be present in the FuncMap", name)
		}

		d, _ := time.ParseDuration("2h45m")
		actual := f.(func(d time.Duration) string)(d)

		expected := "2:45:00"
		if actual != expected {
			t.Fatalf("%s: wanted %s but got %s", name, expected, actual)
		}
	})

	testcases1 := []struct {
		name string
		in   bool
		want string
	}{
		{name: "is datetime", in: true, want: "PT2H45M"},
		{name: "is not datetime", in: false, want: "2h45"},
	}
	for _, tc := range testcases1 {
		t.Run("test fmtDuration "+tc.name, func(t *testing.T) {
			name := "fmtDuration"
			f, ok := fm[name]
			if !ok {
				t.Fatalf("%s function must be present in the FuncMap", name)
			}

			d, _ := time.ParseDuration("2h45m")
			actual := f.(func(d time.Duration, isDatetime bool) string)(d, tc.in)

			expected := tc.want
			if actual != expected {
				t.Fatalf("%s: wanted %s but got %s", name, expected, actual)
			}
		})
	}

	t.Run("test inc", func(t *testing.T) {
		name := "inc"
		f, ok := fm[name]
		if !ok {
			t.Fatalf("%s function must be present in the FuncMap", name)
		}

		actual := f.(func(i int) int)(5)

		if actual != 6 {
			t.Fatalf("%s: wanted 6 but got %d", name, actual)
		}
	})

	testcases2 := []struct {
		name string
		in   string
		want bool
	}{
		{name: "is true", in: "https://www.google.com", want: true},
		{name: "is false", in: "google.com", want: false},
	}
	for _, tc := range testcases2 {
		t.Run("test isUrl "+tc.name, func(t *testing.T) {
			name := "isUrl"
			f, ok := fm[name]
			if !ok {
				t.Fatalf("%s function must be present in the FuncMap", name)
			}

			actual := f.(func(s string) bool)(tc.in)

			expected := tc.want
			if actual != expected {
				t.Fatalf("%s: wanted %v but got %v", name, expected, actual)
			}
		})
	}

	testcases3 := []struct {
		name string
		in   uuid.UUID
		want bool
	}{
		{name: "is invalid", in: uuid.UUID{}, want: false},
		{name: "is valid", in: uuid.New(), want: true},
	}
	for _, tc := range testcases3 {
		t.Run("test isUuidValid "+tc.name, func(t *testing.T) {
			name := "isUuidValid"
			f, ok := fm[name]
			if !ok {
				t.Fatalf("%s function must be present in the FuncMap", name)
			}

			actual := f.(func(u uuid.UUID) bool)(tc.in)

			expected := tc.want
			if actual != expected {
				t.Fatalf("%s: wanted %v but got %v", name, expected, actual)
			}
		})
	}

	testcases4 := []struct {
		name string
		in   []int
		want int
	}{
		{name: "multiple numbers", in: []int{3, 4, 3}, want: 36},
		{name: "nothing", in: []int{}, want: 1},
	}
	for _, tc := range testcases4 {
		t.Run("test mul "+tc.name, func(t *testing.T) {
			name := "mul"
			f, ok := fm[name]
			if !ok {
				t.Fatalf("%s function must be present in the FuncMap", name)
			}

			actual := f.(func(n ...int) int)(tc.in...)

			expected := tc.want
			if actual != expected {
				t.Fatalf("%s: wanted %v but got %v", name, expected, actual)
			}
		})
	}

	testcases5 := []struct {
		name     string
		in       string
		endIndex int
		want     string
	}{
		{name: "empty string", in: "", endIndex: 4, want: ""},
		{name: "part  string", in: "i am a little dog", endIndex: 4, want: "i am..."},
		{name: "all string", in: "i am a little dog", endIndex: -1, want: "i am a little dog"},
	}
	for _, tc := range testcases5 {
		t.Run("test substring "+tc.name, func(t *testing.T) {
			name := "substring"
			f, ok := fm[name]
			if !ok {
				t.Fatalf("%s function must be present in the FuncMap", name)
			}

			actual := f.(func(s string, endIndex int) string)(tc.in, tc.endIndex)

			expected := tc.want
			if actual != expected {
				t.Fatalf("%s: wanted %v but got %v", name, expected, actual)
			}
		})
	}
}
