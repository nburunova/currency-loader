package currency

import (
	"context"

	"github.com/nburunova/currency-loader/src/services/cbr"
	"github.com/pkg/errors"
)

type CbrClient interface {
	Valutes(ctx context.Context) (cbr.Valutes, error)
}

type CbrLoader struct {
	client CbrClient
}

func NewCbrLoader(client CbrClient) CbrLoader {
	return CbrLoader{
		client: client,
	}
}

func (c CbrLoader) Load(ctx context.Context) (Valutes, error) {
	cbrValutes, err := c.client.Valutes(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot load Cbr currency course")
	}
	res := make(Valutes, 0)
	for _, val := range cbrValutes {
		res = append(res, NewValute(ValuteCode(val.Code()), ValuteValue(val.Value())))
	}
	return res, nil
}
