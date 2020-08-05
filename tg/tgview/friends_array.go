package tgview

import (
	"github.com/lodthe/bdaytracker-go/models"
)

// FriendsArray implements methods of the sort.Interface
type FriendsArray struct {
	friends []models.Friend
}

func (a *FriendsArray) Len() int {
	return len(a.friends)
}

func (a *FriendsArray) Less(i, j int) bool {
	if a.friends[i].BMonth != a.friends[j].BMonth {
		return a.friends[i].BMonth < a.friends[j].BMonth
	}
	if a.friends[i].BDay != a.friends[j].BDay {
		return a.friends[i].BDay < a.friends[j].BDay
	}
	return a.friends[i].UUID < a.friends[j].UUID
}

func (a *FriendsArray) Swap(i, j int) {
	a.friends[i], a.friends[j] = a.friends[j], a.friends[i]
}
