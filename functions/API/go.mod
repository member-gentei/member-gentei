module github.com/member-gentei/member-gentei/functions/API

go 1.14

require (
	cloud.google.com/go/firestore v1.3.0
	github.com/deepmap/oapi-codegen v1.3.13
	github.com/getkin/kin-openapi v0.26.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/member-gentei/member-gentei/pkg v0.0.0-20201027004120-4d0e6608b2f8
	github.com/rs/zerolog v1.20.0
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	google.golang.org/api v0.33.0
	google.golang.org/grpc v1.31.1
)

replace github.com/member-gentei/member-gentei/pkg => ../../pkg
