package models

import "time"

// Postgres reference
type AchievementReference struct {
	ID                 string     `json:"id"`
	StudentID          string     `json:"student_id"`
	MongoAchievementID string     `json:"mongo_achievement_id"`
	Status             string     `json:"status"`
	SubmittedAt        *time.Time `json:"submitted_at"`
	VerifiedAt         *time.Time `json:"verified_at"`
	VerifiedBy         *string    `json:"verified_by"`
	RejectionNote      *string    `json:"rejection_note"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`
}

// DTOs for requests/responses
type CreateAchievementRequest struct {
	StudentID string                 `json:"student_id" validate:"required"`
	Doc map[string]interface{} `json:"doc"`
}

type UpdateAchievementRequest struct {
	MongoAchievementID *string `json:"mongo_achievement_id,omitempty"`
	
}

type SubmitRequest struct {
	// empty usually, but kept for extensibility
}

type RejectRequest struct {
	Note string `json:"note" validate:"required"`
}

type AchievementResponse struct {
	ID                 string                 `json:"id"`
	StudentID          string                 `json:"student_id"`
	MongoAchievementID string                 `json:"mongo_achievement_id"`
	Status             string                 `json:"status"`
	Doc                map[string]interface{} `json:"doc,omitempty"` 
	SubmittedAt        *time.Time             `json:"submitted_at"`
	VerifiedAt         *time.Time             `json:"verified_at"`
	VerifiedBy         *string                `json:"verified_by"`
	RejectionNote      *string                `json:"rejection_note"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          *time.Time             `json:"updated_at"`
}

// Mongo document (what we store in Mongo)
type MongoAchievement struct {
	ID          interface{}            `bson:"_id,omitempty" json:"id"`
	StudentID   string                 `bson:"student_id" json:"student_id"`
	Title       string                 `bson:"title,omitempty" json:"title,omitempty"`
	Description string                 `bson:"description,omitempty" json:"description,omitempty"`
	Files       []map[string]interface{} `bson:"files,omitempty" json:"files,omitempty"` // attachments metadata
	Extra       map[string]interface{} `bson:"extra,omitempty" json:"extra,omitempty"`
	CreatedAt   time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt   *time.Time             `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
