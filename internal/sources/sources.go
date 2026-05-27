package sources

import "signalboard/internal/content"

type Source interface {
	Content() []content.Content
}
