package repository

import (
	"context"
	"database/sql"

	"github.com/innovember/real-time-forum/internal/chat"
	"github.com/innovember/real-time-forum/internal/consts"
	"github.com/innovember/real-time-forum/internal/models"
)

type RoomRepository struct {
	dbConn *sql.DB
}

func NewRoomRepository(conn *sql.DB) chat.RoomRepository {
	return &RoomRepository{
		dbConn: conn,
	}
}

func (rr *RoomRepository) InsertRoom(userID1, userID2 int64) (*models.Room, error) {
	var (
		ctx    context.Context
		tx     *sql.Tx
		result sql.Result
		err    error
		room   = &models.Room{}
		user   = &models.User{}
	)
	ctx = context.Background()
	if tx, err = rr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if result, err = tx.Exec(`
	INSERT INTO rooms
	DEFAULT VALUES`); err != nil {
		tx.Rollback()
		return nil, err
	}
	room.ID, err = result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if _, err = tx.Exec(`INSERT INTO room(room_id, user_id)
						VALUES(?,?)`, room.ID, userID1); err != nil {
		tx.Rollback()
		return nil, err
	}
	if _, err = tx.Exec(`INSERT INTO room(room_id, user_id)
						VALUES(?,?)`, room.ID, userID2); err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.QueryRow(`
						SELECT id, nickname
						FROM users
						WHERE id = ?`, userID2).Scan(&user.ID, &user.Nickname)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	room.User = user
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return room, nil
}

func (rr *RoomRepository) SelectRoomByUsers(userID1, userID2 int64) (int64, error) {
	var (
		ctx    context.Context
		tx     *sql.Tx
		err    error
		roomID int64
	)
	ctx = context.Background()
	if tx, err = rr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return 0, err
	}
	if err = tx.QueryRow(`SELECT room_id
							FROM room
							WHERE user_id IN (?,?)
							GROUP BY room_id
							HAVING COUNT (*) > 1;
						 `, userID1, userID2).Scan(
		&roomID); err != nil {
		tx.Rollback()
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return roomID, nil
}

func (rr *RoomRepository) SelectUsersByRoom(roomID int64) ([]models.User, error) {
	var (
		ctx   context.Context
		tx    *sql.Tx
		err   error
		users []models.User
		rows  *sql.Rows
	)
	ctx = context.Background()
	if tx, err = rr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	rows, err = tx.Query(`SELECT id, nickname
							FROM users
							WHERE id IN (
								SELECT user_id
								FROM room
								WHERE room_id = ?
							);
						 `, roomID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Nickname)
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return users, nil
}

func (rr *RoomRepository) SelectAllUsers(userID int64) ([]*models.User, error) {
	var (
		ctx   context.Context
		tx    *sql.Tx
		err   error
		rows  *sql.Rows
		users []*models.User
	)
	ctx = context.Background()
	if tx, err = rr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	rows, err = tx.Query(`SELECT id, nickname
	 FROM users
	  WHERE id != ?
	  ORDER BY nickname ASC`, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Nickname)
		users = append(users, &u)
	}
	if err = rows.Err(); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return users, nil
}

func (rr *RoomRepository) SelectOnlineUsers(userID int64) ([]*models.User, error) {
	var (
		ctx   context.Context
		tx    *sql.Tx
		err   error
		rows  *sql.Rows
		users []*models.User
	)
	ctx = context.Background()
	if tx, err = rr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	rows, err = tx.Query(`SELECT id, nickname
	 FROM users
	  WHERE id != ?
	  AND status = ?
	  ORDER BY nickname ASC`, userID, consts.StatusOnline)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Nickname)
		users = append(users, &u)
	}
	if err = rows.Err(); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return users, nil
}

func (rr *RoomRepository) DeleteRoom(id int64) error {
	var (
		ctx context.Context
		tx  *sql.Tx
		err error
	)
	ctx = context.Background()
	if tx, err = rr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM rooms
		WHERE id = ?`, id); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec(`DELETE FROM room
		WHERE room_id = ?`, id); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec(`DELETE FROM messages
		WHERE room_id = ?`, id); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rr *RoomRepository) InsertMessage(msg *models.Message) (*models.Message, error) {
	var (
		ctx    context.Context
		tx     *sql.Tx
		result sql.Result
		err    error
	)
	ctx = context.Background()
	if tx, err = rr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if result, err = tx.Exec(`INSERT INTO messages(room_id,author_id, message,message_date)
								VALUES(?,?,?,?)`,
		msg.RoomID,
		msg.User.ID,
		msg.Content,
		msg.MessageDate); err != nil {
		tx.Rollback()
		return nil, err
	}
	if msg.ID, err = result.LastInsertId(); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return msg, nil
}

func (rr *RoomRepository) SelectMessages(roomID, lastMessageID, userID int64) ([]models.Message, error) {
	var (
		ctx      context.Context
		tx       *sql.Tx
		err      error
		rows     *sql.Rows
		messages []models.Message
		total    int64
	)
	ctx = context.Background()
	if tx, err = rr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if err = tx.QueryRow(`SELECT count(id) AS total
	 					FROM messages;
						 `).Scan(
		&total); err != nil {
		tx.Rollback()
		return nil, err
	}
	if lastMessageID == 0 {
		lastMessageID = total + 1
	}
	rows, err = tx.Query(`SELECT m.id, m.room_id, m.message,
							m.message_date, m.read,
							u.id, u.nickname
							FROM messages as m
							LEFT JOIN users as u
							ON m.author_id = u.id
							WHERE m.room_id = $1
							AND m.id < $2
							ORDER BY m.id DESC , m.message_date DESC
							LIMIT 10
						 `, roomID, lastMessageID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			m models.Message
			u models.User
		)
		rows.Scan(&m.ID, &m.RoomID, &m.Content, &m.MessageDate,
			&m.Read, &u.ID, &u.Nickname)
		m.User = &u
		if m.User.ID == userID {
			m.IsYourMessage = true
		}
		messages = append(messages, m)
	}
	if err = rows.Err(); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return messages, nil
}

func (rr *RoomRepository) SelectLastMessageDate(roomID int64) (int64, error) {
	var (
		ctx             context.Context
		tx              *sql.Tx
		err             error
		lastMessageDate int64
	)
	ctx = context.Background()
	if tx, err = rr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return 0, err
	}
	if err = tx.QueryRow(`SELECT message_date
							FROM messages
							WHERE room_id = ?
							ORDER BY message_date DESC
							LIMIT 1
						 `, roomID).Scan(
		&lastMessageDate); err != nil {
		tx.Rollback()
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return lastMessageDate, nil
}

func (rr *RoomRepository) SelectRoomByID(roomID int64) (*models.Room, error) {
	var (
		ctx  context.Context
		tx   *sql.Tx
		err  error
		room models.Room
	)
	ctx = context.Background()
	if tx, err = rr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if err = tx.QueryRow(`SELECT id
							FROM rooms
							WHERE id = ?;
						 `, roomID).Scan(&room.ID); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &room, nil
}

func (rr *RoomRepository) SelectUnReadMessages(roomID int64) (int64, error) {
	var (
		ctx   context.Context
		tx    *sql.Tx
		err   error
		total int64
	)
	ctx = context.Background()
	if tx, err = rr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return 0, err
	}
	if err = tx.QueryRow(`SELECT count(id) AS total
	 					FROM messages
						 WHERE room_id = ?
						 AND read = 0;
						 `, roomID).Scan(
		&total); err != nil {
		tx.Rollback()
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return total, nil
}

func (rr *RoomRepository) UpdateMessageStatus(roomID, messageID int64) error {
	var (
		ctx context.Context
		tx  *sql.Tx
		err error
	)
	ctx = context.Background()
	if tx, err = rr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`UPDATE messages
						 SET read = 1
						 WHERE room_id = ?
						 AND id = ?`, roomID, messageID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rr *RoomRepository) UpdateMessagesStatusForReceiver(roomID, userID int64) error {
	var (
		ctx context.Context
		tx  *sql.Tx
		err error
	)
	ctx = context.Background()
	if tx, err = rr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return err
	}
	if _, err = tx.Exec(`UPDATE messages
						 SET read = 1
						 WHERE room_id = ?
						 AND author_id = ?`, roomID, userID); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rr *RoomRepository) GetLastMessage(roomID int64) (*models.Message, error) {
	var (
		ctx         context.Context
		tx          *sql.Tx
		err         error
		lastMessage models.Message
	)
	ctx = context.Background()
	if tx, err = rr.dbConn.BeginTx(ctx, &sql.TxOptions{}); err != nil {
		return nil, err
	}
	if err = tx.QueryRow(`SELECT m.id, m.room_id, m.message,
							m.message_date, m.read
							FROM messages AS m
							WHERE room_id = ?
							ORDER BY m.id DESC , m.message_date DESC 
							LIMIT 1
						 `, roomID).Scan(
		&lastMessage.ID, &lastMessage.RoomID,
		&lastMessage.Content, &lastMessage.MessageDate,
		&lastMessage.Read); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &lastMessage, nil
}
