package repo

import (
	"database/sql"
	"fmt"
	"justTest/internal/models"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}

}
func (r *NotificationRepository) SaveNotification(notification *models.Notification) error {
	query := ` 
	insert into notifications ( user_id, type, title, message , 
	                           data , is_read , priority ) 
	values ($1,$2,$3,$4,$5,$6,$7)
	returning id;
`
	err := r.db.QueryRow(query,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.Data,
		notification.IsRead,
		notification.Priority,
	).Scan(&notification.ID)
	if err != nil {
		return err
	}
	return nil
}
func (r *NotificationRepository) GetNotificationByID(id int64) (*models.Notification, error) {
	query := ` 
	select id, user_id, type, title, message , 
	    data , is_read , priority, created_at, updated_at from notifications where id = $1;

`
	notification := &models.Notification{}
	err := r.db.QueryRow(query, id).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Type,
		&notification.Title,
		&notification.Message,
		&notification.Data,
		&notification.IsRead,
		&notification.Priority,

		&notification.CreatedAt,
		&notification.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (r *NotificationRepository) GetUserNotifications(userID string, limit, offset int) ([]*models.Notification, error) {
	query := `
select id, user_id, type, title, message , 
	    data , is_read , priority, created_at, updated_at 
from notifications where user_id = $1
 ORDER BY created_at DESC
 limit $2 offset $3;
`
	var notifications []*models.Notification
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying notifications: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		notification := &models.Notification{}
		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Type,
			&notification.Title,
			&notification.Message,
			&notification.Data,
			&notification.IsRead,
			&notification.Priority,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning notification: %v", err)
		}
		notifications = append(notifications, notification)

	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notifications: %v", err)
	}
	return notifications, nil
}

func (r *NotificationRepository) GetUnreadNotifications(userID string) ([]*models.Notification, error) {
	query := ` 
select id, user_id, type, title, message ,
data , is_read , priority, created_at, updated_at from notifications where user_id = $1 and is_read = 'false'
`
	var notifications []*models.Notification
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying notifications: %v", err)

	}
	defer rows.Close()
	for rows.Next() {
		notification := &models.Notification{}
		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Type,
			&notification.Title,
			&notification.Message,
			&notification.Data,

			&notification.IsRead,
			&notification.Priority,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning notification: %v", err)
		}
		notifications = append(notifications, notification)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notifications: %v", err)
	}
	return notifications, nil

}

func (r *NotificationRepository) MarkAsRead(id int64) error {
	query := ` 
update notifications set is_read = true where id = $1;`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error updating notifications: %v", err)
	}
	return nil
}
func (r *NotificationRepository) MarkAllAsRead(userID string) error {
	query := ` 
	update notifications set is_read = true where user_id = $1;`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("error updating notifications: %v", err)
	}
	return nil
}
func (r *NotificationRepository) DeleteNotification(id int64) error {
	query := ` 
Delete from notifications where id = $1;
`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting notifications: %v", err)
	}
	return nil

}
