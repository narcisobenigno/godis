package client

type RedisCommand interface {
	ToResp() Resp
}

type SimpleRedisCommand struct {
	name string
}

func (this *SimpleRedisCommand) ToResp() Resp {
	return &RespBulkString{this.name}
}

type RedisCommandWithArgments struct {
	command   *SimpleRedisCommand
	arguments []Argument
}

func (this *RedisCommandWithArgments) ToResp() Resp {
	keysAsBulkString := make([]Resp, len(this.arguments))
	for i, k := range this.arguments {
		keysAsBulkString[i] = &RespBulkString{k}
	}
	return RespArray(append(
		[]Resp{this.command.ToResp()},
		keysAsBulkString...,
	))
}

type GetCommand struct {
	godis Godis
	key   Key
}

func (this *GetCommand) Execute() (string, error) {
	reply, err := this.godis.Execute(
		&RedisCommandWithArgments{
			&SimpleRedisCommand{"GET"},
			[]Argument{this.key},
		},
	)

	if err != nil {
		return "", err
	}

	return reply.String()
}

type SetCommand struct {
	godis Godis
	key   Key
	value Argument
}

func (this *SetCommand) Execute() error {
	command := &RedisCommandWithArgments{
		&SimpleRedisCommand{"SET"},
		[]Argument{this.key, this.value},
	}
	_, err := this.godis.Execute(command)

	if err != nil {
		return err
	}

	return nil
}

type ExistsCommand struct {
	godis Godis
	key   Key
	keys  []Key
}

func (this *ExistsCommand) Execute() (int64, error) {
	reply, err := this.godis.Execute(
		&RedisCommandWithArgments{
			&SimpleRedisCommand{"EXISTS"},
			append([]Argument{this.key}, this.keys...),
		},
	)
	if err != nil {
		return -1, err
	}

	return reply.Integer()
}

type DelCommand struct {
	godis Godis
	key   Key
	keys  []Key
}

func (this *DelCommand) Execute() (int64, error) {
	reply, err := this.godis.Execute(
		&RedisCommandWithArgments{
			&SimpleRedisCommand{"DEL"},
			append([]Argument{this.key}, this.keys...),
		},
	)

	if err != nil {
		return -1, err
	}

	return reply.Integer()
}

type FlushDbCommand struct {
	godis Godis
}

func (this *FlushDbCommand) Execute() error {
	_, err := this.godis.Execute(&SimpleRedisCommand{"FLUSHDB"})
	return err
}
