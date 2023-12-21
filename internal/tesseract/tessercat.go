package tesseract

import (
	"github.com/EmirShimshir/tasker-bot/internal/config"
)

type Nlp struct {
	//client    *gosseract.Client
	//languages []string
}

func NewNlp(cfg config.Tesseract) *Nlp {
	return &Nlp{
		//client:    gosseract.NewClient(),
		//languages: cfg.NlpLanguages,
	}
}

func (n *Nlp) TextRecognition(urlFile string) (string, error) {
	//resp, err := http.Get(urlFile)
	//if err != nil {
	//	return "", err
	//}
	//defer resp.Body.Close()
	//
	//file, err := os.CreateTemp("", "")
	//if err != nil {
	//	return "", err
	//}
	//defer file.Close()
	//defer os.Remove(file.Name())
	//
	//_, err = io.Copy(file, resp.Body)
	//if err != nil {
	//	return "", err
	//}
	//
	//err = n.client.SetLanguage(n.languages...)
	//if err != nil {
	//	return "", err
	//}
	//
	//err = n.client.SetImage(file.Name())
	//if err != nil {
	//	return "", err
	//}
	//
	//text, err := n.client.Text()
	//if err != nil {
	//	return "", err
	//}
	//
	//return text, nil

	return "", nil
}

func (n *Nlp) CloseClient() {
	//err := n.client.Close()
	//if err != nil {
	//	log.Fatal(err)
	//}
}
