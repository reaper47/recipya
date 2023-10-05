package extensions_test

import (
	"github.com/reaper47/recipya/internal/utils/extensions"
	"math"
	"testing"
)

func BenchmarkFloatToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = extensions.FloatToString(3.140000, "%f")
	}
}

func BenchmarkSumStringToFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = extensions.ScaleString("4 3/4 cups of tea", 2.5)
	}
}

func BenchmarkSumString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = extensions.SumString("4 3/4 3 64 34 65 1/2 cups of tea 1/2 3/4")
	}
}

func TestFloatToString(t *testing.T) {
	testcases := []struct {
		name   string
		in     float64
		format string
		want   string
	}{
		{
			name:   "plain float",
			in:     3,
			format: "%f",
			want:   "3",
		},
		{
			name:   "decimal with trailing zeroes",
			in:     3.140000,
			format: "%.2f",
			want:   "3.14",
		},
		{
			name:   "no trailing zeroes",
			in:     3.14159,
			format: "%f",
			want:   "3.14159",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := extensions.FloatToString(tc.in, tc.format)
			assertStringsEqual(t, got, tc.want)
		})
	}
}

func TestScaleString(t *testing.T) {
	testcases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "no numbers",
			in:   "watermelons",
			want: "watermelons",
		},
		{
			name: "fraction",
			in:   "1/2 watermelon",
			want: "1 watermelon",
		},
		{
			name: "another fraction",
			in:   "2/3 tbsp honey",
			want: "1.333333 tbsp honey",
		},
		{
			name: "number with fraction",
			in:   "4 3/4 cups of tea",
			want: "9.5 cups of tea",
		},
		{
			name: "number with fraction no space",
			in:   "41/3 apples",
			want: "27.333333 apples",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := extensions.ScaleString(tc.in, 2)
			assertStringsEqual(t, got, tc.want)
		})
	}
}

func TestSumString(t *testing.T) {
	testcases := []struct {
		name string
		in   string
		want float64
	}{
		{
			name: "many integers",
			in:   "12 5 7 9",
			want: 33.,
		},
		{
			name: "many fractions",
			in:   "1/2 3/4 1/2 4/5",
			want: 2.55,
		},
		{
			name: "many integers sprinkled with fractions",
			in:   "1/2 2 1/2 2 1/2",
			want: 5.5,
		},
		{
			name: "only text",
			in:   "I love apples",
			want: 0.,
		},
		{
			name: "numbers then text",
			in:   "apples 4 1/2",
			want: 4.5,
		},
		{
			name: "numbers then text then numbers",
			in:   "1 1/2 apples 1/2",
			want: 1.5,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := extensions.SumString(tc.in)
			assertFloatsEqual(t, got, tc.want, 1e-3)
		})
	}
}

func assertFloatsEqual(t testing.TB, got, want, threshold float64) {
	t.Helper()
	if math.Abs(got-want) > threshold {
		t.Fatalf("got %g but want %g", got, want)
	}
}

func assertStringsEqual(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("got %q but want %q", got, want)
	}
}
