package einvoice

import "strings"

var digits = [...]string{"không", "một", "hai", "ba", "bốn", "năm", "sáu", "bảy", "tám", "chín"}

// AmountInWords converts an integer amount (VND) to Vietnamese words.
// E.g. 1234567895 → "Một tỷ hai trăm ba mươi bốn triệu năm trăm sáu mươi bảy nghìn tám trăm chín mươi lăm đồng"
func AmountInWords(amount int64) string {
	if amount == 0 {
		return "Không đồng"
	}

	negative := amount < 0
	if negative {
		amount = -amount
	}

	parts := []string{}
	// Process groups of 3 digits from right to left: ones, thousands, millions, billions
	groups := []string{"", "nghìn", "triệu", "tỷ"}

	for i := 0; amount > 0; i++ {
		group := int(amount % 1000)
		amount /= 1000
		if group == 0 {
			continue
		}
		s := readThreeDigits(group)
		if i < len(groups) && groups[i] != "" {
			s += " " + groups[i]
		}
		parts = append([]string{s}, parts...)
	}

	result := strings.Join(parts, " ")
	result = strings.ToUpper(result[:1]) + result[1:]
	if negative {
		result = "Âm " + result
	}
	return result + " đồng"
}

func readThreeDigits(n int) string {
	hundreds := n / 100
	remainder := n % 100
	tens := remainder / 10
	ones := remainder % 10

	var s []string

	if hundreds > 0 {
		s = append(s, digits[hundreds]+" trăm")
	}

	if tens > 1 {
		s = append(s, digits[tens]+" mươi")
		if ones > 0 {
			if ones == 5 {
				s = append(s, "lăm")
			} else {
				s = append(s, digits[ones])
			}
		}
	} else if tens == 1 {
		s = append(s, "mười")
		if ones > 0 {
			if ones == 5 {
				s = append(s, "lăm")
			} else {
				s = append(s, digits[ones])
			}
		}
	} else if hundreds > 0 && ones > 0 {
		if ones == 5 {
			s = append(s, "lăm")
		} else {
			s = append(s, digits[ones])
		}
	} else if ones > 0 {
		if ones == 5 {
			s = append(s, "năm")
		} else {
			s = append(s, digits[ones])
		}
	}

	return strings.Join(s, " ")
}
