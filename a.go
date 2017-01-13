package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	flag.Parse()
	if *datasource == "" {
		*datasource = os.Getenv("DATASOURCE")
		if *datasource == "" {
			log.Fatalln("-datasource is required")
		}
	}
	p := NewGrafanaPanel()
	fill(p)
	m, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(m))
}

var (
	//tagName = flag.String("tagname", "*", "dev*")
	filters    = flag.String("filters", "", "e.g: tag:Name,dev-*,instance-type,m3.large")
	metricName = flag.String("metricName", "CPUUtilization", "CloudWatch MetricName")
	region     = flag.String("region", "ap-northeast-1", "AWS region")
	statistics = flag.String("statistics", "Average", "e.g: Average,Maximum,Minimum,Sum,SampleCount")
	datasource = flag.String("datasource", "", "data source name defined in the grafana")
)

func fill(p *GrafanaPanel) {
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(*region)})

	f := []*ec2.Filter{
		&ec2.Filter{
			Name:   aws.String("instance-state-name"),
			Values: []*string{aws.String("running")},
		},
		//&ec2.Filter{
		//	Name:   aws.String("tag:Name"),
		//	Values: []*string{aws.String(*tagName)},
		//},
	}
	if *filters != "" {
		fs := strings.Split(*filters, ",")
		if len(fs)%2 == 1 {
			log.Fatalln("illegal filters option: it should be even")
		}
		for i := 0; i*2 < len(fs); i++ {
			f = append(f, &ec2.Filter{
				Name:   aws.String(strings.TrimSpace(fs[i*2])),
				Values: []*string{aws.String(strings.TrimSpace(fs[i*2+1]))},
			})
		}
	}

	req := ec2.DescribeInstancesInput{Filters: f}
	res, err := svc.DescribeInstances(&req)
	if err != nil {
		panic(err)
		log.Fatalf("failed to DescribeInstances: %v", err)
	}

	ref := 0
	for _, res := range res.Reservations {
		//fmt.Println(len(res.Instances))
		for _, i := range res.Instances {
			alias := ""
			for _, t := range i.Tags {
				if *t.Key == "Name" {
					alias = *t.Value
				}
			}
			p.Targets = append(p.Targets, Target{
				Dimensions: map[string]string{"InstanceId": *i.InstanceId},
				MetricName: *metricName,
				Namespace:  "AWS/EC2",
				Period:     "",
				Region:     *region,
				Statistics: []string{
					*statistics,
				},
				RefID: fmt.Sprintf("A%d", ref),
				Alias: alias,
			})
			ref += 1
		}
	}

}

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

func NewGrafanaPanel() *GrafanaPanel {
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
		Title:      "EC2 CPU Utilization",
		Datasource: *datasource,
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
