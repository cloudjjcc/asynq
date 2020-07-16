module github.com/cloudjjcc/asynq/tools

go 1.13

require (
	github.com/cloudjjcc/asynq v0.0.0-00010101000000-000000000000
	github.com/go-redis/redis/v7 v7.2.0
	github.com/google/uuid v1.1.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.6.2
)

replace github.com/cloudjjcc/asynq => ./..
