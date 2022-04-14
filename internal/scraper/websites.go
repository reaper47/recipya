package scraper

import (
	"fmt"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeWebsite(doc *html.Node, host string) (models.RecipeSchema, error) {
	switch []rune(host)[0] {
	case 'a':
		switch host {
		case "abril":
			return scrapeLdJSON(doc)
		case "acouplecooks":
			return scrapeGraph(doc)
		case "afghankitchenrecipes":
			return scrapeAfghanKitchen(doc)
		case "allrecipes":
			return scrapeLdJSONs(doc)
		case "amazingribs":
			return scrapeGraph(doc)
		case "ambitiouskitchen":
			return scrapeGraph(doc)
		case "archanaskitchen":
			return scrapeArchanasKitchen(doc)
		case "atelierdeschefs":
			return findRecipeLdJSON(doc)
		case "averiecooks":
			return findRecipeLdJSON(doc)
		}
	case 'b':
		switch host {
		case "bakingmischief":
			return scrapeGraph(doc)
		case "baking-sense":
			return scrapeGraph(doc)
		case "bbc":
			return findRecipeLdJSON(doc)
		case "bbcgoodfood":
			return findRecipeLdJSON(doc)
		case "bettycrocker":
			return scrapeBettyCrocker(doc)
		case "bigoven":
			return scrapeLdJSON(doc)
		case "bonappetit":
			return scrapeLdJSON(doc)
		case "bowlofdelicious":
			return scrapeGraph(doc)
		case "budgetbytes":
			return scrapeGraph(doc)
		}
	case 'c':
		switch host {
		case "castironketo":
			return scrapeGraph(doc)
		case "cdkitchen":
			return scrapeCdKitchen(doc)
		case "chefkoch":
			return findRecipeLdJSON(doc)
		case "comidinhasdochef":
			return findRecipeLdJSON(doc)
		case "cookeatshare":
			return scrapeLdJSON(doc)
		case "cookieandkate":
			return scrapeGraph(doc)
		case "cookinglight":
			return scrapeLdJSONs(doc)
		case "cookstr":
			return findRecipeLdJSON(doc)
		case "copykat":
			return scrapeGraph(doc)
		case "countryliving":
			return scrapeLdJSON(doc)
		case "cuisineaz":
			return scrapeLdJSONs(doc)
		case "cybercook":
			return scrapeLdJSONs(doc)
		}
	case 'd':
		switch host {
		case "delish":
			return scrapeLdJSON(doc)
		case "ditchthecarbs":
			return scrapeGraph(doc)
		case "domesticate-me":
			return scrapeGraph(doc)
		case "downshiftology":
			return scrapeGraph(doc)
		case "dr":
			return scrapeDk(doc)
		}
	case 'e':
		switch host {
		case "eatingbirdfood":
			return scrapeGraph(doc)
		case "eatingwell":
			return scrapeLdJSONs(doc)
		case "eatsmarter":
			return scrapeLdJSON(doc)
		case "eatwhattonight":
			return findRecipeLdJSON(doc)
		case "epicurious":
			return scrapeLdJSON(doc)
		case "expressen":
			return scrapeEspressen(doc)
		}
	case 'f':
		switch host {
		case "fifteenspatulas":
			return scrapeGraph(doc)
		case "finedininglovers":
			return scrapeGraph(doc)
		case "fitmencook":
			return scrapeFitMenCook(doc)
		case "food":
			return scrapeLdJSON(doc)
		case "food52":
			return scrapeLdJSON(doc)
		case "foodandwine":
			return scrapeLdJSONs(doc)
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
			return findRecipeLdJSON(doc)
		case "gimmesomeoven":
			return scrapeGraph(doc)
		case "globo":
			return scrapeGlobo(doc)
		case "gonnawantseconds":
			return scrapeGraph(doc)
		case "greatbritishchefs":
			return scrapeGreatBritishChefs(doc)
		}
	case 'h':
		switch host {
		case "halfbakedharvest":
			return scrapeGraph(doc)
		case "hassanchef":
			return scrapeLdJSON(doc)
		case "headbangerskitchen":
			return findRecipeLdJSON(doc)
		case "hellofresh":
			return scrapeLdJSON(doc)
		case "homechef":
			return scrapeHomeChef(doc)
		case "hostthetoast":
			return scrapeGraph(doc)
		}
	case 'i':
		switch host {
		case "indianhealthyrecipes":
			return scrapeGraph(doc)
		case "innit":
			return scrapeLdJSON(doc)
		case "inspiralized":
			return scrapeGraph(doc)
		}
	case 'j':
		switch host {
		case "jamieoliver":
			return scrapeLdJSON(doc)
		case "jimcooksfoodgood":
			return findRecipeLdJSON(doc)
		case "joyfoodsunshine":
			return scrapeGraph(doc)
		case "justataste":
			return scrapeGraph(doc)
		case "justonecookbook":
			return scrapeGraph(doc)
		}
	case 'k':
		switch host {
		case "kennymcgovern":
			return scrapeGraph(doc)
		case "kingarthurbaking":
			rs, err := scrapeGraph(doc)
			rs.AtContext = "http://schema.org/"
			return rs, err
		case "kochbar":
			return findRecipeLdJSON(doc)
		case "koket":
			return scrapeLdJSONs(doc)
		case "kuchnia-domowa":
			return scrapeKuchniadomova(doc)
		case "kwestiasmaku":
			return scrapeKwestiasmaku(doc)
		}
	case 'l':
		switch host {
		case "lecremedelacrumb":
			return scrapeGraph(doc)
		case "lekkerensimpel":
			return scrapeLekkerenSimpel(doc)
		case "littlespicejar":
			return findRecipeLdJSON(doc)
		case "livelytable":
			return scrapeGraph(doc)
		case "lovingitvegan":
			return scrapeGraph(doc)
		}
	case 'm':
		switch host {
		case "madensverden":
			return scrapeGraph(doc)
		case "marthastewart":
			return scrapeLdJSONs(doc)
		case "matprat":
			return scrapeLdJSON(doc)
		case "melskitchencafe":
			return scrapeGraph(doc)
		case "mindmegette":
			return scrapeMindMegette(doc)
		case "minimalistbaker":
			return scrapeGraph(doc)
		case "misya":
			return scrapeLdJSON(doc)
		case "momswithcrockpots":
			return scrapeGraph(doc)
		case "monsieur-cuisine":
			return scrapeMonsieurCuisine(doc)
		case "motherthyme":
			return scrapeGraph(doc)
		case "mybakingaddiction":
			return findRecipeLdJSON(doc)
		case "mykitchen101":
			return scrapeMyKitchen101(doc)
		case "mykitchen101en":
			return scrapeMyKitchen101(doc)
		case "myplate":
			return scrapeMyPlate(doc)
		case "myrecipes":
			return scrapeLdJSONs(doc)
		}
	case 'n':
		switch host {
		case "nourishedbynutrition":
			return scrapeNourishedByNutrition(doc)
		case "nutritionbynathalie":
			return scrapeNutritionbynathalie(doc)
		case "nytimes":
			return scrapeLdJSON(doc)
		}
	case 'o':
		switch host {
		case "ohsheglows":
			return findRecipeLdJSON(doc)
		case "onceuponachef":
			return findRecipeLdJSON(doc)
		}
	case 'p':
		switch host {
		case "paleorunningmomma":
			return scrapeGraph(doc)
		case "panelinha":
			return scrapePanelinha(doc)
		case "paninihappy":
			return scrapePaniniHappy(doc)
		case "practicalselfreliance":
			return scrapeLdJSON(doc)
		case "primaledgehealth":
			return scrapeGraph(doc)
		case "przepisy":
			return scrapePrzepisy(doc)
		case "purelypope":
			return scrapePurelyPope(doc)
		case "purplecarrot":
			return scrapeLdJSON(doc)
		}
	case 'r':
		switch host {
		case "rachlmansfield":
			return scrapeGraph(doc)
		case "rainbowplantlife":
			return scrapeGraph(doc)
		case "realsimple":
			return scrapeLdJSONs(doc)
		case "recipetineats":
			return scrapeGraph(doc)
		case "redhousespice":
			return scrapeGraph(doc)
		case "reishunger":
			return scrapeReisHunger(doc)
		case "rezeptwelt":
			return scrapeRezeptwelt(doc)
		}
	case 's':
		switch host {
		case "sallysbakingaddiction":
			return scrapeGraph(doc)
		case "sallys-blog":
			return scrapeGraph(doc)
		case "saveur":
			return scrapeSaveur(doc)
		case "seriouseats":
			return scrapeLdJSON(doc)
		case "simplyquinoa":
			return scrapeGraph(doc)
		case "simplyrecipes":
			return scrapeLdJSON(doc)
		case "simplywhisked":
			return scrapeGraph(doc)
		case "skinnytaste":
			return scrapeLdJSON(doc)
		case "southernliving":
			return scrapeSouthernLiving(doc)
		case "spendwithpennies":
			return scrapeGraph(doc)
		case "springlane":
			return findRecipeLdJSON(doc)
		case "steamykitchen":
			return scrapeGraph(doc)
		case "streetkitchen":
			return scrapeStreetKitchen(doc)
		case "sunbasket":
			return scrapeLdJSON(doc)
		case "sundpaabudget":
			return scrapeGraph(doc)
		case "sweetcsdesigns":
			return findRecipeLdJSON(doc)
		case "sweetpeasandsaffron":
			return scrapeGraph(doc)
		}
	case 't':
		switch host {
		case "tasteofhome":
			return findRecipeLdJSON(doc)
		case "tastesoflizzyt":
			return scrapeGraph(doc)
		case "tasty":
			return scrapeLdJSON(doc)
		case "tastykitchen":
			return scrapeTastyKitchen(doc)
		case "tesco":
			return scrapeTesco(doc)
		case "theclevercarrot":
			return scrapeGraph(doc)
		case "thehappyfoodie":
			return scrapeTheHappyFoodie(doc)
		case "thekitchenmagpie":
			return scrapeGraph(doc)
		case "thenutritiouskitchen":
			return scrapeGraph(doc)
		case "thepioneerwoman":
			return scrapeLdJSON(doc)
		case "thespruceeats":
			return scrapeLdJSON(doc)
		case "thevintagemixer":
			return scrapeLdJSON(doc)
		case "thewoksoflife":
			return scrapeGraph(doc)
		case "timesofindia":
			return findRecipeLdJSON(doc)
		case "tine":
			return scrapeLdJSON(doc)
		case "tudogostoso":
			return scrapeLdJSON(doc)
		case "twopeasandtheirpod":
			return scrapeGraph(doc)
		}
	case 'v':
		switch host {
		case "valdemarsro":
			return scrapeValdemarsro(doc)
		case "vanillaandbean":
			return scrapeLdJSON(doc)
		case "vegolosi":
			return findRecipeLdJSON(doc)
		case "vegrecipesofindia":
			return scrapeGraph(doc)
		}
	case 'w':
		switch host {
		case "watchwhatueat":
			return scrapeGraph(doc)
		case "whatsgabycooking":
			return scrapeGraph(doc)
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
			return findRecipeLdJSON(doc)
		}
	default:
		switch host {
		case "101cookbooks":
			root := getElement(doc, "class", "wprm-recipe-container")
			return scrapeLdJSON(root)
		}
	}
	return models.RecipeSchema{}, fmt.Errorf("domain '%s' is not implemented", host)
}
