// +build prod

//go:generate go run github.com/go-bindata/go-bindata/go-bindata -fs -pkg chores -o assets.go templates/... static/...
package chores
