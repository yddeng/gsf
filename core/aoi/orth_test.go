package aoi

import (
	"math/rand"
	"testing"
)

const (
	maxUser = 5000
	width   = int32(10240)
	height  = int32(10240)
	radius  = 10
)

type testUser struct {
	ID  int
	aoi Entity
}

func (u *testUser) OnAOIUpdate(enter, leave []User) {
	//for _, v := range enter {
	//	oth := v.(*testUser)
	//	log.Printf("%v enter %v's aoi\n", oth.ID, u.ID)
	//}
	//
	//for _, v := range leave {
	//	oth := v.(*testUser)
	//	log.Printf("%v leave %v's aoi\n", oth.ID, u.ID)
	//}
}

func TestOrth(t *testing.T) {
	radius := int32(10)
	scene := NewOrth(radius, radius)

	obj1 := &testUser{ID: 1}
	obj2 := &testUser{ID: 2}
	obj3 := &testUser{ID: 3}
	obj4 := &testUser{ID: 4}

	obj1.aoi, _ = scene.Add(obj1.ID, Position{-0, 1}, obj1)
	obj2.aoi, _ = scene.Add(obj2.ID, Position{-2, 4}, obj2)
	obj3.aoi, _ = scene.Add(obj3.ID, Position{-1, 3}, obj3)
	obj4.aoi, _ = scene.Add(obj4.ID, Position{1, 2}, obj4)

	obj2.aoi.Move(Position{0, 12})

	scene.Rem(obj2.aoi)
}

func BenchmarkOrthManager_Add(b *testing.B) {
	var err error
	manager := NewOrth(radius*4, radius*4)
	users := map[interface{}]*testUser{}

	for i := 0; i < maxUser; i++ {
		user := &testUser{ID: i}
		user.aoi, err = manager.Add(user.ID, Position{X: rand.Int31n(width), Y: rand.Int31n(height)}, user)
		if err != nil {
			b.Fatal(err)
		}
		users[user.ID] = user
	}

	for i := 0; i < b.N; i++ {
		id := int(rand.Int31n(maxUser))
		u := users[id]
		if err = manager.Rem(u.aoi); err != nil {
			b.Fatal(err)
		}
		u.aoi, err = manager.Add(u.ID, Position{X: rand.Int31n(width), Y: rand.Int31n(height)}, u)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkOrthEntity_Move(b *testing.B) {
	var err error
	manager := NewOrth(radius*4, radius*4)
	users := make([]*testUser, 0, maxUser)

	for i := 0; i < maxUser; i++ {
		user := &testUser{ID: i}
		user.aoi, err = manager.Add(user.ID, Position{X: rand.Int31n(width), Y: rand.Int31n(height)}, user)
		if err != nil {
			b.Fatal(err)
		}
		users = append(users, user)
	}

	step := int32(radius * 2)
	dir := []int32{-1, 0, 1}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := users[i%maxUser]
		pos := user.aoi.Position()
		pos.X += step * dir[rand.Int31n(3)]
		if pos.X < 0 {
			pos.X = 0
		}
		if pos.X > width {
			pos.X = width
		}
		pos.Y += step * dir[rand.Int31n(3)]
		if pos.Y < 0 {
			pos.Y = 0
		}
		if pos.Y > height {
			pos.Y = height
		}
		if err := user.aoi.Move(pos); err != nil {
			b.Fatal(err)
		}
	}

	// 1700
}
