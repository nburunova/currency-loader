package sql_storage

import (
	"context"
	"testing"
	"time"

	"github.com/nburunova/currency-loader/src/currency"
	"github.com/nburunova/currency-loader/src/services/log"
	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	logger := log.NewLogger("json", "debug")
	lifeTime := time.Duration(600) * time.Millisecond
	strg, err := New(lifeTime, logger)
	assert.Nil(t, err)
	ctx := context.Background()
	strg.Start(ctx, lifeTime/10)

	vals := currency.Valutes{
		currency.NewValute("RUB", 100),
		currency.NewValute("RU2", 200),
		currency.NewValute("RUB3", 300),
	}
	strg.Save(ctx, vals)
	res, err := strg.ValutesOnPage(ctx, 1, 2)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(res))

	time.Sleep(lifeTime * 2)
	res, err = strg.ValutesOnPage(ctx, 1, 2)
	assert.NotNil(t, err)
	assert.Nil(t, res)
}
