package friendship

type Friend struct {
	UUID string

	// First name + surname.
	Name string
	VKID *int

	BDay   int
	BMonth int
}
