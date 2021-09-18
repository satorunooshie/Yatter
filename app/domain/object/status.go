package object

type Status struct {
	ID        int64     `json:"id"`
	AccountID AccountID `json:"-" db:"account_id"`
	Content   string    `json:"content"`
	CreateAt  DateTime  `json:"create_at,omitempty" db:"create_at"`
	DeleteAt  *DateTime `json:"-" db:"delete_at"`

	Account Account `json:"account,omitempty"`
}
