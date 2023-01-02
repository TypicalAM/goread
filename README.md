<p align="center">
    <img width="400" src="assets/cover.png" />
</p>

👋 Hello! This is an RSS reader which allows you to browse the web better! It's accompanied by a beautiful TUI made with [bubble tea](https://github.com/charmbracelet/bubbletea)

## 🌃 Getting started

Getting it up and running is pretty easy! This is how you do it.

### 1. Clone the repository

```sh
git clone --depth=1 https://github.com/TypicalAM/goread && cd goread
```

### 2. Build the executable

```sh
go build -o goread cmd/goread/main.go
```

### 3. Run the program!

```sh
./goread
```

## ✨ Tasks to do

Here are the things that I've not yet implemented

- [ ] Waiting for the window size message before adding the first tab
- [ ] Help interface with the key bindings
- [ ] Adding and removing categories and feeds
- [ ] A main category where all the feeds are aggregated

## 📸 Here is a demo of what it looks like:

<p align="center">
    <img width="700" src="assets/example1.gif" />
</p>

## 💁 Credit where credit is due

### Libraries

The demo was made using [vhs](https://github.com/charmbracelet/vhs/), which is an amazing tool, and you should definitely check it out. The entirety of [charm.sh](https://charm.sh) libraries was essential to the development of this project.

### Fonts & logo

The font in use for the logo is sen-regular designed by "Philatype" and licensed under Open Font License. The icon was designed by throwaway icons.

