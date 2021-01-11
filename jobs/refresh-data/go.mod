module github.com/member-gentei/member-gentei/jobs/refresh-data

go 1.15

require (
	github.com/bwmarrin/discordgo v0.22.0
	github.com/mark-ignacio/zerolog-gcp v0.3.0
	github.com/member-gentei/member-gentei/pkg v0.0.0-20201221020045-018f6414ee45
	github.com/rs/zerolog v1.20.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
)

replace github.com/member-gentei/member-gentei/pkg => ../../pkg
