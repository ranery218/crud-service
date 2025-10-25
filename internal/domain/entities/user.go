package entities

import "github.com/samber/mo"

type User struct {
	ID             string
	Username       string
	Email          string
	HashedPassword string
}

type UserAttrs struct {
	ID             string
	Username       string
	Email          string
	HashedPassword string
}

type UserFilterAttrs struct {
	ID       mo.Option[string]
	Email    mo.Option[string]
	Username mo.Option[string]
}
