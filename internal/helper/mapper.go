package helper

import "strings"

func MapVendor(vendor string) string {
	v := strings.ToUpper(strings.TrimSpace(vendor))

	mapping := map[string]string{
		"PT. SERIKAT HANTAR EKSPEDISI": "SHUTTLE PT. SERIKAT HANTAR EKSPEDISI",
		"PT SERIKAT HANTAR EKSPEDISI":  "SHUTTLE PT. SERIKAT HANTAR EKSPEDISI",

		"PT ANUGERAH BERSATU":  "SHUTTLE PT. ANUGERAH BERSATU",
		"PT. ANUGERAH BERSATU": "SHUTTLE PT. ANUGERAH BERSATU",

		"PT. ASIA PASIFIK LOGISTIK": "SHUTTLE VENDOR ASIA PASIFIK LOGISTIK",
		"PT ASIA PASIFIK LOGISTIK":  "SHUTTLE VENDOR ASIA PASIFIK LOGISTIK",
	}

	if val, ok := mapping[v]; ok {
		return val
	}

	return vendor
}
