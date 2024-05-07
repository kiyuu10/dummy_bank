package gapi

import (
	db "github.com/kiyuu10/simplebank/db/sqlc"
	"github.com/kiyuu10/simplebank/pb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
	}
}
