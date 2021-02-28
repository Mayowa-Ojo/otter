# Otter
[![Open Source Love](https://badges.frapsoft.com/os/v3/open-source.png?v=103)](https://github.com/ellerbrock/open-source-badges/)
[![MIT Licence](https://badges.frapsoft.com/os/mit/mit.svg?v=103)](https://opensource.org/licenses/mit-license.php)


> otter tries to extend features of the heroku cli to create new possibilities when handling your deployments.

### Features
Otter can be considered bleeding edge in many ways. As the project grows, new features will be implemented.

#### Config Vars
Config vars are how heroku lets you manage your app's environment variables. Otter offers a more convinient way to control these variables.

- Set multiple config vars at once from your `.env` file
```
$ otter auth
$ otter env:set -f .env
```