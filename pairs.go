package oauth

import(
)

type Pair struct {
	Key string
	Value string
}

type Params []*Pair

func (p Params) Len() int {  return len(p)}
func (p Params) Less(i, j int) bool {
	if p[i].Key == p[j].Key {
		return p[i].Value < p[j].Value
	}
	return p[i].Key < p[j].Key
}
func (p Params) Swap(i, j int) { p[i], p[j] = p[j], p[i]}


func (p *Params) Add(pair *Pair){
	a := *p
	n := len(a)

	if n+1 > cap(a) {
		s := make([]*Pair, n, 2*n+1)
		copy(s, a)
		a = s
	}
	a = a[0 : n+1]
	a[n] = pair
	*p = a

}


