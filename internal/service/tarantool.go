package service

import (
	"context"
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	dt "github.com/tarantool/go-tarantool/v2/datetime"

	"otus/go-server-project/internal/config"
	"otus/go-server-project/internal/models"
	repo "otus/go-server-project/internal/transport"
)

const spaceName = "dialogs"

type storageRepo interface {
	ListFullDialogs(ctx context.Context) ([]models.DialogDTO, error)
}

type tarantoolRepo interface {
	GetBySecondary(string, uint64) ([]any, error)
	Insert(string, []any) ([]any, error)
	Clean(string) error
}

type TarantoolService struct {
	conn   tarantoolRepo
	source storageRepo
}

func internalConfig() []string {
	// Создание пространства users
	createSpace := `
        box.schema.space.create('dialogs', {
            if_not_exists = true,
            format = {
                {name = 'from_user_id', type = 'string'},
				{name = 'to_user_id', type = 'string'},
                {name = 'message', type = 'string'},
                {name = 'created_at', type = 'datetime'},
				{name = 'user_pair_hash', type = 'unsigned'}
            }
        })`

	// Create a primary index (unique by default)
	createPrimaryIndex := `
		box.space.dialogs:create_index('primary', {
			parts = { 'user_pair_hash', 'created_at' },
			type = 'TREE',
			unique = true
		})`

	// Создание индекса по user_pair_hash
	createFromUserIndex := `
        box.space.dialogs:create_index('secondary', {
            if_not_exists = true,
            parts = {'user_pair_hash'},
			type = 'TREE',
			unique = false
        })
	`

	return []string{createSpace, createPrimaryIndex, createFromUserIndex}
}

func absHashText(s1, s2 string) uint32 {
	// Get lexicographically least and greatest strings
	least := s1
	greatest := s2
	if strings.Compare(s1, s2) > 0 {
		least, greatest = s2, s1
	}

	// Concatenate least and greatest
	combined := least + greatest

	// Compute FNV-1a hash (32-bit)
	hasher := fnv.New32a()
	hasher.Write([]byte(combined))
	hashValue := int32(hasher.Sum32())

	// Return absolute value
	if hashValue < 0 {
		return uint32(-hashValue)
	}
	return uint32(hashValue)
}

func NewTarantoolService(ctx context.Context, cfg config.Tarantool, storage storageRepo) *TarantoolService {
	conn := repo.NewTarantoolConnection(ctx, cfg.Host, cfg.Port, cfg.User, cfg.Password, *repo.NewTarantoolOpts())
	_ = conn.InitializeSpace(internalConfig())
	t := &TarantoolService{conn: conn, source: storage}
	t.LoadDialogs(ctx)
	return t
}

func (r *TarantoolService) ListDialogs(ctx context.Context, user1 string, user2 string,
) ([]models.Dialog, error) {
	dialogs := make([]models.Dialog, 0)

	ds, err := r.conn.GetBySecondary(spaceName, uint64(absHashText(user1, user2)))
	if err != nil {
		return nil, fmt.Errorf("batch get dialogs via Lua function call failed: %w", err)
	}

	for _, d := range ds {
		if len(d.([]any)) == 0 {
			continue
		}
		dialog := models.Dialog{
			From: d.([]any)[0].(string),
			To:   d.([]any)[1].(string),
			Text: d.([]any)[2].(string),
		}
		dialogs = append(dialogs, dialog)
	}
	return dialogs, nil
}

func (r *TarantoolService) SendMessage(
	ctx context.Context, fromUserID, toUserID, message string, createdAt *time.Time,
) error {
	var (
		ct  dt.Datetime
		err error
	)
	if createdAt == nil {
		t := time.Now().UTC()
		createdAt = &t
	}

	ct, err = dt.MakeDatetime(*createdAt)
	if err != nil {
		return fmt.Errorf("failed to convert time to datetime: %w", err)
	}

	args := []any{
		fromUserID,
		toUserID,
		message,
		ct,
		absHashText(fromUserID, toUserID),
	}

	_, err = r.conn.Insert("dialogs", args)
	if err != nil {
		return fmt.Errorf("create dialog via Lua function call failed: %w", err)
	}

	return err
}

func (r *TarantoolService) LoadDialogs(ctx context.Context) {
	r.CleanDialogs()
	dialogs, err := r.source.ListFullDialogs(ctx)
	if err != nil {
		fmt.Printf("failed to load dialogs from DB: %v", err)
		return
	}
	for _, dialog := range dialogs {
		err := r.SendMessage(ctx, dialog.FromUserID, dialog.ToUserID, dialog.Message, dialog.CreatedAt)
		if err != nil {
			fmt.Printf("failed to create dialog in Tarantool: %v", err)
		}
	}
}

func (r *TarantoolService) CleanDialogs() {
	err := r.conn.Clean("dialogs")
	if err != nil {
		fmt.Printf("failed to clean dialogs space in Tarantool: %v", err)
	}
}
