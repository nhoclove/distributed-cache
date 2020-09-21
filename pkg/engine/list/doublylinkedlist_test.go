package list

import "testing"

func TestAddBefore(t *testing.T) {
	ll := List{}
	ll.AddToFront([]byte("1"))
	ll.AddToBack([]byte("2"))
	ll.AddAfter(ll.end, []byte("3"))
	ll.Print()
}

func TestRemove(t *testing.T) {
	ll := List{}
	ll.AddToFront([]byte("1"))
	ll.AddToBack([]byte("2"))
	ll.Remove(ll.end)
	ll.Remove(ll.end)
	ll.Remove(ll.end)
	ll.Print()
}
