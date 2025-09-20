package breed

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type (
	Validator struct {
		apiBase    string
		httpClient *http.Client
	}

	CatAPIReponse []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
)

func NewValidator() Validator {
	return Validator{
		apiBase:    "https://api.thecatapi.com/v1",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// Returns true if the breed name matches one of the breeds from TheCatAPI (case-insensitive)
func (v Validator) IsValid(ctx context.Context, breedName string) (bool, error) {
	url := fmt.Sprintf("%s/breeds/search?q=%s", v.apiBase, breedName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create http request, err: %v", err)
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to do http request, err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("TheCatAPI responded with status %d", resp.StatusCode)
	}

	var result CatAPIReponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to decode response body, err: %v", err)
	}

	for _, b := range result {
		if equalIgnoreCase(b.Name, breedName) {
			return true, nil
		}
	}

	return false, nil
}

func equalIgnoreCase(a, b string) bool {
	return strings.EqualFold(a, b)
}
