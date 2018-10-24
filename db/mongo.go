/*
Copyright 2018 Christian Banse

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/oxisto/titan/model"
	"github.com/sirupsen/logrus"
)

var mongo *mgo.Database
var log *logrus.Entry

// ProductTypeIDBlacklist is a list of blacklisted type IDs
var ProductTypeIDBlacklist = map[int]bool{
	29202: true, 27038: true, 23883: true,
}

func init() {
	log = logrus.WithField("component", "db")
}

func InitMongoDB(mongoAddr string) {
	log.Infof("Using MongoDB @ %s", mongoAddr)

	session, err := mgo.Dial(mongoAddr)
	if err != nil {
		panic(err)
	}

	mongo = session.DB("titan")
}

func GetCategories() ([]model.Category, error) {
	categories := []model.Category{}

	err := mongo.C("categories").Find(bson.M{"published": true}).All(&categories)

	return categories, err
}

func GetProductTypeIDs() ([]int32, error) {
	types := []map[string]interface{}{}

	err := mongo.C("blueprints").Pipe([]bson.M{
		{"$group": bson.M{"_id": "$activities.manufacturing.products.typeID"}},
		{"$unwind": "$_id"},
		{"$lookup": bson.M{
			"from":         "types",
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "type",
		}},
		{"$unwind": "$type"},
		{"$match": bson.M{
			"$and": []bson.M{
				{"type.published": true},
				{"$or": []bson.M{
					{"type.metaGroupID": 1},
					{"type.metaGroupID": 2}}}}}},
	}).All(&types)

	typeIDs := []int32{}

	for _, v := range types {
		if typeID, ok := v["_id"].(int); ok && !ProductTypeIDBlacklist[typeID] {
			typeIDs = append(typeIDs, int32(typeID))
		}
	}

	return typeIDs, err
}

func GetActivityMaterials(activity string, blueprint model.Blueprint, runs int, materialModifier float64, materials interface{}) error {
	pipe := mongo.C("blueprints").Pipe([]bson.M{
		{"$match": bson.M{"_id": blueprint.BlueprintTypeID}},
		{"$unwind": "$activities." + activity + ".materials"},
		{"$group": bson.M{"_id": nil, "materials": bson.M{"$push": "$activities." + activity + ".materials"}}},
		{"$unwind": "$materials"},
		{"$lookup": bson.M{"from": "types",
			"localField":   "materials.typeID",
			"foreignField": "_id",
			"as":           "type"}},
		{"$project": bson.M{"_id": 0,
			"typeID":   "$materials.typeID",
			"typeName": "$type.name",
			"quantity": bson.M{"$ceil": bson.M{"$multiply": []interface{}{"$materials.quantity", runs, materialModifier}}}}},
		{"$unwind": "$typeName"},
		{"$sort": bson.M{"typeName": 1}},
	})

	if err := pipe.All(materials); err != nil {
		return err
	}

	return nil
}

func GetActivitySkills(activity string, blueprint model.Blueprint, skills interface{}) error {
	pipe := mongo.C("blueprints").Pipe([]bson.M{
		{"$match": bson.M{"_id": blueprint.BlueprintTypeID}},
		{"$unwind": "$activities." + activity + ".skills"},
		{"$group": bson.M{"_id": nil, "skills": bson.M{"$push": "$activities." + activity + ".skills"}}},
		{"$unwind": "$skills"},
		{"$lookup": bson.M{"from": "types",
			"localField":   "skills.typeID",
			"foreignField": "_id",
			"as":           "type"}},
		{"$project": bson.M{"_id": 0,
			"skillID":       "$skills.typeID",
			"skillName":     "$type.name",
			"requiredLevel": "$skills.level"}},
		{"$unwind": "$skillName"},
		{"$sort": bson.M{"skillName": 1}},
	})

	if err := pipe.All(skills); err != nil {
		return err
	}

	return nil
}

func GetType(typeID int32) model.Type {
	t := model.Type{}

	pipe := mongo.C("types").Pipe([]bson.M{
		{"$match": bson.M{"_id": typeID}},
		{"$lookup": bson.M{"from": "groups",
			"localField":   "groupID",
			"foreignField": "_id",
			"as":           "group"}},
		{"$unwind": "$group"},
	})

	pipe.One(&t)

	return t
}

func GetBlueprint(typeID int32, basedOn string) model.Blueprint {
	blueprint := model.Blueprint{}

	match := bson.M{}
	match[basedOn] = typeID

	// fetch the blueprint
	mongo.C("blueprints").Find(match).One(&blueprint)

	return blueprint
}

// TODO: use redis
func GetMaterialTypeIDs(activity string) []int32 {
	types := []bson.M{}

	pipe := mongo.C("blueprints").Pipe([]bson.M{
		{"$unwind": "$activities." + activity + ".products"},
		{"$lookup": bson.M{"from": "types",
			"localField":   "activities." + activity + ".products.typeID",
			"foreignField": "_id",
			"as":           "productType"}},
		{"$unwind": "$productType"},
		{"$match": bson.M{"$or": []bson.M{
			{"productType.metaGroupID": bson.M{"$lte": 2}},
			{"productType.metaGroupID": bson.M{"$exists": false}},
		}}},
		{"$group": bson.M{
			"_id": "$activities." + activity + ".materials.typeID"},
		},
		{"$unwind": "$_id"},
		{"$group": bson.M{
			"_id": "$_id",
		}},
	})

	pipe.All(&types)

	typeIDs := []int32{}

	for _, v := range types {
		if typeID, ok := v["_id"].(int); ok {
			typeIDs = append(typeIDs, int32(typeID))
		}
	}

	return typeIDs
}
