package service

import (
	"strconv"
)

func checkOrderNumber(orderNumber string) bool {
	digits := make([]int, len(orderNumber))
	for i, s := range orderNumber {
		digit, err := strconv.Atoi(string(s))
		if err != nil {
			return false
		}
		digits[i] = digit
	}

	sum := 0
	parity := len(orderNumber) % 2
	for i, digit := range digits {
		if i%2 == parity {
			digit = digit * 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}
