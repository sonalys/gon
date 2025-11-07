package goncoder

func isInteger(s string) bool {
	if s == "" {
		return false
	}
	// Allow leading + or -
	if s[0] == '+' || s[0] == '-' {
		s = s[1:]
	}
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func isFloat(s string) bool {
	if s == "" {
		return false
	}
	// Allow leading + or -
	if s[0] == '+' || s[0] == '-' {
		s = s[1:]
	}
	if s == "" {
		return false
	}

	dotSeen := false
	digitSeen := false
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9':
			digitSeen = true
		case r == '.':
			if dotSeen {
				return false
			}
			dotSeen = true
		default:
			return false
		}
	}
	// Must contain at least one digit
	return digitSeen && dotSeen
}
