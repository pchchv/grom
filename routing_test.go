package grom

type Ctx struct{}

type routeTest struct {
	route string
	get   string
	vars  map[string]string
}
