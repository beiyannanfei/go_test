package main
//https://godoc.org/gopkg.in/mgo.v2   文档地址
import (
	"time"
	"labix.org/v2/mgo/bson"
	"labix.org/v2/mgo"
	"fmt"
)

type User61 struct {
	Id_       bson.ObjectId `bson:"_id"`
	Name      string        `bson:"name"`
	Age       int           `bson:"age"`
	JonedAt   time.Time     `bson:"joned_at"`
	Interests []string      `bson:"interests"`
}

func main() {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		fmt.Println("connect mongo err:", err)
		return
	}
	fmt.Println("connect mongo success")
	defer session.Close()

	//connect db
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("go_test")

	//switch collection
	c := db.C("people")

	//insertOne(c)
	//insertMulti(c)
	//findOne(c)
	//findAll(c)
	//update(c)
	remove(c)
}

func remove(c *mgo.Collection) {
	id := bson.NewObjectId()
	user := User61{
		Id_:       id,
		Name:      "test",
		Age:       30,
		JonedAt:   time.Now(),
		Interests: []string{"node", "php", "go"},
	}
	err := c.Insert(&user)

	if err != nil {
		fmt.Println("insertOne err:", err)
		return
	}
	fmt.Println("insertOne success")

	err = c.Remove(bson.M{"name": "test"})
	if err != nil {
		fmt.Println("Remove err:", err)
		return
	}
	fmt.Println("Remove success")

	c.Insert(&user)

	err = c.RemoveId(id)
	if err != nil {
		fmt.Println("RemoveId err:", err)
		return
	}
	fmt.Println("RemoveId success")

	id1 := bson.NewObjectId()
	id2 := bson.NewObjectId()
	user1 := User61{
		Id_:       id1,
		Name:      "test",
		Age:       30,
		JonedAt:   time.Now(),
		Interests: []string{"node", "php", "go"},
	}
	user2 := User61{
		Id_:       id2,
		Name:      "test",
		Age:       30,
		JonedAt:   time.Now(),
		Interests: []string{"node", "php", "go"},
	}
	user3 := User61{
		Id_:       bson.NewObjectId(),
		Name:      "test",
		Age:       30,
		JonedAt:   time.Now(),
		Interests: []string{"node", "php", "go"},
	}
	c.Insert(&user1, &user2, &user3)
	changeInfo, err := c.RemoveAll(bson.M{"name": "test"})
	if err != nil {
		fmt.Println("RemoveAll err:", err)
		return
	}
	fmt.Println("RemoveAll success changeInfo:", changeInfo)
}

func update(c *mgo.Collection) {
	err := c.Update(
		bson.M{"name": "AAAAA"},
		bson.M{"$push": bson.M{
			"interests": "Golang",
		}},
	)
	if err != nil {
		fmt.Println("Update err:", err)
		return
	}
	fmt.Println("update success")

	changeInfo, err := c.UpdateAll(
		nil,
		bson.M{"$pull": bson.M{
			"interests": "a",
		}},
	)
	if err != nil {
		fmt.Println("UpdateAll err:", err)
		return
	}
	fmt.Println("update success changeInfo:", changeInfo)

	err = c.UpdateId(
		bson.ObjectIdHex("5bbd68c5bf9a24146e000001"),
		bson.M{"$inc": bson.M{
			"age": 1,
		}},
	)
	if err != nil {
		fmt.Println("UpdateId err:", err)
		return
	}
	fmt.Println("UpdateId success")

	changeInfo, err = c.Upsert(
		bson.M{"name": "CCCCC"},
		bson.M{"$set": bson.M{
			"age":       20,
			"interests": []string{"10", "20"},
		}},
	)
	if err != nil {
		fmt.Println("Upsert err:", err)
		return
	}
	fmt.Println("Upsert success changeInfo:", changeInfo)

	changeInfo, err = c.UpsertId(
		bson.NewObjectId(),
		bson.M{
			"$set": bson.M{
				"age":       25,
				"interests": []string{"100", "200"},
			},
		},
	)
	if err != nil {
		fmt.Println("UpsertId err:", err)
		return
	}
	fmt.Println("UpsertId success changeInfo:", changeInfo)
}

func findAll(c *mgo.Collection) {
	var users []User61
	c.Find(nil).All(&users)
	fmt.Printf("findAll count: %v, users: %v\n", len(users), users)

	c.Find(bson.M{"name": bson.M{"$in": []string{"AAAAA", "BBBBB"}}}).All(&users)
	fmt.Printf("findAll count: %v, users: %v\n", len(users), users)
}

func findOne(c *mgo.Collection) {
	id := "6bbc8656bf9a240c21000001"
	objectId := bson.ObjectIdHex(id)
	var user User61
	err := c.Find(bson.M{"_id": objectId}).One(&user)
	if err != nil {
		fmt.Println("findOne err:", err)
	} else {
		fmt.Printf("user: %v\n", user)
	}

	err = c.FindId(objectId).One(&user)
	if err != nil {
		fmt.Println("FindId err:", err)
	} else {
		fmt.Printf("user: %v\n", user)
	}

	err = c.Find(bson.M{"name": "AAAAA"}).One(&user)
	if err != nil {
		fmt.Println("findOne err:", err)
	} else {
		fmt.Printf("user: %v\n", user)
	}

}

func insertMulti(c *mgo.Collection) {
	err := c.Insert(
		&User61{
			Id_:       bson.NewObjectId(),
			Name:      "AAAAA",
			Age:       29,
			JonedAt:   time.Now(),
			Interests: []string{"a", "b", "c"},
		},
		&User61{
			Id_:       bson.NewObjectId(),
			Name:      "BBBBB",
			Age:       28,
			JonedAt:   time.Now(),
			Interests: []string{"1", "2", "3"},
		},
	)

	if err != nil {
		fmt.Println("insertMulti err:", err)
	}
	fmt.Println("insertMulti success")
}

func insertOne(c *mgo.Collection) {
	err := c.Insert(&User61{
		Id_:       bson.NewObjectId(),
		Name:      "Jimmy Kuu",
		Age:       33,
		JonedAt:   time.Now(),
		Interests: []string{"Develop", "Movie"},
	})

	if err != nil {
		fmt.Println("insertOne err:", err)
	}
	fmt.Println("insertOne success")
}
