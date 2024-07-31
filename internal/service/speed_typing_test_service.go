package service

import (
	"STTMain/internal/storage"

	sttv1 "github.com/skinkvi/protosSTT/gen/go/stt"
)

type SpeedTypingTestService struct {
	sttv1.UnimplementedSpeedTypingTestServer
	storage *storage.PostgresStorage
}

func NewSpeedTypingTestService(storage *storage.PostgresStorage) (*SpeedTypingTestService, error) {
	return &SpeedTypingTestService{
		storage: storage,
	}, nil
}

// func (s *SpeedTypingTestService) StartTest(ctx context.Context, req *sttv1.StartTestRequest) (*sttv1.StartTestResponse, error) {
// 	words, err := s.storage.GetRandomWords(1)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get random words: %w", err)
// 	}
// 	if len(words) == 0 {
// 		return nil, fmt.Errorf("no words found")
// 	}
// 	return &sttv1.StartTestResponse{
// 		UserId:     req.UserId,
// 		TextToType: words[0].Text,
// 	}, nil
// }

// func (s *SpeedTypingTestService) SubmitTest(ctx context.Context, req *sttv1.SubmitTestRequest) (*sttv1.SubmitTestResponse, error) {
// 	// Предположим, что text_to_type передается в контексте или хранится в сервисе
// 	textToType := "example text" // Замените на реальное значение
// 	isCorrect := req.TypedText == textToType
// 	accuracy := float32(1.0) // Замените на реальное значение
// 	speed := float32(1.0)    // Замените на реальное значение

// 	return &sttv1.SubmitTestResponse{
// 		Accuracy: accuracy,
// 		Speed:    speed,
// 	}, nil
// }

// func (s *SpeedTypingTestService) EndTest(ctx context.Context, req *sttv1.EndTestRequest) (*sttv1.EndTestResponse, error) {
// 	// Реализуйте метод EndTest
// 	return &sttv1.EndTestResponse{}, nil
// }
