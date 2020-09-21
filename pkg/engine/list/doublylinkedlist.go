package list

import (
	"fmt"
)

type node struct {
	data []byte
	pre  *node
	next *node
}

type List struct {
	len        int
	front, end *node
}

func (l *List) AddToFront(data []byte) {
	n := &node{data: data, pre: nil, next: l.front}
	if l.front == nil {
		l.end = n
	} else {
		l.front.pre = n
	}
	l.len++
	l.front = n
}

func (l *List) AddToBack(data []byte) {
	n := &node{data: data, pre: l.end, next: nil}
	if l.end == nil {
		l.front = n
	} else {
		l.end.next = n
	}
	l.len++
	l.end = n
}

func (l *List) AddBefore(n *node, data []byte) {
	tmp := &node{data: data, pre: n.pre, next: n}
	n.pre = tmp
	if tmp.pre != nil {
		tmp.pre.next = tmp
	} else {
		l.front = tmp
	}
	l.len++
}

func (l *List) AddAfter(n *node, data []byte) {
	tmp := &node{data: data, pre: n, next: n.next}
	n.next = tmp
	if tmp.next != nil {
		tmp.next.pre = tmp
	} else {
		l.end = tmp
	}
	l.len++
}

func (l *List) Remove(n *node) {
	if n.pre == nil {
		l.front = l.front.next
		if l.front != nil {
			l.front.pre = nil
		}
	} else if n.next == nil {
		l.end = l.end.pre
		if l.end != nil {
			l.front.pre = nil
		}
	} else {
		n.pre.next = n.next
		n.next.pre = n.pre
	}
	l.len--
	n.pre = nil
	n.next = nil
}

func (l *List) Print() {
	p := l.front
	for {
		if p != nil {
			fmt.Printf(" -> %s", string(p.data))
			p = p.next
		} else {
			break
		}
	}
}
