# Otter
[![Open Source Love](https://badges.frapsoft.com/os/v3/open-source.png?v=103)](https://github.com/ellerbrock/open-source-badges/)
[![MIT Licence](https://badges.frapsoft.com/os/mit/mit.svg?v=103)](https://opensource.org/licenses/mit-license.php)

> otter tries to extend features of the heroku cli to create new possibilities when handling your deployments. [WIP] :construction:

### Commands
auth   [ -r revoke ]                                               -- authorize otter client with your heroku account
config [ -a app ] [ -l list ] [ -f file ] [ -s set ] [ -r remove ] -- manage your production environment variables
help   [ -h help ]                                                 -- show help info

### Features

#### Config Vars
Config vars are how heroku lets you manage your app's environment variables. Otter offers a more convinient way to control these variables.

- Set multiple config vars at once from your `.env` file
```sh
$ otter config --app guarded-savannah-87990 --file env.yaml
```

### Installation