// router.go


package router


type Route struct {
	Address *IPv4
	Masked  *IPv4
	Prefix  uint8
	To      *IPv4
}


func new_route(address *IPv4, prefix uint8, to *IPv4) *Route {
	if prefix > 32 {
		return nil
	}
	return &Route{
		Address: address,
		Masked:  address.mask_with_prefix(prefix),
		Prefix:  prefix,
		To:      to,
	}
}


type Router interface {
	// Add a route to the router
	Add(Route)
	// Deprecate a route
	Drop(Route)
	// Deprecate a gateway
	DropAllTo(IPv4)
	// Get the route with longest match; nil if  missing
	Get(IPv4) *IPv4
}



//
