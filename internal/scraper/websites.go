package scraper

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

// ErrNotImplemented is the error used when the website is not supported.
var ErrNotImplemented = errors.New("domain is not implemented")

func scrapeWebsite(doc *goquery.Document, host string) (models.RecipeSchema, error) {
	switch []rune(host)[0] {
	case 'a':
		switch host {
		case "abril":
			return parseLdJSON(doc)
		case "acouplecooks":
			return parseGraph(doc)
		case "afghankitchenrecipes":
			return scrapeAfghanKitchen(doc)
		case "allrecipes":
			return parseLdJSON(doc)
		case "amazingribs":
			return parseGraph(doc)
		case "ambitiouskitchen": /**/
			return parseGraph(doc)
		case "archanaskitchen":
			return scrapeArchanasKitchen(doc)
		case "atelierdeschefs":
			return parseLdJSON(doc)
		case "averiecooks":
			return parseLdJSON(doc)
		}
	case 'b':
		switch host {
		case "bakingmischief":
			return parseGraph(doc)
		case "baking-sense":
			return parseGraph(doc)
		case "bbc":
			return parseLdJSON(doc)
		case "bbcgoodfood":
			return parseLdJSON(doc)
		case "bettycrocker":
			return scrapeBettyCrocker(doc)
		case "bigoven":
			return parseLdJSON(doc)
		case "bonappetit":
			return parseLdJSON(doc)
		case "bowlofdelicious":
			return parseGraph(doc)
		case "budgetbytes":
			return parseGraph(doc)
		}
	case 'c':
		switch host {
		case "cafedelites":
			return parseGraph(doc)
		case "castironketo":
			return parseGraph(doc)
		case "cdkitchen":
			return scrapeCdKitchen(doc)
		case "chefkoch":
			return parseLdJSON(doc)
		case "comidinhasdochef":
			return parseLdJSON(doc)
		case "cookeatshare":
			return parseLdJSON(doc)
		case "cookieandkate":
			return parseGraph(doc)
		case "copykat":
			return parseGraph(doc)
		case "countryliving":
			return parseLdJSON(doc)
		case "cuisineaz":
			return parseGraph(doc)
		case "cybercook":
			return parseLdJSON(doc)
		}
	case 'd':
		switch host {
		case "delish":
			return parseLdJSON(doc)
		case "ditchthecarbs":
			return parseGraph(doc)
		case "domesticate-me":
			return parseGraph(doc)
		case "dr":
			return scrapeDk(doc)
		}
	case 'e':
		switch host {
		case "eatingbirdfood":
			return parseGraph(doc)
		case "eatingwell":
			return parseLdJSON(doc)
		case "eatsmarter":
			return parseLdJSON(doc)
		case "epicurious":
			return parseLdJSON(doc)
		case "expressen":
			return scrapeEspressen(doc)
		}
	case 'f':
		switch host {
		case "fifteenspatulas":
			return parseGraph(doc)
		case "finedininglovers":
			return parseGraph(doc)
		case "fitmencook":
			return scrapeFitMenCook(doc)
		case "food":
			return parseLdJSON(doc)
		case "food52":
			return parseLdJSON(doc)
		case "foodandwine":
			return parseLdJSON(doc)
		case "foodrepublic":
			return scrapeFoodRepublic(doc)
		case "forksoverknives":
			return scrapeForksOverKnives(doc)
		case "franzoesischkochen":
			return scrapeFranzoesischKochen(doc)
		}
	case 'g':
		switch host {
		case "giallozafferano":
			return parseLdJSON(doc)
		case "gimmesomeoven":
			return parseGraph(doc)
		case "globo":
			return scrapeGlobo(doc)
		case "gonnawantseconds":
			return parseGraph(doc)
		case "greatbritishchefs":
			return scrapeGreatBritishChefs(doc)
		}
	case 'h':
		switch host {
		case "halfbakedharvest":
			return scrapeHalfBakedHarvest(doc)
		case "hassanchef":
			return parseLdJSON(doc)
		case "headbangerskitchen":
			return parseGraph(doc)
		case "hellofresh":
			return parseLdJSON(doc)
		case "hostthetoast":
			return parseGraph(doc)
		}
	case 'i':
		switch host {
		case "indianhealthyrecipes":
			return parseGraph(doc)
		case "innit":
			return parseLdJSON(doc)
		case "inspiralized":
			return parseLdJSON(doc)
		}
	case 'j':
		switch host {
		case "jamieoliver":
			return parseLdJSON(doc)
		case "jimcooksfoodgood":
			return parseLdJSON(doc)
		case "joyfoodsunshine":
			return parseGraph(doc)
		case "justataste":
			return parseGraph(doc)
		case "justonecookbook":
			return parseGraph(doc)
		}
	case 'k':
		switch host {
		case "kennymcgovern":
			return scrapeKennyMcGovern(doc)
		case "kingarthurbaking":
			return parseGraph(doc)
		case "kochbar":
			return parseLdJSON(doc)
		case "koket":
			return parseLdJSON(doc)
		case "kuchnia-domowa":
			return scrapeKuchniadomova(doc)
		case "kwestiasmaku":
			return scrapeKwestiasmaku(doc)
		}
	case 'l':
		switch host {
		case "lecremedelacrumb":
			return parseGraph(doc)
		case "lekkerensimpel":
			return scrapeLekkerenSimpel(doc)
		case "littlespicejar":
			return parseLdJSON(doc)
		case "livelytable":
			return parseLdJSON(doc)
		case "lovingitvegan":
			return parseGraph(doc)
		}
	case 'm':
		switch host {
		case "madensverden":
			return parseGraph(doc)
		case "marthastewart":
			return parseLdJSON(doc)
		case "matprat":
			return parseLdJSON(doc)
		case "melskitchencafe":
			return parseGraph(doc)
		case "mindmegette":
			return scrapeMindMegette(doc)
		case "minimalistbaker":
			return parseGraph(doc)
		case "misya":
			return parseLdJSON(doc)
		case "momsdish":
			return parseGraph(doc)
		case "momswithcrockpots":
			return parseGraph(doc)
		case "monsieur-cuisine":
			return scrapeMonsieurCuisine(doc)
		case "motherthyme":
			return parseLdJSON(doc)
		case "mybakingaddiction":
			return parseLdJSON(doc)
		case "mykitchen101", "mykitchen101en":
			return parseGraph(doc)
		case "myplate":
			return scrapeMyPlate(doc)
		case "myrecipes":
			return parseLdJSON(doc)
		}
	case 'n':
		switch host {
		case "nourishedbynutrition":
			return parseGraph(doc)
		case "nutritionbynathalie":
			return scrapeNutritionByNathalie(doc)
		case "nytimes":
			return parseLdJSON(doc)
		}
	case 'o':
		switch host {
		case "ohsheglows":
			return parseLdJSON(doc)
		case "onceuponachef":
			return parseLdJSON(doc)
		}
	case 'p':
		switch host {
		case "paleorunningmomma":
			return parseGraph(doc)
		case "panelinha":
			return scrapePanelinha(doc)
		case "paninihappy":
			return scrapePaniniHappy(doc)
		case "practicalselfreliance":
			return parseLdJSON(doc)
		case "primaledgehealth":
			return parseGraph(doc)
		case "przepisy":
			return scrapePrzepisy(doc)
		case "purelypope":
			return scrapePurelyPope(doc)
		case "purplecarrot":
			return parseLdJSON(doc)
		}
	case 'r':
		switch host {
		case "rachlmansfield":
			return parseGraph(doc)
		case "rainbowplantlife":
			return parseGraph(doc)
		case "realsimple":
			return parseLdJSON(doc)
		case "recettes":
			return scrapeRecettesDuQuebec(doc)
		case "recipetineats":
			return parseGraph(doc)
		case "redhousespice":
			return parseGraph(doc)
		case "reishunger":
			return scrapeReisHunger(doc)
		case "rezeptwelt":
			return scrapeRezeptwelt(doc)
		}
	case 's':
		switch host {
		case "sallysbakingaddiction":
			return parseGraph(doc)
		case "saveur":
			return parseLdJSON(doc)
		case "seriouseats":
			return parseLdJSON(doc)
		case "simplyquinoa":
			return parseGraph(doc)
		case "simplyrecipes":
			return parseLdJSON(doc)
		case "simplywhisked":
			return parseGraph(doc)
		case "skinnytaste":
			return parseGraph(doc)
		case "southernliving":
			return parseLdJSON(doc)
		case "spendwithpennies":
			return parseGraph(doc)
		case "springlane":
			return parseLdJSON(doc)
		case "steamykitchen":
			return parseLdJSON(doc)
		case "streetkitchen":
			return scrapeStreetKitchen(doc)
		case "sunbasket":
			return parseLdJSON(doc)
		case "sundpaabudget":
			return parseGraph(doc)
		case "sweetcsdesigns":
			return parseLdJSON(doc)
		case "sweetpeasandsaffron":
			return parseGraph(doc)
		}
	case 't':
		switch host {
		case "tasteofhome":
			return parseLdJSON(doc)
		case "tastesoflizzyt":
			return parseGraph(doc)
		case "tasty":
			return parseLdJSON(doc)
		case "tastykitchen":
			return scrapeTastyKitchen(doc)
		case "tesco":
			return scrapeTesco(doc)
		case "theclevercarrot":
			return parseGraph(doc)
		case "thehappyfoodie":
			return scrapeTheHappyFoodie(doc)
		case "thekitchenmagpie":
			return parseGraph(doc)
		case "thenutritiouskitchen":
			return parseGraph(doc)
		case "thepioneerwoman":
			return parseLdJSON(doc)
		case "thespruceeats":
			return parseLdJSON(doc)
		case "thevintagemixer":
			return parseLdJSON(doc)
		case "thewoksoflife":
			return parseGraph(doc)
		case "timesofindia":
			return parseLdJSON(doc)
		case "tine":
			return parseLdJSON(doc)
		case "tudogostoso":
			return parseLdJSON(doc)
		case "twopeasandtheirpod":
			return parseGraph(doc)
		}
	case 'v':
		switch host {
		case "valdemarsro":
			return scrapeValdemarsro(doc)
		case "vanillaandbean":
			return parseGraph(doc)
		case "vegolosi":
			return parseLdJSON(doc)
		case "vegrecipesofindia":
			return parseGraph(doc)
		}
	case 'w':
		switch host {
		case "watchwhatueat":
			return parseGraph(doc)
		case "whatsgabycooking":
			return parseGraph(doc)
		case "wikibooks":
			return scrapeWikiBooks(doc)
		case "woop":
			return scrapeWoop(doc)
		}
	case 'y':
		switch host {
		case "ye-mek":
			return scrapeYemek(doc)
		}
	case 'z':
		switch host {
		case "zenbelly":
			return parseLdJSON(doc)
		}
	default:
		switch host {
		case "101cookbooks":
			return parseLdJSON(doc)
		}
	}
	return models.RecipeSchema{}, ErrNotImplemented
}

func parseUnsupportedWebsite(doc *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseLdJSON(doc)
	if err != nil {
		rs, err = parseGraph(doc)
		if err != nil {
			return models.RecipeSchema{}, ErrNotImplemented
		}
	}
	return rs, nil
}
