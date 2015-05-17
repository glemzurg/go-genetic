package genetic

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
)

// Create a suite.
type UtilsSuite struct{}

var _ = Suite(&UtilsSuite{})

// Add the tests.

func (s *UtilsSuite) Test_InStrings(c *C) {
	// NOTE: the searched string must be sorted!

	// Exercise function.
	c.Check(inStrings(nil, "A"), Equals, false)
	c.Check(inStrings([]string{}, "A"), Equals, false)
	c.Check(inStrings([]string{"A"}, "A"), Equals, true)
	c.Check(inStrings([]string{"A", "B"}, "A"), Equals, true)
	c.Check(inStrings([]string{"A", "B", "C"}, "A"), Equals, true)

	// Exercise function.
	c.Check(inStrings(nil, "B"), Equals, false)
	c.Check(inStrings([]string{}, "B"), Equals, false)
	c.Check(inStrings([]string{"A"}, "B"), Equals, false)
	c.Check(inStrings([]string{"A", "B"}, "B"), Equals, true)
	c.Check(inStrings([]string{"A", "B", "C"}, "B"), Equals, true)

	// Exercise function.
	c.Check(inStrings(nil, "C"), Equals, false)
	c.Check(inStrings([]string{}, "C"), Equals, false)
	c.Check(inStrings([]string{"A"}, "C"), Equals, false)
	c.Check(inStrings([]string{"A", "B"}, "C"), Equals, false)
	c.Check(inStrings([]string{"A", "B", "C"}, "C"), Equals, true)
}
