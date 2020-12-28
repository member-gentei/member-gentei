module github.com/member-gentei/member-gentei/functions/DCYTAccount

go 1.15

require (
	cloud.google.com/go/firestore v1.3.0
	github.com/member-gentei/member-gentei/pkg v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.20.0
	google.golang.org/grpc v1.33.2
)

replace github.com/member-gentei/member-gentei/pkg => ../../pkg
