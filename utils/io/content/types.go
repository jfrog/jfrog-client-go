package content

type SortableContentItem interface {
	GetSortKey() string
}
