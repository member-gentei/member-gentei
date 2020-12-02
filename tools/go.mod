module github.com/member-gentei/member-gentei/tools

go 1.14

require (
	cloud.google.com/go/firestore v1.3.0
	firebase.google.com/go v3.13.0+incompatible
	firebase.google.com/go/v4 v4.1.0
	github.com/member-gentei/member-gentei/pkg v0.0.0-20201115025050-759c1329cf82
	github.com/mitchellh/go-homedir v1.1.0
	github.com/rs/zerolog v1.20.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	google.golang.org/api v0.35.0
	google.golang.org/grpc v1.33.2
)

replace github.com/member-gentei/member-gentei/pkg => ../pkg
