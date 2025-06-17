package config

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"

	"sigs.k8s.io/yaml"

	"github.com/authgear/authgear-server/pkg/util/validation"
)

type ProviderType string

const (
	ProviderTypeTwilio       ProviderType = "twilio"
	ProviderTypeAccessYou    ProviderType = "accessyou"
	ProviderTypeSendCloud    ProviderType = "sendcloud"
	ProviderTypeAccessYouOTP ProviderType = "accessyou_otp"
)

var _ = RootSchema.Add("ProviderType", `
{
	"type": "string",
	"enum": ["twilio", "accessyou", "sendcloud", "accessyou_otp"]
}
`)

type Provider struct {
	Name         string                      `json:"name,omitempty"`
	Type         ProviderType                `json:"type,omitempty"`
	Twilio       *ProviderConfigTwilio       `json:"twilio,omitempty" nullable:"true"`
	AccessYou    *ProviderConfigAccessYou    `json:"accessyou,omitempty" nullable:"true"`
	SendCloud    *ProviderConfigSendCloud    `json:"sendcloud,omitempty" nullable:"true"`
	AccessYouOTP *ProviderConfigAccessYouOTP `json:"accessyou_otp,omitempty" nullable:"true"`
}

type ProviderConfigTwilio struct {
	AccountSID string `json:"account_sid,omitempty"`

	// From and MessagingServiceSID are mutually exclusive.
	From                string `json:"from,omitempty"`
	MessagingServiceSID string `json:"messaging_service_sid,omitempty"`

	// AuthToken and (APIKey and APIKeySecret) are mutually exclusive.
	AuthToken    string `json:"auth_token,omitempty"`
	APIKey       string `json:"api_key,omitempty"`
	APIKeySecret string `json:"api_key_secret,omitempty"`
}

type ProviderConfigAccessYou struct {
	From      string `json:"from,omitempty"`
	BaseUrl   string `json:"base_url,omitempty"`
	AccountNo string `json:"accountno,omitempty"`
	User      string `json:"user,omitempty"`
	Pwd       string `json:"pwd,omitempty"`
}

type ProviderConfigAccessYouOTP struct {
	BaseUrl   string `json:"base_url,omitempty"`
	AccountNo string `json:"accountno,omitempty"`
	User      string `json:"user,omitempty"`
	Pwd       string `json:"pwd,omitempty"`
	A         string `json:"a,omitempty"`
}

var _ = RootSchema.Add("Provider", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"name": { "type": "string" },
		"type": { "$ref": "#/$defs/ProviderType" },
		"twilio": { "$ref": "#/$defs/ProviderConfigTwilio" },
		"accessyou": { "$ref": "#/$defs/ProviderConfigAccessYou" },
		"sendcloud": { "$ref": "#/$defs/ProviderConfigSendCloud" },
		"accessyou_otp": { "$ref": "#/$defs/ProviderConfigAccessYouOTP" }
	},
	"allOf": [
		{
			"if": { "properties": { "type": { "const": "twilio" } }},
			"then": { "required": ["twilio"] }
		},
		{
			"if": { "properties": { "type": { "const": "accessyou" } }},
			"then": { "required": ["accessyou"] }
		},
		{
			"if": { "properties": { "type": { "const": "sendcloud" } }},
			"then": { "required": ["sendcloud"] }
		},
		{
			"if": { "properties": { "type": { "const": "accessyou_otp" } }},
			"then": { "required": ["accessyou_otp"] }
		}
	]
}
`)

var _ = RootSchema.Add("ProviderConfigTwilio", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"account_sid": { "type": "string" },
		"from": { "type": "string" },
		"messaging_service_sid": { "type": "string" },
		"auth_token": { "type": "string" },
		"api_key": { "type": "string" },
		"api_key_secret": { "type": "string" }
	},
	"required": ["account_sid"],
	"allOf": [
		{
			"oneOf": [
				{
					"required": ["from"]
				},
				{
					"required": ["messaging_service_sid"]
				}
			]
		},
		{
			"oneOf": [
				{
					"required": ["auth_token"]
				},
				{
					"required": ["api_key", "api_key_secret"]
				}
			]
		}
	]
}
`)

var _ = RootSchema.Add("ProviderConfigAccessYou", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"from": { "type": "string" },
		"base_url": { "type": "string" },
		"accountno": { "type": "string" },
		"user": { "type": "string" },
		"pwd": {"type": "string"}
	},
	"required": ["from", "accountno", "user", "pwd"]
}
`)

var _ = RootSchema.Add("ProviderConfigAccessYouOTP", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"base_url": { "type": "string" },
		"accountno": { "type": "string" },
		"user": { "type": "string" },
		"pwd": {"type": "string"},
		"a": {"type": "string"}
	},
	"required": ["accountno", "user", "pwd", "a"]
}
`)

type ProviderSelectorSwitchType string

const (
	ProviderSelectorSwitchTypeMatchPhoneNumberAlpha2         ProviderSelectorSwitchType = "match_phone_number_alpha2"
	ProviderSelectorSwitchTypeMatchAppID                     ProviderSelectorSwitchType = "match_app_id"
	ProviderSelectorSwitchTypeMatchAppIDAndPhoneNumberAlpha2 ProviderSelectorSwitchType = "match_app_id_and_phone_number_alpha2"
	ProviderSelectorSwitchTypeDefault                        ProviderSelectorSwitchType = "default"
)

var _ = RootSchema.Add("ProviderSelectorSwitchType", `
{
	"type": "string",
	"enum": ["match_phone_number_alpha2", "match_app_id", "match_app_id_and_phone_number_alpha2", "default"]
}
`)

type ProviderSelectorSwitchRule struct {
	Type              ProviderSelectorSwitchType `json:"type,omitempty"`
	UseProvider       string                     `json:"use_provider,omitempty"`
	PhoneNumberAlpha2 string                     `json:"phone_number_alpha2,omitempty"`
	AppID             string                     `json:"app_id,omitempty"`
}

var _ = RootSchema.Add("ProviderSelectorSwitchRule", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"type": { "$ref": "#/$defs/ProviderSelectorSwitchType" },
		"use_provider": { "type": "string" },
		"phone_number_alpha2": { "type": "string" },
		"app_id": { "type": "string" }
	},
	"allOf": [
		{
			"if": { "properties": { "type": { "const": "match_phone_number_alpha2" } }},
			"then": { "required": ["phone_number_alpha2"] }
		},
		{
			"if": { "properties": { "type": { "const": "match_app_id" } }},
			"then": { "required": ["app_id"] }
		},
		{
			"if": { "properties": { "type": { "const": "match_app_id_and_phone_number_alpha2" } }},
			"then": { "required": ["phone_number_alpha2", "app_id"] }
		}
	],
	"required": ["type", "use_provider"]
}
`)

type ProviderSelector struct {
	Switch []*ProviderSelectorSwitchRule `json:"switch,omitempty"`
}

var _ = RootSchema.Add("ProviderSelector", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"switch": {
			"type": "array",
			"minItems": 1,
			"items": { "$ref": "#/$defs/ProviderSelectorSwitchRule" }
		}
	},
	"required": ["switch"]
}
`)

type RootConfig struct {
	Providers        []*Provider       `json:"providers,omitempty"`
	ProviderSelector *ProviderSelector `json:"provider_selector,omitempty"`
}

var _ validation.Validator = (*RootConfig)(nil)

var _ = RootSchema.Add("RootConfig", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"providers": {
			"type": "array",
			"minItems": 1,
			"items": { "$ref": "#/$defs/Provider" }
		},
		"provider_selector": { "$ref": "#/$defs/ProviderSelector" }
	},
	"required": ["providers", "provider_selector"]
}
`)

func (c *RootConfig) Validate(ctx context.Context, validationCtx *validation.Context) {
	c.ValidateProviderSelectorUseProvider(validationCtx)
	c.ValidateProviderSelectorDefault(validationCtx)
	c.ValidateSendCloudConfigs(validationCtx)
}

func (c *RootConfig) ValidateProviderSelectorUseProvider(ctx *validation.Context) {
	providers := c.Providers
	for i, switchCase := range c.ProviderSelector.Switch {
		useProvider := switchCase.UseProvider
		idx := slices.IndexFunc(providers, func(p *Provider) bool { return p.Name == useProvider })
		if idx == -1 {
			ctx.Child("provider_selector", "switch", strconv.Itoa(i), "use_provider").EmitErrorMessage(fmt.Sprintf("provider %s not found", useProvider))
		}
	}
}

func (c *RootConfig) ValidateProviderSelectorDefault(ctx *validation.Context) {
	for _, switchCase := range c.ProviderSelector.Switch {
		if switchCase.Type == ProviderSelectorSwitchTypeDefault {
			return
		}
	}
	ctx.Child("provider_selector", "switch").EmitErrorMessage("provider selector default not found")
}

func (c *RootConfig) ValidateSendCloudConfigs(ctx *validation.Context) {
	for i, provider := range c.Providers {
		if provider.Type == ProviderTypeSendCloud {
			c.ValidateSendCloudConfig(ctx.Child("providers", strconv.Itoa(i), "sendcloud"), provider.SendCloud)
		}
	}
}

func (c *RootConfig) ValidateSendCloudConfig(ctx *validation.Context, sendCloudConfig *ProviderConfigSendCloud) {
	templates := sendCloudConfig.Templates
	for i, template := range templates {
		ctxTemplates := ctx.Child("templates", strconv.Itoa(i))
		if len(template.TemplateVariableKeyMappings) == 0 {
			ctxTemplates.Child("template_variable_key_mappings").EmitErrorMessage("missing template_variable_key_mappings")
		}
	}

	for i, templateAssignment := range sendCloudConfig.TemplateAssignments {
		ctxTemplateAssignment := ctx.Child("template_assignments", strconv.Itoa(i))
		defaultTemplateID := templateAssignment.DefaultTemplateID
		idx := slices.IndexFunc(templates, func(t *SendCloudTemplate) bool { return t.TemplateID == defaultTemplateID })

		if idx == -1 {
			ctxTemplateAssignment.Child("default_template_id").EmitErrorMessage(fmt.Sprintf("template_id %v not found", defaultTemplateID))
		}

		for j, byLanguage := range templateAssignment.ByLanguages {
			ctxByLanguage := ctxTemplateAssignment.Child("by_languages", strconv.Itoa(j))
			templateID := byLanguage.TemplateID
			idx = slices.IndexFunc(templates, func(t *SendCloudTemplate) bool { return t.TemplateID == templateID })
			if idx == -1 {
				ctxByLanguage.Child("template_id").EmitErrorMessage(fmt.Sprintf("template_id %v not found", templateID))
			}
		}
	}

}

func ParseRootConfigFromYAML(ctx context.Context, inputYAML []byte) (*RootConfig, error) {
	const validationErrorMessage = "invalid configuration"

	jsonData, err := yaml.YAMLToJSON(inputYAML)
	if err != nil {
		return nil, err
	}

	err = RootSchema.Validator().ValidateWithMessage(ctx, bytes.NewReader(jsonData), validationErrorMessage)
	if err != nil {
		return nil, err
	}

	var config RootConfig
	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	err = validation.ValidateValueWithMessage(ctx, &config, validationErrorMessage)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
