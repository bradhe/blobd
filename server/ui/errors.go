package ui

import "errors"

var (
	ErrAssetNotFound            = errors.New("ui: asset not found")
	ErrTemplateParsingFailed    = errors.New("ui: failed to parse asset template")
	ErrTemplateProcessingFailed = errors.New("ui: failed to process asset template")
)
