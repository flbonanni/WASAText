package datamodels

type Conversation struct {
	ConversationID string   `json:"conversation_id"`
	Participants   []string `json:"participants"`
	LastMessage    string   `json:"last_message,omitempty"` // omitempty allows the field to be optional
}

func (c *Conversation) ConvFromDatabase(conv database.Conversation) {
	c.ConversationID = conv.ConversationID
	c.Participants = conv.Participants
	c.LastMessage = conv.LastMessage
}
