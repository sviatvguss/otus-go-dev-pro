package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	if l.len == 0 {
		l.front = new(ListItem)
		l.back = l.front
		l.front.Value = v
	} else {
		newFront := new(ListItem)
		newFront.Next = l.front
		l.front.Prev = newFront
		newFront.Value = v
		l.front = newFront
	}
	l.len++
	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.len == 0 {
		l.back = new(ListItem)
		l.front = l.back
		l.back.Value = v
	} else {
		newBack := new(ListItem)
		newBack.Prev = l.back
		l.back.Next = newBack
		newBack.Value = v
		l.back = newBack
	}
	l.len++
	return l.back
}

func (l *list) Remove(i *ListItem) {
	switch {
	case i.Next != nil && i.Prev != nil: // not front, not back
		i.Prev.Next, i.Next.Prev = i.Next, i.Prev
	case i.Next == nil && i.Prev == nil: // front is back
		l.front, l.back = nil, nil
	case i.Next == nil: // back
		i.Prev.Next = nil
		l.back = i.Prev
	case i.Prev == nil: // front
		i.Next.Prev = nil
		l.front = i.Next
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	i.Next = l.front
	l.front.Prev = i
	l.front = i
	l.len++
}
