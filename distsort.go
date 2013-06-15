package spacegoo

import (
	"sort"
)

type planetsByDist struct {
	Planets Planets
	Floats  bool
	X       int
	Y       int
	FX      float64
	FY      float64
}

func (sorter planetsByDist) Len() int {
	return len(sorter.Planets)
}

func (sorter planetsByDist) Less(i, j int) bool {
	pi := sorter.Planets[i]
	pj := sorter.Planets[j]
	if sorter.Floats {
		return pi.FDist(sorter.FX, sorter.FY) < pj.FDist(sorter.FX, sorter.FY)
	}
	return pi.Dist(sorter.X, sorter.Y) < pj.Dist(sorter.X, sorter.Y)
}

func (sorter planetsByDist) Swap(i, j int) {
	t := sorter.Planets[i]
	sorter.Planets[j] = sorter.Planets[i]
	sorter.Planets[i] = t
}

func (pl Planets) SortByFDist(X, Y float64) Planets {
	plcpy := make(Planets, len(pl))
	sorter := planetsByDist{plcpy, false, 0, 0, X, Y}
	sort.Sort(sorter)
	return sorter.Planets
}

func (pl Planets) SortByDist(X, Y int) Planets {
	plcpy := make(Planets, len(pl))
	sorter := planetsByDist{plcpy, false, X, Y, 0.0, 0.0}
	sort.Sort(sorter)
	return sorter.Planets
}
