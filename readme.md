# Connect to GCE Win 
[![CircleCI](https://circleci.com/gh/mpppk/connect-to-gce-win.svg?style=svg)](https://circleci.com/gh/mpppk/connect-to-gce-win)

## Installation

Download from [GitHub Releases](https://github.com/mpppk/connect-to-gce-win/releases)

## Setup
`connect-to-gce-win` requires that `GOOGLE_APPLICATION_CREDENTIALS` environment variable is set correctly.  
(If you are not ready, see [this tutorial](https://cloud.google.com/docs/authentication/getting-started))

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
(example gif uses [Microsoft Remote Desktop Beta](https://rink.hockeyapp.net/apps/5e0c144289a51fca2d3bfa39ce7f2b06/), but you can use any clients)

```
$ ./connect-to-gce-win
```