module github.com/member-gentei/member-gentei/bot

go 1.14

require (
	cloud.google.com/go/firestore v1.3.0
	cloud.google.com/go/pubsub v1.3.1
	github.com/Lukaesebrot/dgc v1.0.7-0.20200816224117-b4b1b682649a
	github.com/bwmarrin/discordgo v0.22.0
	github.com/deepmap/oapi-codegen v1.3.13
	github.com/getkin/kin-openapi v0.23.0
	github.com/karrick/tparse/v2 v2.8.1
	github.com/mark-ignacio/zerolog-gcp v0.2.0
	github.com/member-gentei/member-gentei/pkg v0.0.0-20201031063345-01759886af5b
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/rs/zerolog v1.20.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/zekroTJA/timedmap v0.0.0-20200518230343-de9b879d109a
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43 // indirect
	google.golang.org/api v0.35.0
	google.golang.org/grpc v1.33.2
)

replace github.com/member-gentei/member-gentei/pkg => ../pkg
