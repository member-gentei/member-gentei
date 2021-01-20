module github.com/member-gentei/member-gentei/bot

go 1.14

require (
	cloud.google.com/go/firestore v1.3.0
	cloud.google.com/go/pubsub v1.3.1
	github.com/BurntSushi/toml v0.3.1
	github.com/Lukaesebrot/dgc v1.0.7-0.20200816224117-b4b1b682649a
	github.com/bwmarrin/discordgo v0.22.1-0.20201217190221-8d6815dde7ed
	github.com/deepmap/oapi-codegen v1.3.13
	github.com/getkin/kin-openapi v0.23.0
	github.com/google/go-cmp v0.5.2
	github.com/lthibault/jitterbug v0.0.0-20200313035244-37ff5f417161
	github.com/mark-ignacio/zerolog-gcp v0.2.0
	github.com/member-gentei/member-gentei/pkg v0.0.0-20201031063345-01759886af5b
	github.com/nicksnyder/go-i18n/v2 v2.1.1
	github.com/rs/zerolog v1.20.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	golang.org/x/text v0.3.3
	google.golang.org/api v0.35.0
	google.golang.org/grpc v1.33.2
)

replace github.com/member-gentei/member-gentei/pkg => ../pkg
