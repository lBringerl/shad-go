//go:build !solution

package speller

import "fmt"

func spellOnes(n int64) string {
	if n > 10 {
		panic(fmt.Sprintf("spellOnes: too big number: %d", n))
	}
	switch n {
	case 0:
		return ""
	case 1:
		return "one"
	case 2:
		return "two"
	case 3:
		return "three"
	case 4:
		return "four"
	case 5:
		return "five"
	case 6:
		return "six"
	case 7:
		return "seven"
	case 8:
		return "eight"
	case 9:
		return "nine"
	default:
		return "unknown"
	}
}

func spellLT20(n int64) string {
	if n > 20 {
		panic(fmt.Sprintf("spellLT20: too big number: %d", n))
	}
	switch n {
	case 10:
		return "ten"
	case 11:
		return "eleven"
	case 12:
		return "twelve"
	case 13:
		return "thirteen"
	case 14:
		return "fourteen"
	case 15:
		return "fifteen"
	case 16:
		return "sixteen"
	case 17:
		return "seventeen"
	case 18:
		return "eighteen"
	case 19:
		return "nineteen"
	default:
		return spellOnes(n)
	}
}

func spellTensNoOnes(tens int64) string {
	switch tens {
	case 0:
		return ""
	case 10:
		return "ten"
	case 20:
		return "twenty"
	case 30:
		return "thirty"
	case 40:
		return "forty"
	case 50:
		return "fifty"
	case 60:
		return "sixty"
	case 70:
		return "seventy"
	case 80:
		return "eighty"
	case 90:
		return "ninety"
	default:
		return "unknown"
	}
}

func spellTens(n int64) string {
	if n > 99 {
		panic(fmt.Sprintf("spellTens: too big number: %d", n))
	}
	ones := n % 10
	tens := n - ones
	if ones == 0 {
		return spellTensNoOnes(tens)
	}
	switch tens {
	case 0:
		return spellOnes(n)
	case 10:
		return spellLT20(n)
	case 20:
		return fmt.Sprintf("twenty-%s", spellOnes(ones))
	case 30:
		return fmt.Sprintf("thirty-%s", spellOnes(ones))
	case 40:
		return fmt.Sprintf("forty-%s", spellOnes(ones))
	case 50:
		return fmt.Sprintf("fifty-%s", spellOnes(ones))
	case 60:
		return fmt.Sprintf("sixty-%s", spellOnes(ones))
	case 70:
		return fmt.Sprintf("seventy-%s", spellOnes(ones))
	case 80:
		return fmt.Sprintf("eighty-%s", spellOnes(ones))
	case 90:
		return fmt.Sprintf("ninety-%s", spellOnes(ones))
	default:
		return "unknown"
	}
}

func spellHundreds(n int64) string {
	if n > 999 {
		panic(fmt.Sprintf("spellHundreds: too big number: %d", n))
	}
	if n < 100 {
		return spellTens(n)
	}
	hundreds := n - (n % 100)
	hundredOnes := hundreds / 100
	tens := n - hundreds
	if tens == 0 {
		return fmt.Sprintf("%s hundred", spellOnes(hundredOnes))
	}
	return fmt.Sprintf("%s hundred %s", spellOnes(hundredOnes), spellTens(tens))
}

func spellThousands(n int64) string {
	if n > 999999 {
		panic(fmt.Sprintf("spellThousands: too big number: %d", n))
	}
	if n < 1000 {
		return spellHundreds(n)
	}
	hundreds := n % 1000
	thousands := n - hundreds
	thousandHundreds := thousands / 1000
	if hundreds == 0 {
		return fmt.Sprintf("%s thousand", spellHundreds(thousandHundreds))
	}
	return fmt.Sprintf("%s thousand %s", spellHundreds(thousandHundreds), spellHundreds(hundreds))
}

func spellMillions(n int64) string {
	if n > 999999999 {
		panic(fmt.Sprintf("spellMillions: too big number: %d", n))
	}
	if n < 1000000 {
		return spellThousands(n)
	}
	thousands := n % 1000000
	millions := n - thousands
	millionsThousands := millions / 1000000
	if thousands == 0 {
		return fmt.Sprintf("%s million", spellHundreds(millionsThousands))
	}
	return fmt.Sprintf("%s million %s", spellHundreds(millionsThousands), spellThousands(thousands))
}

func spellBillions(n int64) string {
	if n > 999999999999 {
		panic(fmt.Sprintf("spellBillions: too big number: %d", n))
	}
	if n < 1000000000 {
		return spellMillions(n)
	}
	millions := n % 1000000000
	billions := n - millions
	billionsThousands := billions / 1000000000
	if millions == 0 {
		return fmt.Sprintf("%s billion", spellHundreds(billionsThousands))
	}
	return fmt.Sprintf("%s billion %s", spellHundreds(billionsThousands), spellMillions(millions))
}

func Spell(n int64) string {
	if n == 0 {
		return "zero"
	}
	if n < 0 {
		return fmt.Sprintf("minus %s", spellBillions(-n))
	}
	return spellBillions(n)
}
