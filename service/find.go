package service

import (
	"context"
	"encoding/json"
	"fin_im/conf"
	"fin_im/model/ws"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SendSortMsg struct {
	Content  string `json:"content"`
	Read     uint   `json:"read"`
	CreateAt int64  `json:"create_at"`
}

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
	sendIDTimeCurcor, _ := sendCollection.Find(
		context.TODO(),
		bson.D{},
		// options.Find().SetSort(bson.D{{"startTime", -1}}),
		options.Find().SetLimit(int64(pageSize)),
	)
	idTimeCurcor, _ := idCollection.Find(
		context.TODO(),
		bson.D{},
		// options.Find().SetSort(bson.D{{"startTime", -1}}),
		options.Find().SetLimit(int64(pageSize)),
	)

	err = sendIDTimeCurcor.All(context.TODO(), &resultYou)
	err = idTimeCurcor.All(context.TODO(), &resultMe)
	results, err = AppendAndSort(resultMe, resultYou)
	return
}

func AppendAndSort(resultMe, resultYou []ws.Trainer) (results []ws.Result, err error) {
	for _, r := range resultMe {
		sendSortMsg := SendSortMsg{
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		jsonString, _ := json.Marshal(sendSortMsg)
		result := ws.Result{
			StartTime: r.StartTime,
			Msg:       string(jsonString),
			From:      "me",
		}
		results = append(results, result)
	}

	for _, r := range resultYou {
		sendSortMsg := SendSortMsg{
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		jsonString, _ := json.Marshal(sendSortMsg)
		result := ws.Result{
			StartTime: r.StartTime,
			Msg:       string(jsonString),
			From:      "other",
		}
		results = append(results, result)
	}
	return
}
