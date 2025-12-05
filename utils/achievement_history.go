package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/Lutfania/ekrp/config"
)

// userID may be interface{} from c.Locals; caller should cast
func InsertAchievementHistory(achievementID, oldStatus, newStatus string, changedBy interface{}) error {
	var changedByID *string
	if v, ok := changedBy.(string); ok && v != "" {
		changedByID = &v
	}
	_, err := config.DB.Exec(context.Background(),
		`INSERT INTO achievement_reference_history (achievement_ref_id, old_status, new_status, changed_by, note, changed_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		achievementID, oldStatus, newStatus, changedByID, nil, time.Now())
	if err != nil {
		// not fatal; return error to caller
		return fmt.Errorf("history insert: %w", err)
	}
	return nil
}
