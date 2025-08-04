package interfaces

type Builder interface {
	Config() ConfigBuilder
	//Registry() RegistryBuilder
	//Service() ServerBuilder
	//Middleware() MiddlewareBuilder
}
