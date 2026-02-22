package config

import (
	"context"
	"time"

	"github.com/S42yt/tuubaa-bot/core"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type guildConfig struct {
	GuildID string            `bson:"guild_id"`
	Roles   map[string]string `bson:"roles"`
}

func GetRole(guildID, key string) (string, error) {
	db := core.NewMongoHandler()
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()
	if err := db.Connect(ctx); err != nil {
		return "", err
	}
	defer db.Disconnect(ctx)

	coll := db.Collection("guild_configs")
	var cfg guildConfig
	if err := coll.FindOne(ctx, bson.M{"guild_id": guildID}).Decode(&cfg); err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil
		}
		return "", err
	}
	if cfg.Roles == nil {
		return "", nil
	}
	if v, ok := cfg.Roles[key]; ok {
		return v, nil
	}
	return "", nil
}

func GetRoles(guildID string) (map[string]string, error) {
	db := core.NewMongoHandler()
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()
	if err := db.Connect(ctx); err != nil {
		return nil, err
	}
	defer db.Disconnect(ctx)

	coll := db.Collection("guild_configs")
	var cfg guildConfig
	if err := coll.FindOne(ctx, bson.M{"guild_id": guildID}).Decode(&cfg); err != nil {
		if err == mongo.ErrNoDocuments {
			return map[string]string{}, nil
		}
		return nil, err
	}
	if cfg.Roles == nil {
		return map[string]string{}, nil
	}
	return cfg.Roles, nil
}
