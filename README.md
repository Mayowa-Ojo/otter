# Otter
[![Open Source Love](https://badges.frapsoft.com/os/v3/open-source.png?v=103)](https://github.com/ellerbrock/open-source-badges/)
[![MIT Licence](https://badges.frapsoft.com/os/mit/mit.svg?v=103)](https://opensource.org/licenses/mit-license.php)

> otter tries to extend features of the heroku cli to create new possibilities when handling your deployments. [WIP :construction:]

### Commands
auth   [ -r revoke ]                                               -- authorize otter client with your heroku account

config [ -a app ] [ -l list ] [ -f file ] [ -s set ] [ -r remove ] -- manage your production environment variables

help   [ -h help ]                                                 -- show help info

### Documentation

#### Auth
Otter requires authorization via heroku oauth
- authorize otter client: `$ otter auth`
- revoke authorization: `$ otter auth --revoke`

#### Config Vars
Config vars are how heroku lets you manage your app's environment variables. Otter offers a more convinient way to control these variables.

- set single variable: `$ otter config --app guarded-savannah-87990 --set PORT:8870`
- Set multiple variables at once from your `.env` file [also supports `json` and `yaml`]: `$ otter config --app guarded-savannah-87990 --file env.yaml`
- list variables: `$ otter config --app guarded-savannah-87990 --list`

### Installation
If you have go installed [v1.13+], you can clone this repository and run go install or go build <path/to/executable>.

You can also download a pre-built binary from [releases](https://github.com/Mayowa-Ojo/otter/releases)
#### Build from source
```sh
$ git clone https://github.com/Mayowa-Ojo/otter
$ cd otter
$ go install 
```

#### Releases
Linux
```sh
$ wget -P <path/to/downloads> https://github.com/Mayowa-Ojo/otter/releases/download/v0.1/otter.v0.1
$ mv <path/to/downloads>/otter.v0.1 <$GOPATH>
$ otter --help
```