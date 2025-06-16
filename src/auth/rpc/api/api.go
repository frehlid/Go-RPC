package api

import "local/auth/types"

type SignupArgs struct {
	Id       string
	Password string
}

type SignupResult struct {
	Result bool
}

type LoginArgs struct {
	Id       string
	Password string
}

type LoginResult struct {
	Cap types.UserCap
}

type GetIdArgs struct {
	Cap types.UserCap
}

type GetIdResult struct {
	Id string
}
