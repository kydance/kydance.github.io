# 终端神器：一篇文章玩转终端神器 Tmux，ZSH，Oh-My-Zsh


{{< admonition type=abstract title="导语" open=true >}}
在现代开发环境中，高效的终端管理是提升工作效率的关键。
Tmux 作为一款强大的终端复用工具，不仅能让你在一个终端窗口中同时操作多个会话，还能实现窗口分割、会话保持等高级功能。
无论是本地开发还是远程服务器管理，Tmux 都能让你的终端操作更加得心应手。
本文将带你全面了解 Tmux 的各项功能，从基础操作到高级配置，让你的终端使用效率得到质的飞跃。
{{< /admonition >}}

<!--more-->

## TMux

### Feature

- 强劲的、易于使用的命令行界面
- 可以横向、纵向分割窗口
- 窗格可以自由移动和调整大小，或直接利用四个预设布局之一
- 支持 UTF-8 编码及 256 色终端
- 可在多个缓冲区进行复制和粘贴
- 可通过交互式菜单来选择窗口、会话及客户端
- 支持跨窗口搜索
- 支持自动及手动锁定窗口
- 可以自由配置绑定快捷键

### Tmux 中的 server, session, window 和 Pane

在 Tmux 系统中，存在以下极其重要的大小层级: `Server` -> `Session` -> `Window` -> `Pane`.

- **`Server`**: 整个 tmux 的后台服务. NOTE: 当配置文件不生效时，就需要使用 `tmux kill-server` 来重启 Tmux
- **`Session`**: 可以理解为 workplace
- **`Window`**: 相当于 VIM 中的buffer
- **`Pane`**: 窗口中的小分屏，相当于 VIM 中的 `split` / `vsplit`

### Installation

Require: version >= 2.1

#### Linux

```shell
yum install -y tmux
yay -S tmux
apt-get install tmux
```

#### Mac

```shell
brew install tmux
```

### 常用命令

Tmux 的默认 prefix-key 是 `<C-b>`

#### 启动新 session

```shell
$ tmux [new -s sessionName -n windowName]
# e.g. tmux new -s kyden -n nvim
```

#### 恢复 Session

```shell
tmux at[-t sessionName]
```

#### Session List

```shell
tmux ls
```

#### 关闭 Session

```shell
tmux kill-session -t sessionName
```

#### 关闭整个 tmux 服务器

```shell
tmux kill-server
```

#### Session Command

| prefix-key | command | description |
| :--- | :--- | :--- |
| `<C-b>` | `?` | 显示快捷键帮助文档 |
| `<C-b>` | `d` | 断开当前 Session |
| `<C-b>` | `r` | 强制重载当前 Session |
| `<C-b>` | `:` | 进入命令模式，可直接输入命令 |

#### Window Command

| prefix-key | command | description |
| :--- | :--- | :--- |
| `<C-b>` | `c` | 新建窗口 |
| `<C-b>` | `&` | 关闭当前窗口 |
| `<C-b>` | `p / n / <number>` | 切换到上一个 / 下一个 / 指定窗口 |
| `<C-b>` | `w` | 打开窗口列表，用于切换窗口 |
| `<C-b>` | `,` | 重命名当前窗口 |
| `<C-b>` | `.` | 修改窗口编号 |

#### Pane Command

| prefix-key | command | description |
| :--- | :--- | :--- |
| `<C-b>` | `"` / `%` | 新建上下 / 左右 pane |
| `<C-b>` | `x` | 关闭当前 pane |
| `<C-b>` | `z` | 最大化当前 pane(重复取消最大化) |
| `<C-b>` | `q` | 显示面板编号，在编号消失前输入对应的数字可切换到相应的面板 |
| `<C-b>` | `<left>` / `<right>` / `up` / `down` | 移动光标切换面板 |
| `<C-b>` | `o` | 选择下一 pane |
| `<C-b>` | `<space>` | 在自带的面板布局中循环切换 |

### 配置文件

配置文件 `.tmux.conf` 通常位于 `~/.tmux.conf` 处，可输入 `restart tmux` 进行 mtux 重启

```conf
# recover colorful terminal
set -g default-terminal "xterm-256color"

# 窗口面板起始序列号
set -g base-index 1
set -g pane-base-index 1

# 从tmux v1.6版起，支持设置第二个指令前缀，使用 ` 作为第二指令前缀
# set-option -g prefix2 `

# (Tmux v2.1) 支持鼠标: 选取文本、调整面板大小、选中并切换面板
set-option -g mouse on

# 状态栏窗口名称格式
set -wg window-status-format " #I #W "
# 状态栏当前窗口名称格式(#I：序号，#w：窗口名称，#F：间隔符)
set -wg window-status-current-format " #I:#W#F "
# 状态栏窗口名称之间的间隔
set -wg window-status-separator ""

# 开启vi风格后，支持vi的C-d、C-u、hjkl等快捷键
setw -g mode-keys vi
# 绑定 Escape 进入 复制 模式
bind Escape copy-mode

setw -g automatic-rename off
setw -g allow-rename off
```

## zsh

### zsh 介绍与安装

[ZSH](https://www.zsh.org/) 是一个兼容 bash 的 shell，相较于 bash 具有以下优点：

- Tab 补全功能强大。命令、命令参数、文件路径均可以补全
- 插件丰富。快速输入以前使用过的命令、快速跳转文件夹、显示系统负载这些都可以通过插件实现
- 主题丰富
- 可定制性高

Installation

```shell
# macos
brew install zsh

# Arch Linux
pacman -S zsh
```

**设置 zsh 为默认 shell**: `chsh -s /bin/zsh`

### oh-my-zsh

`cURL` 下载并安装 oh-my-zsh

```zsh
sh -c "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
```

#### powerlevel10k theme

```zsh
# 下载
git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ${ZSH_CUSTOM:-$HOME/.oh-my-zsh/custom}/themes/powerlevel10k

# 设置 in .zshrc
ZSH_THEME="powerlevel10k/powerlevel10k"
```

#### zsh-autosuggestions

[**zsh-autosuggestions**](https://github.com/zsh-users/zsh-autosuggestions) 是一个命令提示插件，当你输入命令时，会自动推测你可能需要输入的命令，按下右键可以快速采用建议

```zsh
# Download
git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions

# Set in .zshrc
plugins=(... zsh-autosuggestions ...)
```

#### zsh-syntax-highlighting

[zsh-syntax-highlighting](https://github.com/zsh-users/zsh-syntax-highlighting) 是一个命令语法校验插件，在输入命令的过程中，若指令不合法，则指令显示为红色，若指令合法就会显示为绿色。

```zsh
# Download
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting

# Set in .zshrc
plugins=(... zsh-syntax-highlighting ...)
```

#### z

oh-my-zsh 内置了 z 插件。z 是一个文件夹快捷跳转插件，对于曾经跳转过的目录，只需要输入最终目标文件夹名称，就可以快速跳转，避免再输入长串路径，提高切换文件夹的效率。

```zsh
# Set in .zshrc
plugins=(... z ...)
```

#### extract

oh-my-zsh 内置了 extract 插件。extract 用于解压任何压缩文件，不必根据压缩文件的后缀名来记忆压缩软件

使用 x 命令即可解压文件

```zsh
# Set in .zshrc
plugins=(... extract ...)
```

##### Supported file extensions

| Extension         | Description                          |
| :---------------- | :----------------------------------- |
| `7z`              | 7zip file                            |
| `Z`               | Z archive (LZW)                      |
| `apk`             | Android app file                     |
| `aar`             | Android library file                 |
| `bz2`             | Bzip2 file                           |
| `cab`             | Microsoft cabinet archive            |
| `cpio`            | Cpio archive                         |
| `deb`             | Debian package                       |
| `ear`             | Enterprise Application aRchive       |
| `exe`             | Windows executable file              |
| `gz`              | Gzip file                            |
| `ipa`             | iOS app package                      |
| `ipsw`            | iOS firmware file                    |
| `jar`             | Java Archive                         |
| `lrz`             | LRZ archive                          |
| `lz4`             | LZ4 archive                          |
| `lzma`            | LZMA archive                         |
| `obscpio`         | cpio archive used on OBS             |
| `rar`             | WinRAR archive                       |
| `rpm`             | RPM package                          |
| `sublime-package` | Sublime Text package                 |
| `tar`             | Tarball                              |
| `tar.bz2`         | Tarball with bzip2 compression       |
| `tar.gz`          | Tarball with gzip compression        |
| `tar.lrz`         | Tarball with lrzip compression       |
| `tar.lz`          | Tarball with lzip compression        |
| `tar.lz4`         | Tarball with lz4 compression         |
| `tar.xz`          | Tarball with lzma2 compression       |
| `tar.zma`         | Tarball with lzma compression        |
| `tar.zst`         | Tarball with zstd compression        |
| `tbz`             | Tarball with bzip compression        |
| `tbz2`            | Tarball with bzip2 compression       |
| `tgz`             | Tarball with gzip compression        |
| `tlz`             | Tarball with lzma compression        |
| `txz`             | Tarball with lzma2 compression       |
| `tzst`            | Tarball with zstd compression        |
| `vsix`            | VS Code extension zip file           |
| `war`             | Web Application archive (Java-based) |
| `whl`             | Python wheel file                    |
| `xpi`             | Mozilla XPI module file              |
| `xz`              | LZMA2 archive                        |
| `zip`             | Zip archive                          |
| `zlib`            | zlib archive                         |
| `zst`             | Zstandard file (zstd)                |
| `zpaq`            | Zpaq file                            |

#### web-search

oh-my-zsh 内置了 web-search 插件。web-search 能让我们在命令行中使用搜索引擎进行搜索。使用搜索引擎关键字+搜索内容 即可自动打开浏览器进行搜索

使用 web-search 命令即可搜索

```zsh
# Set in .zshrc
plugins=(... web-search ...)
```

例如，这两个是等价的:

```zsh
web_search google oh-my-zsh
google oh-my-zsh
```

一些常用的搜索上下文如下：

| Context               | URL                                             |
| --------------------- | ----------------------------------------------- |
| `bing`                | `https://www.bing.com/search?q=`                |
| `google`              | `https://www.google.com/search?q=`              |
| `github`              | `https://github.com/search?q=`                  |
| `baidu`               | `https://www.baidu.com/s?wd=`                   |
| `youtube`             | `https://www.youtube.com/results?search_query=` |
| `chatgpt`             | `https://chatgpt.com/?q=`                       |

#### Uninstall

```zsh
uninstall_oh_my_zsh
```

### Q&A

1. Kitty SSH 远程字符问题，例如`删除(Backspace)键` => `空格(space)`，无法正常工作？

> 检查终端 TERM 设置，确保在 Kitty 本地终端的 `~/.zshrc` 中添加以下内容，
> 【Kitty 默认使用 xterm-kitty，有时远程服务器不支持它，可以尝试更改为 xterm-256color】
>
> ```bash
> export TERM=xterm-256color
> ```


---

> Author: [kyden](https://github.com/kydance)  
> URL: http://kydance.github.io/posts/terminal-configure-guide/  

