package service

import (
	"bufio"
	"context"
	"math/rand"
	"os"
	"time"

	sttv1 "github.com/skinkvi/protosSTT/gen/go/stt"
)

//TODO: Сейчас тут используется файл words.txt для хранения слов, они потом записываются в слайс и потом выводятся это все вообще не круто нужно придумать какой то другой способ для того что бы хранить слова и выдалвать эти слова прользователю

//TODO: Сделать что то нормальное  с ошибками мне не нравиться реализация по типу
// if err != nil {
// 	return nil, err
// }
// вот как по мне это прям гавнище так что нужно что то сделать

type SpeedTypingTestService struct {
	sttv1.UnimplementedSpeedTypingTestServer
	words []string
	// здесь нужно записать зависимости к примеру storage
}

func NewSpeedTypingTestService() (*SpeedTypingTestService, error) {
	words, err := readWordsFromFile("words.txt")
	if err != nil {
		return nil, err
	}

	return &SpeedTypingTestService{
		words: words,
	}, nil
}

func (s *SpeedTypingTestService) StartTest(ctx context.Context, req *sttv1.StartTestRequest) (*sttv1.StartTestResponse, error) {
	rand.Seed(time.Now().UnixNano())
	wordIndex := rand.Intn(len(s.words))
	word := s.words[wordIndex]

	return &sttv1.StartTestResponse{
		UserId:     req.UserId,
		TextToType: word,
	}, nil
}

func (s *SpeedTypingTestService) SubmitTest(ctx context.Context, req *sttv1.SubmitTestRequest) (*sttv1.SubmitTestResponse, error) {
	// Реализуйте метод SubmitTest
	return nil, nil
}

func (s *SpeedTypingTestService) EndTest(ctx context.Context, req *sttv1.EndTestRequest) (*sttv1.EndTestResponse, error) {
	// Реализуйте метод EndTest
	return nil, nil
}

func readWordsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	return words, scanner.Err()
}
