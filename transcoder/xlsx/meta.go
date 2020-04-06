package xlsx

type meta struct {
	firstRow      int
	index         int
	lastRow       int
	cells         int
	nonEmptyCells int
}

func (m *meta) hasMore() bool {
	return m.index <= m.lastRow
}

type metaSlice []*meta

func (a metaSlice) Len() int           { return len(a) }
func (a metaSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a metaSlice) Less(i, j int) bool {
	if a[i].cells == a[j].cells {
		return a[i].nonEmptyCells < a[j].nonEmptyCells
	}
	return a[i].cells < a[j].cells }
