package main

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "errors"
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
  Username string `bson:"username"`
  Fullname string `bson:"fullname"`
  Role string `bson:"role"`
  Points int32 `bson:"points"`
  Avatar string `bson:"avatar"`
}

func (db *Database) QueryUser(username string, password string) (*User, error) {
  users := db.session.DB("world_cup").C("users")
  if users != nil {
    return nil, errors.New("users collection not found")
  }

  var user User
  if users.Find(bson.M{"username":username, "password":password}).One(&user) != nil {
    return nil, errors.New("user not found or password incorrect")
  }
  return &user, nil
}
