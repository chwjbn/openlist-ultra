package jellyfin

type JellyfinDataItem struct {
	Name         string
	ServerId     string
	Id           string
	ChannelId    string
	IsFolder     bool
	Type         string
	LocationType string
	MediaType    string
}

type JellyfinRespItems struct {
	Items            []JellyfinDataItem
	TotalRecordCount int
	StartIndex       int
}
