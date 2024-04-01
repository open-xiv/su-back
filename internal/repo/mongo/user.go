package rmongo

import (
	"errors"
	"github.com/open-xiv/su-back/internal/tools"
	"github.com/open-xiv/su-back/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

func InitUser(coll *mongo.Collection, user model.User) (primitive.ObjectID, error) {
	// check name unique
	_, err := PullUserByName(coll, user.Person.Name)
	if err == nil {
		zap.L().Debug("user name already exists")
		return primitive.NilObjectID, errors.New("user name already exists")
	}

	// create new user
	user.ID = primitive.NewObjectID()
	user.ServerRecord.Update = time.Now().Unix()
	// random string with length 12
	user.Person.Key = tools.RandStringBytes(12)

	// insert
	rst, err := coll.InsertOne(nil, user)
	if err != nil {
		zap.L().Debug("failed to insert user", zap.Error(err))
		return primitive.NilObjectID, err
	}
	zap.L().Debug("new user inserted", zap.Any("id", rst.InsertedID))
	return rst.InsertedID.(primitive.ObjectID), err
}

func PullUser(coll *mongo.Collection, id primitive.ObjectID) (model.User, error) {
	var user model.User
	err := coll.FindOne(nil, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		zap.L().Debug("failed to find user", zap.Error(err))
		return model.User{}, err
	}
	zap.L().Debug("user found", zap.Any("id", user.ID))
	return user, err
}

func PullUserByName(coll *mongo.Collection, name string) (model.User, error) {
	var user model.User
	err := coll.FindOne(nil, bson.M{"person.name": name}).Decode(&user)
	if err != nil {
		zap.L().Debug("failed to find user", zap.Error(err))
		return model.User{}, err
	}
	zap.L().Debug("user found", zap.Any("id", user.ID))
	return user, err
}

func PullUserByKey(coll *mongo.Collection, key string) (model.User, error) {
	var user model.User
	err := coll.FindOne(nil, bson.M{"person.key": key}).Decode(&user)
	if err != nil {
		zap.L().Debug("failed to find user", zap.Error(err))
		return model.User{}, err
	}
	zap.L().Debug("user found", zap.Any("id", user.ID))
	return user, err
}

func PushUser(coll *mongo.Collection, user model.User) (model.User, error) {
	user.ServerRecord.Update = time.Now().Unix()
	_, err := coll.ReplaceOne(nil, bson.M{"_id": user.ID}, user)
	if err != nil {
		zap.L().Debug("failed to update user", zap.Error(err))
		return model.User{}, err
	}
	zap.L().Debug("user updated", zap.Any("id", user.ID))
	return user, err
}

func RemoveUser(coll *mongo.Collection, id primitive.ObjectID) error {
	_, err := coll.DeleteOne(nil, bson.M{"_id": id})
	if err != nil {
		zap.L().Debug("failed to delete user", zap.Error(err))
		return err
	}
	zap.L().Debug("user deleted", zap.Any("id", id))
	return err
}

func PatchUser(coll *mongo.Collection, id primitive.ObjectID, patch bson.M) error {
	// update record (concat with patch)
	patch["record.update"] = time.Now().Unix()
	// update user
	rst, err := coll.UpdateOne(nil, bson.M{"_id": id}, bson.M{"$set": patch})
	if err != nil {
		zap.L().Debug("failed to patch user", zap.Error(err))
		return err
	}
	zap.L().Debug("user patched", zap.Any("id", id), zap.Any("rst", rst))
	return err
}

func InsertFight(coll *mongo.Collection, id primitive.ObjectID, fightId primitive.ObjectID) error {
	// update server record
	patch := bson.M{"$set": bson.M{"record.update": time.Now().Unix()}}
	rst, err := coll.UpdateOne(nil, bson.M{"_id": id}, patch)
	if err != nil {
		zap.L().Debug("failed to patch user", zap.Error(err))
		return err
	}

	// update fights
	patch = bson.M{"$push": bson.M{"fight_ids": fightId}}
	// update user
	rst, err = coll.UpdateOne(nil, bson.M{"_id": id}, patch)
	if err != nil {
		zap.L().Debug("failed to insert fight", zap.Error(err))
		return err
	}
	zap.L().Debug("fight inserted", zap.Any("id", id), zap.Any("rst", rst))
	return err
}
