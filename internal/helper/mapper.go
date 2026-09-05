package helper

import "strings"

func MapVendor(vendor string) string {
	v := strings.ToUpper(strings.TrimSpace(vendor))

	// mapping := map[string]string{
	// 	"PT. SERIKAT HANTAR EKSPEDISI": "PT. SERIKAT HANTAR EKSPEDISI",
	// 	"PT SERIKAT HANTAR EKSPEDISI":  "PT. SERIKAT HANTAR EKSPEDISI",

	// 	"PT ANUGERAH BERSATU":  "PT. ANUGERAH BERSATU",
	// 	"PT. ANUGERAH BERSATU": "PT. ANUGERAH BERSATU",

	// 	"PT. ASIA PASIFIK LOGISTIK": "PT. ASIA PASIFIK LOGISTIK",
	// 	"PT ASIA PASIFIK LOGISTIK":  "PT. ASIA PASIFIK LOGISTIK",
	// }

	// if val, ok := mapping[v]; ok {
	// 	return val
	// }

	if strings.HasPrefix(v, "PT ") {
		return strings.Replace(v, "PT ", "PT. ", 1)
	}

	return v
}
