package stat

import (
	"math"

	"github.com/khoanguyen96/WarofRuneterra/server/dragon"
)

const PointStart = 100

func WithRegion(r Regions) *Reward {
	re := &Reward{
		regions: r,
	}

	re.begin()

	return re
}

func (r *Reward) begin() {
	r.set = make(PointSet, len(r.regions))

	for _, region := range r.regions {
		r.set[region] = PointStart
	}
}

func (r *Reward) tallyInit() map[string]int64 {
	tally := make(map[string]int64)

	for _, region := range r.regions {
		tally[region] = 0
	}

	return tally
}

func (r *Reward) Calc(win bool, cards []*dragon.Card) *PointSet {
	tally := r.tallyInit()
	total := len(cards)

	for _, c := range cards {
		if t, ok := tally[c.Region]; ok {
			tally[c.Region] = t + 1
		}
	}

	// todo: win or lose?
	for region, hits := range tally {
		score := (float64(hits) / float64(total)) * 100.0
		score = math.Round(score)

		i := int64(score)

		if win {
			r.set[region] += i * 10
		} else {
			r.set[region] -= (1 - i) * 10
		}
	}

	return &r.set
}
