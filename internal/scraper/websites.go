package scraper

import (
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
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
		case "aldi":
			return scrapeAldi(doc)
		case "alittlebityummy":
			return scrapeALittleBitYummy(doc)
		case "all-clad":
			return scrapeAllClad(doc)
		case "angielaeats":
			return scrapeAngielaEats(doc)
		case "aniagotuje":
			return scrapeAniagotuje(doc)
		case "antilliaans-eten":
			return scrapeAntilliaansEten(doc)
		case "argiro":
			return scrapeArgiro(doc)
		case "archanaskitchen":
			return scrapeArchanasKitchen(doc)
		default:
			return parseWebsite(doc)
		}
	case 'b':
		switch host {
		case "bakels":
			return scrapeBritishBakels(doc)
		case "barefootcontessa":
			return scrapeBarefootcontessa(doc)
		case "bbcgoodfood":
			return scrapeBbcgoodfood(doc)
		case "bettybossi":
			return scrapeBettybossi(doc)
		case "bettycrocker":
			return scrapeBettyCrocker(doc)
		case "bingingwithbabish":
			return scrapeBingingWithBabish(doc)
		case "blueapron":
			return scrapeBlueapron(doc)
		case "bottomlessgreens":
			return scrapeBottomLessGreens(doc)
		case "brianlagerstrom":
			return scrapeBrianLagerstrom(doc)
		case "briceletbaklava":
			return scrapeBriceletbaklava(doc)
		case "britishbakels":
			return scrapeBritishBakels(doc)
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
		case "cestmafournee":
			return scrapeCestMaFournee(doc)
		case "chefnini":
			return scrapeChefnini(doc)
		case "chetnamakan":
			return scrapeChetnamakan(doc)
		case "chinesecookingdemystified":
			return scrapeChinesecookingdemystified(doc)
		case "chuckycruz":
			return scrapeChuckycruz(doc)
		case "closetcooking":
			return scrapeClosetcooking(doc)
		case "cook-talk":
			return scrapeCooktalk(doc)
		case "coop":
			return scrapeCoop(doc)
		case "costco":
			return scrapeCostco(doc)
		case "culy":
			return scrapeCuly(doc)
		default:
			return parseWebsite(doc)
		}
	case 'd':
		switch host {
		case "dherbs":
			return scrapeDherbs(doc)
		case "dish":
			return scrapeDish(doc)
		case "donnahay":
			return scrapeDonnaHay(doc)
		case "dr":
			return scrapeDk(doc)
		case "drinkoteket":
			return scrapeDrinkoteket(doc)
		default:
			return parseWebsite(doc)
		}
	case 'e':
		switch host {
		case "eatwell101":
			return scrapeEatwell101(doc)
		case "epicurious":
			return scrapeEpicurious(doc)
		case "etenvaneefke":
			return scrapEetenvaneefke(doc)
		case "expressen":
			return scrapeEspressen(doc)
		default:
			return parseWebsite(doc)
		}
	case 'f':
		switch host {
		case "farmhousedelivery":
			return scrapeFarmhousedelivery(doc)
		case "felix":
			return scrapeFelix(doc)
		case "finedininglovers":
			return scrapeFineDiningLovers(doc)
		case "fitmencook":
			return scrapeFitMenCook(doc)
		case "foodnetwork":
			return scrapeFoodNetwork(doc)
		case "foodrepublic":
			return scrapeFoodRepublic(doc)
		case "forksoverknives":
			return scrapeForksOverKnives(doc)
		case "francescakookt":
			return scrapeFrancescakookt(doc)
		case "franzoesischkochen":
			return parseWebsite(doc)
		default:
			return parseWebsite(doc)
		}
	case 'g':
		switch host {
		case "gazoakleychef":
			return scrapeGazoakleychef(doc)
		case "gesund-aktiv":
			return scrapeGesundAktiv(doc)
		case "giallozafferano":
			return scrapeGiallozafferano(doc)
		case "globo":
			return scrapeGlobo(doc)
		case "glutenfreetables":
			return scrapeGlutenFreeTables(doc)
		case "goodeatings":
			return scrapeGoodEatings(doc)
		case "grandfrais":
			return scrapeGrandfrais(doc)
		case "greatbritishchefs":
			return scrapeGreatBritishChefs(doc)
		case "grimgrains":
			return scrapeGrimGrains(doc)
		case "grouprecipes":
			return scrapeGrouprecipes(doc)
		case "gutekueche":
			return scrapeGutekueche(doc)
		default:
			return parseWebsite(doc)
		}
	case 'h':
		switch host {
		case "halfbakedharvest":
			return scrapeHalfBakedHarvest(doc)
		case "heatherchristo":
			return scrapeHeatherChristo(doc)
		case "homebrewanswers":
			return scrapeHomebrewAnswers(doc)
		case "homechef":
			return scrapeHomechef(doc)
		default:
			return parseWebsite(doc)
		}
	case 'i':
		switch host {
		case "instantpot":
			return scrapeInstantPot(doc)
		default:
			return parseWebsite(doc)
		}
	case 'j':
		switch host {
		case "jamieoliver":
			return scrapeJamieOliver(doc)
		case "jaimyskitchen":
			return scrapeJaimysKitchen(doc)
		case "juliegoodwin":
			return scrapeJuliegoodwin(doc)
		case "justbento":
			return scrapeJustbento(doc)
		default:
			return parseWebsite(doc)
		}
	case 'k':
		switch host {
		case "keepinitkind":
			return scrapeKeepinItKind(doc)
		case "kennymcgovern":
			return scrapeKennyMcGovern(doc)
		case "kingarthurbaking":
			return scrapeKingArthurBaking(doc)
		case "kitchenaid":
			return scrapeKitchenaid(doc)
		case "kochbucher":
			return scrapeKochbucher(doc)
		case "kookjij":
			return scrapeKookjij(doc)
		case "kptncook":
			return scrapeKptncook(doc)
		case "kuchnia-domowa":
			return scrapeKuchniadomova(doc)
		case "kuchynalidla":
			return scrapeKuchynalidla(doc)
		case "kwestiasmaku":
			return scrapeKwestiasmaku(doc)
		default:
			return parseWebsite(doc)
		}
	case 'l':
		switch host {
		case "lahbco":
			return scrapeLahbco(doc)
		case "latelierderoxane":
			return scrapeLatelierderoxane(doc)
		case "lekkerensimpel":
			return scrapeLekkerenSimpel(doc)
		case "lidl":
			return scrapeLidl(doc)
		case "lidl-kochen":
			return scrapeLidlKochen(doc)
		case "livingthegreenlife":
			return scrapeLivingTheGreenLife(doc)
		case "lithuanianintheusa":
			return scrapeLithuanianInTheUSA(doc)
		case "loveandlemons":
			return scrapeLoveAndLemons(doc)
		default:
			return parseWebsite(doc)
		}
	case 'm':
		switch host {
		case "maangchi":
			return scrapeMaangchi(doc)
		case "meljoulwan":
			return scrapeMeljoulwan(doc)
		case "mexicanmademeatless":
			return scrapeMexicanMadeMeatless(doc)
		case "mindmegette":
			return scrapeMindMegette(doc)
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
		case "nigella":
			return scrapeNigella(doc)
		case "ninjatestkitchen":
			return scrapeNinjatestkitchen(doc)
		default:
			return parseWebsite(doc)
		}
	case 'o':
		switch host {
		case "okokorecepten":
			return scrapeOkokorecepten(doc)
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
		case "pickuplimes":
			return scrapePickupLimes(doc)
		case "plentyvegan":
			return scrapePlentyVegan(doc)
		case "ploetzblog":
			return scrapePloetzblog(doc)
		case "popsugar":
			return scrapePopsugar(doc)
		case "potatorolls":
			return scrapePotatoRolls(doc)
		case "projectgezond":
			return scrapeProjectgezond(doc)
		case "przepisy":
			return scrapePrzepisy(doc)
		case "purelypope":
			return scrapePurelyPope(doc)
		case "purewow":
			return scrapePureWow(doc)
		case "puurgezond":
			return scrapePuurgezond(doc)
		default:
			return parseWebsite(doc)
		}
	case 'q':
		switch host {
		case "quitoque":
			return scrapeQuitoque(doc)
		default:
			return parseWebsite(doc)
		}
	case 'r':
		switch host {
		case "radiofrance":
			return scrapeRadioFrance(doc)
		case "recettes":
			return scrapeRecettesDuQuebec(doc)
		case "recipecommunity":
			return scrapeRecipeCommunity(doc)
		case "reddit":
			return scrapeReddit(doc)
		case "reishunger":
			return scrapeReisHunger(doc)
		case "rezeptwelt":
			return scrapeRezeptwelt(doc)
		case "robinasbell":
			return scrapeRobinasBell(doc)
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
		case "seriouseats":
			return scrapeSeriousEats(doc)
		case "smittenkitchen":
			return scrapeSmittenKitchen(doc)
		case "southerncastiron":
			return scrapeSoutherncastiron(doc)
		case "spiceboxtravels":
			return scrapeSpiceBoxTravels(doc)
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
		case "thatvegandad":
			return scrapeThatVeganDad(doc)
		case "thecookingguy":
			return scrapeTheCookingGuy(doc)
		case "thefoodflamingo":
			return scrapeTheFoodFlamingo(doc)
		case "theguccha":
			return scrapeTheGuccha(doc)
		case "thehappyfoodie":
			return scrapeTheHappyFoodie(doc)
		case "theheartysoul":
			return scrapeTheHeartySoul(doc)
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
		case "vegan-pratique":
			return scrapeVeganPratique(doc)
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
		case "yumelise":
			return scrapeYumelise(doc)
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
		return models.RecipeSchema{}, ErrNotImplemented
	}

	if rs.Yield == nil {
		rs.Yield = &models.Yield{Value: 1}
	} else if rs.Yield.Value == 0 {
		rs.Yield.Value = 1
	}

	return rs, nil
}

func parseLdJSON(root *goquery.Document) (models.RecipeSchema, error) {
	for _, node := range root.Find("script[type='application/ld+json']").Nodes {
		if node.FirstChild == nil {
			continue
		}

		var rs = models.NewRecipeSchema()
		err := json.Unmarshal([]byte(strings.ReplaceAll(node.FirstChild.Data, "\n", "")), &rs)

		var found bool
		if len(rs.AtGraph) > 0 {
			for _, schema := range rs.AtGraph {
				if schema.AtType.Value == "Recipe" {
					rs = *schema
					found = true
					break
				}
			}
		}

		if !found && err != nil || rs.Equal(models.NewRecipeSchema()) {
			var xrs []models.RecipeSchema
			err = json.Unmarshal([]byte(node.FirstChild.Data), &xrs)
			if err != nil {
				continue
			}

			for _, rs = range xrs {
				if rs.AtType != nil && rs.AtType.Value == "Recipe" {
					return rs, nil
				}
			}
			continue
		}

		if rs.AtType.Value != "Recipe" {
			continue
		}
		return rs, nil
	}
	return models.RecipeSchema{}, ErrNotImplemented
}
