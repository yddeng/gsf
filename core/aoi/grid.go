package aoi

import "errors"

/*
  网格
*/
type grid struct {
	r, c     int32                       // row/col coordinate.
	entities map[interface{}]*GridEntity // AOI entities inside.
}

// Grid AOI manager.
type GridManager struct {
	lb, rt        Position  // left-bottom, right-top.
	gWidth, gHigh int32     // grid width.
	grids         [][]*grid // all grids.
}

func NewGrid(lb, rt Position, gWidth, gHigh int32) *GridManager {
	if gWidth <= 0 || gHigh <= 0 {
		panic(errors.New("invalid width or high"))
	}

	width, height := rt.X-lb.X, rt.Y-lb.Y
	r, c := width/gWidth, height/gHigh
	if width%gWidth > 0 {
		r += 1
	}
	if height%gHigh > 0 {
		c += 1
	}

	g := &GridManager{
		lb:     lb,
		rt:     rt,
		gWidth: gWidth,
		gHigh:  gHigh,
	}
	g.grids = make([][]*grid, r)
	for i := int32(0); i < r; i++ {
		g.grids[i] = make([]*grid, c)
		for j := int32(0); j < c; j++ {
			g.grids[i][j] = &grid{
				r:        i,
				c:        j,
				entities: map[interface{}]*GridEntity{},
			}
		}
	}

	return g
}

func (m *GridManager) Add(key interface{}, pos Position, user User) (Entity, error) {
	fixPos := m.fixPos(pos)
	grid, vector := m.getGridVectorByPos(fixPos)
	if user == nil {
		return nil, errors.New("nil user")
	}

	if _, ok := grid.entities[key]; ok {
		return nil, errors.New("entity already exist")
	}

	entity := &GridEntity{
		entity:  newEntity(key, fixPos, user),
		manager: m,
		grid:    grid,
		vector:  vector,
	}

	var entities []User
	enter := []User{user}
	grids := m.getNearGridsByPos(fixPos)
	for _, g := range grids {
		for _, e := range g.entities {
			entities = append(entities, e.user)
			e.user.OnAOIUpdate(enter, nil)
		}
	}

	grid.entities[key] = entity
	entity.user.OnAOIUpdate(entities, nil)

	return entity, nil
}

func (m *GridManager) Rem(e Entity) error {
	ee, ok := e.(*GridEntity)
	if !ok {
		return errors.New("invalid entity")
	}

	if ee.manager != m {
		return errors.New("object not belongs to manager")
	}

	delete(ee.grid.entities, ee.key)

	leave := []User{e.User()}
	grids := m.getNearGridsByPos(e.Position())
	for _, g := range grids {
		for _, v := range g.entities {
			v.user.OnAOIUpdate(nil, leave)
		}
	}

	ee.manager = nil
	ee.grid = nil
	return nil
}

func (m *GridManager) PosNearAOI(pos Position, distance int32) []User {

	lbp := m.fixPos(Position{
		X: pos.X - distance,
		Y: pos.Y - distance,
	})
	rtp := m.fixPos(Position{
		X: pos.X + distance,
		Y: pos.Y + distance,
	})

	lbg, _ := m.getGridVectorByPos(lbp)
	rtg, _ := m.getGridVectorByPos(rtp)

	grids := []*grid{}
	for r := lbg.r; r <= rtg.r; r++ {
		for c := lbg.c; c <= rtg.c; c++ {
			grids = append(grids, m.grids[r][c])
		}
	}

	ret := make([]User, 0)
	for _, g := range grids {
		for _, e := range g.entities {
			ret = append(ret, e.user)
		}
	}
	return ret
}

func (m *GridManager) fixPos(pos Position) Position {
	if pos.X < m.lb.X {
		pos.X = m.lb.X
	}
	if pos.X > m.rt.X {
		pos.X = m.rt.X
	}
	if pos.Y < m.lb.X {
		pos.Y = m.lb.X
	}
	if pos.Y > m.rt.Y {
		pos.Y = m.rt.Y
	}
	return pos
}

func (m *GridManager) getGridVectorByPos(pos Position) (*grid, int) {
	r := (pos.X - m.lb.X) / m.gWidth
	c := (pos.Y - m.lb.Y) / m.gHigh
	if r == int32(len(m.grids)) {
		r -= 1
	}
	if c == int32(len(m.grids[0])) {
		c -= 1
	}

	// 位置判断
	x := (pos.X - m.lb.X) % m.gWidth
	y := (pos.Y - m.lb.Y) % m.gHigh
	vectorX := x - m.gWidth/2
	vectorY := y - m.gHigh/2

	if vectorX == 0 && vectorY == 0 { // 0
		return m.grids[r][c], 0
	} else if vectorX > 0 && vectorY > 0 { // 1
		return m.grids[r][c], 1
	} else if vectorX < 0 && vectorY > 0 { // 2
		return m.grids[r][c], 2
	} else if vectorX < 0 && vectorY < 0 { // 3
		return m.grids[r][c], 3
	} else if vectorX > 0 && vectorY < 0 { // 4
		return m.grids[r][c], 4
	} else if vectorX > 0 && vectorY == 0 { // 5
		return m.grids[r][c], 5
	} else if vectorX == 0 && vectorY > 0 { // 6
		return m.grids[r][c], 6
	} else if vectorX < 0 && vectorY == 0 { // 7
		return m.grids[r][c], 7
	} else { //if vectorX == 0 && vectorY < 0 { // 8
		return m.grids[r][c], 8
	}

}

/*
	    6
	  2 | 1
	7 --0-- 5
	  3 | 4
		8
	判断 pos 所在格子的位置
*/

func (m *GridManager) getNearGridsByPos(pos Position) []*grid {

	curGrid, vector := m.getGridVectorByPos(pos)
	rc := []struct{ r, c int32 }{}
	switch vector {
	case 0:
		rc = []struct{ r, c int32 }{{curGrid.r, curGrid.c}}
	case 1:
		rc = []struct{ r, c int32 }{{curGrid.r, curGrid.c}, {curGrid.r + 1, curGrid.c}, {curGrid.r, curGrid.c + 1}, {curGrid.r + 1, curGrid.c + 1}}
	case 2:
		rc = []struct{ r, c int32 }{{curGrid.r, curGrid.c}, {curGrid.r - 1, curGrid.c}, {curGrid.r, curGrid.c + 1}, {curGrid.r - 1, curGrid.c + 1}}
	case 3:
		rc = []struct{ r, c int32 }{{curGrid.r, curGrid.c}, {curGrid.r - 1, curGrid.c}, {curGrid.r, curGrid.c - 1}, {curGrid.r - 1, curGrid.c - 1}}
	case 4:
		rc = []struct{ r, c int32 }{{curGrid.r, curGrid.c}, {curGrid.r + 1, curGrid.c}, {curGrid.r, curGrid.c - 1}, {curGrid.r + 1, curGrid.c - 1}}
	case 5:
		rc = []struct{ r, c int32 }{{curGrid.r, curGrid.c}, {curGrid.r + 1, curGrid.c}}
	case 6:
		rc = []struct{ r, c int32 }{{curGrid.r, curGrid.c}, {curGrid.r, curGrid.c + 1}}
	case 7:
		rc = []struct{ r, c int32 }{{curGrid.r, curGrid.c}, {curGrid.r - 1, curGrid.c}}
	default: // 8
		rc = []struct{ r, c int32 }{{curGrid.r, curGrid.c}, {curGrid.r, curGrid.c - 1}}
	}

	grids := []*grid{}
	maxR, maxC := int32(len(m.grids))-1, int32(len(m.grids[0]))-1
	for _, v := range rc {
		if v.r < 0 || v.r > maxR || v.c < 0 || v.c > maxC {
			continue
		}
		grids = append(grids, m.grids[v.r][v.c])
	}

	return grids
}

type GridEntity struct {
	entity
	manager *GridManager
	grid    *grid
	vector  int
}

func (e *GridEntity) Move(pos Position) error {
	fixPos := e.manager.fixPos(pos)
	if fixPos == e.pos {
		return nil
	}

	newGrid, newVector := e.manager.getGridVectorByPos(fixPos)
	if (newGrid.r == e.grid.r && newGrid.c == e.grid.c) && newVector == e.vector {
		return nil
	}

	oldGrid := e.grid
	// 不在当前格子
	if newGrid.r != e.grid.r || newGrid.c != e.grid.c {
		delete(oldGrid.entities, e.key)
	}

	oldGrids := e.manager.getNearGridsByPos(e.pos)
	newGrids := e.manager.getNearGridsByPos(fixPos)
	leaveGrids, enterGrids := oldNewGrids(oldGrids, newGrids)

	selfLeave := []User{e.user}
	othLeave := make([]User, 0)
	for _, g := range leaveGrids {
		for _, v := range g.entities {
			v.user.OnAOIUpdate(nil, selfLeave)
			othLeave = append(othLeave, v.user)
		}

	}

	selfEnter := []User{e.user}
	othEnter := make([]User, 0)
	for _, g := range enterGrids {
		for _, v := range g.entities {
			v.user.OnAOIUpdate(selfEnter, nil)
			othEnter = append(othEnter, v.user)
		}

	}

	newGrid.entities[e.key] = e
	e.grid = newGrid
	e.vector = newVector
	e.pos = fixPos
	e.user.OnAOIUpdate(othEnter, othLeave)

	return nil
}

func (e *GridEntity) TraverseAOI(fn func(u User) error) error {
	if fn == nil {
		return errors.New("nil fn")
	}

	grids := e.manager.getNearGridsByPos(e.pos)
	for _, g := range grids {
		if e.grid != g {
			for _, v := range g.entities {
				if err := fn(v.user); err != nil {
					return err
				}
			}
		} else {
			for _, v := range g.entities {
				if e != v {
					if err := fn(v.user); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func oldNewGrids(oldGrids, newGrids []*grid) ([]*grid, []*grid) {
redo:
	for i, og := range oldGrids {
		for j, ng := range newGrids {
			if og.c == ng.c && og.r == ng.r {
				oldGrids = append(oldGrids[:i], oldGrids[i+1:]...)
				newGrids = append(newGrids[:j], newGrids[j+1:]...)
				goto redo
			}
		}
	}
	return oldGrids, newGrids
}
