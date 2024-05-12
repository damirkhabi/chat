package chat

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/arifullov/chat-server/internal/client/db"
	"github.com/arifullov/chat-server/internal/model"
	"github.com/arifullov/chat-server/internal/repository"
)

const (
	tableName         = "chats"
	messagesTableName = "messages"

	idColumn        = "id"
	usernamesColumn = "usernames"

	fromColumn      = "\"from\""
	textColumn      = "text"
	chatIDColumn    = "chat_id"
	timestampColumn = "timestamp"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, chat *model.Chat) (int64, error) {
	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(usernamesColumn).
		Values(chat.Usernames).
		Suffix("RETURNING id")
	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "chat_repository.Create",
		QueryRaw: query,
	}

	var chatID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chatID)
	if err != nil {
		return 0, err
	}
	return chatID, nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builderDelete := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return err
	}
	q := db.Query{
		Name:     "chat_repository.Delete",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return model.ErrChatNotFound
	}

	return nil
}

func (r *repo) CreateMessage(ctx context.Context, message *model.Message) error {
	builderInsert := sq.Insert("messages").
		PlaceholderFormat(sq.Dollar).
		Columns(fromColumn, textColumn, chatIDColumn, timestampColumn).
		Values(message.From, message.Text, message.ChatID, message.Timestamp).
		Suffix("RETURNING id")
	query, args, err := builderInsert.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chat_repository.CreateMessage",
		QueryRaw: query,
	}

	var chatID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chatID)
	if err != nil {
		return err
	}
	return nil
}
