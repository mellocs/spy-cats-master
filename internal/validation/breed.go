package validation

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const catApiURL = "https://api.thecatapi.com/v1/breeds"

func ValidateBreed(breed string) (bool, error) {
	res, err := http.Get(catApiURL)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("validation error: %s", res.Status)
	}

	var breeds []struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(res.Body).Decode(&breeds); err != nil {
		return false, err
	}

	for _, b := range breeds {
		if b.Name == breed {
			return true, nil
		}
	}

	return false, nil
}
