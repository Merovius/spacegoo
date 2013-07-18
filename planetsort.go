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

type planetsByShips Planets

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
	sorter.Planets[i] = sorter.Planets[j]
	sorter.Planets[j] = t
}

func (sorter planetsByShips) Len() int {
	return len(sorter)
}

func (sorter planetsByShips) Less(i, j int) bool {
	return sorter[i].Ships.Sum() < sorter[j].Ships.Sum()
}

func (sorter planetsByShips) Swap(i, j int) {
	t := sorter[i]
	sorter[i] = sorter[j]
	sorter[j] = t
}

// SortByFDist sorts the planets by distance from a given point (in
// float-coordinates)
func (pl Planets) SortByFDist(X, Y float64) Planets {
	plcpy := make(Planets, len(pl))
	copy(plcpy, pl)
	sorter := planetsByDist{plcpy, false, 0, 0, X, Y}
	sort.Sort(sorter)
	return sorter.Planets
}

// SortByDist sorts the planets by distance from a given point (in
// int-coordinates)
func (pl Planets) SortByDist(X, Y int) Planets {
	plcpy := make(Planets, len(pl))
	copy(plcpy, pl)
	sorter := planetsByDist{plcpy, false, X, Y, 0.0, 0.0}
	sort.Sort(sorter)
	return sorter.Planets
}

// SortByShips sorts the planets in ascending order by power of fleetsize.
func (pl Planets) SortByShips() Planets {
	plcpy := make(Planets, len(pl))
	copy(plcpy, pl)
	sort.Sort(planetsByShips(plcpy))
	return plcpy
}
