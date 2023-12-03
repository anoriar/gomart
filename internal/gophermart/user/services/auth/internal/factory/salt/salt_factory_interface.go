package salt

//go:generate mockgen -source=salt_factory_interface.go -destination=mock_salt_factory/salt_factory.go -package=mock_salt_factory
type SaltFactoryInterface interface {
	GenerateSalt() ([]byte, error)
}
