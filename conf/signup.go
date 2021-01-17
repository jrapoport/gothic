package conf

type SignupConfig struct {
	// Code if true, require a signup code.
	Code bool `json:"code"`
	// Defaults if true, generate a random username & color on signup if a username is not provided.
	Defaults bool `json:"defaults"`
}
