// prefix_tree_test.go


package router


import (
	"testing"
)


func test_comprehensive(t *testing.T) {
	// Initialize routes in table
	pt := new_prefix_tree()
	rt1 := &IPv4{1, 1, 1, 1}
	rt2 := &IPv4{2, 2, 2, 2}
	rt3 := &IPv4{3, 3, 3, 3}
	rt4 := &IPv4{4, 4, 4, 4}
	to_add := []*Route{
		new_route(&IPv4{192, 168, 0, 0}, 16, rt1),
		new_route(&IPv4{192, 168, 0, 0}, 18, rt2),
		new_route(&IPv4{192, 168, 64, 0}, 18, rt3),
		new_route(&IPv4{192, 168, 128, 0}, 18, rt4),
	}
	for _, r := range to_add {
		pt.Add(*r)
	}
	// Check that queries against the table are answered correctly
	cases := []struct {
		Name     string
		Input    *IPv4
		Expected *IPv4
	}{
		{
			"192.168.0.1 goes into the /18 not the /16",
			&IPv4{192, 168, 0, 1},
			rt2,
		},
		{
			"192.168.127.255 goes to second block of /18",
			&IPv4{192, 168, 127, 255},
			rt3,
		},
		{
			"Last subnet of /18 unrouted; goes to /16",
			&IPv4{192, 168, 192, 1},
			rt1,
		},
		{
			"Unknown gets unrouted without default gateway",
			&IPv4{10, 0, 0, 0},
			nil,
		},
	}
	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			got := pt.Get(*test.Input)
			if got == nil && test.Expected == nil {
				return
			} else if got != nil && test.Expected == nil {
				t.Errorf("Expected nil but got %v", *got)
			} else if got == nil && test.Expected != nil {
				t.Errorf("Expected %v but got nil", *test.Expected)
			} else if *got != *test.Expected {
				t.Errorf("expected %v got %v", *test.Expected, *got)
			}
		})
	}
	// Add a default gateway and drop a route
	default_gateway := *new_route(&IPv4{0, 0, 0, 0}, 0, &IPv4{9, 0, 0, 0})
	pt.Add(default_gateway)
	// Remove first subnet of /18
	pt.Drop(*to_add[1])
	cases = []struct {
		Name     string
		Input    *IPv4
		Expected *IPv4
	}{
		{
			"192.168.0.1 now goes into the /16",
			&IPv4{192, 168, 0, 1},
			rt2,
		},
		{
			"192.168.127.255 still goes to second block of /18",
			&IPv4{192, 168, 127, 255},
			rt3,
		},
		{
			"Last subnet of /18 unrouted; goes to /16",
			&IPv4{192, 168, 192, 1},
			rt1,
		},
		{
			"Unknown now gets through default gateway",
			&IPv4{10, 0, 0, 0},
			default_gateway.To,
		},
	}
	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			got := pt.Get(*test.Input)
			if got == nil && test.Expected == nil {
				return
			} else if got != nil && test.Expected == nil {
				t.Errorf("Expected nil but got %v", *got)
			} else if got == nil && test.Expected != nil {
				t.Errorf("Expected %v but got nil", *test.Expected)
			} else if *got != *test.Expected {
				t.Errorf("expected %v got %v", *test.Expected, *got)
			}
		})
	}
}



//
