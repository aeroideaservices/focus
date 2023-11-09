package access_control

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Role struct {
	name       string
	privileges *Privileges
}

func NewRole(name string, privilegesStrings []string) *Role {
	privileges := &Privileges{}
	for _, ps := range privilegesStrings {
		privileges.append(ps)
	}
	return &Role{name: name, privileges: privileges}
}

func RolesFromFile(filename string) (roles []*Role) {
	if filepath.Ext(filename) != ".json" {
		panic("wrong roles file extension")
	}

	fileData, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	rawRoles := make(map[string][]string)
	err = json.Unmarshal(fileData, &rawRoles)
	if err != nil {
		panic(err)
	}

	for roleName, actions := range rawRoles {
		roles = append(roles, NewRole(roleName, actions))
	}

	return roles
}

func (r Role) HasAccess(action *Action) bool {
	return r.privileges.accessByAction(action.String()) == Accessed
}
