package util

func HasHave(count int) string {
	if count <= 1 {
		return "has"
	} else {
		return "have"
	}
}

func SingularPlural(count int) string {
	if count <= 1 {
		return ""
	} else {
		return "s"
	}
}
