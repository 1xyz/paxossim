package v1

import (
	"container/list"
)

type typeToProcessMap map[ProcessType]*list.List

func newMultiValueMap() *typeToProcessMap {
	v := make(typeToProcessMap)
	return &v
}

func (mm typeToProcessMap) get(k ProcessType) (*list.List, bool) {
	v, ok := mm[k]
	return v, ok
}

func (mm typeToProcessMap) put(p ProcessInbox) {
	k := p.Type()
	_, ok := mm[k]
	if !ok {
		mm[k] = list.New()
	}
	mm[k].PushBack(p)
}

func (mm typeToProcessMap) remove(k ProcessInbox) {
	v, ok := mm[k.Type()]
	if !ok {
		return
	}
	for e := v.Front(); e != nil; e = e.Next() {
		if e.Value.(ProcessInbox).ID() == k.ID() {
			v.Remove(e)
			return
		}
	}
}

type AddrSet map[Addr]bool

func (a AddrSet) Contains(addr Addr) bool {
	_, ok := a[addr]
	return ok
}

func (a AddrSet) Add(addr Addr) {
	a[addr] = true
}

func (a AddrSet) Remove(addr Addr) {
	if a.Contains(addr) {
		delete(a, addr)
	}
}

func (a AddrSet) Len() int {
	return len(a)
}
