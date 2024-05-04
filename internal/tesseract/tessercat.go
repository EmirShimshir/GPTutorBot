package tesseract

import (
	"github.com/EmirShimshir/tasker-bot/internal/config"
	"github.com/otiai10/gosseract/v2"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"sync"
)

type Nlp struct {
	client    *gosseract.Client
	languages []string
	mutex     *sync.Mutex
}

func NewNlp(cfg config.Tesseract) *Nlp {
	return &Nlp{
		client:    gosseract.NewClient(),
		languages: cfg.NlpLanguages,
		mutex:     &sync.Mutex{},
	}
}

func (n *Nlp) TextRecognition(urlFile string) (string, error) {
	resp, err := http.Get(urlFile)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	file, err := os.CreateTemp("", "")
	if err != nil {
		return "", err
	}

	defer file.Close()
	defer os.Remove(file.Name())

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	err = n.client.SetLanguage(n.languages...)
	if err != nil {
		return "", err
	}

	err = n.client.SetImage(file.Name())
	if err != nil {
		return "", err
	}

	n.mutex.Lock()
	log.WithFields(log.Fields{}).Info("start text")
	text, err := n.client.Text()
	if err != nil {
		return "", err
	}
	log.WithFields(log.Fields{}).Info("end text")
	n.mutex.Unlock()

	return text, nil
}

func (n *Nlp) CloseClient() {
	err := n.client.Close()
	if err != nil {
		log.Fatal(err)
	}
}
