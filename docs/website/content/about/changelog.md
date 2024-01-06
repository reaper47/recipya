---
title: Changelog
weight: 1
---

## v1.0.1 (release date TBD)

### Refresh Scraper

The scraper module has been refreshed to keep it up-to-date. 
 
The `mob.co.uk` website has been removed because the recipe's script tag in the scraped HTML contains invalid JSON.
The `json` and `yaml` libraries unfortunately cannot decode the JSON.

The `nutritionbynathalie.com` website has been removed because it cannot be scraped anymore.