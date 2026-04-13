package types

type RouteStatus struct {
	CutOff     string `json:"cut_off"`
	StatusRute string `json:"status_rute"`
}

type GatewayRoutes map[string][]RouteStatus
