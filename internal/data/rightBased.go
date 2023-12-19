package data

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

type RightsBased struct {
	Id     string
	Login  string
	Name   string
	Rights map[string]bool
	Token  string
}

func (b RightsBased) CheckEither(rights []string) error {
	has := false

	for _, right := range rights {
		_, ok := b.Rights[right]

		if ok {
			has = true
			break
		}
	}

	var err error

	if !has {
		err = fmt.Errorf("does not have at least one of the rigths: %s", rights)
	}

	if err == nil {
		log.Err(err).Msg("error check either")
	} else {
		log.Info().Msg("success check either")
	}

	return err
}

func (b RightsBased) CheckAll(rights []string) error {
	has := true

	for _, right := range rights {
		_, ok := b.Rights[right]

		if !ok {
			has = false
			break
		}
	}

	var err error

	if !has {
		err = fmt.Errorf("must have all the rights: %s", rights)
	}

	if err == nil {
		log.Err(err).Msg("error check either")
	} else {
		log.Info().Msg("success check either")
	}

	return err
}
