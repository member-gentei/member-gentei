module github.com/member-gentei/member-gentei/jobs/member-check

go 1.15

require (
	cloud.google.com/go/firestore v1.3.0
	github.com/j0hnsmith/connspy v0.0.0-20200203145744-79259064872c // indirect
	github.com/mark-ignacio/zerolog-gcp v0.2.0
	github.com/member-gentei/member-gentei/pkg v0.0.0-20201201174628-adf757afdb38
	github.com/mitchellh/go-homedir v1.1.0
	github.com/rs/zerolog v1.20.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
)

replace github.com/member-gentei/member-gentei/pkg => ../../pkg
