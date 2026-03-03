package metadata

// --- Engine Metadata Definitions ---

// Scope defines the isolation level for components.
type Scope string

const (
	GlobalScope Scope = "global"
	ServerScope Scope = "server"
	ClientScope Scope = "client"
)

// Category defines the functional type of component.
type Category string

const (
	CategoryInfrastructure Category = "infrastructure"
	CategoryLogger         Category = "logger"
	CategoryRegistry       Category = "registry"
	CategoryClient         Category = "client"
	CategoryServer         Category = "server"
	CategoryMiddleware     Category = "middleware"
	CategoryDatabase       Category = "database"
	CategoryCache          Category = "cache"
	CategoryObjectStore    Category = "objectstore"
	CategoryQueue          Category = "queue"
	CategoryTask           Category = "task"
	CategoryMail           Category = "mail"
	CategoryStorage        Category = "storage"
)

// Standard Priorities (Lower values execute earlier)
const (
	PriorityInfrastructure int = 100
	PriorityRegistry       int = 200
	PriorityStorage        int = 300
	PriorityClientStack    int = 400
	PriorityServerStack    int = 500
)
