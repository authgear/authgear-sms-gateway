package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"

	"sigs.k8s.io/yaml"

	"github.com/authgear/authgear-server/pkg/util/validation"
)

type ProviderType string

const (
	ProviderTypeTwilio    ProviderType = "twilio"
	ProviderTypeNexmo     ProviderType = "nexmo"
	ProviderTypeAccessYou ProviderType = "accessyou"
	ProviderTypeSendCloud ProviderType = "sendcloud"
	ProviderTypeInfobip   ProviderType = "infobip"
)

var _ = SMSProviderConfigSchema.Add("ProviderType", `
{
	"type": "string",
	"enum": ["twilio", "nexmo", "accessyou", "sendcloud", "infobip"]
}
`)

type Provider struct {
	Name      string                   `json:"name,omitempty"`
	Type      ProviderType             `json:"type,omitempty"`
	Twilio    *ProviderConfigTwilio    `json:"twilio,omitempty" nullable:"true"`
	Nexmo     *ProviderConfigNexmo     `json:"nexmo,omitempty" nullable:"true"`
	AccessYou *ProviderConfigAccessYou `json:"accessyou,omitempty" nullable:"true"`
	SendCloud *ProviderConfigSendCloud `json:"sendcloud,omitempty" nullable:"true"`
	Infobip   *ProviderConfigInfobip   `json:"infobip,omitempty" nullable:"true"`
}

type ProviderConfigTwilio struct {
	AccountSID          string `json:"account_sid,omitempty"`
	AuthToken           string `json:"auth_token,omitempty"`
	MessagingServiceSID string `json:"message_service_sid,omitempty"`
}

type ProviderConfigNexmo struct {
	APIKey    string `json:"api_key,omitempty"`
	APISecret string `json:"api_secret,omitempty"`
}

type ProviderConfigAccessYou struct {
	AccountNo string `json:"accountno,omitempty"`
	Pwd       string `json:"pwd,omitempty"`
}

type ProviderConfigSendCloud struct {
	SMSUser string `json:"sms_user,omitempty"`
	SMSKey  string `json:"sms_key,omitempty"`
}

type ProviderConfigInfobip struct {
	APIKey string `json:"api_key,omitempty"`
}

var _ = SMSProviderConfigSchema.Add("Provider", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"name": { "type": "string" },
		"type": { "$ref": "#/$defs/ProviderType" },
		"twilio": { "$ref": "#/$defs/ProviderConfigTwilio" },
		"nexmo": { "$ref": "#/$defs/ProviderConfigNexmo" },
		"accessyou": { "$ref": "#/$defs/ProviderConfigAccessYou" },
		"sendcloud": { "$ref": "#/$defs/ProviderConfigSendCloud" },
		"infobip": { "$ref": "#/$defs/ProviderConfigInfobip" }
	},
	"allOf": [
		{
			"if": { "properties": { "type": { "const": "twilio" } }},
			"then": { "required": ["twilio"] }
		},
		{
			"if": { "properties": { "type": { "const": "nexmo" } }},
			"then": { "required": ["nexmo"] }
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
			"if": { "properties": { "type": { "const": "infobip" } }},
			"then": { "required": ["infobip"] }
		}
	]
}
`)

var _ = SMSProviderConfigSchema.Add("ProviderConfigTwilio", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"account_sid": { "type": "string" },
		"auth_token": {"type": "string"},
		"message_service_sid": {"type": "string"}
	},
	"required": ["account_sid", "auth_token", "message_service_sid"]
}
`)

var _ = SMSProviderConfigSchema.Add("ProviderConfigNexmo", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"api_key": { "type": "string" },
		"api_secret": {"type": "string"}
	},
	"required": ["api_key", "api_secret"]
}
`)

var _ = SMSProviderConfigSchema.Add("ProviderConfigAccessYou", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"accountno": { "type": "string" },
		"pwd": {"type": "string"}
	},
	"required": ["accountno", "pwd"]
}
`)

var _ = SMSProviderConfigSchema.Add("ProviderConfigSendCloud", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"sms_user": { "type": "string" },
		"sms_key": {"type": "string"}
	},
	"required": ["sms_user", "sms_key"]
}
`)

var _ = SMSProviderConfigSchema.Add("ProviderConfigInfobip", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"api_key": { "type": "string" }
	},
	"required": ["api_key"]
}
`)

type ProviderSelectorSwitchType string

const (
	ProviderSelectorSwitchTypeMatchPhoneNumberAlpha2 ProviderSelectorSwitchType = "match_phone_number_alpha2"
	ProviderSelectorSwitchTypeDefault                ProviderSelectorSwitchType = "default"
)

var _ = SMSProviderConfigSchema.Add("ProviderSelectorSwitchType", `
{
	"type": "string",
	"enum": ["match_phone_number_alpha2", "default"]
}
`)

type ProviderSelectorSwitchRule struct {
	Type              ProviderSelectorSwitchType `json:"type,omitempty"`
	UseProvider       string                     `json:"use_provider,omitempty"`
	PhoneNumberAlpha2 string                     `json:"phone_number_alpha2,omitempty"`
}

var _ = SMSProviderConfigSchema.Add("ProviderSelectorSwitchRule", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"type": { "$ref": "#/$defs/ProviderSelectorSwitchType" },
		"use_provider": { "type": "string" },
		"phone_number_alpha2": { "type": "string" }
	},
	"allOf": [
		{
			"if": { "properties": { "type": { "const": "phone_number_alpha2" } }},
			"then": { "required": ["phone_number_alpha2"] }
		}
	],
	"required": ["type", "use_provider"]
}
`)

type ProviderSelector struct {
	Switch []*ProviderSelectorSwitchRule `json:"switch,omitempty"`
}

var _ = SMSProviderConfigSchema.Add("ProviderSelector", `
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

type SMSProviderConfig struct {
	Providers        []*Provider       `json:"providers,omitempty"`
	ProviderSelector *ProviderSelector `json:"provider_selector,omitempty"`
}

var _ = SMSProviderConfigSchema.Add("SMSProviderConfig", `
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

func (c *SMSProviderConfig) Validate(ctx *validation.Context) {
	c.ValidateProvider(ctx)
}

func (c *SMSProviderConfig) ValidateProvider(ctx *validation.Context) {
	providers := c.Providers
	for i, switchCase := range c.ProviderSelector.Switch {
		useProvider := switchCase.UseProvider
		idx := slices.IndexFunc(providers, func(p *Provider) bool { return p.Name == useProvider })
		if idx == -1 {
			ctx.Child("provider_selector", "switch", strconv.Itoa(i), "use_provider").EmitErrorMessage(fmt.Sprintf("provider %s not found", useProvider))
		}
	}
}

func ParseSMSProviderConfigFromYAML(inputYAML []byte) (*SMSProviderConfig, error) {
	const validationErrorMessage = "invalid configuration"

	jsonData, err := yaml.YAMLToJSON(inputYAML)
	if err != nil {
		return nil, err
	}

	err = SMSProviderConfigSchema.Validator().ValidateWithMessage(bytes.NewReader(jsonData), validationErrorMessage)
	if err != nil {
		return nil, err
	}

	var config SMSProviderConfig
	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	err = validation.ValidateValueWithMessage(&config, validationErrorMessage)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
