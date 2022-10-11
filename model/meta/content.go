package meta

type MetaContent struct {
	Classes   []ClassMeta    `json:"entities"`
	Relations []RelationMeta `json:"relations"`
	Codes     []CodeMeta     `json:"codes"`
	Packages  []PackageMeta  `json:"packages"`
	Diagrams  []interface{}  `json:"diagrams"`
	X6Nodes   []interface{}  `json:"x6Nodes"`
	X6Edges   []interface{}  `json:"x6Edges"`
}
