# Themes

The default theme is based on the GNOME Adwaita colors. 
These can be customized as follows, using css color values (hex, rgb or rgba):
- `background` - the background color for the page
- `element_background` - the background of the page elements (buttons, search, etc.)
- `foreground` - text color
- `foreground_secondary` - secondary text color (descriptions), also used for icons
- `accent_color` - accent color - used on hover over elements
- `success_color` - color shown on health indicators for success
- `warning_color` - color shown on health indicators for waiting
- `error_color` - color shown on health indicators for error

## Example themes

### Solarized

Original colors by [Ethan Schoonover](https://ethanschoonover.com/solarized/)

```yaml
  theme:
    dark:
      background: rgb(0, 43, 54)
      element_background: rgb(7, 54, 66)
      foreground: rgb(131, 148, 150)
      foreground_secondary: rgb(88, 110, 117)
      accent_color: rgb(38, 139, 210)
      success_color: rgb(133, 153, 0)
      warning_color: rgb(181, 137, 0)
      error_color: rgb(220, 50, 47)
    light:
      background: rgb(253, 246, 227)
      element_background: rgb(238, 232, 213)
      foreground: rgb(101, 123, 131)
      foreground_secondary: rgba(88, 110, 117)
      accent_color: rgb(38, 139, 210)
      success_color: rgb(133, 153, 0)
      warning_color: rgb(181, 137, 0)
      error_color: rgb(220, 50, 47)
```