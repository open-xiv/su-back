package rmongo

import (
	"github.com/open-xiv/su-back/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

func InitFights(coll *mongo.Collection, fight model.Fight) (primitive.ObjectID, error) {
	fight.ID = primitive.NewObjectID()
	fight.ServerRecord.Update = time.Now().Unix()
	rst, err := coll.InsertOne(nil, fight)
	if err != nil {
		zap.L().Debug("failed to insert fight", zap.Error(err))
		return primitive.NilObjectID, err
	}
	zap.L().Debug("new fight inserted", zap.Any("id", rst.InsertedID))
	return rst.InsertedID.(primitive.ObjectID), err
}

func PullFight(coll *mongo.Collection, id primitive.ObjectID) (model.Fight, error) {
	var fight model.Fight
	err := coll.FindOne(nil, bson.M{"_id": id}).Decode(&fight)
	if err != nil {
		zap.L().Debug("failed to find fight", zap.Error(err))
		return model.Fight{}, err
	}
	zap.L().Debug("fight found", zap.Any("id", fight.ID))
	return fight, err
}

func PushFight(coll *mongo.Collection, fight model.Fight) (model.Fight, error) {
	fight.ServerRecord.Update = time.Now().Unix()
	_, err := coll.ReplaceOne(nil, bson.M{"_id": fight.ID}, fight)
	if err != nil {
		zap.L().Debug("failed to update fight", zap.Error(err))
		return model.Fight{}, err
	}
	zap.L().Debug("fight updated", zap.Any("id", fight.ID))
	return fight, err
}

func RemoveFight(coll *mongo.Collection, id primitive.ObjectID) error {
	_, err := coll.DeleteOne(nil, bson.M{"_id": id})
	if err != nil {
		zap.L().Debug("failed to delete fight", zap.Error(err))
		return err
	}
	zap.L().Debug("fight deleted", zap.Any("id", id))
	return err
}

func _(coll *mongo.Collection, id primitive.ObjectID, patch bson.M) error {
	// update record
	patch["record.update"] = time.Now().Unix()
	// update fight
	rst, err := coll.UpdateOne(nil, bson.M{"_id": id}, bson.M{"$set": patch})
	if err != nil {
		zap.L().Debug("failed to patch fight", zap.Error(err))
		return err
	}
	zap.L().Debug("fight patched", zap.Any("id", id), zap.Any("rst", rst))
	return err
}
