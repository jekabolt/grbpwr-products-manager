module github.com/jekabolt/grbpwr-manager

go 1.16

require (
	github.com/caarlos0/env/v6 v6.4.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/httprate v0.5.0 // indirect
	github.com/go-chi/render v1.0.1
	github.com/go-kivik/couchdb/v4 v4.0.0-20200818191020-c997633e0a27
	github.com/rs/cors v1.7.0
	github.com/rs/zerolog v1.19.0
	github.com/tidwall/buntdb v1.2.4
	github.com/tidwall/gjson v1.8.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/text v0.3.2 // indirect
)

//replace gitlab.com/dvision/go-cri => ../../dvision/go-cri
