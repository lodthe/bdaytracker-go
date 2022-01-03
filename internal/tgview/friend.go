package tgview

import (
	"fmt"
	"strconv"

	friendship "github.com/lodthe/bdaytracker-go/internal/friendship"
)

// formatFriend returns a formatted representation of the given friend.
// It looks as belows:
// Name — 20.04
// If the VK ID of the friend is known, Name is a hyperlink to the VK profile.
func formatFriend(friend friendship.Friend) string {
	name := fmt.Sprintf("<code>%s</code>", friend.Name)
	if friend.VKID != nil {
		name = fmt.Sprintf("<a href=\"vk.com/id%d\">%s</a>", *friend.VKID, friend.Name)
	}

	if friend.BMonth == 0 || friend.BDay == 0 {
		return name
	}

	return fmt.Sprintf("%s — %02d.%02d", name, friend.BDay, friend.BMonth)
}

// formatFriendWithIndex returns a formatted representation of the given friend.
// It looks as belows:
// 013. Name — 20.04
// It takes friendsNumber to add enough leading zeroes before the index.
func formatFriendWithIndex(friend friendship.Friend, index, friendsNumber int) string {
	maxIndexLength := len(strconv.Itoa(friendsNumber))
	format := fmt.Sprintf("<b>%%0%dd</b>. %%s", maxIndexLength)
	return fmt.Sprintf(format, index, formatFriend(friend))
}
