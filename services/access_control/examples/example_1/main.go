package main

import (
	"fmt"
	"github.com/aeroideaservices/focus/services/access_control"
	"os"
	"path/filepath"
)

func main() {
	modulePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	roles := access_control.RolesFromFile(filepath.Join(modulePath, "examples", "example_1", "roles.json"))
	role := roles[0]
	actions := []*access_control.Action{
		access_control.NewAction("/cart", "POST"),
		access_control.NewAction("/catalog.categories.main-category", "GET"),
	}

	for _, action := range actions {
		switch role.HasAccess(action) {
		case true:
			fmt.Println(fmt.Sprintf("accessed to %s\n", action))
		default:
			fmt.Println(fmt.Sprintf("denied to %s\n", action))
		}
	}
}
