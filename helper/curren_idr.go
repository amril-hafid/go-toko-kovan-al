package helper

import (
	"fmt"
	"strings"
)

func FormatRupiah(amount int) string {
	// Mengubah angka menjadi string dengan format ribuan
	formatted := fmt.Sprintf("%d", amount)
	n := len(formatted)

	// Menambahkan titik setiap tiga digit dari belakang
	var result strings.Builder
	for i, digit := range formatted {
		if i > 0 && (n-i)%3 == 0 {
			result.WriteString(".")
		}
		result.WriteRune(digit)
	}

	// Menambahkan prefix "Rp. "
	return "IDR " + result.String()
}
