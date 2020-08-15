package state

type ID int

const (
	None ID = iota
	AddFriend
	RemoveFriend
	ImportFromVK
)
