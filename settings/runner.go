package settings

type Runner struct {
	Image     string `json:"image"`
	TagSuffix string `json:"tag_suffix"`
	Shell     string `json:"shell"`
	Command   string `json:"command"`
}
