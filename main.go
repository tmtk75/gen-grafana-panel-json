package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jawher/mow.cli"
)

var (
	versionShort = "" // Expected given by -ldflags
	versionLong  = "" // Expected given by -ldflags
)

func main() {
	app := cli.App("gen-grafana-panel-json", "JSON Generator for Grafana CloudWatch datasource")
	app.Version("version", versionShort)
	app.Command("ec2", "EC2", func(c *cli.Cmd) {
		filters := c.String(cli.StringOpt{Name: "filters", Desc: `e.g: "tag:Name,dev-*", "instance-type,m3.large"`})
		opts := newCloudWatchOpts(c)
		c.Spec = "DATASOURCE_NAME [OPTIONS]"
		c.Action = func() {
			p := NewGrafanaPanel(*opts.dsName, "EC2 "+*opts.metricName)
			p.Targets = NewTargetsEC2(opts, *filters)
			PrintGrafanaPanelJSON(p)
		}
	})
	app.Command("sqs", "SQS", func(c *cli.Cmd) {
		opts := newCloudWatchOpts(c)
		px := c.String(cli.StringArg{Name: "PREFIX", Desc: "Prefix to filter"})
		rp := c.Bool(cli.BoolOpt{Name: "remove-prefix", Desc: "Prefix"})
		c.Spec = "[OPTIONS] DATASOURCE_NAME PREFIX"
		c.Action = func() {
			var qs []string
			if terminal.IsTerminal(int(os.Stdin.Fd())) {
				qs = ListQueues(*opts.region, *px)
			} else {
				bytes, err := ioutil.ReadAll(os.Stdin)
				if err != nil {
					log.Fatalf("failed to read stdin: %v", err)
				}
				qs = strings.Split(strings.Trim(string(bytes), "\n"), "\n")
			}
			p := NewGrafanaPanel(*opts.dsName, fmt.Sprintf("SQS %v-* %v", *px, *opts.metricName))
			p.Targets = NewTargetsSQS(opts, qs, *px, *rp)
			PrintGrafanaPanelJSON(p)
		}
	})
	app.Command("list-queues", "SQS", func(c *cli.Cmd) {
		var (
			p = c.String(cli.StringArg{Name: "PREFIX", Desc: "Prefix to filter"})
			r = c.String(cli.StringArg{Name: "REGION", Value: "ap-northeast-1", Desc: "ap-northeast-1 by default"})
		)
		c.Spec = "PREFIX [REGION]"
		c.Action = func() {
			for _, q := range ListQueues(*r, *p) {
				fmt.Printf("%v\n", q)
			}
		}
	})
	app.Command("cloudwatch", "CloudWatch", func(c *cli.Cmd) {
		c.Command("list-metrics", "", func(c *cli.Cmd) {
			var (
				r  = c.String(cli.StringArg{Name: "REGION", Value: "ap-northeast-1", Desc: "ap-northeast-1 by default"})
				ns = c.String(cli.StringArg{Name: "NAMESPACE", Desc: "CloudWatch namespace e.g) AWS/EC2"})
			)
			c.Spec = "NAMESPACE [REGION]"
			c.Action = func() {
				ListMetrics(*ns, *r, nil)
			}
		})
	})
	app.Run(os.Args)
}

func ListQueues(region, prefix string) []string {
	svc := sqs.New(session.New(), &aws.Config{Region: aws.String(region)})
	req := sqs.ListQueuesInput{
		QueueNamePrefix: aws.String(prefix),
	}
	res, err := svc.ListQueues(&req)
	if err != nil {
		log.Fatalf("%v", err)
	}
	qs := make([]string, 0)
	for _, q := range res.QueueUrls {
		qs = append(qs, *q)
	}
	return qs
}

func queueName(q string) string {
	s := strings.Split(q, "/")
	return s[len(s)-1] // queue name
}

type cloudWatchOpts struct {
	dsName     *string
	metricName *string
	region     *string
	statistics *string
}

func ListMetrics(ns, region string, nextToken *string) {
	svc := cloudwatch.New(session.New(), &aws.Config{Region: aws.String(region)})
	req := cloudwatch.ListMetricsInput{
		Namespace:  aws.String(ns),
		Dimensions: []*cloudwatch.DimensionFilter{
		//&cloudwatch.DimensionFilter{
		//	Name:  aws.String("QueueName"),
		//	Value: aws.String("stg-jp"),
		//},
		},
		NextToken: nextToken,
	}
	res, err := svc.ListMetrics(&req)
	if err != nil {
		log.Fatalf("%v", err)
	}
	for _, m := range res.Metrics {
		fmt.Printf("%v %v\n", *m.Dimensions[0].Value, *m.MetricName)
	}
	if res.NextToken != nil && *res.NextToken != "" {
		ListMetrics(ns, region, res.NextToken)
	}
}

func newCloudWatchOpts(c *cli.Cmd) *cloudWatchOpts {
	return &cloudWatchOpts{
		dsName:     c.String(cli.StringArg{Name: "DATASOURCE_NAME", Desc: "Grafana datasource name"}),
		metricName: c.String(cli.StringOpt{Name: "metricName m", Value: "CPUUtilization", Desc: "CloudWatch MetricName"}),
		region:     c.String(cli.StringOpt{Name: "region r", Value: "ap-northeast-1", Desc: "AWS region"}),
		statistics: c.String(cli.StringOpt{Name: "statistics s", Value: "Average", Desc: "e.g: Average,Maximum,Minimum,Sum,SampleCount"}),
	}
}

func PrintGrafanaPanelJSON(p *GrafanaPanel) {
	m, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("failed to marshal: %v", err)
	}
	fmt.Println(string(m))
}

func alias(s string) string {
	a := regexp.MustCompile("Of").ReplaceAllString(s, "")
	re := regexp.MustCompile("[a-z]+")
	return re.ReplaceAllString(a, "")
}

func removePrefix(prefix, s string) string {
	return strings.TrimLeft(s, prefix)
}

func NewTargetsSQS(opts *cloudWatchOpts, urls []string, prefix string, rp bool) []Target {
	targets := make([]Target, 0)
	for i, q := range urls {
		qn := queueName(q)
		if rp {
			qn = removePrefix(prefix, qn)
		}
		t := Target{
			Dimensions: map[string]string{
				"QueueName": queueName(q),
			},
			MetricName: *opts.metricName,
			Namespace:  "AWS/SQS",
			Period:     "300",
			RefID:      fmt.Sprintf("ID-%v", i),
			Region:     *opts.region,
			Statistics: []string{
				*opts.statistics,
			},
			Alias: alias(*opts.metricName) + ":" + qn,
		}
		targets = append(targets, t)
	}
	return targets //[0:20]
}

func NewTargetsEC2(opts *cloudWatchOpts, filters string) []Target {
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
	if filters != "" {
		fs := strings.Split(filters, ",")
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

	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(*opts.region)})
	req := ec2.DescribeInstancesInput{Filters: f}
	res, err := svc.DescribeInstances(&req)
	if err != nil {
		log.Fatalf("failed to describe-instances: %v", err)
		log.Fatalf("failed to DescribeInstances: %v", err)
	}

	ref := 0
	targets := make([]Target, 0)
	for _, res := range res.Reservations {
		//fmt.Println(len(res.Instances))
		for _, i := range res.Instances {
			alias := ""
			for _, t := range i.Tags {
				if *t.Key == "Name" {
					alias = *t.Value
				}
			}
			targets = append(targets, Target{
				Dimensions: map[string]string{"InstanceId": *i.InstanceId},
				MetricName: *opts.metricName,
				Namespace:  "AWS/EC2",
				Period:     "",
				Region:     *opts.region,
				Statistics: []string{
					*opts.statistics,
				},
				RefID: fmt.Sprintf("A%d", ref),
				Alias: alias,
			})
			ref += 1
		}
	}

	sort.Sort(Targets(targets))
	return targets
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
