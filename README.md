# alertService
AlertManager service for Prometheus

 - [x] lark notification alarm message
 - [ ] enterprise weChat notification alarm message
 - [x] alarm message persistence

## How to use
### 1. compile service
```shell
$ git clone https://github.com/gagraler/alert-service.git

$ cd alert-service

$ make -f build/Makefile
```

### 2. run service
Please modify the contents in the configuration file before running the service

```shell
$ mysql -u root -p'123456' < alert-service.sql

$ export LARK_BOT_SIGN_KEY=secret_key
$ export LARK_BOT_URL=lark_url
$ nohup ./alert-service > alert-service.log 2>&1 &
```

#### build docker images
```shell
$ docker pull ghcr.io/gagraler/alert-service:latest

$ docker run -d \
    -e LARK_BOT_SIGN_KEY=secret_key \
    -e LARK_BOT_URL=lark_url \
    -p 8588:8588 \
    alert-service
    -v /etc/alert-service.toml:/opt/alert-service.toml
```

## HTTP API
HTTP Api The unified prefix defaults to `/api/v1/alertService`

### test api
WebHook API `alertMessage/hook`
```shell
$ curl --location 'http://localhost:8588/api/v1/alertService/alertMessage/hook' \
--header 'Content-Type: application/json' \
--data '{
    "receiver":"test",
    "status":"firing",
    "alerts":[
        {
            "status":"firing",
            "labels":{
                "alertname":"TestAlert",
                "instance":"Grafana"
            },
            "annotations":{
                "summary":"Notification test"
            },
            "startsAt":"2024-01-11T15:04:03.277535848Z",
            "endsAt":"0001-01-01T00:00:00Z",
            "generatorURL":"",
            "fingerprint":"57c6d9296de2ad39",
            "silenceURL":"http://localhost:3000/alerting/silence/new?alertmanager=grafana&matcher=alertname%3DTestAlert&matcher=instance%3DGrafana",
            "dashboardURL":"",
            "panelURL":"",
            "values":null,
            "valueString":"[ metric='\''foo'\'' labels={instance=bar} value=10 ]"
        }
    ],
    "groupLabels":{
        "alertname":"TestAlert",
        "instance":"Grafana"
    },
    "commonLabels":{
        "alertname":"TestAlert",
        "instance":"Grafana"
    },
    "commonAnnotations":{
        "summary":"Notification test"
    },
    "externalURL":"http://localhost:3000/",
    "version":"1",
    "groupKey":"test-57c6d9296de2ad39-1704985443",
    "truncatedAlerts":0,
    "orgId":1,
    "title":"[FIRING:1] TestAlert Grafana ",
    "state":"alerting",
    "message":"**Firing**\n\nValue: [no value]\nLabels:\n - alertname = TestAlert\n - instance = Grafana\nAnnotations:\n - summary = Notification test\nSilence: http://localhost:3000/alerting/silence/new?alertmanager=grafana&matcher=alertname%3DTestAlert&matcher=instance%3DGrafana\n"
}'
```

## License and Copyright
[MIT](https://choosealicense.com/licenses/mit/)

- Email: gagral@sina.com
