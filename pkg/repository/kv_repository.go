package repository

import (
	"context"
	"time"

	"github.com/benlocal/lai-panel/pkg/database"
	"github.com/hertz-contrib/cache/persist"
	"github.com/jmoiron/sqlx"
)

type KvRepository struct {
	db    *sqlx.DB
	cache *persist.MemoryStore
}

func NewKvRepository() *KvRepository {
	memoryStore := persist.NewMemoryStore(1 * time.Minute)
	return &KvRepository{
		db:    database.GetDB(),
		cache: memoryStore,
	}
}

func (r *KvRepository) Create(key string, value string, subKey *string) error {
	query := `INSERT INTO kv (key, value) VALUES (:key, :value)`
	if subKey != nil {
		query = `INSERT INTO kv (key, sub_key, value) VALUES (:key, :sub_key, :value)`
	}

	r.cache.Delete(context.Background(), key)
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"key":     key,
		"value":   value,
		"sub_key": subKey,
	})
	return err
}

func (r *KvRepository) Update(key string, value string) error {
	_, err := r.db.Exec(`UPDATE kv SET value = ? WHERE key = ?`, value, key)
	if err == nil {
		r.cache.Set(context.Background(), key, value, time.Duration(1)*time.Minute)
	}
	return err
}

func (r *KvRepository) Get(key string) (string, error) {
	c := ""
	err := r.cache.Get(context.Background(), key, &c)
	if err == nil && c != "" {
		return c, nil
	}

	query := `SELECT value FROM kv WHERE key = ?`
	var value string
	err = r.db.Get(&value, query, key)
	if err == nil {
		r.cache.Set(context.Background(), key, value, time.Duration(1)*time.Minute)
	}
	return value, err
}

func (r *KvRepository) GetWithSubKey(subKey string) (map[string]string, error) {
	query := `SELECT key, value FROM kv WHERE sub_key = ?`
	var values []struct {
		Key   string `db:"key"`
		Value string `db:"value"`
	}
	err := r.db.Select(&values, query, subKey)

	result := make(map[string]string)
	if err == nil {
		for _, v := range values {
			result[v.Key] = v.Value
			r.cache.Set(context.Background(), v.Key, v.Value, time.Duration(1)*time.Minute)
		}
	}

	return result, err
}

func (r *KvRepository) Delete(key string) error {
	_, err := r.db.Exec(`DELETE FROM kv WHERE key = ?`, key)
	if err == nil {
		r.cache.Delete(context.Background(), key)
	}
	return err
}

func (r *KvRepository) DeleteWithSubKey(subKey string) error {
	_, err := r.db.Exec(`DELETE FROM kv WHERE sub_key = ?`, subKey)
	if err == nil {
		r.cache.Cache.Purge()
	}
	return err
}
