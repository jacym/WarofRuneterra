package stat

type (
	PointSet map[string]int64

	Regions []string

	Reward struct {
		Regions Regions
		Set     PointSet
	}
)
