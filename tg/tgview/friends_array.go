package tgview

import (
	"github.com/lodthe/bdaytracker-go/models"
)

// FriendsArray implements methods of the sort.Interface
type FriendsArray struct {
	Friends []models.Friend
}

func (a *FriendsArray) Len() int {
	return len(a.Friends)
}

func (a *FriendsArray) Less(i, j int) bool {
	if a.Friends[i].BMonth != a.Friends[j].BMonth {
		return a.Friends[i].BMonth < a.Friends[j].BMonth
	}
	if a.Friends[i].BDay != a.Friends[j].BDay {
		return a.Friends[i].BDay < a.Friends[j].BDay
	}
	return a.Friends[i].UUID < a.Friends[j].UUID
}

func (a *FriendsArray) Swap(i, j int) {
	a.Friends[i], a.Friends[j] = a.Friends[j], a.Friends[i]
}
