package utils

import (
	"fmt"
	"math"
	"strings"
)

type Money struct {
	Amount   float64 `bson:"amount"`
	Currency string  `bson:"currency;omitempty"`
}

func (m *Money) Round() Money {
	newMoney := NewMoney(m.Amount, m.Currency)
	newMoney.Amount = math.Round(newMoney.Amount)
	return newMoney
}

func (m *Money) Ceil() Money {
	newMoney := NewMoney(m.Amount, m.Currency)
	newMoney.Amount = math.Ceil(newMoney.Amount)
	return newMoney
}

func (m *Money) DivideMoney(money Money) Money {
	newMoney := NewMoney(m.Amount, m.Currency)
	newMoney.DivideWithMoney(money)
	return newMoney
}

func (m *Money) Divide(amount float64) Money {
	newMoney := NewMoney(m.Amount, m.Currency)
	newMoney.DivideWith(amount)
	return newMoney
}

func (m *Money) DivideWithMoney(money Money) {
	m.validateCurrency(money)
	m.Amount /= money.Amount
}

func (m *Money) DivideWith(amount float64) { m.Amount /= amount }

func (m *Money) MultiplyMoney(money Money) Money {
	newMoney := NewMoney(m.Amount, m.Currency)
	newMoney.MultiplyWithMoney(money)
	return newMoney
}

func (m *Money) Multiply(amount float64) Money {
	newMoney := NewMoney(m.Amount, m.Currency)
	newMoney.MultiplyWith(amount)
	return newMoney
}

func (m *Money) MultiplyWithMoney(money Money) {
	m.validateCurrency(money)
	m.Amount *= money.Amount
}

func (m *Money) MultiplyWith(amount float64) { m.Amount *= amount }

func (m *Money) SubstractMoney(money Money) Money {
	newMoney := NewMoney(m.Amount, m.Currency)
	newMoney.SubstractWithMoney(money)
	return newMoney
}

func (m *Money) Substract(amount float64) Money {
	newMoney := NewMoney(m.Amount, m.Currency)
	newMoney.SubstractWith(amount)
	return newMoney
}

func (m *Money) SubstractWithMoney(money Money) {
	m.validateCurrency(money)
	m.Amount -= money.Amount
}

func (m *Money) SubstractWith(amount float64) { m.Amount -= amount }

func (m *Money) Compare(other Money) int {
	m.validateCurrency(other)
	if m.Amount > other.Amount {
		return 1
	} else if m.Amount < other.Amount {
		return -1
	}
	return 0
}

func (m *Money) AddMoney(money Money) Money {
	newMoney := NewMoney(m.Amount, m.Currency)
	newMoney.AddWithMoney(money)
	return newMoney
}

func (m *Money) Add(amount float64) Money {
	newMoney := NewMoney(m.Amount, m.Currency)
	newMoney.AddWith(amount)
	return newMoney
}

func (m *Money) AddWithMoney(money Money) {
	m.validateCurrency(money)
	m.Amount += money.Amount
}

func (m *Money) validateCurrency(money Money) {
	if strings.ToLower(m.Currency) != strings.ToLower(money.Currency) {
		panic(fmt.Errorf("Currency not same"))
	}
}

func (m *Money) AddWith(amount float64) { m.Amount += amount }

func NewMoney(amount float64, currency string) Money {
	return Money{Amount: amount, Currency: currency}
}
