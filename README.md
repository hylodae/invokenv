# Invokenv

Invokenv is a command-line tool to help you to manage your environment variables across several environments.

Inspired by [direnv](https://direnv.net/) and [summon](https://github.com/cyberark/summon), Invokenv combines the flexibility of the former and the user experience of the latter.

## Install

### From sources

```
$ go install github.com/hylodae/invokenv@latest
```

## Usage

To use Invokenv, you must first create an `invokenv.yml` file containing the list of your environment variables:

```yaml
---
- name: OS_AUTH_URL
  provider: file
  params:
    path: project-config.yaml
    format: yaml
    key: openstack.config.auth_url

- name: OS_USERNAME
  provider: raw
  params:
    value: my-username

- name: OS_PASSWORD
  provider: exec
  params:
    command: gopass show -o my-project/os-password
```

Then, every time you will have to execute a program using these environment variables, just run it as a subcommand of Invokenv:

```
$ invokenv openstack server list
```

### Providers

#### Raw

List of supported parameters:

- `value`: raw value of the environment variable.

Example:

```yaml
# invokenv.yml
---
- name: FOO
  provider: raw
  params:
    value: bar
```

```
$ invokenv echo $FOO
bar
```

#### File

List of supported parameters:

- `path`: path to the file to read.
- `format` (optional): file's content format. The following formats are supported: `yaml`, `json`, `toml`.
- `key` (optional): key corresponding to the value to retrieve from the file (JMESPath expression is supported).

Example:

```yaml
# invokenv.yml
---
- name: FOO
  provider: file
  params:
    path: config.yml
    format: yaml
    key: openstack.projects[0].foo
```

```yaml
# config.yml
---
openstack:
  projects:
    - foo: bar
```

```
$ invokenv echo $FOO
bar
```

#### Exec

List of supported parameters:

- `command`: command to execute, output will be used as value.

Example:

```yaml
# invokenv.yml
---
- name: FOO
  provider: exec
  params:
    command: echo "bar"
```

```
$ invokenv echo $FOO
bar
```

## Troubleshooting

If you encounter any issues while using Invokenv, we recommend checking the documentation and reviewing the existing [Github issues](https://github.com/hylodae/invokenv/issues) for assistance.

If you think you've identified a bug and can't find a related issue, don't hesitate to [submit a new one](https://github.com/hylodae/invokenv/issues/new)! Make sure to provide as much information as possible about your environment.

## Contributing

We gladly welcome [pull requests](https://github.com/hylodae/invokenv/pulls)! Please feel free to suggest improvements or changes to enhance its stability and reliability.

