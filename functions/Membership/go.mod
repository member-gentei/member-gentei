module github.com/mark-ignacio/member-gentei/functions/Membership

go 1.14

require (
	cloud.google.com/go/firestore v1.3.0
	firebase.google.com/go v3.13.0+incompatible
	github.com/mark-ignacio/member-gentei/pkg v0.0.0-20201027004120-4d0e6608b2f8
	github.com/rs/zerolog v1.20.0
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
)

replace github.com/mark-ignacio/member-gentei/pkg => ../../pkg
