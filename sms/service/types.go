package service

import "context"

type Service interface {
	Send(ctx context.Context, tpl string, args []NamedArg, numbers ...string) error
}

type NamedArg struct {
	Name  string
	Value string
}
