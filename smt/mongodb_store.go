package smt

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongodbStore struct {
	ctx      context.Context
	client   *mongo.Client
	database string
}

func NewMongoStore(ctx context.Context, client *mongo.Client, database string) *MongodbStore {
	return &MongodbStore{ctx: ctx, client: client, database: database}
}

func (m *MongodbStore) Client() *mongo.Client {
	return m.client
}

func (m *MongodbStore) GetBranch(key BranchKey) (*BranchNode, error) {
	keyHash := key.GetHash()
	var node BranchNode
	collection := m.client.Database(m.database).Collection(key.SmtName)
	err := collection.FindOne(m.ctx, bson.M{"_id": keyHash}).Decode(&node)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (m *MongodbStore) InsertBranch(key BranchKey, node BranchNode) error {
	keyHash := key.GetHash()
	collection := m.client.Database(m.database).Collection(key.SmtName)
	update := bson.M{"$set": node}
	updateOpts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": keyHash}, update, updateOpts)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongodbStore) RemoveBranch(key BranchKey) error {
	keyHash := key.GetHash()
	collection := m.client.Database(m.database).Collection(key.SmtName)
	_, err := collection.DeleteOne(m.ctx, bson.M{"_id": keyHash})
	if err != nil {
		return err
	}
	return nil
}
