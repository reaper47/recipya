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
		case "afghankitchenrecipes":
			return scrapeAfghanKitchen(doc)
		case "ah":
			return scrapeAh(doc)
		case "argiro":
			return scrapeArgiro(doc)
		case "archanaskitchen":
			return scrapeArchanasKitchen(doc)
		default:
			return parseWebsite(doc)
		}
	case 'b':
		switch host {
		case "barefootcontessa":
			return scrapeBarefootcontessa(doc)
		case "bbcgoodfood":
			return scrapeBbcgoodfood(doc)
		case "bettycrocker":
			return scrapeBettyCrocker(doc)
		case "blueapron":
			return scrapeBlueapron(doc)
		case "brianlagerstrom":
			return scrapeBrianLagerstrom(doc)
		case "briceletbaklava":
			return scrapeBriceletbaklava(doc)
		case "bodybuilding":
			return scrapeBodybuilding(doc)
		case "bongeats":
			return scrapeBongeats(doc)
		default:
			return parseWebsite(doc)
		}
	case 'c':
		switch host {
		case "cdkitchen":
			return scrapeCdKitchen(doc)
		case "chefnini":
			return scrapeChefnini(doc)
		case "closetcooking":
			return scrapeClosetcooking(doc)
		case "cook-talk":
			return scrapeCooktalk(doc)
		case "coop":
			return scrapeCoop(doc)
		case "costco":
			return scrapeCostco(doc)
		default:
			return parseWebsite(doc)
		}
	case 'd':
		switch host {
		case "dr":
			return scrapeDk(doc)
		default:
			return parseWebsite(doc)
		}
	case 'e':
		switch host {
		case "eatwell101":
			return scrapeEatwell101(doc)
		case "expressen":
			return scrapeEspressen(doc)
		default:
			return parseWebsite(doc)
		}
	case 'f':
		switch host {
		case "farmhousedelivery":
			return scrapeFarmhousedelivery(doc)
		case "fitmencook":
			return scrapeFitMenCook(doc)
		case "foodnetwork":
			return scrapeFoodNetwork(doc)
		case "foodrepublic":
			return scrapeFoodRepublic(doc)
		case "forksoverknives":
			return scrapeForksOverKnives(doc)
		case "franzoesischkochen":
			return parseWebsite(doc)
		default:
			return parseWebsite(doc)
		}
	case 'g':
		switch host {
		case "gesund-aktiv":
			return scrapeGesundAktiv(doc)
		case "globo":
			return scrapeGlobo(doc)
		case "grandfrais":
			return scrapeGrandfrais(doc)
		case "greatbritishchefs":
			return scrapeGreatBritishChefs(doc)
		case "grimgrains":
			return scrapeGrimGrains(doc)
		case "grouprecipes":
			return scrapeGrouprecipes(doc)
		default:
			return parseWebsite(doc)
		}
	case 'h':
		switch host {
		case "halfbakedharvest":
			return scrapeHalfBakedHarvest(doc)
		case "heatherchristo":
			return scrapeHeatherChristo(doc)
		case "homechef":
			return scrapeHomechef(doc)
		default:
			return parseWebsite(doc)
		}
	case 'i':
		switch host {
		default:
			return parseWebsite(doc)
		}
	case 'j':
		switch host {
		case "juliegoodwin":
			return scrapeJuliegoodwin(doc)
		case "justbento":
			return scrapeJustbento(doc)
		default:
			return parseWebsite(doc)
		}
	case 'k':
		switch host {
		case "kennymcgovern":
			return scrapeKennyMcGovern(doc)
		case "kochbucher":
			return scrapeKochbucher(doc)
		case "kptncook":
			return scrapeKptncook(doc)
		case "kuchnia-domowa":
			return scrapeKuchniadomova(doc)
		case "kwestiasmaku":
			return scrapeKwestiasmaku(doc)
		default:
			return parseWebsite(doc)
		}
	case 'l':
		switch host {
		case "latelierderoxane":
			return scrapeLatelierderoxane(doc)
		case "lekkerensimpel":
			return scrapeLekkerenSimpel(doc)
		case "livingthegreenlife":
			return scrapeLivingTheGreenLife(doc)
		default:
			return parseWebsite(doc)
		}
	case 'm':
		switch host {
		case "maangchi":
			return scrapeMaangchi(doc)
		case "meljoulwan":
			return scrapeMeljoulwan(doc)
		case "mindmegette":
			return scrapeMindMegette(doc)
		/*case "monsieur-cuisine":
		return scrapeMonsieurCuisine(doc)*/
		case "moulinex":
			return scrapeMoulinex(doc)
		case "mundodereceitasbimby":
			return scrapeMundodereceitasbimby(doc)
		case "myplate":
			return scrapeMyPlate(doc)
		default:
			return parseWebsite(doc)
		}
	case 'n':
		switch host {
		case "ninjatestkitchen":
			return scrapeNinjatestkitchen(doc)
		default:
			return parseWebsite(doc)
		}
	case 'o':
		switch host {
		case "owen-han":
			return scrapeOwenhan(doc)
		default:
			return parseWebsite(doc)
		}
	case 'p':
		switch host {
		case "panelinha":
			return scrapePanelinha(doc)
		case "paninihappy":
			return scrapePaniniHappy(doc)
		case "ploetzblog":
			return scrapePloetzblog(doc)
		case "projectgezond":
			return scrapeProjectgezond(doc)
		case "przepisy":
			return scrapePrzepisy(doc)
		case "purelypope":
			return scrapePurelyPope(doc)
		default:
			return parseWebsite(doc)
		}
	case 'r':
		switch host {
		case "recettes":
			return scrapeRecettesDuQuebec(doc)
		case "recipecommunity":
			return scrapeRecipeCommunity(doc)
		case "reishunger":
			return scrapeReisHunger(doc)
		case "rezeptwelt":
			return scrapeRezeptwelt(doc)
		case "rosannapansino":
			return scrapeRosannapansino(doc)
		default:
			return parseWebsite(doc)
		}
	case 's':
		switch host {
		case "saboresajinomoto":
			return scrapeSaboresajinomoto(doc)
		case "sallys-blog":
			return scrapeSallysblog(doc)
		case "southerncastiron":
			return scrapeSoutherncastiron(doc)
		case "streetkitchen":
			return scrapeStreetKitchen(doc)
		case "sunset":
			return scrapeSunset(doc)
		default:
			return parseWebsite(doc)
		}
	case 't':
		switch host {
		case "tastykitchen":
			return scrapeTastyKitchen(doc)
		case "tesco":
			return scrapeTesco(doc)
		case "thecookingguy":
			return scrapeTheCookingGuy(doc)
		case "thehappyfoodie":
			return scrapeTheHappyFoodie(doc)
		default:
			return parseWebsite(doc)
		}
	case 'u':
		switch host {
		case "uitpaulineskeuken":
			return scrapeUitpaulineskeuken(doc)
		case "usapears":
			return scrapeUsapears(doc)
		default:
			return parseWebsite(doc)
		}
	case 'v':
		switch host {
		case "valdemarsro":
			return scrapeValdemarsro(doc)
		default:
			return parseWebsite(doc)
		}
	case 'w':
		switch host {
		case "waitrose":
			return scrapeWaitrose(doc)
		case "wholefoodsmarket":
			return scrapeWholefoodsmarket(doc)
		case "wikibooks":
			return scrapeWikiBooks(doc)
		case "woop":
			return scrapeWoop(doc)
		default:
			return parseWebsite(doc)
		}
	case 'y':
		switch host {
		case "ye-mek":
			return scrapeYemek(doc)
		}
	case 'z':
		switch host {
		case "zeit":
			return scrapeZeit(doc)
		default:
			return parseWebsite(doc)
		}
	default:
		switch host {
		case "15gram":
			return scrape15gram(doc)
		}
		return parseWebsite(doc)
	}
	return models.RecipeSchema{}, ErrNotImplemented
}

func parseWebsite(doc *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseLdJSON(doc)
	if err != nil {
		rs, err = parseGraph(doc)
		if err != nil {
			return models.RecipeSchema{}, ErrNotImplemented
		}
	}

	if rs.Yield.Value == 0 {
		rs.Yield.Value = 1
	}

	return rs, nil
}
