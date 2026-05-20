package resolver

import "net/http"

type MultiResolver struct {
	resolvers map[string]Resolver
}

func NewMultiResolver() *MultiResolver {
	return &MultiResolver{
		resolvers: map[string]Resolver{},
	}
}

func (dr *MultiResolver) Register(source string, resolver Resolver) {
	dr.resolvers[source] = resolver
}

func (dr *MultiResolver) Resolve(r *http.Request, key string) (any, bool) {

	s, _, err := getSourceAndKey(key)
	if err != nil {
		return nil, false
	}

	resolver, ok := dr.resolvers[s]
	if !ok {
		return nil, false
	}

	v, ok := resolver.Resolve(r, key)

	return v, ok
}
