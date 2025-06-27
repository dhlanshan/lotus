package snowflake

type iSnowWorker interface {
	NextId() int64
}
