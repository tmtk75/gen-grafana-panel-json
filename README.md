# README
Generate Grafana panel JSON for CloudWatch datasource.

Grafana has [HTTP API](http://docs.grafana.org/reference/http_api) but at a glance it doens't have
APIs for panel operations. We con configure panels with JSON. This command generates
JSONs for panels which datasource is CloudWatch.

# How to use
You need to configure `~/.aws/credentials` first.

Generate a Grafana panel JSON for all EC2 instances in ap-northeast-1.
```
$ gen-grafana-panel-json -datasource CloudWatch -region ap-northeast-1
```

Switching profile.
```
$ AWS_PROFILE=staging gen-grafana-panel-json -datasource CloudWatch -region us-west-2
```

To filter EC2, use `-filters`.
```
$ gen-grafana-panel-json -datasource CloudWatch -filters tag:Name,dev-*,instance-type,t2.small \
  | jq . | head
{
  "aliasColors": {},
  "bars": false,
  "datasource": "CloudWatch",
  "fill": 1,
  "id": 2,
  "legend": {
    "alignAsTable": false,
    "avg": false,
    "current": false,
```

To exclude some elements in targets, use `jq` and `-stdin` option.
```
$ gen-grafana-panel-json -datasource CloudWatch \
  | jq '[.targets[]|select((.alias|test("dev-.*"))|not)]' \
  | gen-grafana-panel-json -stdin
```
