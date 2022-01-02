package vk

import (
	"fmt"

	"github.com/SevereCloud/vksdk/object"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/lodthe/bdaytracker-go/models"
)

func friendObjectToFriend(obj *object.FriendsUserXtrLists) models.Friend {
	result := models.Friend{
		UUID: fmt.Sprint(uuid.New()),
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
