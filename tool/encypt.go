package tool

type Args struct {
	Path         string   `json:"path"`
	Policy       []string `json:"policy"`
	PlainIndices string   `json:"plainIndices"`
}

func NewArgs(p, i string, plc []string) *Args {
	return &Args{
		Path:         p,
		Policy:       plc,
		PlainIndices: i,
	}
}
