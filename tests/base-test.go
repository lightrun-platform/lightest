package tests

import (
	. "prerequisites-tester/communication/response-handling"
	. "prerequisites-tester/tests/logging"
)

type TestBase struct {
	Context Context
	TestLogger
}

func (test *TestBase) GetContext() Context {
	return test.Context
}

func (test *TestBase) Initialise(name string) {
	test.Context = Context{
		Store:           make(map[string]interface{}),
		HandlingContext: &HandlingContext{},
	}
	test.TestName = name
}

type Context struct {
	Store           map[string]interface{}
	HandlingContext *HandlingContext
}
