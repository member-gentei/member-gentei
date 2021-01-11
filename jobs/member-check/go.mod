module github.com/member-gentei/member-gentei/jobs/member-check

go 1.15

require (
	cloud.google.com/go/firestore v1.4.0 // indirect
	cloud.google.com/go/pubsub v1.9.1
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/mark-ignacio/zerolog-gcp v0.3.0
	github.com/member-gentei/member-gentei/pkg v0.0.0-20201216034645-22b33a8d8da8
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.4.0 // indirect
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/rs/zerolog v1.20.0
	github.com/spf13/afero v1.5.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.7.1
	golang.org/x/net v0.0.0-20201216054612-986b41b23924 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/member-gentei/member-gentei/pkg => ../../pkg
