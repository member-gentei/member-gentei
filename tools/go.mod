module github.com/member-gentei/member-gentei/tools

go 1.14

require (
	cloud.google.com/go/firestore v1.3.0
	cloud.google.com/go/pubsub v1.3.1
	firebase.google.com/go v3.13.0+incompatible
	github.com/VividCortex/ewma v1.1.1 // indirect
	github.com/golang/protobuf v1.5.1 // indirect
	github.com/member-gentei/member-gentei/pkg v0.0.0-20201115025050-759c1329cf82
	github.com/rs/zerolog v1.20.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/vbauerster/mpb v3.4.0+incompatible
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/net v0.0.0-20210316092652-d523dce5a7f4 // indirect
	golang.org/x/oauth2 v0.0.0-20210313182246-cd4f82c27b84
	golang.org/x/sys v0.0.0-20210320140829-1e4c9ba3b0c4 // indirect
	google.golang.org/api v0.42.0
	google.golang.org/genproto v0.0.0-20210319143718-93e7006c17a6 // indirect
	google.golang.org/grpc v1.36.0
)

replace github.com/member-gentei/member-gentei/pkg => ../pkg
