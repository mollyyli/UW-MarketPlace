package sessions

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

//RedisStore represents a session.Store backed by redis.
type RedisStore struct {
	//Redis client used to talk to redis server.
	Client *redis.Client
	//Used for key expiry time on redis.
	SessionDuration time.Duration
}

//NewRedisStore constructs a new RedisStore
func NewRedisStore(client *redis.Client, sessionDuration time.Duration) *RedisStore {
	//initialize and return a new RedisStore struct
	return &RedisStore{
		client,
		sessionDuration,
	}
}

//Store implementation

//Save saves the provided `sessionState` and associated SessionID to the store.
//The `sessionState` parameter is typically a pointer to a struct containing
//all the data you want to associated with the given SessionID.
func (rs *RedisStore) Save(sid SessionID, sessionState interface{}) error {
	//TODO: marshal the `sessionState` to JSON and save it in the redis database,
	//using `sid.getRedisKey()` for the key.
	//return any errors that occur along the way.
	newSession, err := json.Marshal(sessionState)

	if err != nil {
		return err
	}
	key := sid.getRedisKey()
	cmd := rs.Client.Set(key, newSession, rs.SessionDuration)

	if cmd.Err() != nil {
		return ErrStateNotFound
	}
	return nil
}

//Get populates `sessionState` with the data previously saved
//for the given SessionID
func (rs *RedisStore) Get(sid SessionID, sessionState interface{}) error {
	//TODO: get the previously-saved session state data from redis,
	//unmarshal it back into the `sessionState` parameter
	//and reset the expiry time, so that it doesn't get deleted until
	//the SessionDuration has elapsed.

	prev := rs.Client.Get(sid.getRedisKey())

	err := prev.Err()
	if err != nil {
		return ErrStateNotFound
	}
	buffer, err := prev.Result()
	if err != nil {
		return ErrStateNotFound
	}
	if err := json.Unmarshal([]byte(buffer), sessionState); err != nil {
		return fmt.Errorf("error JSON unmarshal: %v", err)
	}
	log.Println("session state marshal", sessionState)
	cmd := rs.Client.Expire(sid.getRedisKey(), rs.SessionDuration)
	if cmd.Err() != nil {
		return ErrStateNotFound
	}
	return nil
}

//Delete deletes all state data associated with the SessionID from the store.
func (rs *RedisStore) Delete(sid SessionID) error {
	//TODO: delete the data stored in redis for the provided SessionID
	getErr := rs.Client.Get(sid.getRedisKey()).Err()
	if getErr != nil {
		return getErr
	}
	deleted := rs.Client.Del(sid.getRedisKey())
	deleteErr := deleted.Err()
	if deleteErr != nil {
		return deleteErr
	}
	return nil
}

//getRedisKey() returns the redis key to use for the SessionID
func (sid SessionID) getRedisKey() string {
	//convert the SessionID to a string and add the prefix "sid:" to keep
	//SessionID keys separate from other keys that might end up in this
	//redis instance
	return "sid:" + sid.String()
}
