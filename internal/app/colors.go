package app

func greenText(s string) string {
	return "\033[32m" + s + "\033[0m"
}

func redText(s string) string {
	return "\033[31m" + s + "\033[0m"
}
