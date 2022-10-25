package business

type PluginInfo struct {
	Id          uint64 `json:"id"`
	Title       string `json:"title"`
	Url         string `json:"url"`
	PluginId    string `json:"pluginId"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Version     string `json:"version"`
}
