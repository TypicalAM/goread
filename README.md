<p align="center">
    <img width="500" src="assets/cover.png" />
</p>

👋 Hello! goread is an RSS/Atom feed reader for the terminal. It allows you to categorize and follow feeds and read articles right in the commandline! It's accompanied by a beautiful TUI made with [bubble tea](https://github.com/charmbracelet/bubbletea). Features include:

- Categorizing feeds
- Downloading articles for later use
- Offline mode
- Customizable colorschemes
- Stunning TUI

## ❤️ Getting started

### Installing with `go install`

Install with

```
go install github.com/TypicalAM/goread@latest
```

### Installing with `homebrew`

Add repository

```
brew tap TypicalAM/goread
```

Install

```
brew install goread
```

## 📸 What does it look like?

Here is a gif of some basic usage:

<p align="center">
    <img width="700" src="assets/example1.gif" />
</p>

You can use a colorscheme from `pywal` to create a goread colorscheme!

<p align="center">
    <img width="700" src="assets/example2.gif" />
</p>

## ⚙️ Configuration

### 📝 The urls file

The urls file contains the categories and feeds that you are subscribed to! This file is generated by the program in the config directory (usually `~/.config/goread/urls.yml`) and looks similar to this:

```yaml
categories:
  - name: News
    desc: News from around the world
    subscriptions:
      - name: BBC
        desc: News from the BBC
        url: http://feeds.bbci.co.uk/news/rss.xml
  - name: Tech
    desc: Tech news
    subscriptions:
      - name: Wired
        desc: News from the wired team
        url: https://www.wired.com/feed/rss
      - name: Chris Titus Tech (virtualization)
        desc: Chris Titus Tech on virtualization
        url: https://christitus.com/categories/virtualization/index.xml
```

You can edit this file to change the app's contents in an automated manner (remember that you can also edit entries in the TUI!).

### 🌃 The colorscheme file

The colorscheme file contains the colorscheme of your application! It can be generated by hand or using
the `--dump_colors` flag. The colorscheme file is usually at `~/.config/goread/colorscheme.json` - here is how it looks like!

```json
{
  "bg_dark": "#161622",
  "bg_darker": "#11111a",
  "text": "#FFFFFF",
  "text_dark": "#47485b",
  "color1": "#c29fec",
  "color2": "#ddbec0",
  "color3": "#89b4fa",
  "color4": "#e06c75",
  "color5": "#98c379",
  "color6": "#fab387",
  "color7": "#f1c1e4"
}
```

You can use the `--get_colors` flag to generate a colorscheme from pywal. For that you have to supply it with the
pywal `colors.json` file which is usually located at `~/.cache/wal/colors.json`. To generate the `colors.json` file you can run `wal -stni ~/wallpapers/example.png`.

## ✨ Contributing

### TODOs

Here are the things that I've not yet implemented, contributions and suggestions are very welcome!

- [x] URL highlighting and opening
- [x] Automatically theming the glamour viewer
- [ ] AI-Generated feed suggestions
- [ ] Adding customizable keybinds

### Issues

If something doesn't work feel free to create an issue and include:

- Output of `goread --version` if applicable
- Logs are usually located at `/tmp/goread.log` on linux and `%TMP%\goread.log` on Windows

## 💁 Credit where credit is due

### Libraries

The demo was made using [vhs](https://github.com/charmbracelet/vhs/), which is an amazing tool, and you should definitely check it out. The entirety of [charm.sh](https://charm.sh) libraries was essential to the development of this project. The [cobra](https://github.com/spf13/cobra/) library helped to make the cli flags and settings.

### Fonts & logo

The font in use for the logo is sen-regular designed by "Philatype" and licensed under Open Font License. The icon was designed by throwaway icons.

