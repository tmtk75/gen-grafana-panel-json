# README
Generate Grafana panel JSON for CloudWatch datasource.

Grafana has [HTTP API](http://docs.grafana.org/reference/http_api) but at a glance it doens't have
APIs for panel operations. We con configure panels with JSON. This command generates
JSONs for panels which datasource is CloudWatch.

# Getting Started
You need to configure `~/.aws/credentials` first.

Generate a Grafana panel JSON for all EC2 instances in ap-northeast-1 for CPUUtilization.
```
$ gen-grafana-panel-json ec2 CloudWatch
```

Switching profile.
```
$ AWS_PROFILE=staging gen-grafana-panel-json ec2 CloudWatch
```

Use `AWS_SDK_LOAD_CONFIG=1` if you use STS assume role in your config file.
```
$ AWS_SDK_LOAD_CONFIG=1 AWS_PROFILE=staging gen-grafana-panel-json ec2 CloudWatch
```

# EC2
To filter EC2, use `--filters`.
```
$ gen-grafana-panel-json ec2 --filters tag:Name,dev-*,instance-type,t2.* CloudWatch
```

# SQS
Give metric name, datasource name and prefix.
```
$ gen-grafana-panel-json sqs --metric-name ApproximateNumberOfMessagesVisible CloudWatch dev-
```

# DynamoDB
```
$ AWS_PROFILE=prod1 gen-grafana-panel-json dynamodb \
    --region ue-west-1 \
    --metric-name ConsumedReadCapacityUnits \
    --statistics Sum \
    CloudWatch-prod
```
