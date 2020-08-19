package aoi

import (
	"math/rand"
	"testing"
)

func TestGrid(t *testing.T) {
	radius := int32(10)
	aoiManager := NewGrid(Position{0, 0}, Position{1024, 1024}, radius, radius)

	obj1 := &testUser{ID: 1}
	obj2 := &testUser{ID: 2}
	obj3 := &testUser{ID: 3}
	obj4 := &testUser{ID: 4}

	obj1.aoi, _ = aoiManager.Add(obj1.ID, Position{5, 15}, obj1)
	obj2.aoi, _ = aoiManager.Add(obj2.ID, Position{25, 15}, obj2)
	obj3.aoi, _ = aoiManager.Add(obj3.ID, Position{45, 5}, obj3)
	obj4.aoi, _ = aoiManager.Add(obj4.ID, Position{5, 30}, obj4)

	obj4.aoi.Move(Position{15, 20})
	obj4.aoi.Move(Position{55, 19})
	aoiManager.Rem(obj4.aoi)
}

func BenchmarkGridManager_Add(b *testing.B) {
	var err error
	manager := NewGrid(Position{0, 0}, Position{width, height}, radius*2, radius*2)

	users := map[interface{}]*testUser{}
	for i := 0; i < maxUser; i++ {
		user := &testUser{ID: i}
		user.aoi, err = manager.Add(user.ID, Position{X: rand.Int31n(width), Y: rand.Int31n(height)}, user)
		if err != nil {
			b.Fatal(err)
		}
		users[user.ID] = user
	}

	b.ResetTimer()
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

func BenchmarkGridEntity_Move(b *testing.B) {
	var err error
	manager := NewGrid(Position{0, 0}, Position{width, height}, radius*2, radius*2)
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
		if err = user.aoi.Move(pos); err != nil {
			b.Fatal(err)
		}
	}
	//1200
}
