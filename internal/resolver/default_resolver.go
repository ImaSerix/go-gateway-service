package resolver

import "net/http"

type MultiResolver struct {
	resolvers map[string]Resolver
}

func NewMultiResolver() *MultiResolver {

	r := &MultiResolver{
		resolvers: map[string]Resolver{},
	}

	r.Register("context", NewContextResolver())
	r.Register("route", NewRouterResolver())
	r.Register("query", NewQueryResolver())
	r.Register("header", NewHeaderResolver())

	return &MultiResolver{
		resolvers: map[string]Resolver{},
	}
}

func (dr *MultiResolver) Register(source string, resolver Resolver) {
	dr.resolvers[source] = resolver
}

func (dr *MultiResolver) Resolve(r *http.Request, source string, key string) (any, bool) {

	resolver, ok := dr.resolvers[source]
	if !ok {
		return nil, false
	}

	v, ok := resolver.Resolve(r, key)

	return v, ok
}
