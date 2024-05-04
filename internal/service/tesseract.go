package service

import log "github.com/sirupsen/logrus"

func (s *Service) ProcessFile(urlFile string) (string, error) {
	log.WithFields(log.Fields{
		"urlFile": urlFile,
	}).Info("start Process file")

	text, err := s.nlp.TextRecognition(urlFile)
	if err != nil {
		return "", err
	}

	log.WithFields(log.Fields{
		"text": text,
	}).Info("TextRecognized")

	return text, nil
}
