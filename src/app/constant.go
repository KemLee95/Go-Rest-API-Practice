package app

type MESSAGE_LANGUAGE string

const (
	ENGLISH               MESSAGE_LANGUAGE = "en"
	CHINA                 MESSAGE_LANGUAGE = "ch"
	LANGUAGE_HEADER_PARAM                  = "x-language-code"
)
