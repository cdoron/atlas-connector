module github.com/cdoron/atlas-connector

go 1.17

require github.com/cdoron/datacatalog-go v0.0.0-20220601124730-cdb9bca2d7db

require (
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/go-resty/resty/v2 v2.7.0 // indirect
	github.com/gorilla/handlers v1.5.1 // indirect
	github.com/gorilla/mux v1.7.3 // indirect
	golang.org/x/net v0.0.0-20211029224645-99673261e6eb // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/cdoron/datacatalog-go => ./api
