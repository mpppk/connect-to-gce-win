# Connect to GCE Win 
[![CircleCI](https://circleci.com/gh/mpppk/connect-to-gce-win.svg?style=svg)](https://circleci.com/gh/mpppk/connect-to-gce-win)

## Installation

Download from [GitHub Releases](https://github.com/mpppk/connect-to-gce-win/releases)

## Setup
This setup assume GOOGLE_APPLICATION_CREDENTIALS environment variable is set correctly.
(If you are not ready, see this tutorial)

Edit `~/.config/connect-to-gce-win/.connect-to-gce-win.yaml` as below.

```yaml
userName: your_windows_user_name
project: your_gcp_project_name 
zone: your_gcp_zone
instanceName: target_instance_name(optional)
```

example:

```yaml
userName: mpppk
project: sample-project 
zone: asia-northeast1-b
```

Then, execute binary and your RDP Client will be started.

```
$ ./connect-to-gce-win
```