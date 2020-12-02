module github.com/mark-ignacio/member-gentei/functions/DCYTAccount

go 1.15

require (
	cloud.google.com/go/firestore v1.3.0
	firebase.google.com/go v3.13.0+incompatible
	github.com/rs/zerolog v1.20.0
	google.golang.org/grpc v1.30.0
)

replace github.com/mark-ignacio/member-gentei/pkg => ../../pkg
