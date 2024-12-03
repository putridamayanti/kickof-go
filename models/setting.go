package models

const ConfigCollection = "configs"
const SettingCollection = "settings"

type Configuration struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Values      []string `json:"values"`
	Type        string   `json:"type"`
}

type Setting struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Currency Currency `json:"currency"`
	Address  string   `json:"address"`
	Tax      string   `json:"tax"`
}

type Currency struct {
	Code     string `json:"code"`
	Currency string `json:"currency"`
	Symbol   string `json:"symbol"`
}

type Tax struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	ServicePercentage string `json:"servicePercentage"`
	TaxPercentage     string `json:"taxPercentage"`
	UseService        bool   `json:"useService"`
}
