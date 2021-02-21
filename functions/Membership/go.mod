module github.com/member-gentei/member-gentei/functions/Membership

go 1.14

require (
	cloud.google.com/go/firestore v1.4.0
	github.com/member-gentei/member-gentei/pkg v0.0.0-20210213071734-d54e68722d0e
	github.com/rs/zerolog v1.20.0
	golang.org/x/oauth2 v0.0.0-20210210192628-66670185b0cd
	google.golang.org/api v0.40.0
)

replace github.com/member-gentei/member-gentei/pkg => ../../pkg
