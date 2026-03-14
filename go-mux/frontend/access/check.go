package access

// Level represents the access level of a user.
type Level uint32

const (
	LevelPublic Level = 0
	LevelEditor Level = 1 << iota
	LevelAdmin  Level = LevelEditor | (1 << iota)
)
