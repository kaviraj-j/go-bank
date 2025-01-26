package util

const (
	USD = "USD"
	CAD = "CAD"
	EUR = "EUR"
	INR = "INR"
)

func IsValidCurrency(currency string) bool {
	switch currency {
	case USD, CAD, EUR, INR:
		return true
	}
	return false
}
