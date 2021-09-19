package object

type (
	StatusID = int64

	Status struct {
		ID        StatusID  `json:"id"`
		AccountID AccountID `json:"-" db:"account_id"`
		Content   string    `json:"content"`
		CreateAt  DateTime  `json:"create_at,omitempty" db:"create_at"`
		DeleteAt  *DateTime `json:"-" db:"delete_at"`

		Account         *Account           `json:"account,omitempty"`
		MediaAttachment []*MediaAttachment `json:"media_attachments,omitempty"`
	}
)
