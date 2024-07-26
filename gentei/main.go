package main

import (
	"github.com/member-gentei/member-gentei/gentei/cmd"
	"github.com/rs/zerolog/log"
)

func main() {
	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			if pErr, ok := r.(error); ok {
				log.Err(pErr).Msg("panic recovery")
			} else {
				log.Error().Any("recover", r).Msg("panic recovery")
			}
		}
	}()
	cmd.Execute()
}
