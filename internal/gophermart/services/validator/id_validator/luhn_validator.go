package id_validator

import (
	"strconv"
	"strings"
)

type LuhnValidator struct {
}

func NewLuhnValidator() *LuhnValidator {
	return &LuhnValidator{}
}

func (validator *LuhnValidator) Validate(number string) bool {
	// Удаление пробелов из номера
	number = strings.Replace(number, " ", "", -1)

	// Проверка, что номер состоит только из цифр
	for _, char := range number {
		if char < '0' || char > '9' {
			return false
		}
	}

	// Разбор номера на цифры
	digits := make([]int, len(number))
	for i, char := range number {
		digits[i], _ = strconv.Atoi(string(char))
	}

	// Расчет контрольной суммы
	sum := 0
	double := false
	for i := len(digits) - 1; i >= 0; i-- {
		digit := digits[i]
		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		double = !double
	}

	// Проверка, делится ли сумма на 10
	return sum%10 == 0
}
