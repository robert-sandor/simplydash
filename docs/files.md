# File providers

These providers are simple YAML files, that contain a list of categories and links.

Example:
```yaml
- name: bookmarks
  items:
    - name: youtube
      url: https://youtube.com
      icon: https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/youtube.png
      description: youtube.com
```

The structure of the file is as follows:
- List of Categories, each containing:
  - `name` - this is the name of the category
    - does not need to be unique, as items from the categories with the same name will be merged
  - `items`
    - List of items in the category - each item will be displayed as its own link in the dashboard
      - `name` - this is the name of the application the item links to
      - `url` - the link to the application itself
      - `icon` - a link to an icon for the application
        - if missing, will attempt to find and icon from [walkxcode/Dashboard-Icons](https://github.com/walkxcode/Dashboard-Icons)
        - to add a custom icon, set this to the relative path to the icon like `icons/mycustomicon.png`
      - `description` - a description of what the application does
        - this defaults to the value of the `url`
