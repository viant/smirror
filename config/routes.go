package config

//Routes represents route slice
type Routes []*Route

//HasMatch returns the first match route
func (r Routes) HasMatch(URL string) *Route {
	for i := range r {
		if r[i].HasMatch(URL) {
			return r[i]
		}
	}
	return nil
}
