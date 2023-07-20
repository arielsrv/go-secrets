package secrets

type Builder interface {
	Build() Store
}
