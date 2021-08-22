package common

import "strconv"

func AsInt(s string) (int, error) {
	code, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return code, nil
}
