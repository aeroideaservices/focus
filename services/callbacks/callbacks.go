package callbacks

import "github.com/google/uuid"

type Callbacks struct {
	AfterCreate func(ids ...uuid.UUID)
	AfterUpdate func(ids ...uuid.UUID)
	AfterDelete func(ids ...uuid.UUID)
}

func (c Callbacks) GoAfterCreate(ids ...uuid.UUID) {
	if c.AfterCreate != nil {
		go c.AfterCreate(ids...)
	}
}

func (c Callbacks) GoAfterUpdate(ids ...uuid.UUID) {
	if c.AfterUpdate != nil {
		go c.AfterUpdate(ids...)
	}
}

func (c Callbacks) GoAfterDelete(ids ...uuid.UUID) {
	if c.AfterDelete != nil {
		go c.AfterDelete(ids...)
	}
}
