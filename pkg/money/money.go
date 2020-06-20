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
		amount:   a,
		currency: EUR,
	}
}

func Usd(a int) Money {
	return Money{
		amount:   a,
		currency: USD,
	}
}

type Currency string

func (c Currency) Symbol() string {
	switch c {
	case EUR:
		return "?"
	case USD:
		return "$"
	}

	return ""
}

type Money struct {
	amount   int
	currency Currency
}

func (m Money) Add(o Money) (Money, error) {
	if err := m.sameCurrency(o); err != nil {
		return Money{}, err
	}

	return Money{amount: m.amount + o.amount, currency: m.currency}, nil
}

func (m Money) Substract(o Money) (Money, error) {
	if err := m.sameCurrency(o); err != nil {
		return Money{}, err
	}

	return Money{amount: m.amount - o.amount, currency: m.currency}, nil
}

func (m Money) Multiply(f int) Money {
	return Money{amount: m.amount * f, currency: m.currency}
}

func (m Money) Divide(f int) Money {
	return Money{amount: m.amount / f, currency: m.currency}
}

func (m Money) sameCurrency(o Money) error {
	if m.currency != o.currency {
		return fmt.Errorf("currency %s does not match with %s", m.currency, o.currency)
	}

	return nil
}

func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", float64(m.amount)/100, m.currency.Symbol())
}
