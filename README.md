![facebook_cover_photo_2](https://github.com/JunNishimura/Goit/assets/28744711/919d78ca-52bf-481d-883e-b17cf0b9ea69)

<p align='center'>
  <img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/JunNishimura/Goit">
  <img alt="GitHub" src="https://img.shields.io/github/license/JunNishimura/Goit">
  <a href="https://github.com/JunNishimura/Goit/actions/workflows/test.yml"><img src="https://github.com/JunNishimura/Goit/actions/workflows/test.yml/badge.svg" alt="test"></a>
  <a href="https://goreportcard.com/report/github.com/JunNishimura/Goit"><img src="https://goreportcard.com/badge/github.com/JunNishimura/Goit" alt="Go Report Card"></a>
</p>

# Goit - Git made by Golang

## 📖 Overview
Goit is version control tool just like Git made by Golang.

## 💻 Installation
### Homebrew Tap
```
brew install JunNishimura/tap/Goit
```

### go intall
```
go install github.com/JunNishimura/Goit@latest
```

## 🔨 Commands
### Available
- [x] `init` - initialize Goit, make .goit directory where you init
- [x] `add` - make goit object and register to index
- [x] `commit` - make commit object
- [x] `log` - show commit history
- [x] `config` - set config. e.x.) name, email
- [x] `cat-file` - show goit object data
- [x] `ls-files` - show index
- [x] `hash-object` - show hash of file
- [x] `rev-parse` - show hash of reference such as branch, HEAD
- [x] `update-ref` - update reference
- [x] `write-tree` - write tree object

### Future
- [ ] tag
- [ ] rm
- [ ] switch
- [ ] checkout
- [ ] branch
- [ ] merge
- [ ] stash
- [ ] reset
- [ ] revert
- [ ] restore
- [ ] diff
- [ ] read-tree
- [ ] symbolic-ref

## 👀 How to use
### 0. Install Goit
see Installation above.

### 1. Move to the directory you want to use goit
```
cd /home/usr/sample 
```

### 2. Initialize Goit
```
goit init
```

### 3. add test.txt
```
echo "Hello, World" > test.txt
```

### 4. run `add` command to convert text file to goit object
```
goit add test.txt
```

### 5. run `config` command to set config
```
goit config --global user.name 'Goit'
goit config --global uesr.email 'goit@example.com'
```

### 6. run `commit` command to make commit object
```
goit commit -m "init"
```

### 7. look at inside commit object
```
goit rev-parse main
goit cat-file -p `hash`
```

Perfect!! You now understand the basic of goit🎉


## 🪧 License
Goit is released under MIT License. See [MIT](https://raw.githubusercontent.com/JunNishimura/Goit/main/LICENSE)
