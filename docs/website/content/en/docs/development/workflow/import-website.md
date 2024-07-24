---
title: Support a Website
---

You first need to understand how the scraper works to support a website. 
Then, we will guide you through adding a website to the supported list with an example.

## The Scraper

Recipya developed its own recipe scraper, which resides in the [internal/scraper](https://github.com/reaper47/recipya/tree/main/internal/scraper) 
package. This scraper uses [goquery](https://github.com/PuerkitoBio/goquery) to extract information from web pages.
Its main file is [scraper.go](https://github.com/reaper47/recipya/blob/main/internal/scraper/scraper.go). You will find a single exposed function named `Scrape`, which
takes a URL and a [files service](https://github.com/reaper47/recipya/blob/main/internal/services/service.go#L102) as parameters. The files services is an interface with functions to
manipulate files in the OS. The use of an interface simplifies the process of mocking file operations during testing.

You can read how the function works, but essentially it involves fetching the web page using Go's HTTP client,
creating a `goquery` document from the response, extracting into a 
[models.RecipeSchema](https://github.com/reaper47/recipya/blob/main/internal/models/schema-recipe.go) struct, uploading 
the image to the server, and finally returning the recipe schema model. The image is compressed and resized.
Whether compression is too high remains is subject to evaluation.

## Workflow

Let's assume a user requests https://www.example.com/recipes/declicious-bbq-steak to be supported.
This section will help you understand how to add this website to the list of supported sites.

### Database

Initially, a SQLite migration file needs to be created using Goose to insert the desired website into the 
websites table. To do so, open a terminal and navigate to the root of the project. Then, generate
the migration file.

```bash
task new-migration name=support_website
```

The `support_website` is the name of the migration. It can be anything else. The command will create a new file of the 
form `timestamp_name_of_migration.sql` under `internal/services/migrations`. It will be embedded into the executable on build 
and will be executed when the user starts the server. 

The final step involves inserting the website into the database:

```sql {filename="internal/services/migrations/timestamp_support_website.sql"}
-- +goose Up
INSERT INTO websites (host, url) 
VALUES ('example.com', 'https://www.example.com/recipes/declicious-bbq-steak');

-- +goose Down
DELETE FROM websites
WHERE host IN ('example.com');
```

The host field could eventually be removed because we can determine it from Go using 
the [net/url](https://pkg.go.dev/net/url#URL.Hostname) package.

### Test

Setting up a test involves accessing the website and creating a test case within `internal/scraper/scraper_{letter}_test.go`.
In our case, open [internal/scraper/scraper_E_test.go](https://github.com/reaper47/recipya/blob/main/internal/scraper/scraper_E_test.go)
because `example` begins with `E`. The tests within the `testcases` slice are listed alphabetically so insert your `name: "example.com"`
test where appropriate. You can use an existing struct as a template.

Next, alternate between the recipe web page and the test to modify the `models.RecipeSchema` of 
the `want` field. You can proceed to writing code once the setup is done. 

Executing the test by clicking the green play button to the left it should confirm its failure.
If you notice the test returns a `models.RecipeSchema` that looks valid, then replace the empty schema
of the test with the one from the output and make the test go green. Otherwise, continue to the next section.

### The Go Code

The initial step is to include the `example.com` case within the list of supported websites. To achieve this, open 
[internal/scraper/websites.go](https://github.com/reaper47/recipya/blob/main/internal/scraper/websites.go). This file contains the `scrapeWebsite` function, which executes the relevant 
scrape function for the parsed HTML web page. Your task involves adding the host within the switch-case block. 
Therefore, add `case "example"` to the switch-case block
of [case 'e'](https://github.com/reaper47/recipya/blob/main/internal/scraper/websites.go#L64). 

Following this, the body of the case must be added by calling a custom HTML parser function.
Its naming convention is `scrape{Host}`. In your case, it would be `return scrapeExample(doc)`.
Then, create a new file named `example.go` under `internal/scraper` and add the 
`func scrapeKuchniadomova(root *goquery.Document) (models.RecipeSchema, error)` function. Please check any 
custom scraper file to understand how to write your own.

Congratulations! That is all you need to know to support a website. Feel free to open a PR once your scrape function is 
written and the tests pass.
