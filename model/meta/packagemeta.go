package meta

const (
	PACKAGE_NORMAL = "Normal"
	PACKAGE_THIRD  = "ThirdParty"
)

type PackageMeta struct {
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`
	Label       string `json:"label"`
	System      bool   `json:"system"`
	Sharable    bool   `json:"sharable"`
	StereoType  string `json:"stereoType"`
	TokenScript string `json:"tokenScript"`
	AppId       uint64
}
