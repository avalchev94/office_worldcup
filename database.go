package main

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Database struct {
	session *mgo.Session
}

func NewDB() (*Database, error) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		return nil, err
	}

	return &Database{session}, nil
}

func (db *Database) Close() {
	db.session.Close()
}

type User struct {
	ID       bson.ObjectId `bson:"_id"`
	Username string        `bson:"username"`
	Fullname string        `bson:"fullname"`
	Password string        `bson:"password"`
	Role     string        `bson:"role"`
	Points   int32         `bson:"points"`
	Avatar   string        `bson:"avatar"`
}

func (db *Database) GetUser(username string) (User, error) {
	users := db.session.DB("world_cup").C("users")
	if users == nil {
		return User{}, errors.New("users collection not found")
	}

	var user User
	if users.Find(bson.M{"username": username}).One(&user) != nil {
		return User{}, errors.New("user not found")
	}
	return user, nil
}

func (db *Database) AddUser(user User) error {
	users := db.session.DB("world_cup").C("users")
	if users == nil {
		return errors.New("users collection not found")
	}

	if count, _ := users.Find(bson.M{"username": user.Username}).Count(); count > 0 {
		return errors.New("already have such user")
	}

	users.Insert(&user)
	return nil
}

type Team struct {
	ID   bson.ObjectId `bson:"_id"`
	Name string        `bson:"name"`
}

func (db *Database) GetTeamByID(id bson.ObjectId) (Team, error) {
	teams := db.session.DB("world_cup").C("teams")
	if teams == nil {
		return Team{}, errors.New("teams not found")
	}

	var t Team
	if teams.Find(bson.M{"_id": id}).One(&t) != nil {
		return Team{}, errors.New("team not found")
	}

	return t, nil
}

type Match struct {
	ID     bson.ObjectId `bson:"_id"`
	Host   bson.ObjectId `bson:"host"`
	Guest  bson.ObjectId `bson:"guest"`
	Result string        `bson:"result"`
	Date   time.Time     `bson:"date"`
}

func (db *Database) GetTodayMatches() ([]Match, error) {
	matches := db.session.DB("world_cup").C("matches")
	if matches == nil {
		return nil, errors.New("matches not found")
	}

	var result []Match
	matches.Find(bson.M{
		"date": bson.M{
			"$gt": time.Now(),
		},
	}).All(&result)

	return result, nil
}

type Prediction struct {
	ID        bson.ObjectId `bson:"_id"`
	Match     bson.ObjectId `bson:"match"`
	User      bson.ObjectId `bson:"user"`
	Predicted string        `bson:"predicted"`
}

func (db *Database) GetPrediction(match, user bson.ObjectId) (Prediction, error) {
	predictions := db.session.DB("world_cup").C("predictions")
	if predictions == nil {
		return Prediction{}, errors.New("predictions not found")
	}

	var p Prediction
	if predictions.Find(bson.M{"match": match, "user": user}).One(&p) != nil {
		return Prediction{}, errors.New("prediction not found")
	}

	return p, nil
}
