package currency

type ValuteCode string

type ValuteValue float64

type Valute struct {
	code  ValuteCode
	value ValuteValue
}

func (v Valute) Code() ValuteCode {
	return v.code
}

func (v Valute) Value() ValuteValue {
	return v.value
}

type Valutes []Valute

func NewValute(code ValuteCode, value ValuteValue) Valute {
	return Valute{
		code:  code,
		value: value,
	}
}

func (vals Valutes) SplitByPages(pageSize int) []Valutes {
	pages := make([]Valutes, 0)
	for i := 0; i < len(vals); i = i + pageSize {
		end := i + pageSize
		if end >= len(vals) {
			end = len(vals)
		}
		pages = append(pages, vals[i:end])
	}
	return pages
}
