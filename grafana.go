package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type GrafanaPanel struct {
	AliasColors     struct{}      `json:"aliasColors"`
	Bars            bool          `json:"bars"`
	Datasource      string        `json:"datasource"`
	Fill            int           `json:"fill"`
	ID              int           `json:"id"`
	Legend          Legend        `json:"legend"`
	Lines           bool          `json:"lines"`
	Linewidth       int           `json:"linewidth"`
	Links           []interface{} `json:"links"`
	NullPointMode   string        `json:"nullPointMode"`
	Percentage      bool          `json:"percentage"`
	Pointradius     int           `json:"pointradius"`
	Points          bool          `json:"points"`
	Renderer        string        `json:"renderer"`
	SeriesOverrides []interface{} `json:"seriesOverrides"`
	Span            int           `json:"span"`
	Stack           bool          `json:"stack"`
	SteppedLine     bool          `json:"steppedLine"`
	Targets         []Target      `json:"targets"`
	Thresholds      []interface{} `json:"thresholds"`
	TimeFrom        interface{}   `json:"timeFrom"`
	TimeShift       interface{}   `json:"timeShift"`
	Title           string        `json:"title"`
	Tooltip         Tooltip       `json:"tooltip"`
	Transparent     bool          `json:"transparent"`
	Type            string        `json:"type"`
	Xaxis           Xaxis         `json:"xaxis"`
	Yaxes           []Yaxis       `json:"yaxes"`
}

type Xaxis struct {
	Mode   string        `json:"mode"`
	Name   interface{}   `json:"name"`
	Show   bool          `json:"show"`
	Values []interface{} `json:"values"`
}

type Yaxis struct {
	Format  string      `json:"format"`
	Label   interface{} `json:"label"`
	LogBase int         `json:"logBase"`
	Max     interface{} `json:"max"`
	Min     interface{} `json:"min"`
	Show    bool        `json:"show"`
}

type Tooltip struct {
	Shared    bool   `json:"shared"`
	Sort      int    `json:"sort"`
	ValueType string `json:"value_type"`
}

type Legend struct {
	AlignAsTable bool `json:"alignAsTable"`
	Avg          bool `json:"avg"`
	Current      bool `json:"current"`
	HideEmpty    bool `json:"hideEmpty"`
	HideZero     bool `json:"hideZero"`
	Max          bool `json:"max"`
	Min          bool `json:"min"`
	RightSide    bool `json:"rightSide"`
	Show         bool `json:"show"`
	Total        bool `json:"total"`
	Values       bool `json:"values"`
}

type Target struct {
	Alias      string            `json:"alias"`
	Dimensions map[string]string `json:"dimensions"`
	MetricName string            `json:"metricName"`
	Namespace  string            `json:"namespace"`
	Period     string            `json:"period"`
	RefID      string            `json:"refId"`
	Region     string            `json:"region"`
	Statistics []string          `json:"statistics"`
}

func NewGrafanaPanel(ds, title string) *GrafanaPanel {
	return &GrafanaPanel{
		Type:            "graph",
		Links:           []interface{}{},
		NullPointMode:   "null",
		SeriesOverrides: []interface{}{},
		Thresholds:      []interface{}{},
		Xaxis:           Xaxis{Mode: "time", Show: true, Values: []interface{}{}},
		Yaxes: []Yaxis{
			{Format: "short", Show: true, LogBase: 1},
			{Format: "short", Show: true, LogBase: 1},
		},
		Tooltip: Tooltip{
			Shared:    true,
			Sort:      0,
			ValueType: "individual",
		},
		Title:      title,
		Datasource: ds,
		Fill:       1,
		ID:         2,
		Legend: Legend{
			Show: true,
		},
		Lines:       true,
		Linewidth:   1,
		Pointradius: 5,
		Renderer:    "flot",
		Span:        8,
		Targets:     []Target{},
	}
}

type Targets []Target

func (f Targets) Len() int {
	return len(f)
}

func (f Targets) Less(i, j int) bool {
	return strings.Compare(f[i].Alias, f[j].Alias) < 0
}

func (f Targets) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func PrintGrafanaPanelJSON(p *GrafanaPanel) {
	m, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("failed to marshal: %v", err)
	}
	fmt.Println(string(m))
}
