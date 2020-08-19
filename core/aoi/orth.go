package aoi

import (
	"container/list"
	"errors"
	"math"
)

/*
  十字链表
*/
// Orthogonal List AOI manager.
type OrthManager struct {
	entities      map[interface{}]*OrthEntity
	gWidth, gHigh int32
	xAxis         list.List // x axis link-list. y sort
}

func NewOrth(gWidth, gHigh int32) *OrthManager {
	return &OrthManager{
		entities: map[interface{}]*OrthEntity{},
		gHigh:    gHigh,
		gWidth:   gWidth,
	}
}

func (m *OrthManager) Add(key interface{}, pos Position, user User) (Entity, error) {
	if _, ok := m.entities[key]; ok {
		return nil, errors.New("entity already exist")
	}

	if user == nil {
		return nil, errors.New("nil user")
	}

	var (
		entity = &OrthEntity{
			entity:   newEntity(key, pos, user),
			manager:  m,
			entities: map[interface{}]*OrthEntity{},
		}
		xNode *list.Element
	)

	node := m.xAxis.Front()
	for node != nil && pos.X > node.Value.(*OrthEntity).Position().X {
		node = node.Next()
	}
	if node == nil {
		xNode = m.xAxis.PushBack(entity)
		//log.Printf("%v insert at back of x\n", entity.Key())
	} else {
		xNode = m.xAxis.InsertBefore(entity, node)
		//log.Printf("%v insert before %v at x\n", entity.Key(), node.Value.(*OrthEntity).Key())
	}

	entity.xNode = xNode
	m.entities[entity.Key()] = entity
	entity.updateAOI()

	return entity, nil
}

func (m *OrthManager) Rem(e Entity) error {
	ee, ok := e.(*OrthEntity)
	if !ok {
		return errors.New("invalid entity")
	}

	if ee.manager != m {
		return errors.New("object not belongs to manager")
	}

	for _, v := range ee.entities {
		v.remEntity(ee)
	}
	ee.entities = nil

	m.xAxis.Remove(ee.xNode)
	ee.xNode = nil

	ee.manager = nil
	delete(m.entities, e.Key())
	return nil
}

func (m *OrthManager) PosNearAOI(pos Position, distance int32) []User {
	ret := make([]User, 0)
	node := m.xAxis.Front()
	for node != nil {
		if math.Abs(float64(pos.X-node.Value.(*OrthEntity).pos.X)) <= float64(distance) &&
			math.Abs(float64(pos.Y-node.Value.(*OrthEntity).pos.Y)) <= float64(distance) {
			ret = append(ret, node.Value.(*OrthEntity).user)
		}
		node = node.Next()
	}

	return ret
}

// Orthogonal List AOI entity.
type OrthEntity struct {
	entity
	manager  *OrthManager
	xNode    *list.Element               // 链表节点
	entities map[interface{}]*OrthEntity // AOI Entities
}

func (e *OrthEntity) Move(pos Position) error {
	oldPos := e.Position()

	// 坐标未改变
	if oldPos == pos {
		return nil
	}

	xNode := e.xNode

	// 更新 x 链
	if pos.X != oldPos.X {
		node := xNode
		if pos.X < oldPos.X {
			for node.Prev() != nil && node.Prev().Value.(*OrthEntity).Position().X > pos.X {
				node = node.Prev()
			}
			if node != xNode {
				e.manager.xAxis.MoveBefore(xNode, node)
			}
		} else {
			for node.Next() != nil && node.Next().Value.(*OrthEntity).Position().X < pos.X {
				node = node.Next()
			}
			if node != xNode {
				e.manager.xAxis.MoveAfter(xNode, node)
			}
		}
	}

	e.pos = pos
	e.updateAOI()

	return nil
}

func (e *OrthEntity) TraverseAOI(fn func(u User) error) error {
	if fn == nil {
		return errors.New("nil fn")
	}

	for _, v := range e.entities {
		if err := fn(v.User()); err != nil {
			return nil
		}
	}

	return nil
}

func (e *OrthEntity) updateAOI() {
	xNode := e.xNode
	pos := e.pos
	width := e.manager.gWidth
	high := e.manager.gHigh

	entities := map[interface{}]*OrthEntity{}
	enter := make([]User, 0)

	// left
	node := xNode.Prev()
	for node != nil {
		oth := node.Value.(*OrthEntity)
		if pos.X-oth.Position().X > width {
			break
		}

		if math.Abs(float64(pos.Y-oth.Position().Y)) <= float64(high) {
			if _, ok := entities[oth.Key()]; !ok {
				entities[oth.Key()] = oth
				if _, ok := e.entities[oth.key]; !ok {
					enter = append(enter, oth.user)
					oth.addEntity(e)
				}
				delete(e.entities, oth.key)
			} else {
				panic("repeatedly")
			}
		}

		node = node.Prev()
	}

	// right
	node = xNode.Next()
	for node != nil {
		oth := node.Value.(*OrthEntity)
		if oth.Position().X-pos.X > width {
			break
		}

		if math.Abs(float64(pos.Y-oth.Position().Y)) <= float64(high) {
			if _, ok := entities[oth.Key()]; !ok {
				entities[oth.Key()] = oth
				if _, ok := e.entities[oth.key]; !ok {
					enter = append(enter, oth.user)
					oth.addEntity(e)
				}
				delete(e.entities, oth.key)
			} else {
				panic("repeatedly")
			}
		}

		node = node.Next()
	}

	leave := make([]User, len(e.entities))
	i := 0
	for _, v := range e.entities {
		leave[i] = v.user
		v.remEntity(e)
		i++
	}

	e.entities = entities
	e.user.OnAOIUpdate(enter, leave)
}

func (e *OrthEntity) addEntity(oth *OrthEntity) {
	if _, ok := e.entities[oth.key]; !ok {
		e.entities[oth.key] = oth
		enter := []User{oth.user}
		e.user.OnAOIUpdate(enter, nil)
	} else {
		panic(errors.New("entity repeated"))
	}
}

func (e *OrthEntity) remEntity(oth *OrthEntity) {
	if _, ok := e.entities[oth.key]; ok {
		delete(e.entities, oth.key)
		leave := []User{oth.user}
		e.user.OnAOIUpdate(nil, leave)
	} else {
		panic(errors.New("entity not exist"))
	}
}
