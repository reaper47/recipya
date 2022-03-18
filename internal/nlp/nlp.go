package nlp

import (
	"log"
	"sync"
	"unicode"

	"github.com/jdkato/prose/v2"
)

func CapitalizeParagraphs(paragraphs []string) []string {
	output := make([]string, len(paragraphs))

	var wg sync.WaitGroup
	for i, p := range paragraphs {
		wg.Add(1)
		go func(i int, text string, buffer []string) {
			defer wg.Done()

			doc, err := prose.NewDocument(text)
			if err != nil {
				log.Printf("could not create a document from '%s': %s", text, err)
				return
			}

			for _, s := range doc.Sentences() {
				if len(s.Text) == 0 {
					continue
				}

				r := []rune(s.Text)
				if buffer[i] != "" {
					buffer[i] += " "
				}
				buffer[i] += string(append([]rune{unicode.ToUpper(r[0])}, r[1:]...))
			}
		}(i, p, output)
	}
	wg.Wait()

	return output
}
