package mattermost

// LogError is helper for error logging
func (c *Client) LogError(items ...interface{}) {
	if c.Log != nil {
		c.Log.Error(items)
	}
}

// LogInfo is a helper for info logging
func (c *Client) LogInfo(items ...interface{}) {
	if c.Log != nil {
		c.Log.Info(items)
	}
}
