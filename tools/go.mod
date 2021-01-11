module github.com/member-gentei/member-gentei/tools

go 1.14

require (
	cloud.google.com/go/pubsub v1.3.1
	firebase.google.com/go v3.13.0+incompatible
	github.com/member-gentei/member-gentei/pkg v0.0.0-20201115025050-759c1329cf82
	github.com/rs/zerolog v1.20.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	google.golang.org/api v0.35.0
)

replace github.com/member-gentei/member-gentei/pkg => ../pkg
