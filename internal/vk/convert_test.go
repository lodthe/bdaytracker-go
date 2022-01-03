package vk

import (
	"testing"

	"github.com/SevereCloud/vksdk/object"
	"github.com/stretchr/testify/assert"
)

func TestFriendObjectToFriend(t *testing.T) {
	out := friendObjectToFriend(&object.FriendsUserXtrLists{
		UsersUser: object.UsersUser{
			FirstName: "Pavel",
			LastName:  "Durov",
			ID:        123123,
			Bdate:     "20.01.2001",
		},
	})

	assert.Equal(t, "Pavel Durov", out.Name)
	assert.Equal(t, 123123, *out.VKID)
	assert.Equal(t, 20, out.BDay)
	assert.Equal(t, 1, out.BMonth)

	out = friendObjectToFriend(&object.FriendsUserXtrLists{
		UsersUser: object.UsersUser{
			FirstName: "Nikolay",
			LastName:  "Durov",
			ID:        222882828,
			Bdate:     "01.08",
		},
	})

	assert.Equal(t, "Nikolay Durov", out.Name)
	assert.Equal(t, 222882828, *out.VKID)
	assert.Equal(t, 1, out.BDay)
	assert.Equal(t, 8, out.BMonth)
}
