package meta

type PackageMeta struct {
	Uuid     string `json:"uuid"`
	Name     string `json:"name"`
	Label    string `json:"label"`
	System   bool   `json:"system"`
	Sharable bool   `json:"sharable"`
	AppId    uint64
}
