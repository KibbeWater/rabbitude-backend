package structures

type LibraryLoad func() bool
type LibraryFree func() bool
type LibraryInterface struct {
	Name string

	IsAvailable bool
	Load        LibraryLoad
	Free        LibraryFree
}
