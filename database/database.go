package database

import (
	"errors"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"log"
	"os"
	"time"
)

func ObjectIdHex(id string) bson.ObjectId {
	return bson.ObjectIdHex(id)
}

type Database struct {
	session *mgo.Session
}

func New() (*Database, error) {
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
	ID        bson.ObjectId `bson:"_id"`
	Username  string        `bson:"username"`
	Fullname  string        `bson:"fullname"`
	Favourite bson.ObjectId `bson:"favourite"`
	Password  string        `bson:"password"`
	Role      string        `bson:"role"`
	Points    int32         `bson:"points"`
	Avatar    string        `bson:"avatar"`
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

func (db *Database) GetUsers() ([]User, error) {
	users := db.session.DB("world_cup").C("users")
	if users == nil {
		return nil, errors.New("users collection not found")
	}

	var result []User
	err := users.Find(nil).All(&result)
	return result, err
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

const (
	GroupStage    = "group"
	KnockoutStage = "knockout"
)

type Match struct {
	ID     bson.ObjectId `bson:"_id"`
	Host   bson.ObjectId `bson:"host"`
	Guest  bson.ObjectId `bson:"guest"`
	Stage  string        `bson:"stage"`
	Result string        `bson:"result"`
	Date   time.Time     `bson:"date"`
}

func (db *Database) GetMatch(id bson.ObjectId) (Match, error) {
	matches := db.session.DB("world_cup").C("matches")
	if matches == nil {
		return Match{}, errors.New("matches not found")
	}

	var match Match
	if (matches.Find(bson.M{"_id": id}).One(&match) != nil) {
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

func (db *Database) GetUnfinishedMatches() ([]Match, error) {
	matches := db.session.DB("world_cup").C("matches")
	if matches == nil {
		return nil, errors.New("matches not found")
	}

	var result []Match
	matches.Find(bson.M{"result": ""}).All(&result)

	return result, nil
}

func (db *Database) UpdateMatch(m Match) error {
	matches := db.session.DB("world_cup").C("matches")
	if matches == nil {
		return errors.New("matches not found")
	}

	return matches.Update(bson.M{"_id": m.ID}, &m)
}

func (db *Database) AddMatch(m Match) error {
	matches := db.session.DB("world_cup").C("matches")
	if matches == nil {
		return errors.New("matches not found")
	}

	m.ID = bson.NewObjectId()
	return matches.Insert(&m)
}

type Prediction struct {
	ID        bson.ObjectId `bson:"_id"`
	Match     bson.ObjectId `bson:"match"`
	User      bson.ObjectId `bson:"user"`
	Predicted string        `bson:"predicted"`
	Score     int32         `bson:"score"`
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

func (db *Database) GetPredictionsByMatch(match bson.ObjectId) ([]Prediction, error) {
	predictions := db.session.DB("world_cup").C("predictions")
	if predictions == nil {
		return nil, errors.New("predictions not found")
	}

	var result []Prediction
	if predictions.Find(bson.M{"match": match}).All(&result) != nil {
		return nil, errors.New("no predictions for that match")
	}

	return result, nil
}

func (db *Database) GetPredictionsByUser(user bson.ObjectId) ([]Prediction, error) {
	predictions := db.session.DB("world_cup").C("predictions")
	if predictions == nil {
		return nil, errors.New("predictions not found")
	}

	var result []Prediction
	if predictions.Find(bson.M{"user": user}).All(&result) != nil {
		return nil, errors.New("no predictions for that match")
	}

	return result, nil
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

func (db *Database) UpdatePrediction(p Prediction) error {
	predictions := db.session.DB("world_cup").C("predictions")
	if predictions == nil {
		return errors.New("predictions not found")
	}

	return predictions.Update(bson.M{"_id": p.ID}, &p)
}
