package models

import (
	"bytes"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/conf"
	"github.com/jrapoport/gothic/storage"
	"github.com/vcraescu/go-paginator/v2"
	"github.com/vcraescu/go-paginator/v2/adapter"
)

type AuditAction string
type auditLogType string

const (
	LoginAction                     AuditAction = "login"
	LogoutAction                    AuditAction = "logout"
	InviteAcceptedAction            AuditAction = "invite_accepted"
	UserSignedUpAction              AuditAction = "user_signed_up"
	UserInvitedAction               AuditAction = "user_invited"
	UserDeletedAction               AuditAction = "user_deleted"
	UserModifiedAction              AuditAction = "user_modified"
	UserRecoveryRequestedAction     AuditAction = "user_recovery_requested"
	UserConfirmationRequestedAction AuditAction = "user_confirmation_requested"
	TokenRevokedAction              AuditAction = "token_revoked"
	TokenRefreshedAction            AuditAction = "token_refreshed"

	account auditLogType = "account"
	team    auditLogType = "team"
	token   auditLogType = "token"
	user    auditLogType = "user"
)

var actionLogTypeMap = map[AuditAction]auditLogType{
	LoginAction:                     account,
	LogoutAction:                    account,
	InviteAcceptedAction:            account,
	UserSignedUpAction:              team,
	UserInvitedAction:               team,
	UserDeletedAction:               team,
	TokenRevokedAction:              token,
	TokenRefreshedAction:            token,
	UserModifiedAction:              user,
	UserRecoveryRequestedAction:     user,
	UserConfirmationRequestedAction: user,
}

func init() {
	storage.AddMigration(&AuditLogEntry{})
}

// AuditLogEntry is the database model for audit log entries.
type AuditLogEntry struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:char(36)"`
	Payload   Map       `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

func (a AuditLogEntry) TableName() string {
	c := conf.Current()
	return storage.Namespace(c) + "audit_log"
}

func NewAuditLogEntry(tx *storage.Connection, actor *User, action AuditAction, traits map[string]interface{}) error {
	id := uuid.New()
	l := AuditLogEntry{
		ID: id,
		Payload: Map{
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
			"actor_id":    actor.ID,
			"actor_email": actor.Email,
			"action":      action,
			"log_type":    actionLogTypeMap[action],
		},
	}

	if name, ok := actor.UserMetaData["full_name"]; ok {
		l.Payload["actor_name"] = name
	}

	if traits != nil {
		l.Payload["traits"] = traits
	}

	if err := tx.Create(&l).Error; err != nil {
		err = fmt.Errorf("%w creating audit log entry", err)
		return err
	}
	return nil
}

func FindAuditLogEntries(tx *storage.Connection, filterColumns []string, filterValue string, pageParams *Pagination) ([]*AuditLogEntry, error) {
	q := tx.Model(AuditLogEntry{}).Order("created_at desc")

	if len(filterColumns) > 0 && filterValue != "" {
		lf := "%" + filterValue + "%"

		builder := bytes.NewBufferString("(")
		values := make([]interface{}, len(filterColumns))

		for idx, _ := range filterColumns {
			builder.WriteString("payload LIKE ?")
			values[idx] = lf

			if idx+1 < len(filterColumns) {
				builder.WriteString(" OR ")
			}
		}
		builder.WriteString(")")
		fmt.Println(builder.String())
		q = q.Where(builder.String(), values...)
	}

	var logs []*AuditLogEntry
	var err error
	if pageParams != nil {
		p := paginator.New(adapter.NewGORMAdapter(q), int(pageParams.PerPage))
		p.SetPage(int(pageParams.Page))
		if err = p.Results(&logs); err != nil {
			return nil, err
		}
		var cnt int
		if cnt, err = p.PageNums(); err != nil {
			return nil, err
		}
		pageParams.Count = uint64(cnt)
	} else {
		err = q.Find(&logs).Error
	}

	return logs, err
}