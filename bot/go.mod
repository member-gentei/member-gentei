module github.com/member-gentei/member-gentei/bot

go 1.14

require (
	cloud.google.com/go v0.75.0
	cloud.google.com/go/firestore v1.3.0
	cloud.google.com/go/pubsub v1.3.1
	github.com/BurntSushi/toml v0.3.1
	github.com/Lukaesebrot/dgc v1.0.7-0.20200816224117-b4b1b682649a
	github.com/bwmarrin/discordgo v0.22.1-0.20201217190221-8d6815dde7ed
	github.com/deepmap/oapi-codegen v1.3.13
	github.com/getkin/kin-openapi v0.23.0
	github.com/google/go-cmp v0.5.4
	github.com/lthibault/jitterbug v0.0.0-20200313035244-37ff5f417161
	github.com/mark-ignacio/zerolog-gcp v0.2.0
	github.com/member-gentei/member-gentei/pkg v0.0.0-20201031063345-01759886af5b
	github.com/nicksnyder/go-i18n/v2 v2.1.1
	github.com/rs/zerolog v1.20.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	go.opencensus.io v0.22.6 // indirect
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	golang.org/x/oauth2 v0.0.0-20210126194326-f9ce19ea3013 // indirect
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	golang.org/x/text v0.3.5
	google.golang.org/api v0.38.0
	google.golang.org/genproto v0.0.0-20210126160654-44e461bb6506
	google.golang.org/grpc v1.35.0
	google.golang.org/protobuf v1.25.0
)

replace github.com/member-gentei/member-gentei/pkg => ../pkg
