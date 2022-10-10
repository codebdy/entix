package meta

const (
	CLASSS_ENTITY      string = "Entity"
	CLASSS_ENUM        string = "Enum"
	CLASSS_ABSTRACT    string = "Abstract"
	CLASS_VALUE_OBJECT string = "ValueObject"
	CLASS_SERVICE      string = "Service"
)

type ClassMeta struct {
	Uuid        string          `json:"uuid"`
	InnerId     uint64          `json:"innerId"`
	Name        string          `json:"name"`
	Label       string          `json:"label"`
	StereoType  string          `json:"stereoType"`
	Attributes  []AttributeMeta `json:"attributes"`
	Methods     []MethodMeta    `json:"methods"`
	Root        bool            `json:"root"`
	Description string          `json:"description"`
	SoftDelete  bool            `json:"softDelete"`
	System      bool            `json:"system"`
	IdNoShift   bool            `json:"idNoShift"`
	PackageUuid string          `json:"packageUuid"`
	AppId       uint64
}
