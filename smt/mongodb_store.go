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
	smtName  string
}

func NewMongoStore(ctx context.Context, client *mongo.Client, database, smtName string) *MongodbStore {
	return &MongodbStore{ctx: ctx, client: client, database: database, smtName: smtName}
}

func (m *MongodbStore) Collection() *mongo.Collection {
	return m.client.Database(m.database).Collection(m.smtName)
}

type MongodbRoot struct {
	Root H256 `json:"r" bson:"r"`
}

func (m *MongodbStore) UpdateRoot(root H256) error {
	var data MongodbRoot
	data.Root = root
	update := bson.M{"$set": data}
	updateOpts := options.Update().SetUpsert(true)
	_, err := m.Collection().UpdateOne(context.Background(), bson.M{"_id": "root"}, update, updateOpts)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongodbStore) Root() (H256, error) {
	var root MongodbRoot
	err := m.Collection().FindOne(m.ctx, bson.M{"_id": "root"}).Decode(&root)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return H256Zero(), nil
		}
		return nil, err
	}
	return root.Root, nil
}

func (m *MongodbStore) GetBranch(key BranchKey) (*BranchNode, error) {
	keyHash := key.GetHash()
	var node BranchNode
	err := m.Collection().FindOne(m.ctx, bson.M{"_id": keyHash}).Decode(&node)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, StoreErrorNotExist
		}
		return nil, err
	}
	return &node, nil
}

func (m *MongodbStore) InsertBranch(key BranchKey, node BranchNode) error {
	keyHash := key.GetHash()
	update := bson.M{"$set": node}
	updateOpts := options.Update().SetUpsert(true)
	_, err := m.Collection().UpdateOne(context.Background(), bson.M{"_id": keyHash}, update, updateOpts)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongodbStore) RemoveBranch(key BranchKey) error {
	keyHash := key.GetHash()
	_, err := m.Collection().DeleteOne(m.ctx, bson.M{"_id": keyHash})
	if err != nil {
		return err
	}
	return nil
}
