package spacegoo

import (
	"sort"
)

type Fleets []Fleet

type fleetByEtaSorter struct {
	f Fleets
}

func (fs fleetByEtaSorter) Len() int {
	return len(fs.f)
}

func (fs fleetByEtaSorter) Less(i, j int) bool {
	return fs.f[i].Eta < fs.f[j].Eta
}

func (fs fleetByEtaSorter) Swap(i, j int) {
	t := fs.f[i]
	fs.f[j] = fs.f[i]
	fs.f[i] = t
}

// Sort the Fleets ascending by ETA
func (f Fleets) Sort() {
	sort.Sort(fleetByEtaSorter{f})
}

type fleetByShipsSorter struct {
	f Fleets
}

func (fs fleetByShipsSorter) Len() int {
	return len(fs.f)
}

func (fs fleetByShipsSorter) Less(i, j int) bool {
	return fs.f[i].Ships.Sum() < fs.f[j].Ships.Sum()
}

func (fs fleetByShipsSorter) Swap(i, j int) {
	t := fs.f[i]
	fs.f[j] = fs.f[i]
	fs.f[i] = t
}

// Sort the Fleets ascending by ETA
func (f Fleets) SortByStrength() {
	sort.Sort(fleetByShipsSorter{f})
}
