package money

import (
	"fmt"
)

const (
	EUR Currency = "EUR"
	USD Currency = "USD"
)

func Eur(a int) Money {
	return Money{
		Amount:   a,
		Currency: EUR,
	}
}

func Usd(a int) Money {
	return Money{
		Amount:   a,
		Currency: USD,
	}
}

func FromCurrency(a int, c Currency) Money {
	return Money{
		Amount:   a,
		Currency: c,
	}
}

type Currency string

func (c Currency) Symbol() string {
	switch c {
	case EUR:
		return "\u20ac"
	case USD:
		return "$"
	}

	return ""
}

type Money struct {
	Amount   int
	Currency Currency
}

func (m Money) Add(o Money) (Money, error) {
	if err := m.sameCurrency(o); err != nil {
		return Money{}, err
	}

	return Money{Amount: m.Amount + o.Amount, Currency: m.Currency}, nil
}

func (m Money) Substract(o Money) (Money, error) {
	if err := m.sameCurrency(o); err != nil {
		return Money{}, err
	}

	return Money{Amount: m.Amount - o.Amount, Currency: m.Currency}, nil
}

func (m Money) Multiply(f int) Money {
	return Money{Amount: m.Amount * f, Currency: m.Currency}
}

func (m Money) Divide(f int) Money {
	return Money{Amount: m.Amount / f, Currency: m.Currency}
}

func (m Money) Discount(p int) Money {
	discount := int(float64(m.Amount) * (float64(p) / 100))

	return Money{Amount: discount, Currency: m.Currency}
}

func (m Money) sameCurrency(o Money) error {
	if m.Currency != o.Currency {
		return fmt.Errorf("currency %s does not match with %s", m.Currency, o.Currency)
	}

	return nil
}

func (m Money) String() string {
	return fmt.Sprintf("%.2f%s", float64(m.Amount)/100, m.Currency.Symbol())
}
