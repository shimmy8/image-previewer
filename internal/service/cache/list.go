package cache

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
	length    int
	firstItem *ListItem
	lastItem  *ListItem
}

func (lst *list) Len() int {
	return lst.length
}

func (lst *list) Front() *ListItem {
	return lst.firstItem
}

func (lst *list) Back() *ListItem {
	return lst.lastItem
}

func (lst *list) PushFront(v interface{}) *ListItem {
	listItem := &ListItem{Value: v}
	if lst.firstItem == nil {
		lst.firstItem = listItem
		lst.lastItem = listItem
	} else {
		listItem.Next = lst.firstItem
		lst.firstItem.Prev = listItem
		lst.firstItem = listItem
	}
	lst.length++
	return listItem
}

func (lst *list) PushBack(v interface{}) *ListItem {
	listItem := &ListItem{Value: v}
	if lst.lastItem == nil { // lst is empty
		lst.lastItem = listItem
		lst.firstItem = listItem
	} else {
		listItem.Prev = lst.lastItem
		lst.lastItem.Next = listItem
		lst.lastItem = listItem
	}
	lst.length++
	return listItem
}

func (lst *list) Remove(i *ListItem) {
	if i.Prev == nil {
		lst.firstItem = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		lst.lastItem = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	lst.length--
	i = nil
}

func (lst *list) MoveToFront(i *ListItem) {
	if lst.firstItem == i {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
		if i.Next == nil { // it was a last item
			lst.lastItem = i.Prev
		}
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	lst.firstItem.Prev = i
	i.Next = lst.firstItem
	i.Prev = nil
	lst.firstItem = i
}

func NewList() List {
	return new(list)
}
