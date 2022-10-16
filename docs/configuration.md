# Configuration

## The configuration file

The configuration file is split into multiple sections, as folows.
An example of a simple configuration can be found in [deploy/samples](../deploy/samples).

## Settings

```yaml
settings:
  name: simplydash
  enable_health_indicators: true
  theme:
    dark: ...
    light: ...
```

This section controls the settings of the application:
- `name` - this is the name of the dashboard, shown in the UI
- `enable_health_indicators` - this enables the health indicators, which will poll the url of the link for availability
  - the behavior of the health indicators is as follows:
    - green / success if the status is not 404 or 5xx
    - red / error if the status is 404 or 5xx, or the connection has failed
    - yellow / waiting if the status has not yet been determined
- `theme` - this controls the dark and light mode themes for the UI
  - for more information about this section see [Themes](themes.md)

## File providers

```yaml
files:
  path: config/items.yml
  watch: true
```

This section contains a list of files to read for links.
Each item in the list contains:
- `path` - the path to the file, relative to the application directory or absolute
- `watch` - whether to watch this file for changes

See [Files](files.md) for more information about the structure of these files

## Docker

TODO - add section

## Kubernetes

TODO - add section
