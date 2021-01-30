package currency

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaging(t *testing.T) {
	vals := Valutes{
		NewValute(ValuteCode("RUB"), ValuteValue(100)),
		NewValute(ValuteCode("RUB2"), ValuteValue(200)),
		NewValute(ValuteCode("RUB3"), ValuteValue(300)),
	}
	pages := vals.SplitByPages(2)
	assert.Equal(t, 2, len(pages))
	assert.Equal(t, 2, len(pages[0]))
	assert.Equal(t, 1, len(pages[1]))

	pages = vals.SplitByPages(3)
	assert.Equal(t, 1, len(pages))
	assert.Equal(t, 3, len(pages[0]))

	pages = vals.SplitByPages(10)
	assert.Equal(t, 1, len(pages))
	assert.Equal(t, 3, len(pages[0]))

	pages = vals.SplitByPages(1)
	assert.Equal(t, 3, len(pages))
	assert.Equal(t, 1, len(pages[0]))
	assert.Equal(t, 1, len(pages[1]))
	assert.Equal(t, 1, len(pages[2]))
}
