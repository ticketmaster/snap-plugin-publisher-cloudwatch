# Snap CloudWatch Publisher Plugin
## Changlog
Please see [changlog](CHANGELOG.md).

## Getting Started
To get started, you will need snap (v0.14.0-beta or above) running to receive and aggregate sampling data points.
Of course, you will need a working Amazon Web Service (AWS) account so the data points can be published onto CloudWatch.

## System Requirements
* [golang 1.5+](https://golang.org/dl/) (needed only for building)
* [snap](https://github.com/intelsdi-x/snap)

### Operating Systems
All OSs currently supported by snap:
* Linux/amd64
* Darwin/amd64


### Installation
#### Download Snap CloudWatch Publisher Plugin binary:
There is no pre-built binary avaiable yet.

#### To build the plugin binary:
Fork https://github.com/Ticketmaster/snap-plugin-publisher-cloudwatch.git
Clone repo into `$GOPATH/src/github.com/ticketmaster/`:

```
$ git clone https://github.com/Ticketmaster/snap-plugin-publisher-cloudwatch.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `/build/rootfs/`

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* Ensure `$SNAP_PATH` is exported  
`export SNAP_PATH=$GOPATH/src/github.com/intelsdi-x/snap/build`

### Configure Amazon Web Service (AWS)
* Install aws command line and [configure](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-config-files) it correctly.
* This plugin supports aws_access_key_id, aws_secret_access_key, and aws_session_token.  It will also works with EC2 instance using IAM Roles/Policies (CloudWatchFullAccess).

### Example Task
```json
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "1s"
    },
    "workflow": {
        "collect": {
            "metrics": {
                "/intel/psutil/load/load1": {},
                "/intel/psutil/load/load5": {},
                "/intel/psutil/load/load15": {},
                "/intel/psutil/vm/available": {},
                "/intel/psutil/vm/free": {},
                "/intel/psutil/vm/used": {}
            },
            "config": {
                "/intel/mock": {
                    "password": "secret",
                    "user": "root"
                }
            },
            "tags": {
                "/intel": {
                   "cluster": "cluster_xyz",
                   "region": "us-east-1"
                }
            },
            "process": [
                {
                    "plugin_name": "passthru",
                    "plugin_version": 1,
                    "process": null,
                    "publish": [
                        {
                            "plugin_name": "cloudwatch",
                            "plugin_version": 1,
                            "config": {
                                "region": "us-east-1",
                                "namespace": "snap"                    
                            }
                        }
                    ],
                    "config": null
                }
            ],
            "publish": null
        }
    }
}
```
Create task:
```
$ snapctl task create -t sample-task.json
```
