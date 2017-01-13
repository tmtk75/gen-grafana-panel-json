package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	p := GrafanaPanel{
		Type:  "graph",
		Xaxis: Xaxis{Mode: "time", Show: true},
		Yaxes: []Yaxis{
			{Format: "short", Show: true, LogBase: 1},
			{Format: "short", Show: true, LogBase: 1},
		},
		Tooltip: Tooltip{
			Shared:    true,
			Sort:      0,
			ValueType: "individual",
		},
	}
	m, err := json.Marshal(&p)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(m))
}

func list() {
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String("ap-northeast-1")})
	res, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}
	for _, res := range res.Reservations {
		//fmt.Println(len(res.Instances))
		for _, i := range res.Instances {
			fmt.Println(i)
		}
	}

}

type GrafanaPanel struct {
	AliasColors struct{} `json:"aliasColors"`
	Bars        bool     `json:"bars"`
	Datasource  string   `json:"datasource"`
	Fill        int      `json:"fill"`
	ID          int      `json:"id"`
	Legend      struct {
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
	} `json:"legend"`
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
	Targets         []struct {
		Alias      string `json:"alias"`
		Dimensions struct {
			InstanceID string `json:"InstanceId"`
		} `json:"dimensions"`
		MetricName string   `json:"metricName"`
		Namespace  string   `json:"namespace"`
		Period     string   `json:"period"`
		RefID      string   `json:"refId"`
		Region     string   `json:"region"`
		Statistics []string `json:"statistics"`
	} `json:"targets"`
	Thresholds  []interface{} `json:"thresholds"`
	TimeFrom    interface{}   `json:"timeFrom"`
	TimeShift   interface{}   `json:"timeShift"`
	Title       string        `json:"title"`
	Tooltip     Tooltip       `json:"tooltip"`
	Transparent bool          `json:"transparent"`
	Type        string        `json:"type"`
	Xaxis       Xaxis         `json:"xaxis"`
	Yaxes       []Yaxis       `json:"yaxes"`
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
