func (c *Conversation) ConvFromDatabase(conv database.Conversation) {
	c.ConversationID = conv.ConversationID
	c.Participants = conv.Participants
	c.LastMessage = conv.LastMessage
}
