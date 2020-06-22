package money

import (
	"fmt"
)

const (
	EUR Currency = "EUR"
	USD Currency = "USD"
)

// Eur creates a new money instance with Eur currency
func Eur(a int) Money {
	return Money{
		Amount:   a,
		Currency: EUR,
	}
}

// Usd creates a new money instance with Eur currency
func Usd(a int) Money {
	return Money{
		Amount:   a,
		Currency: USD,
	}
}

// FromCurrency creates a new money instance with the given currency
func FromCurrency(c Currency) Money {
	return Money{
		Amount:   0,
		Currency: c,
	}
}

// Currency defines money currency
type Currency string

//Symbol return symbol representation of a currency if exists
func (c Currency) Symbol() string {
	switch c {
	case EUR:
		return "\u20ac"
	case USD:
		return "$"
	}

	return ""
}

//Money represents the value of an amount of money
type Money struct {
	Amount   int
	Currency Currency
}

// Adds returns the value of the money plus the given money amount if the currency matchs
func (m Money) Add(o Money) (Money, error) {
	if err := m.sameCurrency(o); err != nil {
		return Money{}, err
	}

	return Money{Amount: m.Amount + o.Amount, Currency: m.Currency}, nil
}

// Substract returns the value of the money minus the given money amount if the currency matchs
func (m Money) Substract(o Money) (Money, error) {
	if err := m.sameCurrency(o); err != nil {
		return Money{}, err
	}

	return Money{Amount: m.Amount - o.Amount, Currency: m.Currency}, nil
}

// Multiply returns the value of the money multiplied by the factor
func (m Money) Multiply(f int) Money {
	return Money{Amount: m.Amount * f, Currency: m.Currency}
}

// Divide returns the value of the money divided by the factor
func (m Money) Divide(f int) Money {
	return Money{Amount: m.Amount / f, Currency: m.Currency}
}

// Discount returns given percent of the money
func (m Money) Discount(p int) Money {
	discount := int(float64(m.Amount) * (float64(p) / 100))

	return Money{Amount: discount, Currency: m.Currency}
}

// IsZero returns if the money is initialized that means has currency
func (m Money) IsZero() bool {
	return m.Currency == ""
}

//sameCurrency returns error if the money currencies does not match
func (m Money) sameCurrency(o Money) error {
	if m.Currency != o.Currency {
		return fmt.Errorf("currency %s does not match with %s", m.Currency, o.Currency)
	}

	return nil
}

//String prints the value and currency formatted
func (m Money) String() string {
	return fmt.Sprintf("%.2f%s", float64(m.Amount)/100, m.Currency.Symbol())
}
