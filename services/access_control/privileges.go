package access_control

import "strings"

const (
	Nil      Access = iota // Не установлено
	Denied                 // Доступ запрещен
	Accessed               // Доступ разрешен
)

type Access uint

func (a Access) String() string {
	switch a {
	case Nil:
		return "null"
	case Denied:
		return "denied"
	case Accessed:
		return "accessed"
	default:
		return ""
	}
}

type Privilege struct {
	access     Access
	privileges *Privileges
}

func NewPrivilege(access Access) *Privilege {
	return &Privilege{access: access}
}

type Privileges map[string]*Privilege

func (p Privileges) getPrivilege(code string) (privilege *Privilege) {
	privilege = p[code]
	if privilege == nil {
		privilege = p["*"]
	}

	return privilege
}

// Получение уровня доступа к действию
// Если доступ к действию не описан, происходит попытка получения уровня доступа для "*"
func (p Privileges) access(code string) Access {
	privilege := p[code]
	if privilege == nil || privilege.access == Nil {
		privilege = p["*"]
	}

	if privilege == nil {
		return Nil
	}

	return privilege.access
}

func (p Privileges) accessByAction(action string) Access {
	actionParts := strings.SplitN(action, ".", 2)
	switch len(actionParts) {
	case 1:
		return p.access(actionParts[0])
	case 2:
		// Получаем привилегии текущего действия
		privilege := p.getPrivilege(actionParts[0])
		access := Nil
		if privilege != nil {
			access = privilege.access
		}
		if privilege != nil && privilege.privileges != nil {
			// Получаем уровень доступа дочерней привилегии. Если непустой - возвращаем его
			if childAccess := privilege.privileges.accessByAction(actionParts[1]); childAccess != Nil {
				return childAccess
			}
		}

		// Получаем привилегии *. Если дочерних привилегий нет - возвращаем текущий доступ
		allPrivilege := p.getPrivilege("*")
		// Если не найдена - возвращаем доступ привилегии текущего действия
		if allPrivilege == nil {
			return access
		}
		if allPrivilege.privileges != nil {
			// Получаем уровень доступа дочерней привилегии. Если непустой - возвращаем его
			if allChildAccess := allPrivilege.privileges.accessByAction(actionParts[1]); allChildAccess != Nil {
				return allChildAccess
			}
		}

		// Если доступ привилегии текущего действия непустой - возвращаем его
		if access != Nil {
			return access
		}

		// Возвращаем доступ привилегии действия "*"
		return allPrivilege.access
	default:
		return Denied
	}
}

func (p *Privileges) append(action string) {
	access := Accessed
	if strings.HasPrefix(action, "!") {
		access = Denied
		action = strings.TrimPrefix(action, "!")
	}

	actionParts := strings.SplitN(action, ".", 2)
	switch len(actionParts) {
	case 1:
		if privilege, ok := (*p)[actionParts[0]]; ok {
			privilege.access = access
			return
		}
		(*p)[actionParts[0]] = NewPrivilege(access)
	case 2:
		childPrivilegeString := actionParts[1]
		if access == Denied {
			childPrivilegeString = "!" + childPrivilegeString
		}
		var ok bool
		var privilege *Privilege
		if privilege, ok = (*p)[actionParts[0]]; !ok {
			privilege = NewPrivilege(Nil)
			(*p)[actionParts[0]] = privilege
		}
		(*p)[actionParts[0]] = privilege

		if privilege.privileges == nil {
			privilege.privileges = &Privileges{}
		}
		privilege.privileges.append(childPrivilegeString)
	}
}
