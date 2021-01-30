package cbr

import (
	"encoding/xml"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type rawValute struct {
	XMLName  xml.Name `xml:"Valute"`
	ID       string   `xml:"ID,attr"`
	NumCode  string   `xml:"NumCode"`
	CharCode string   `xml:"CharCode"`
	Nominal  string   `xml:"Nominal"`
	Name     string   `xml:"Name"`
	Value    string   `xml:"Value"`
}

type Valute struct {
	id       string
	numCode  string
	charCode string
	nominal  int
	name     string
	value    float64
}

func (v Valute) Value() float64 {
	return v.value / float64(v.nominal)
}

func (v Valute) Code() string {
	return v.charCode
}

type Valutes []Valute

func (r rawValute) Valute() (*Valute, error) {
	var nominal int
	var value float64
	var err error
	if r.Nominal != "" {
		nominal, err = strconv.Atoi(r.Nominal)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot convert valute nominal to int %v", r.Nominal)
		}
	}
	if r.Value != "" {
		value, err = strconv.ParseFloat(strings.Replace(r.Value, ",", ".", -1), 64)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot convert value to float %v", r.Value)
		}
	}
	return &Valute{
		id:       r.ID,
		numCode:  r.NumCode,
		charCode: r.CharCode,
		nominal:  nominal,
		name:     r.Name,
		value:    value,
	}, nil
}

type rawValCurs struct {
	XMLName    xml.Name    `xml:"ValCurs"`
	Date       string      `xml:"Date,attr"`
	Name       string      `xml:"name,attr"`
	RawValutes []rawValute `xml:"Valute"`
}

func (r rawValCurs) Valutes() (Valutes, error) {
	valutes := make(Valutes, 0)
	for _, rValute := range r.RawValutes {
		val, err := rValute.Valute()
		if err != nil {
			return nil, errors.Wrap(err, "cannot convert valute")
		}
		valutes = append(valutes, *val)
	}
	return valutes, nil
}
