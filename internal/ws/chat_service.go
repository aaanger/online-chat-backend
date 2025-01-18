package ws

type IChatService interface {
	CreateChat(name string, userID int) (*Chat, error)
	SaveChat(chatID, userID int, message *Message) error
	JoinChat(userID, roomID int) (*Chat, error)
}

type ChatService struct {
	repo IChatRepository
}

func NewChatService(repo IChatRepository) *ChatService {
	return &ChatService{
		repo: repo,
	}
}

func (s *ChatService) CreateChat(name string, userID int) (*Chat, error) {
	return s.repo.CreateChat(name, userID)
}

func (s *ChatService) SaveChat(chatID, userID int, message *Message) error {
	return s.repo.SaveChat(chatID, userID, message)
}

func (s *ChatService) JoinChat(userID, roomID int) (*Chat, error) {
	return s.repo.JoinChat(userID, roomID)
}
