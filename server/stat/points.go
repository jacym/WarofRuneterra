package stat

import (
	"math"

	"github.com/jacym/WarofRuneterra/server/dragon"
)

const PointStart = 100

func WithRegion(r Regions) *Reward {
	re := &Reward{
		Regions: r,
	}

	re.begin()

	return re
}

func (r *Reward) begin() {
	r.Set = make(PointSet, len(r.Regions))

	for _, region := range r.Regions {
		r.Set[region] = PointStart
	}
}

func (r *Reward) tallyInit() map[string]int64 {
	tally := make(map[string]int64)

	for _, region := range r.Regions {
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
			r.Set[region] += i * 10
		} else {
			r.Set[region] -= (1 - i) * 10
		}
	}

	return &r.Set
}
