package ws

type ChatService struct {
	repo *ChatRepository
}

func (s *ChatService) CreateChat(name string, userID int) (*Chat, error) {
	return s.repo.CreateChat(name, userID)
}

func (s *ChatService) SaveChat(userID int, message *Message) error {
	return s.repo.SaveChat(userID, message)
}

func (s *ChatService) JoinChat(userID, roomID int) (*Chat, error) {
	return s.repo.JoinChat(userID, roomID)
}
