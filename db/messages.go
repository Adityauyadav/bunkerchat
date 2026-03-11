package db

import "github.com/adityauyadav/bunkerchat/models"

func SaveMessage(sentFromID int, sentToID int, content string) error {
	query := `INSERT INTO messages (sent_from_id, sent_to_id, content) VALUES(?,?,?)`
	_, err := DB.Exec(query, sentFromID, sentToID, content)
	return err
}

func GetConversation(sentFromID int, sentToID int) ([]models.Message, error) {
	query := `SELECT * FROM messages WHERE (sent_from_id = ? AND sent_to_id = ?) OR (sent_from_id = ? AND sent_to_id = ?) ORDER BY created_at ASC`
	rows, err := DB.Query(query, sentFromID, sentToID, sentToID, sentFromID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		err := rows.Scan(&message.ID, &message.SentFromID, &message.SentToID, &message.Content, &message.Read, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
