package api

type PutArgs struct {
	Key   string
	Value string
}

type PutResult struct {
	Result bool
}

type GetArgs struct {
	Key string
}

type GetResult struct {
	Result string
}
