package mattermost

// Subscriber is any probe with a subcribe function
type Subscriber interface {
	GetSubscription() *WebSocketSubscription
}
