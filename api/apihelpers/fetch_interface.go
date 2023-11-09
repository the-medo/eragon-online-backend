package apihelpers

type FetchInterface struct {
	ModuleIds    []int32 `json:"module_ids,omitempty"`
	WorldIds     []int32 `json:"world_ids,omitempty"`
	SystemIds    []int32 `json:"system_ids,omitempty"`
	CharacterIds []int32 `json:"character_ids,omitempty"`
	QuestIds     []int32 `json:"quest_ids,omitempty"`

	EntityIds   []int32 `json:"entity_ids,omitempty"`
	PostIds     []int32 `json:"post_ids,omitempty"`
	ImageIds    []int32 `json:"image_ids,omitempty"`
	LocationIds []int32 `json:"location_ids,omitempty"`
	MapIds      []int32 `json:"map_ids,omitempty"`
}
