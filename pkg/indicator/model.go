package indicator

type Indicator struct {
	Open   float32 `json:"open"`
	High   float32 `json:"high"`
	Low    float32 `json:"low"`
	Close  float32 `json:"close"`
	Group  int64   `json:"-"`
	Period string  `json:"-"`
	Start  string  `json:"period_start"`
	End    string  `json:"period_end"`
}
