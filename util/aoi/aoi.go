package aoi

import (
	"errors"
)

type Position struct {
	X int32
	Y int32
}

// Define AOI user interface.
type User interface {
	// Call if when AOI has updated.
	OnAOIUpdate(enter, level []User)
}

// Define AOI entity interface.
type Entity interface {
	// The key to identify entity.
	Key() interface{}
	// The AOI position.
	Position() Position
	// The user of AOI entity.
	User() User
	// Move AOI Entity.
	Move(pos Position) error
	// Traverse the AOI.
	TraverseAOI(fn func(u User) error) error
}

type entity struct {
	key  interface{}
	pos  Position
	user User
}

func newEntity(key interface{}, pos Position, user User) entity {
	if user == nil {
		panic(errors.New("nil user"))
	}
	return entity{
		key:  key,
		pos:  pos,
		user: user,
	}
}

func (e *entity) Key() interface{} { return e.key }

func (e *entity) Position() Position { return e.pos }

func (e *entity) User() User { return e.user }

// Define AOI manager interface.
type Manager interface {
	// Add one entity.
	Add(key interface{}, pos Position, user User) (Entity, error)
	// Remove the entity.
	Rem(e Entity) error
	// Get Near Entity
	PosNearAOI(pos Position, distance int32) []User
}
