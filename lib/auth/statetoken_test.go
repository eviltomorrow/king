package auth

// func init() {
// 	redis.DSN = "redis://:admin123@localhost:6379/0?protocol=3"
// 	if err := redis.Connect(); err != nil {
// 		log.Fatal(err)
// 	}
// }

// var (
// 	jwtToken      = "this.is.jwt-token"
// 	id            = "shepard"
// 	oldStateToken = "202e5a4a0164e068c729b23303c754d2951f7543e942bb6d55ad51e272692e17"
// 	newStateToken = "7d78c29795e5d565e5b4a5f9c9a09c6e6de38149d03a1493ad9b7ce8996460f9"
// )

// func init() {
// 	c := redis.RDB.Keys(context.Background(), "*")
// 	keys, err := c.Result()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if len(keys) != 0 {
// 		i := redis.RDB.Del(context.Background(), keys...)
// 		if err := i.Err(); err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// }

// func TestStateTokenWithParseJwtToken(t *testing.T) {
// 	_assert := assert.New(t)

// 	token, err := StateTokenWithParseJwtToken(jwtToken)
// 	_assert.Nil(err)
// 	_assert.Equal(token, oldStateToken)
// }

// func TestStateTokenWithCount(t *testing.T) {
// 	_assert := assert.New(t)

// 	count, err := StateTokenWithCount(context.Background(), id)
// 	_assert.Nil(err)
// 	_assert.Equal(int64(0), count)

// 	err = StateTokenWithRenew(context.Background(), "", newStateToken, id, 3*time.Second)
// 	_assert.Nil(err)

// 	count, err = StateTokenWithCount(context.Background(), id)
// 	_assert.Nil(err)
// 	_assert.Equal(int64(1), count)

// 	time.Sleep(5 * time.Second)
// 	count, err = StateTokenWithCount(context.Background(), id)
// 	_assert.Nil(err)
// 	_assert.Equal(int64(0), count)

// 	err = StateTokenWithRenew(context.Background(), "", newStateToken, id, 10*time.Second)
// 	_assert.Nil(err)

// 	count, err = StateTokenWithCount(context.Background(), id)
// 	_assert.Nil(err)
// 	_assert.Equal(int64(1), count)

// 	time.Sleep(5 * time.Second)
// 	count, err = StateTokenWithCount(context.Background(), id)
// 	_assert.Nil(err)
// 	_assert.Equal(int64(1), count)
// }

// func TestStateTokenWithSearchList(t *testing.T) {
// 	_assert := assert.New(t)

// 	err := StateTokenWithRenew(context.Background(), "", oldStateToken, id, 3*time.Second)
// 	_assert.Nil(err)

// 	result, err := StateTokenWithSearchList(context.Background(), id)
// 	_assert.Nil(err)
// 	_assert.Equal(1, len(result))
// 	t.Logf("%v", result)

// 	err = StateTokenWithRenew(context.Background(), "", newStateToken, id, 3*time.Second)
// 	_assert.Nil(err)

// 	result, err = StateTokenWithSearchList(context.Background(), id)
// 	_assert.Nil(err)
// 	_assert.Equal(2, len(result))
// 	t.Logf("%v", result)

// 	time.Sleep(5 * time.Second)

// 	result, err = StateTokenWithSearchList(context.Background(), id)
// 	_assert.Nil(err)
// 	_assert.Equal(0, len(result))
// 	t.Logf("%v", result)

// 	result, err = StateTokenWithSearchList(context.Background(), id)
// 	_assert.Nil(err)
// 	_assert.Equal(0, len(result))
// 	t.Logf("%v", result)
// }

// func TestStateTokenWithRenew(t *testing.T) {
// 	_assert := assert.New(t)

// 	err := StateTokenWithRenew(context.Background(), "", oldStateToken, id, 3*time.Second)
// 	_assert.Nil(err)
// 	ok, err := stateTokenWithSearch(context.Background(), oldStateToken)
// 	_assert.Nil(err)
// 	_assert.Equal(true, ok)
// 	time.Sleep(5 * time.Second)

// 	ok, err = stateTokenWithSearch(context.Background(), oldStateToken)
// 	_assert.Nil(err)
// 	_assert.Equal(false, ok)

// 	err = StateTokenWithRenew(context.Background(), oldStateToken, newStateToken, id, 3*time.Second)
// 	_assert.NotNil(err)

// 	err = StateTokenWithRenew(context.Background(), "", oldStateToken, id, 30*time.Second)
// 	_assert.Nil(err)
// 	err = StateTokenWithRenew(context.Background(), oldStateToken, newStateToken, id, 3*time.Second)
// 	_assert.Nil(err)
// 	ok, err = stateTokenWithSearch(context.Background(), oldStateToken)
// 	_assert.Nil(err)
// 	_assert.Equal(false, ok)
// 	ok, err = stateTokenWithSearch(context.Background(), newStateToken)
// 	_assert.Nil(err)
// 	_assert.Equal(true, ok)
// 	time.Sleep(5 * time.Second)
// 	ok, err = stateTokenWithSearch(context.Background(), newStateToken)
// 	_assert.Nil(err)
// 	_assert.Equal(false, ok)
// }

// func TestStateTokenWithRevokeAll(t *testing.T) {
// 	_assert := assert.New(t)

// 	err := StateTokenWithRenew(context.Background(), "", oldStateToken, id, 30*time.Second)
// 	_assert.Nil(err)

// 	ok, err := stateTokenWithSearch(context.Background(), oldStateToken)
// 	_assert.Nil(err)
// 	_assert.Equal(true, ok)

// 	err = StateTokenWithRevokeAll(context.Background(), id)
// 	_assert.Nil(err)

// 	ok, err = stateTokenWithSearch(context.Background(), oldStateToken)
// 	_assert.Nil(err)
// 	_assert.Equal(false, ok)

// 	err = StateTokenWithRenew(context.Background(), "", oldStateToken, id, 30*time.Second)
// 	_assert.Nil(err)
// 	err = StateTokenWithRenew(context.Background(), "", newStateToken, id, 30*time.Second)
// 	_assert.Nil(err)
// 	err = StateTokenWithRevokeAll(context.Background(), id)
// 	_assert.Nil(err)

// 	ok, err = stateTokenWithSearch(context.Background(), oldStateToken)
// 	_assert.Nil(err)
// 	_assert.Equal(false, ok)
// 	ok, err = stateTokenWithSearch(context.Background(), newStateToken)
// 	_assert.Nil(err)
// 	_assert.Equal(false, ok)
// }

// func TestStateTokenWithRevoke(t *testing.T) {
// 	_assert := assert.New(t)

// 	err := StateTokenWithRenew(context.Background(), "", oldStateToken, id, 30*time.Second)
// 	_assert.Nil(err)
// 	ok, err := stateTokenWithSearch(context.Background(), oldStateToken)
// 	_assert.Nil(err)
// 	_assert.Equal(true, ok)

// 	err = StateTokenWithRevoke(context.Background(), oldStateToken)
// 	_assert.Nil(err)
// 	ok, err = stateTokenWithSearch(context.Background(), oldStateToken)
// 	_assert.Nil(err)
// 	_assert.Equal(false, ok)

// 	err = StateTokenWithRenew(context.Background(), "", newStateToken, id, 30*time.Second)
// 	_assert.Nil(err)
// 	ok, err = stateTokenWithSearch(context.Background(), newStateToken)
// 	_assert.Nil(err)
// 	_assert.Equal(true, ok)
// }

// func stateTokenWithSearch(ctx context.Context, token string) (bool, error) {
// 	key := fmt.Sprintf("%s%s", tokenPrefix, token)
// 	c := redis.RDB.Exists(ctx, key)
// 	if err := c.Err(); err != nil {
// 		return false, err
// 	}
// 	ok, err := c.Result()
// 	if err != nil {
// 		return false, err
// 	}
// 	return ok == 1, nil
// }
