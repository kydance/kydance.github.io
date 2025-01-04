# Tmux 使用指南


{{&lt; admonition type=abstract title=&#34;导语&#34; open=true &gt;}}
Tmux，一款优秀的终端复用工具，使用它最直观的好处就是，通过一个终端登录远程主机并运行tmux后，在其中可以开启多个控制台而无需再“浪费”多余的终端来连接这台远程主机
{{&lt; /admonition &gt;}}

&lt;!--more--&gt;

## Feature

- 强劲的、易于使用的命令行界面
- 可以横向、纵向分割窗口
- 窗格可以自由移动和调整大小，或直接利用四个预设布局之一
- 支持 UTF-8 编码及 256 色终端
- 可在多个缓冲区进行复制和粘贴
- 可通过交互式菜单来选择窗口、会话及客户端
- 支持跨窗口搜索
- 支持自动及手动锁定窗口
- 可以自由配置绑定快捷键

## Tmux 中的 server, session, window 和 Pane

在 Tmux 系统中，存在以下极其重要的大小层级: `Server` -&gt; `Session` -&gt; `Window` -&gt; `Pane`.

- **`Server`**: 整个 tmux 的后台服务. NOTE: 当配置文件不生效时，就需要使用 `tmux kill-server` 来重启 Tmux
- **`Session`**: 可以理解为 workplace
- **`Window`**: 相当于 VIM 中的buffer
- **`Pane`**: 窗口中的小分屏，相当于 VIM 中的 `split` / `vsplit`

## Installation

Require: version &gt;= 2.1

### Linux

```shell
yum install -y tmux
yay -S tmux
apt-get install tmux
```

### Mac

```shell
brew install tmux
```

## 常用命令

Tmux 的默认 prefix-key 是 `&lt;C-b&gt;`

### 启动新 session

```shell
$ tmux [new -s sessionName -n windowName]
# e.g. tmux new -s kyden -n nvim
```

### 恢复 Session

```shell
tmux at[-t sessionName]
```

### Session List

```shell
tmux ls
```

### 关闭 Session

```shell
tmux kill-session -t sessionName
```

### 关闭整个 tmux 服务器

```shell
tmux kill-server
```

### Session Command

| prefix-key | command | description |
| :--- | :--- | :--- |
| `&lt;C-b&gt;` | `?` | 显示快捷键帮助文档 |
| `&lt;C-b&gt;` | `d` | 断开当前 Session |
| `&lt;C-b&gt;` | `r` | 强制重载当前 Session |
| `&lt;C-b&gt;` | `:` | 进入命令模式，可直接输入命令 |

### Window Command

| prefix-key | command | description |
| :--- | :--- | :--- |
| `&lt;C-b&gt;` | `c` | 新建窗口 |
| `&lt;C-b&gt;` | `&amp;` | 关闭当前窗口 |
| `&lt;C-b&gt;` | `p / n / &lt;number&gt;` | 切换到上一个 / 下一个 / 指定窗口 |
| `&lt;C-b&gt;` | `w` | 打开窗口列表，用于切换窗口 |
| `&lt;C-b&gt;` | `,` | 重命名当前窗口 |
| `&lt;C-b&gt;` | `.` | 修改窗口编号 |

### Pane Command

| prefix-key | command | description |
| :--- | :--- | :--- |
| `&lt;C-b&gt;` | `&#34;` / `%` | 新建上下 / 左右 pane |
| `&lt;C-b&gt;` | `x` | 关闭当前 pane |
| `&lt;C-b&gt;` | `z` | 最大化当前 pane(重复取消最大化) |
| `&lt;C-b&gt;` | `q` | 显示面板编号，在编号消失前输入对应的数字可切换到相应的面板 |
| `&lt;C-b&gt;` | `&lt;left&gt;` / `&lt;right&gt;` / `up` / `down` | 移动光标切换面板 |
| `&lt;C-b&gt;` | `o` | 选择下一 pane |
| `&lt;C-b&gt;` | `&lt;space&gt;` | 在自带的面板布局中循环切换 |

## 配置文件

配置文件 `.tmux.conf` 通常位于 `~/.tmux.conf` 处，可输入 `restart tmux` 进行 mtum 重启

```conf
# recover colorful terminal
set -g default-terminal &#34;xterm-256color&#34;

# 窗口面板起始序列号
set -g base-index 1
set -g pane-base-index 1

# 从tmux v1.6版起，支持设置第二个指令前缀，使用 ` 作为第二指令前缀
# set-option -g prefix2 `

# (Tmux v2.1) 支持鼠标: 选取文本、调整面板大小、选中并切换面板
set-option -g mouse on

# 状态栏窗口名称格式
set -wg window-status-format &#34; #I #W &#34;
# 状态栏当前窗口名称格式(#I：序号，#w：窗口名称，#F：间隔符)
set -wg window-status-current-format &#34; #I:#W#F &#34;
# 状态栏窗口名称之间的间隔
set -wg window-status-separator &#34;&#34;

# 开启vi风格后，支持vi的C-d、C-u、hjkl等快捷键
setw -g mode-keys vi
# 绑定 Escape 进入 复制 模式
bind Escape copy-mode

setw -g automatic-rename off
setw -g allow-rename off
```


---

> : [kyden](https:github.com/kydance)  
> URL: http://kyden.us.kg/posts/tmux-guide/  

