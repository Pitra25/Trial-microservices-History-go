package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr        string        `yaml:"addr"`
	Password    string        `yaml:"password"`
	User        string        `yaml:"user"`
	DB          int           `yaml:"db"`
	MaxRetries  int           `yaml:"max_retries"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
}

func NewClient(ctx context.Context, cfg Config) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		Username:     cfg.User,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to redis server: %s\n", err.Error())
		return nil, err
	}

	return db, nil
}

// ======================================

var cfg = Config{
	Addr:        "127.0.0.1:6379",
	Password:    "",
	User:        "",
	DB:          0,
	MaxRetries:  5,
	DialTimeout: 10 * time.Second,
	Timeout:     5 * time.Second,
}

type Recording struct {
	ID          int8   `json:"id"`
	Calculation string `json:"calculation"`
	CreatedAt   string `json:"createdAt"`
}

func GetRecordFromHash(key string) ([]Recording, error) {
	ctx := context.Background()
	rdb, err := NewClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %v", err)
	}

	if key != "" {
		// Get single record
		val, err := getRecord(key)
		if err != nil {
			return nil, fmt.Errorf("Failed to get recording: %v", err)
		}

		return []Recording{val}, nil
	} else {
		// Get all records
		var (
			recordings []Recording
			cursor     uint64
			keys       []string
		)

		for {
			const numberOfRecords = 10
			result, nextCursor, err := rdb.Scan(ctx, cursor, "*", numberOfRecords).Result()
			if err != nil {
				return nil, fmt.Errorf("REDIS | scan error: %v", err)
			}

			keys = append(keys, result...)
			cursor = nextCursor

			if cursor == 0 {
				break
			}
		}

		for _, key := range keys {
			recording, err := getRecord(key)
			if err != nil {
				return nil, fmt.Errorf("Failed to get recording: %v", err)
			}

			recordings = append(recordings, recording)
		}

		if len(recordings) == 0 {
			return nil, fmt.Errorf("REDIS | no records found")
		}

		return recordings, nil
	}
}

func CreatrRecording(record Recording) {

	IDstr := fmt.Sprint(record.ID)

	result, err := getRecord(IDstr)
	if err != nil {
		panic("Error get recording")
	} else {
		if result.Calculation != "" {
			return
		}
	}

	dbR, err := NewClient(context.Background(), cfg)
	if err != nil {
		panic(err)
	}

	userJSON, err := json.Marshal(record)
	if err != nil {
		log.Fatal(err)
	}

	const timeOfLife = 30 * time.Second
	if err := dbR.Set(context.Background(), IDstr, userJSON, timeOfLife).Err(); err != nil {
		log.Fatalf("REDIS | Failed to set data, error: %s", err.Error())
	}
}

func getRecord(key string) (Recording, error) {
	ctx := context.Background()
	rdb, err := NewClient(ctx, cfg)
	if err != nil {
		return Recording{}, fmt.Errorf("failed to connect to redis: %v", err)
	}

	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return Recording{}, nil
	} else if err != nil {
		return Recording{}, fmt.Errorf("REDIS | failed to get value: %v", err)
	}

	var recording Recording
	err = json.Unmarshal([]byte(val), &recording)
	if err != nil {
		return Recording{}, fmt.Errorf("REDIS | failed to unmarshal recording: %v", err)
	}

	return recording, nil
}
