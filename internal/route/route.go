package route

import (
	"fmt"
	"strings"

	"github.com/umardev500/jnt-report/internal/types"
)

func GetRitase(routes types.GatewayRoutes, tujuan, cutOff string) *types.RouteStatus {
	// Check if tujuan exists
	routeList, ok := routes[tujuan]
	if !ok {
		fmt.Println(tujuan)
		return nil // empty slice if tujuan not found
	}

	// Filter by clock
	for _, rt := range routeList {
		if rt.CutOff == cutOff {
			return &rt
		}
	}

	return &types.RouteStatus{
		CutOff:     cutOff,
		StatusRute: "ADD",
	}
}

func GetJenisPaket(status string) string {
	s := strings.ToLower(strings.TrimSpace(status))

	if strings.Contains(s, "reguler") || strings.Contains(s, "wajib") {
		return "REG"
	}

	// For balikan
	if strings.Contains(s, "balikan") {
		return "BALIKAN"
	}

	return "ADD"
}
