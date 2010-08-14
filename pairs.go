package oauth

import(
)


/*
type KeyValuePair struct {
	Key string
	Value string
}

type KeyPairs struct{
	Items []*KeyValuePair
}

func (kp KeyPairs) Len() int {  return len(kp.Items)}
func (kp KeyPairs) Less(i, j int) bool {
	if kp.Items[i].Key == kp.Items[j].Key {
		return kp.Items[i].Value < kp.Items[j].Value
	}
	return kp.Items[i].Key < kp.Items[j].Key
}
func (kp KeyPairs) Swap(i, j int) { kp.Items[i], kp.Items[j] = kp.Items[j], kp.Items[i]}


func (kp *KeyPairs) Add( kvp *KeyValuePair){

	if kp.Items == nil { kp.Items = make([]*KeyValuePair, 0, 4) }

	n := len(kp.Items)

	if n+1 > cap(kp.Items) {
		s := make([]*KeyValuePair, n, 2*n+1)
		copy(s, kp.Items)
		kp.Items = s
	}
	kp.Items = kp.Items[0 : n+1]
	kp.Items[n] = kvp

}

func (kp KeyPairs) Get( key string) string { 
	for i := range kp.Items{
		if kp.Items[i].Key == key {
			return  kp.Items[i].Value
		}
	}
	return ""
}

*/



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


