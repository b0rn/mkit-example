package config

// Constants related to the API
type apiConstants struct {
	REST string
	HTTP string
	NATS string
}

type dataStoreConstants struct {
	// We create the connection with nats in the datastore
	NATS  string
	SQL   string
	MONGO string
}

type dataServiceConstants struct {
	DB_DATA string
	MQ_DATA string
}

type useCaseConstants struct {
	MANAGE_USERS string
}

var API = apiConstants{
	REST: "rest",
	HTTP: "http",
	NATS: "nats",
}

var DATASERVICE = dataServiceConstants{
	DB_DATA: "dbData",
	MQ_DATA: "mqData",
}

var USECASE = useCaseConstants{
	MANAGE_USERS: "manageUsers",
}

var DATASTORE = dataStoreConstants{
	NATS:  "nats",
	SQL:   "sql",
	MONGO: "mongo",
}
