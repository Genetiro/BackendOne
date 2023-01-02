package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	User     string
	Password string
	DbName   string
	Host     string
	Port     int
	PoolMax  int
}

//Выбор фреймворка повторять необязательно, я обычно работаю с pgx, поэтому пример именно на нем. Суть коннекта одинакова и в других.

/*
Строка соединения для базы данных постгрес требует указать параметры подключения
Для этого используем структуру Config. Строка (URL для подключения) описана в документации:
https://pkg.go.dev/github.com/jackc/pgx/v4/pgxpool@v4.17.2#ParseConfig

# Example URL
postgres://{User}:{Password}@{Host}:{Port}/{DbName}?sslmode=verify-ca&pool_max_conns={PoolMax}

Юзер и пароль в случае с контейнерами обычно пишем в docker-compose.
DbName это имя базы данных, к  которой хотим подключиться. Это требуется, т.к. по одному адресу могут жить базы разных команд, например.
Host, Port нужны соответственно для того чтобы указать по какому адресу живет база.
PoolMax это максимальное количество соединений, которое мы можем положить в пул соединений *pgxpool.Pool. Этот параметр
связан с параметром базы данных max_connections - максимальное число коннектов к серверу базы данных. PoolMax не может
быть больше, чем max_connections.

Упрощенно это выглядит так: где-то в коде вызываем Connection - значит создается пул соединений к базе данных,
которая описана в конфиге. На первую переменную ctx context.Context можно пока не смотреть.
*/
func Connection(ctx context.Context, pgCfg Config) (*pgxpool.Pool, error) {
	cfgStr := fmt.Sprintf(
		"user=%s password=%s host=%s dbname=%s port=%d pool_max_conns=%d",
		pgCfg.User, pgCfg.Password, pgCfg.Host, pgCfg.DbName, pgCfg.Port, pgCfg.PoolMax)

	cfg, err := pgxpool.ParseConfig(cfgStr)
	if err != nil {
		return nil, err
	}

	pgConn, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return pgConn, nil
}
