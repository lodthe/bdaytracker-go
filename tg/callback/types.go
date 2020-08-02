package callback

func Init() {
	addCallback(OpenMenu{})
	addCallback(AddFriend{})
	addCallback(AddFromVK{})
	addCallback(FriendsList{})
	addCallback(Settings{})
}

type OpenMenu struct {
}

type AddFriend struct {
}

type AddFromVK struct {
}

type FriendsList struct {
}

type Settings struct {
}
