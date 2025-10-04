package graphql

type rootResolver struct {
}

type Initiator func(r *rootResolver) *rootResolver

func New() Initiator {
	return func(r *rootResolver) *rootResolver {
		return r
	}
}

func (i Initiator) Build() *rootResolver {
	return i(&rootResolver{})
}
