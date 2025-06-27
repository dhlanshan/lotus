package idgen

import snowflake "github.com/dhlanshan/lotus/idgen/snowflake_inter"

func GenSnowflakeId() int64 {
	return snowflake.NextId()
}
