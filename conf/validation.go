package conf

type ValidationConfig struct {
	UsernameRegex string `json:"username_regex" split_words:"true" default:"^[a-zA-Z0-9_]{2,255}$"`
	PasswordRegex string `json:"password_regex" split_words:"true" default:"^[a-zA-Z0-9[:punct:]]{8,28}$"`
}
