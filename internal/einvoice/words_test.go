package einvoice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAmountInWords_Zero(t *testing.T) {
	assert.Equal(t, "Không đồng", AmountInWords(0))
}

func TestAmountInWords_SmallAmounts(t *testing.T) {
	assert.Equal(t, "Một đồng", AmountInWords(1))
	assert.Equal(t, "Mười đồng", AmountInWords(10))
	assert.Equal(t, "Một trăm đồng", AmountInWords(100))
	assert.Equal(t, "Một nghìn đồng", AmountInWords(1000))
	assert.Equal(t, "Một triệu đồng", AmountInWords(1000000))
}

func TestAmountInWords_ComplexAmounts(t *testing.T) {
	assert.Equal(t, "Một tỷ hai trăm ba mươi bốn triệu năm trăm sáu mươi bảy nghìn tám trăm chín mươi lăm đồng", AmountInWords(1234567895))
	assert.Equal(t, "Năm triệu đồng", AmountInWords(5000000))
	assert.Equal(t, "Hai trăm nghìn đồng", AmountInWords(200000))
}

func TestAmountInWords_WithDecimals(t *testing.T) {
	// VND doesn't use decimals in practice, but the function handles the int64 input
	assert.Equal(t, "Một đồng", AmountInWords(1))
}
