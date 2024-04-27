package port

type ClientAI interface {
	CreateChatCompletion(msg string) (string, error)
	NewToken(token string)
}
