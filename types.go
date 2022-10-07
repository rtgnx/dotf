package godots

type Dots []Dot

type Variables struct {
	Global map[string]string
}

type Dot struct {
	Name      string            `json:"name"`
	FileMap   map[string]string `json:"fmap"`
	Templates []string          `json:"templates"`
}
