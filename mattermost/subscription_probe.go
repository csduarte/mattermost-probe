package mattermost

// SubscriptionProbe is any probe with a subcribe function
type SubscriptionProbe interface {
	GetSubscription() *WebSocketSubscription
}
