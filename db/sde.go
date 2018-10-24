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
	"fmt"
	"time"

	"io/ioutil"

	"github.com/oxisto/titan/model"
	"gopkg.in/yaml.v2"
)

type StaticDataExport struct {
	Version int32
	Server  string
}

func (t StaticDataExport) ID() int32 {
	return t.Version
}

func (t StaticDataExport) ExpiresOn() *time.Time {
	return nil
}

func (t StaticDataExport) SetExpire(time *time.Time) {

}

func (t StaticDataExport) HashKey() string {
	return fmt.Sprintf("sde:%d", t.ID())
}

func ImportSDE(version int32, server string) {
	log.Infof("Importing SDE %d...", version)

	files := map[string]string{
		"sde/fsd/blueprints.yaml":  "blueprints",
		"sde/fsd/typeIDs.yaml":     "types",
		"sde/fsd/groupIDs.yaml":    "groups",
		"sde/fsd/categoryIDs.yaml": "categories",
	}

	objects := map[string]interface{}{
		"blueprints": make(map[int32]model.Blueprint),
		"types":      make(map[int32]model.Type),
		"groups":     make(map[int32]model.Group),
		"categories": make(map[int32]model.Category),
	}

	for k, v := range files {
		if err := ImportSDEFile(k, v, objects[v]); err != nil {
			log.Errorf("An error occured while importing blueprints: %v", err)
		}
	}
}

func ImportSDEFile(fileName string, objectType string, in interface{}) error {
	log.Infof("Reading %s ...", objectType)

	err := UnmarshalYAMLFromFile(fileName, in)

	if err != nil {
		return err
	}

	log.Infof("Inserting %s ...", objectType)

	collection := mongo.C(objectType)
	collection.RemoveAll(nil)

	typeArray := []model.MetaGroup{}
	metaTypes := make(map[int32]model.MetaGroup)

	err = UnmarshalYAMLFromFile("sde/bsd/invMetaTypes.yaml", &typeArray)

	for _, entry := range typeArray {
		metaTypes[entry.TypeID] = entry

		// set the parent type's meta group to 1 if the metaType is 2
		if entry.MetaGroupID == 2 {
			metaTypes[entry.ParentTypeID] = model.MetaGroup{
				TypeID:      entry.ParentTypeID,
				MetaGroupID: 1,
			}
		}
	}

	i := 0

	switch v := in.(type) {
	case map[int32]model.Blueprint:
		for objectID, blueprint := range v {
			blueprint.ObjectID = objectID

			if err = collection.Insert(blueprint); err != nil {
				return err
			}
			i++
		}
	case map[int32]model.Type:
		for typeID, t := range v {
			t.TypeID = typeID

			if metaType, ok := metaTypes[typeID]; ok {
				t.MetaGroupID = metaType.MetaGroupID
			}

			if err = collection.Insert(t); err != nil {
				return err
			}
			i++
		}
	case map[int32]model.Group:
		for groupID, group := range v {
			group.GroupID = groupID

			if err = collection.Insert(group); err != nil {
				return err
			}
			i++
		}
	case map[int32]model.Category:
		for categoryID, category := range v {
			category.CategoryID = categoryID

			if err = collection.Insert(category); err != nil {
				return err
			}
			i++
		}
	}

	log.Infof("Successfully inserted %d %s.", i, objectType)

	return err
}

func UnmarshalYAMLFromFile(fileName string, out interface{}) (err error) {
	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, out)
}
