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

package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/doublerebel/bellows"
	esi "github.com/evecentral/esiapi/client"
	esiCharacter "github.com/evecentral/esiapi/client/character"
	esiCorporation "github.com/evecentral/esiapi/client/corporation"
	esiIndustry "github.com/evecentral/esiapi/client/industry"
	esiSkills "github.com/evecentral/esiapi/client/skills"
	"github.com/fatih/structs"
	"github.com/go-openapi/runtime/client"
	"github.com/go-redis/redis"
	"github.com/mitchellh/mapstructure"
	"github.com/oxisto/titan/db"
	"github.com/oxisto/titan/model"
	"github.com/sirupsen/logrus"
)

var cache *redis.Client
var log *logrus.Entry

var corporationMap map[int32]int32

func init() {
	log = logrus.WithField("component", "cache")

	corporationMap = map[int32]int32{}
}

func InitCache(redisAddr string) {
	log.Infof("Using Redis cache @ %s", redisAddr)

	cache = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
}

type FetchFuncType func(callerID int32, objectID int32, object model.CachedObject) error

func GetCharacter(characterID int32, character *model.Character) error {
	hashKey := fmt.Sprintf("character:%d", characterID)
	return GetCachedObject(hashKey, characterID, characterID, character, FetchCharacter)
}

func GetCorporation(callerID int32, corporationID int32, corporation *model.Corporation) error {
	hashKey := fmt.Sprintf("corporation:%d", corporationID)
	return GetCachedObject(hashKey, callerID, corporationID, corporation, FetchCorporation)
}

func GetIndustryJobs(callerID int32, corporationID int32, jobs *model.IndustryJobs) error {
	hashKey := fmt.Sprintf("industry-jobs:%d", corporationID)
	return GetCachedObject(hashKey, callerID, corporationID, jobs, FetchCorporationIndustryJobs)
}

func GetAccessToken(characterID int32, accessToken *model.AccessToken) error {
	hashKey := fmt.Sprintf("accesstoken:%d", characterID)
	return GetCachedObject(hashKey, characterID, characterID, accessToken, FetchAccessToken)
}

func GetAccessTokenForCorporation(corporationID int32, accessToken *model.AccessToken) error {
	characterID, ok := corporationMap[corporationID]
	if !ok {
		return fmt.Errorf("Could not find a character for corporation %d", corporationID)
	}

	hashKey := fmt.Sprintf("accesstoken:%d", characterID)
	return GetCachedObject(hashKey, characterID, characterID, accessToken, FetchAccessToken)
}

func GetCachedObject(hashKey string, callerID int32, objectID int32, object model.CachedObject, funcType FetchFuncType) (err error) {
	exists, err := cache.Exists(hashKey).Result()
	if err != nil {
		return err
	}

	// if it exists, read if from cache
	if exists == 1 {
		return ReadCachedObject(hashKey, object)
	}

	log.Debugf("Fetching %s from source...", hashKey)

	// otherwise, fetch it
	if err = funcType(callerID, objectID, object); err == nil {
		// update the cache if no error was found
		return WriteCachedObject(object)
	}

	return err
}

func ReadCachedObject(hashKey string, object model.CachedObject) error {
	var (
		d   *mapstructure.Decoder
		err error
	)

	var fields []string

	if fields, err = cache.HKeys(hashKey).Result(); err != nil {
		return err
	}

	m := map[string]interface{}{}

	if values, err := cache.HMGet(hashKey, fields...).Result(); err != nil {
		return err
	} else {
		for k, v := range values {
			field := fields[k]
			m[field] = v
		}
	}

	m = bellows.Expand(m)

	config := mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           object,
	}

	if d, err = mapstructure.NewDecoder(&config); err != nil {
		return err
	}

	return d.Decode(m)
}

func WriteCachedObject(object model.CachedObject) error {
	m := structs.Map(object)
	m = bellows.Flatten(m)

	log.Debugf("Writing %v to cache...", object.HashKey())

	for k, v := range m {
		if reflect.TypeOf(v).Kind() == reflect.Map {
			// remove maps for now
			delete(m, k)
		}
	}

	_, err := cache.HMSet(object.HashKey(), m).Result()

	// set expiry if necessary
	if tm := object.ExpiresOn(); tm != nil {
		cache.ExpireAt(object.HashKey(), *tm)
	}

	return err
}

func FetchCorporationIndustryJobs(callerID int32, corporationID int32, object model.CachedObject) error {
	jobs, ok := object.(*model.IndustryJobs)
	if !ok {
		return errors.New("passing invalid type to FetchCorporationIndustryJobs function")
	}

	// find access token for corporation
	accessToken := model.AccessToken{}
	err := GetAccessTokenForCorporation(corporationID, &accessToken)
	if err != nil {
		return err
	}

	jobsParams := esiIndustry.NewGetCorporationsCorporationIDIndustryJobsParams()
	jobsParams.CorporationID = corporationID
	jobsResponse, err := esi.Default.Industry.GetCorporationsCorporationIDIndustryJobs(jobsParams, client.BearerToken(accessToken.Token))
	if err != nil {
		return err
	}

	t, err := time.Parse(time.RFC1123, jobsResponse.Expires)
	if err != nil {
		return err
	}

	jobs.SetExpire(&t)

	jobs.CorporationID = corporationID

	jobs.Jobs = map[string]model.IndustryJob{}

	for _, v := range jobsResponse.Payload {
		blueprintTypeID := model.SafeInt32(v.BlueprintTypeID)

		job := model.IndustryJob{
			ActivityID: model.SafeInt32(v.ActivityID),
			Blueprint:  db.GetType(blueprintTypeID),
		}

		jobs.Jobs[strconv.Itoa(int(model.SafeInt32(v.JobID)))] = job
	}

	return nil
}

func FetchCorporation(callerID int32, corporationID int32, object model.CachedObject) error {
	corporation, ok := object.(*model.Corporation)
	if !ok {
		return errors.New("passing invalid type to FetchCorporation function")
	}

	corporationParams := esiCorporation.NewGetCorporationsCorporationIDParams()
	corporationParams.CorporationID = corporationID
	corporationResponse, err := esi.Default.Corporation.GetCorporationsCorporationID(corporationParams)
	if err != nil {
		return err
	}

	corporation.CorporationID = corporationID
	corporation.CorporationName = model.SafeString(corporationResponse.Payload.Name)
	corporation.AllianceID = corporationResponse.Payload.AllianceID
	corporation.Ticker = model.SafeString(corporationResponse.Payload.Ticker)
	corporation.CEOID = model.SafeInt32(corporationResponse.Payload.CeoID)

	expireTime, err := time.Parse(time.RFC1123, corporationResponse.Expires)
	if err != nil {
		return err
	}
	corporation.SetExpire(&expireTime)

	// find access token for character
	accessToken := model.AccessToken{}
	err = GetAccessToken(callerID, &accessToken)
	if err != nil {
		return err
	}

	membersParams := esiCorporation.NewGetCorporationsCorporationIDMembersParams()
	membersParams.CorporationID = corporationID
	membersResponse, err := esi.Default.Corporation.GetCorporationsCorporationIDMembers(membersParams, client.BearerToken(accessToken.Token))
	if err != nil {
		return err
	}

	corporation.Members = map[string]int32{}
	for _, v := range membersResponse.Payload {
		corporation.Members[strconv.Itoa(int(v))] = v
	}

	return nil
}

func FetchCharacter(callerID int32, characterID int32, object model.CachedObject) error {
	character, ok := object.(*model.Character)
	if !ok {
		return errors.New("passing invalid type to FetchCharacter function")
	}

	characterParams := esiCharacter.NewGetCharactersCharacterIDParams()
	characterParams.CharacterID = characterID
	characterResponse, err := esi.Default.Character.GetCharactersCharacterID(characterParams)
	if err != nil {
		return err
	}

	character.CharacterID = characterID
	character.CharacterName = model.SafeString(characterResponse.Payload.Name)
	character.CorporationID = model.SafeInt32(characterResponse.Payload.CorporationID)
	character.AllianceID = characterResponse.Payload.AllianceID

	corporationParams := esiCorporation.NewGetCorporationsCorporationIDParams()
	corporationParams.CorporationID = character.CorporationID

	corporationResponse, err := esi.Default.Corporation.GetCorporationsCorporationID(corporationParams)
	if err != nil {
		return err
	}

	character.CorporationName = model.SafeString(corporationResponse.Payload.Name)

	expireTime, err := time.Parse(time.RFC1123, characterResponse.Expires)
	if err != nil {
		return err
	}
	character.SetExpire(&expireTime)

	// find access token for character
	accessToken := model.AccessToken{}
	err = GetAccessToken(callerID, &accessToken)
	if err != nil {
		return err
	}

	skillsParams := esiSkills.NewGetCharactersCharacterIDSkillsParams()
	skillsParams.CharacterID = characterID
	skillsResponse, err := esi.Default.Skills.GetCharactersCharacterIDSkills(skillsParams, client.BearerToken(accessToken.Token))
	if err != nil {
		return err
	}

	character.Skills = map[string]model.Skill{}
	for _, v := range skillsResponse.Payload.Skills {
		s := model.Skill{
			SkillID:     model.SafeInt32(v.SkillID),
			SkillPoints: model.SafeInt64(v.SkillpointsInSkill),
			Level:       model.SafeInt32(v.TrainedSkillLevel),
		}
		character.Skills[strconv.Itoa(int(s.SkillID))] = s
	}

	// update the corporation map
	corporationMap[character.CorporationID] = characterID
	log.Infof("Setting character %d as token for corporation %d", characterID, character.CorporationID)

	return nil
}

func FetchType(callerID int32, typeID int32, object model.CachedObject) error {
	t, ok := object.(*model.Type)
	if !ok {
		return errors.New("passing invalid type to FetchType function")
	}

	*t = db.GetType(typeID)

	return nil
}

func FetchAccessToken(callerID int32, characterID int32, object model.CachedObject) error {
	accessToken, ok := object.(*model.AccessToken)
	if !ok {
		return errors.New("passing invalid type to FetchAccessToken function")
	}

	// check, if a refresh token exists, otherwise we cannot fetch an access token
	refreshToken := model.RefreshToken{}
	refreshToken.CharacterID = characterID
	hashKey := refreshToken.HashKey()
	exists, err := cache.Exists(hashKey).Result()
	if exists != 1 {
		err = fmt.Errorf("No access or refresh tokens exist for character %d", characterID)
	}
	if err != nil {
		return err
	}

	ReadCachedObject(hashKey, &refreshToken)

	// fetch a new access token
	tokenResponse, expiryTime, _, characterName, err := SSO.AccessToken(refreshToken.Token, true)
	if err != nil {
		return err
	}

	// create a new access token cache object
	*accessToken = model.AccessToken{
		CharacterID:   int32(characterID),
		CharacterName: characterName,
		Token:         tokenResponse.AccessToken,
	}
	accessToken.SetExpire(&expiryTime)

	return nil
}

func GetPrices(regionID int, types []int32) (prices map[int32]model.Price, err error) {
	prices = make(map[int32]model.Price)

	typesToFetch := []int32{}

	for _, typeID := range types {
		var (
			cached  int64
			price   model.Price
			hashKey = fmt.Sprintf("price:%d", typeID)
		)

		// check, if prices are somehow cached
		if cached, err = cache.Exists(hashKey).Result(); err != nil {
			return nil, err
		}

		if cached == int64(1) {
			// read it from cache
			if err = ReadCachedObject(hashKey, &price); err != nil {
				return nil, err
			}

			// add it to the result map
			prices[typeID] = price
		} else {
			typesToFetch = append(typesToFetch, typeID)
		}
	}

	if len(typesToFetch) == 0 {
		return
	} else {
		log.Infof("Need to fetch %d price(s)...", len(typesToFetch))
	}

	chunkSize := 150

	for chunkStart := 0; chunkStart < len(typesToFetch); chunkStart += chunkSize {
		chunkEnd := chunkStart + chunkSize

		if chunkEnd > len(typesToFetch) {
			chunkEnd = len(typesToFetch)
		}

		// deliberately ignore errors here, because fuzzwork json objects are sometimes not properly formatted
		results, _ := FetchPrices(regionID, typesToFetch[chunkStart:chunkEnd])

		for k, price := range results {
			typeID, err := strconv.Atoi(k)
			expireDate := time.Now().Add(time.Hour * time.Duration(1))

			if err != nil {
				return nil, err
			}

			price.TypeID = int32(typeID)
			price.SetExpire(&expireDate)

			WriteCachedObject(&price)

			prices[int32(typeID)] = price
		}
	}

	return
}

func FetchPrices(regionID int, types []int32) (prices map[string]model.Price, err error) {
	typesParam := []string{}

	for _, typeID := range types {
		typesParam = append(typesParam, strconv.Itoa(int(typeID)))
	}

	log.Infof("Requesting %d types from fuzzwork", len(types))

	url := "https://market.fuzzwork.co.uk/aggregates/?region=" + strconv.Itoa(regionID) + "&types=" + strings.Join(typesParam, ",")

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(res.Body).Decode(&prices)

	return
}
