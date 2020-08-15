package helpers

import (
	"sort"

	"github.com/lodthe/bdaytracker-go/models"
)

func SortFriends(friends []models.Friend) []models.Friend {
	sort.Slice(friends, func(i, j int) bool {
		if friends[i].BMonth != friends[j].BMonth {
			return friends[i].BMonth < friends[j].BMonth
		}
		if friends[i].BDay != friends[j].BDay {
			return friends[i].BDay < friends[j].BDay
		}
		return friends[i].UUID < friends[j].UUID
	})
	return friends
}

func RemoveVKFriends(friends []models.Friend) []models.Friend {
	var result []models.Friend
	for _, friend := range friends {
		if friend.VKID == nil {
			result = append(result, friend)
		}
	}
	return result
}
