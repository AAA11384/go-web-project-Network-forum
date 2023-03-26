package snowFlake

import (
	"time"

	"github.com/spf13/viper"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func Init() (err error) {
	var st time.Time
	st, err = time.ParseInLocation("2006-01-02 15:04:05", viper.GetString("snow.start_time"), time.Local)
	if err != nil {
		return
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(viper.GetInt64("snow.machine_id"))
	if err != nil {
		return
	}
	return
}

func GenId() int64 {
	return node.Generate().Int64()
}
