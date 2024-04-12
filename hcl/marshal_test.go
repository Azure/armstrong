package hcl_test

import (
	"testing"

	"github.com/azure/armstrong/hcl"
)

func Test_MarshalIndent(t *testing.T) {
	testcases := []struct {
		input  interface{}
		expect string
	}{
		{
			input:  nil,
			expect: "null",
		},
		{
			input:  "test",
			expect: `"test"`,
		},
		{
			input:  1,
			expect: "1",
		},
		{
			input:  true,
			expect: "true",
		},
		{
			input: []interface{}{"test", 1, true},
			expect: `[
  "test",
  1,
  true,
]`,
		},
		{
			input: map[string]interface{}{
				"test": "test",
				"test1": map[string]interface{}{
					"test2": "test2",
				},
			},
			expect: `{
  test = "test"
  test1 = {
    test2 = "test2"
  }
}`,
		},
		{
			input: map[string]interface{}{
				"/test": "test",
				"2test": map[string]interface{}{
					"${local.test}": "${local.value}",
				},
			},
			expect: `{
  "/test" = "test"
  "2test" = {
    (local.test) = local.value
  }
}`,
		},
	}

	for _, tc := range testcases {
		t.Logf("input: %v", tc.input)
		output := hcl.MarshalIndent(tc.input, "", "  ")
		if tc.expect != output {
			t.Fatalf("expect %s but got %s", tc.expect, output)
		}
	}
}
