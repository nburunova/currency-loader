package api

import "strings"

type CurrenciesResponse struct {
	Page       int              `json:"page"`
	TotalPages int              `json:"total_pages"`
	Valutes    []ValuteResponse `json:"valutes"`
}

type Valute struct {
	code  string
	value float64
}

type Valutes []Valute

func NewValute(code string, value float64) Valute {
	return Valute{
		code:  strings.ToUpper(code),
		value: value,
	}
}

func NewCurrenciesResponse(page, totalPages int, vals Valutes) CurrenciesResponse {
	respVals := make([]ValuteResponse, 0)
	for _, val := range vals {
		respVals = append(respVals, NewValuteResponse(val))
	}
	return CurrenciesResponse{
		Page:       page,
		TotalPages: totalPages,
		Valutes:    respVals,
	}
}
