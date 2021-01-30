package cbr

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCurses(t *testing.T) {
	data := []byte(`
<?xml version="1.0" encoding="UTF-8"?>
<ValCurs Date="30.01.2021" name="Foreign Currency Market">
   <Valute ID="R01010">
      <NumCode>036</NumCode>
      <CharCode>AUD</CharCode>
      <Nominal>1</Nominal>
      <Name>Австралийский доллар</Name>
      <Value>58,3333</Value>
   </Valute>
   <Valute ID="R01020A">
      <NumCode>944</NumCode>
      <CharCode>AZN</CharCode>
      <Nominal>1</Nominal>
      <Name>Азербайджанский манат</Name>
      <Value>44,8809</Value>
   </Valute>
</ValCurs>
  `)
	var res rawValCurs
	err := xml.Unmarshal(data, &res)
	assert.Nil(t, err)
	valutes, err := res.Valutes()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(valutes))
	assert.Equal(t, float64(58.3333), valutes[0].value)
	assert.Equal(t, int(1), valutes[0].nominal)
	assert.Equal(t, "036", valutes[0].numCode)
	assert.Equal(t, "AUD", valutes[0].charCode)
	assert.Equal(t, "Австралийский доллар", valutes[0].name)
}
