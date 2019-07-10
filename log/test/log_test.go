package test

import (
	"testing"

	"github.com/jkaveri/goabs/log"
)

func TestLog(t *testing.T) {
	type testCase struct {
		name string
	}

	testFunc := func(c *testCase) func(t *testing.T) {
		return func(t1 *testing.T) {
			log.Log(
				log.Level(log.LevelDebug),
				log.Message("log\" message"),
				log.Fields(log.LogFields{
					"a": 10,
				}),
			)
		}
	}

	testCases := []*testCase{
		{
			name: "simple test case",
		},
	}

	for i := 0; i < len(testCases); i++ {
		c := testCases[i]
		t.Run(c.name, testFunc(c))
	}
}
