package zerotier

type Member struct {
	Id            string   `json:"id"`
	Address       string   `json:"address"`
	NetworkId     string   `json:"nwid"`
	Authorized    bool     `json:"authorized"`
	ActiveBridge  bool     `json:"activeBridge"`
	IPAssignments []string `json:"ipAssignments"`
}
