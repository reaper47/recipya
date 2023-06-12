package regex_test

import (
	"github.com/reaper47/recipya/internal/utils/regex"
	"golang.org/x/exp/slices"
	"strings"
	"testing"
)

func TestRegex_Email(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		emails := []string{
			"james@bond.com",
			"hello@hello.ca",
			"slave@ukrainia.ua",
			"norway@rocks.no",
		}
		for _, email := range emails {
			t.Run("regex is valid", func(t *testing.T) {
				if !regex.Email.MatchString(email) {
					t.Fatal("got false when want true")
				}
			})
		}
	})

	t.Run("invalid", func(t *testing.T) {
		emails := []string{
			"xyzGmail.com",
			"@gmail.com",
			"email",
			"a@.com",
			".com@",
			"a@",
			"norway@rocks",
		}
		for _, email := range emails {
			t.Run("regex is invalid "+email, func(t *testing.T) {
				if regex.Email.MatchString(email) {
					t.Error("got true when want false")
				}
			})
		}
	})
}

func TestRegex_Quantity(t *testing.T) {
	quantities := []string{
		"1ml", "1mL", "15 ml", "16 mL",
		"1l", "1L", "15 l", "16 L",
		"1°c", "1°f", "15 °c", "16 °f",
		"1°C", "1°F", "15 °C", "16 °F",
	}
	for _, q := range quantities {
		t.Run("regex is valid "+q, func(t *testing.T) {
			if !regex.Quantity.MatchString(q) {
				t.Error("got false when want true for")
			}
		})
	}

	quantities = []string{
		"ml", "mL",
		"l", "L",
		"°c", "°f",
		"°C", "°F",
		"15 mX", "15mx",
	}
	for _, q := range quantities {
		t.Run("regex is invalid "+q, func(t *testing.T) {
			if regex.Quantity.MatchString(q) {
				t.Errorf("got true when want false for %q", q)
			}
		})
	}
}

func TestRegex_Anchor(t *testing.T) {
	t.Run("anchor is valid", func(t *testing.T) {
		a := `<a slot="guide-links-primary" href="https://www.youtube.com/about/press/" style="display: none;">`

		if !regex.Anchor.MatchString(a) {
			t.Error("got false when want true")
		}
	})

	xa := []string{
		`<aa slot="guide-links-primary" href="https://www.youtube.com/about/press/" style="display: none;">`,
		`<aslot="guide-links-primary" href="https://www.youtube.com/about/press/" style="display: none;">`,
		`<a slot="guide-links-primary" href="https://www.youtube.com/about/press/" style="display: none;"`,
	}
	for _, a := range xa {
		t.Run("anchor is invalid ", func(t *testing.T) {
			if regex.Anchor.MatchString(a) {
				t.Errorf("got true when want true for %q", a)
			}
		})
	}
}

func TestRegex_HourMinutes(t *testing.T) {
	xs := []string{
		"45:23",
		"45:00",
		"120:59",
	}
	for _, s := range xs {
		t.Run("regex is valid "+s, func(t *testing.T) {
			if !regex.HourMinutes.MatchString(s) {
				t.Error("got false but want true")
			}
		})
	}

	xs = []string{
		":00",
		"4500",
		"120:60",
		"120 43",
		"120-43",
		"120:",
		"-1:43",
		"10:-43",
		"10:80",
	}
	for _, s := range xs {
		t.Run("regex is invalid "+s, func(t *testing.T) {
			if regex.HourMinutes.MatchString(s) {
				t.Error("got true when want false")
			}
		})
	}
}

func TestRegex_Sentences(t *testing.T) {
	paragraph := "this is a test with 4.7. " +
		"it should be capitalized. " +
		"this is a test. " +
		"can you capitalize this? " +
		"Barbara had been waiting at the table for twenty minutes. " +
		"It had been twenty long and excruciating minutes. " +
		"David had promised that he would be on time today. " +
		"He never was, but he had promised this one time. " +
		"She had made him repeat the promise multiple times over the last week until she'd believed his promise. " +
		"Now she was paying the price."

	got := regex.Sentences.FindAllString(paragraph, -1)
	want := []string{
		"this is a test with 4.7. ",
		"it should be capitalized. ",
		"this is a test. ",
		"can you capitalize this? ",
		"Barbara had been waiting at the table for twenty minutes. ",
		"It had been twenty long and excruciating minutes. ",
		"David had promised that he would be on time today. ",
		"He never was, but he had promised this one time. ",
		"She had made him repeat the promise multiple times over the last week until she'd believed his promise. ",
		"Now she was paying the price.",
	}

	if !slices.EqualFunc(got, want, func(s string, s2 string) bool {
		return strings.TrimSpace(s) == strings.TrimSpace(s2)
	}) {
		t.Fatalf("got\n%v\nbut want\n%v", got, want)
	}
}

func TestRegex_Website(t *testing.T) {
	urls := []string{
		"http://subdomain.example.com",
		"http://www.example.com/path/to/page/subpage",
		"https://www.example.com/path/to/page/subpage?key1=value1&key2=value2",
		"http://username:password@subdomain.example.com:8080/path/to/page/subpage?key1=value1&key2=value2#fragment",
		"http://www.example.com/path/to/page/subpage/subsubpage?key1=value1&key2=value2#fragment",
		"https://subdomain.example.com/path/to/page/subpage/subsubpage?key1=value1&key2=value2#fragment",
		"https://realfood.tesco.com/recipes/salted-honey-and-rosemary-lamb-with-roasties-and-rainbow-carrots.html",
	}

	for _, url := range urls {
		if !regex.URL.MatchString(url) {
			t.Errorf("%v should be matched", url)
		}
	}
}
