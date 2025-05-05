package api

import (
	"effectiveMobile/internal"
	"effectiveMobile/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"math/rand"
)

type External struct {
	r *resty.Client
}

func NewExternal() *External {
	r := resty.New()
	return &External{r}
}

func (e *External) GetAge(name string) (int, error) {
	resp, err := e.r.R().SetQueryParam("name", name).
		Get("https://api.agify.io")
	if err != nil {
		return 0, fmt.Errorf("error getting age: %w", err)
	}

	if resp.StatusCode() != 200 {
		return 0, fmt.Errorf("error getting age: %s", resp.Status())
	}

	var data struct {
		Age int `json:"age"`
	}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return 0, fmt.Errorf("error getting age: %w", err)
	}

	utils.DebugLog(name, "% Age received", data.Age)

	return data.Age, nil
}

func (e *External) GetGender(name string) (string, error) {
	resp, err := e.r.R().SetQueryParams(map[string]string{
		"name": name,
	}).Get("https://api.genderize.io")
	if err != nil {
		return "", fmt.Errorf("error getting gender: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("error getting gender: %s", resp.Status())
	}

	var data struct {
		Gender string `json:"gender"`
	}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return "", fmt.Errorf("error unmarshall gender: %w", err)
	}

	utils.DebugLog(name, "% Gender received", data.Gender)

	return data.Gender, nil
}

func (e *External) GetNationality(name string) (string, error) {
	resp, err := e.r.R().SetQueryParam("name", name).
		Get("https://api.nationalize.io")
	if err != nil {
		return "", fmt.Errorf("error getting nationality: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("error getting nationality: %s", resp.Status())
	}

	var data struct {
		Countries []internal.Country `json:"country"`
	}

	if err = json.Unmarshal(resp.Body(), &data); err != nil {
		return "", fmt.Errorf("error getting nationality: %w", err)
	}

	if len(data.Countries) == 0 {
		return "", errors.New("nationality not found")
	}

	randIndex := rand.Intn(len(data.Countries))
	randCountry := data.Countries[randIndex-1].Country

	utils.DebugLog(name, "% Nationality received", randCountry)

	return randCountry, nil
}
