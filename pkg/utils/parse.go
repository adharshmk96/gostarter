package utils

func ParseNullString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ParseNullInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}
