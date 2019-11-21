package stat

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
	tally := make(map[string]int64, len(r.regions))

	for _, region := range r.regions {
		tally[region] = 0
	}

	return tally
}

func (r *Reward) Calc(win bool, regions []string) *PointSet {
	tally := r.tallyInit()
	total := len(regions)

	for _, name := range regions {
		if t, ok := tally[name]; ok {
			tally[name] = t + 1
		}
	}

	// todo: win or lose?
	for region, hits := range tally {
		score := (hits / int64(total)) * 100

		if win {
			r.set[region] += score * 10
		} else {
			r.set[region] += (1 - score) * 10
		}
	}

	return &r.set
}
