package salt

//go:generate mockgen -source=salt_factory_interface.go -destination=mock/salt_factory.go -package=mock
type SaltFactoryInterface interface {
	GenerateSalt() ([]byte, error)
}
