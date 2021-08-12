from bs4 import BeautifulSoup
import requests


def main():
    page = requests.get("https://github.com/hhursev/recipe-scrapers")
    soup = BeautifulSoup(page.content, "html.parser")

    websites = (
        soup.find(id="user-content-scrapers-available-for")
        .parent.find_next_sibling("ul")
        .findAll("a")
    )
    websites = ["'" + website["href"] + "'" for website in websites]

    with open("./websites.txt", "w") as f:
        f.write("\n".join(websites))


if __name__ == "__main__":
    main()
