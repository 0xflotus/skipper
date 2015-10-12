package routing

import (
	"github.com/zalando/skipper/eskip"
	"github.com/zalando/skipper/filters"
	"net/http"
	"time"
)

type MatchingOptions uint

const (
	MatchingOptionsNone MatchingOptions = 0
	IgnoreTrailingSlash MatchingOptions = 1 << iota
)

func (o MatchingOptions) ignoreTrailingSlash() bool {
	return o&IgnoreTrailingSlash > 0
}

type DataClient interface {
	GetInitial() ([]*eskip.Route, error)
	GetUpdate() ([]*eskip.Route, []string, error)
}

type Options struct {
	FilterRegistry  filters.Registry
	MatchingOptions MatchingOptions
	PollTimeout     time.Duration
	DataClients     []DataClient
	UpdateBuffer    int
}

type Route struct {
	eskip.Route
	Scheme, Host string
	Filters      []filters.Filter
}

type Routing struct {
	getMatcher <-chan *matcher
}

func feedMatchers(updateBuffer int, current *matcher) (chan<- *matcher, <-chan *matcher) {
	// todo: use updateBuffer, when benchmarks show that it matters
	in := make(chan *matcher)
	out := make(chan *matcher, 0)

	go func() {
		for {
			select {
			case current = <-in:
			case out <- current:
			}
		}
	}()

	return in, out
}

func New(o Options) *Routing {
	initialMatcher, _ := newMatcher(nil, MatchingOptionsNone)
	matchersIn, matchersOut := feedMatchers(o.UpdateBuffer, initialMatcher)
	go receiveRouteMatcher(o, matchersIn)
	return &Routing{matchersOut}
}

func (r *Routing) Route(req *http.Request) (*Route, map[string]string) {
	m := <-r.getMatcher
	return m.match(req)
}
