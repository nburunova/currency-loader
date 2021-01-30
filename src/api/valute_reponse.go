package api

type ValuteResponse struct {
	Code  string  `json:"code"`
	Value float64 `json:"value"`
}

func NewValuteResponse(val Valute) ValuteResponse {
	return ValuteResponse{
		Code:  val.code,
		Value: val.value,
	}
}
