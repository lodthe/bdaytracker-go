package vk

import (
	"fmt"

	"github.com/SevereCloud/vksdk/object"
	"github.com/google/uuid"
	friendship "github.com/lodthe/bdaytracker-go/internal/friendship"
	"github.com/sirupsen/logrus"
)

func friendObjectToFriend(obj *object.FriendsUserXtrLists) friendship.Friend {
	result := friendship.Friend{
		UUID: uuid.New().String(),
		Name: fmt.Sprintf("%s %s", obj.FirstName, obj.LastName),
		VKID: &obj.ID,
	}

	if obj.Bdate != "" {
		_, err := fmt.Sscanf(obj.Bdate, "%d.%d", &result.BDay, &result.BMonth)
		if err != nil {
			logrus.WithField("obj", obj).WithError(err).Error("failed to parse bdate")
		}
	}

	return result
}
