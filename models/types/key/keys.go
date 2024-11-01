package key

// Key type
type Key string

// Defined keys
const (
	Unknown        Key = ""
	AccountID          = "account_id"
	Action             = "action"
	AdminID            = "admin_id"
	AvatarURL          = "avatar_url"
	Class              = "class"
	Code               = "code"
	Color              = "color"
	ConfirmedAt        = "confirmed_at"
	Count              = "count"
	Data               = "data"
	Description        = "description"
	Email              = "email"
	Event              = "event"
	ExpirationDate     = "expiration_date"
	Fields             = "fields"
	Filters            = "filters"
	FirstName          = "first_name"
	Hard               = "hard"
	Hostname           = "hostname"
	ID                 = "id"
	IPAddress          = "ip_address"
	Issued             = "issued"
	JWT                = "jwt"
	LastName           = "last_name"
	LastUsed           = "last_used"
	Metadata           = "metadata"
	Name               = "name"
	Nickname           = "nickname"
	Page               = "page"
	PageSize           = "page_size"
	Password           = "password"
	PerPage            = "per_page"
	Provider           = "provider"
	ReCaptcha          = "recaptcha"
	Revoked            = "revoked"
	Role               = "role"
	Service            = "service"
	Session            = "session"
	Sort               = "sort"
	State              = "state"
	Status             = "status"
	Timestamp          = "timestamp"
	Token              = "token"
	Type               = "type"
	UserID             = "user_id"
	Username           = "username"
	Uses               = "uses"
	Valid              = "valid"
)

var Reserved = map[string]struct{}{
	IPAddress: {},
}
