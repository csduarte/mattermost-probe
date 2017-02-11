package config

// Credentials store the email and password for a user in the system
type Credentials struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

// Valid will return true if the credentils have some value for email and password
func (c *Credentials) Valid() bool {
	return len(c.Email) > 0 && len(c.Password) > 0
}
