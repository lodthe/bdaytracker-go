package friendship

import (
	"sort"
)

func SortFriends(friends []Friend) []Friend {
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

func RemoveVKFriends(friends []Friend) []Friend {
	var result []Friend
	for _, friend := range friends {
		if friend.VKID == nil {
			result = append(result, friend)
		}
	}
	return result
}
