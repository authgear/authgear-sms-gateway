package config

type ProviderConfigCMCOM struct {
	From         string `json:"from,omitempty"`
	ProductToken string `json:"product_token,omitempty"`
}

var _ = RootSchema.Add("ProviderConfigCMCOM", `{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"from": { "type": "string" },
		"product_token": { "type": "string" }
	},
	"required": ["from", "product_token"]
}`)
