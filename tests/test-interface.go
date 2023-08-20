package tests

// Tester is an interface that describes objects that can perform a test.
type Tester interface {
	Test() bool
	Name() string
	GetContext() Context
	Initialise()
}
