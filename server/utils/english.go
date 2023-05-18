package utils

func HasHave(count int) string {
	if count == 1 || count == -1 {
		return "has"
	}

	return "have"
}

func SingularPlural(count int) string {
	if count >= -1 && count <= 1 {
		return ""
	}

	return "s"
}
