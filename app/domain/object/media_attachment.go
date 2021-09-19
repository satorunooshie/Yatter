package object

const (
	TypeImage = 1
)

type (
	MediaAttachmentID = int64

	MediaAttachment struct {
		ID          MediaAttachmentID `json:"-"`
		StatusID    StatusID          `json:"id" db:"status_id"`
		Type        int64             `json:"-" db:"type"`
		URL         string            `json:"url,omitempty" db:"url"`
		Description string            `json:"description,omitempty"`
		CreateAt    DateTime          `json:"create_at,omitempty" db:"create_at"`
		DeleteAt    *DateTime         `json:"-" db:"delete_at"`

		MediaType *string `json:"type,omitempty"`
	}
)

func (m *MediaAttachment) SetMediaType() {
	var (
		image = "image"
		other = "other"
	)
	switch m.Type {
	case TypeImage:
		m.MediaType = &image
	default:
		m.MediaType = &other
	}
}
