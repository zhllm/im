package service

import (
	"context"
	"fin_im/conf"
	"fin_im/model/ws"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertMsg(database string, id string, content string, read int, expire int64) error {
	// 插入到mongodb
	collection := conf.MongoDBClient.Database(database).Collection(id)
	comment := ws.Trainer{
		Content:   content,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix() + expire,
		Read:      uint(read),
	}

	_, err := collection.InsertOne(context.TODO(), comment)
	return err
}

func FindMany(database string, sendId string, id string, time int64, pageSize int) (results []ws.Result, err error) {
	var resultMe []ws.Trainer
	var resultYou []ws.Trainer

	sendCollection := conf.MongoDBClient.Database(database).Collection(sendId)
	idCollection := conf.MongoDBClient.Database(database).Collection(id)
	sendIDTimeCurcor, err := sendCollection.Find(context.TODO(), options.Find().SetSort(bson.D{{"startTime", -1}}), options.Find().SetLimit(int64(pageSize)))

}
