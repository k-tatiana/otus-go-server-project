package transport

import (
	"context"
	"fmt"

	"github.com/tarantool/go-tarantool/v2"
)

type Logger struct{}

func (l *Logger) Report(event tarantool.ConnLogKind, conn *tarantool.Connection, v ...interface{}) {
	// Implement logging logic here, e.g., log to console or a file
	// For simplicity, we'll just print to the console
	fmt.Printf("Tarantool Log - Event: %v, Connection: %v, Details: %v\n", event, conn, v)
}

func NewTarantoolLogger() *Logger {
	return &Logger{}
}

func NewTarantoolOpts() *tarantool.Opts {
	return &tarantool.Opts{
		Logger:      NewTarantoolLogger(),
		SkipSchema:  false,
		Concurrency: 2,
	}
}

type TarantoolRepository struct {
	conn *tarantool.Connection
}

func NewTarantoolConnection(ctx context.Context, host, port, user, pwd string, opts tarantool.Opts) *TarantoolRepository {
	dialer := tarantool.NetDialer{
		Address:  host + ":" + port,
		User:     user,
		Password: pwd,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		// Use log.Fatal to print the error and exit
		fmt.Printf("Failed to connect to Tarantool: %v\n", err)
		// Optionally, you can use log.Fatal(err) if you import "log"
		panic(err)
	}

	// Проверка подключения
	_, err = conn.Ping()
	if err != nil {
		fmt.Printf("failed to ping Tarantool: %w", err)
		panic(err)
	}

	return &TarantoolRepository{conn: conn}
}

// InitializeSpace создает пространство и индексы
func (r *TarantoolRepository) InitializeSpace(queries []string) error {

	for _, query := range queries {
		_, err := r.conn.Eval(query, []interface{}{})
		if err != nil {
			return fmt.Errorf("failed to execute query %s: %w", query, err)
		}
	}

	return nil
}

func (r *TarantoolRepository) Clean(space string) error {
	_, err := r.conn.Eval(fmt.Sprintf("box.space.%s:truncate()", space), []interface{}{})
	if err != nil {
		return fmt.Errorf("failed to truncate space %s: %w", space, err)
	}
	return nil
}

// Вызов Lua функций
func (r *TarantoolRepository) Insert(functionName string, args []any) ([]any, error) {
	req := tarantool.NewInsertRequest(functionName).Tuple(args)
	fut := r.conn.Do(req)
	resp, err := fut.Get()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Пакетные операции
func (r *TarantoolRepository) GetBySecondary(spaceName string, whereClause uint64) ([]any, error) {
	op := tarantool.NewSelectRequest(spaceName).
		Index("secondary").
		Key([]interface{}{whereClause}).
		Iterator(tarantool.IterEq)

	fut := r.conn.Do(op)
	resp, err := fut.Get()
	if err != nil {
		return nil, err
	}
	return resp, nil
}
