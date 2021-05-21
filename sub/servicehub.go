package sub

import "github.com/ngobach/subdl/sub/services"

var Hub map[string]services.Service

func init() {
	Hub = map[string]services.Service{
		"subscene": services.NewSubSceneService(),
	}
}
