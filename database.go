package main

import (
	"errors"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"time"
  "os"
  "log"
)

type Database struct {
	session *mgo.Session
}

func NewDB() (*Database, error) {
  session, err := mgo.Dial(os.Getenv("MongoServer"))
  if err != nil {
    log.Fatalln(err)
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
  Favourite bson.ObjectId `bson:"favourite"`
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

  user.ID = bson.NewObjectId()
	err := users.Insert(&user)
	return err
}

type Team struct {
	ID   bson.ObjectId `bson:"_id"`
	Name string        `bson:"name"`
}

func (db *Database) GetTeams() ([]Team, error) {
  teams := db.session.DB("world_cup").C("teams")
	if teams == nil {
		return nil, errors.New("teams not found")
	}

  var result []Team
  return result, teams.Find(nil).All(&result)
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

func (db *Database) GetMatch(id bson.ObjectId) (Match, error) {
  matches := db.session.DB("world_cup").C("matches")
	if matches == nil {
		return Match{}, errors.New("matches not found")
	}

  var match Match
  if (matches.Find(bson.M{"_id":id}).One(&match) != nil) {
    return Match{}, errors.New("match not found")
  }

  return match, nil
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

func (db *Database) AddPrediction(p Prediction) error {
  if _, err := db.GetPrediction(p.Match, p.User); err == nil {
    return errors.New("user has predicted this game already")
  }

  predictions := db.session.DB("world_cup").C("predictions")
  if predictions == nil {
    return errors.New("predictions not found")
  }

  p.ID = bson.NewObjectId()
  err := predictions.Insert(&p)
  return err
}
