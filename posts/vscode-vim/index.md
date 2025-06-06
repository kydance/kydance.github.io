# 【最佳实践】VSCode + Vim = 效率之神


{{< admonition type=abstract title="导语" open=true >}}
想要在现代编辑器中获得极致的编码体验？VSCode 与 Vim 的强强联合将带给你意想不到的效率提升。
本文将详细介绍如何通过 VSCodeVim 插件，在 VSCode 中完美复刻 Vim 的操作方式，让你既能享受 VSCode 强大的功能生态，
又能保持 Vim 快速高效的编辑体验。无论你是 Vim 老手还是新手，这份完整指南都能帮你打造一个更高效的编码环境。
{{< /admonition >}}

<!--more-->

## I. VSCodeVim

<!-- {{< figure src="/posts/vscode-vim/vim.png" title="" height=0.5 width=0.5 >}} -->

[VSCodeVim](https://github.com/VSCodeVim/Vim/#key-remapping) 是一款 vim 模拟器，
它将 vim 的大部分功能都集成在了 VSCode 中，即一个嵌入在 VSCode 中的 vim。

正是由于 VSCodeVim 本质上只是一个 Vim 模拟器，而非真正的 Vim，导致原生 Vim 中的有些功能并不支持，具体支持情况见 [roadmap](https://github.com/VSCodeVim/Vim/blob/master/ROADMAP.md)。
尽管它现在还无法完全模拟 Vim，但这依然不妨碍它的优秀。

| Status             | Command                |
| ------------------ | ---------------------- |
| ✅ | Normal Mode            |
| ✅ | Insert Mode            |
| ✅ | Visual Mode            |
| ✅ | Visual Line Mode       |
| ✅ | Number Prefixes        |
| ✅ | . Operator             |
| ✅ | Searching with / and ? |
| ✅ | Correct Undo/Redo      |
| ⚠ | Command Remapping      |
| ⚠️ | Marks                  |
| ✅ | Text Objects           |
| ✅ | Visual Block Mode      |
| ✅ | Replace Mode           |
| ✅ | Multiple Select Mode   |
| ⚠ | Macros                 |
| ⚠ | Buffer/Window/Tab      |

✅ - command done

⚠️ - some variations of the command are not supported

---

## II. 安装

只需在 VSCode 的插件商店搜索 `vim` 就能找到该插件.

{{< figure src="/posts/vscode-vim/vscode-extension-vim.png" title="" height="256" width="256" >}}

{{< admonition type=tip title="关闭 Mac 的重复键" open=true >}}

当使用 Mac 时，需要输入以下代码，用以关闭 Mac 的重复键

```bash
# For VS Code
$ defaults write com.microsoft.VSCode ApplePressAndHoldEnabled -bool false
# For VS Codium
$ defaults write com.vscodium ApplePressAndHoldEnabled -bool false
# If necessary, reset global default
$ defaults delete -g ApplePressAndHoldEnabled
```

{{< /admonition >}}

---

## III. 文件配置详解

VSCodeVim 的相关配置文件是放在 `settings.json` 中，而不是 `vimrc` 文件.

- 对于**非代码编辑区**的热键将其定义在 `keybindings.json` 中
- 对于**代码编辑区**且属于 vim 的热键将其定义在 `settings.json` 文件中

{{< admonition type=tip title="个人配置参考" open=true >}}
[个人 vscode 配置文件参考：https://github.com/kydance/dotfiles](https://github.com/kydance/dotfiles)
{{< /admonition >}}

### 1. 基础配置

```json
// leader - prefix key
"vim.leader": "<space>",
// To improve performance
"extensions.experimental.affinity": { 
    "vscodevim.vim": 1 
},
// Easy motion
"vim.easymotion": true,
// Use system clipboard
"vim.useSystemClipboard": true,
// 由vim接管ctrl+any的按键，而不是vscode
"vim.useCtrlKeys": true,
"vim.replaceWithRegister": true,
// 忽略大小写
"vim.ignorecase": true,
"vim.smartcase": true,
// 智能行号
"vim.smartRelativeLine": true,
"vim.foldfix": true,
// Highlight search
"vim.hlsearch": true,
"vim.highlightedyank.enable": true,
"vim.highlightedyank.duration": 500,

// 由vscode进行处理，而不是vscode-vim插件
"vim.handleKeys": {
    "<C-d>": true, // 向下滚动半页
    "<C-f>": true, // 向下滚动一页
    "<C-e>": true, // 向下滚动一行
    "<C-s>": true,
    "<C-z>": false,
    "<C-a>": true,
    "<C-c>": true,
    "<C-v>": true,
    "<C-x>": true,
},
```

---

### 2. 快捷键配置

#### NORMAL Mode

```json
"vim.normalModeKeyBindingsNonRecursive": [
    {   // 聚集 terminal
        "before": ["<C-j>"], 
        "commands": ["workbench.action.terminal.focus"]
    },
    {   // 语义级 重命名
        "before": ["leader", "r"],  
        "commands": [ "editor.action.rename"] 
    },
    {   // 
        "before": ["g", "h"],  
        "commands": [ "editor.action.showHover"]
    },
    {
        "before": ["g", "d"],
        "commands": ["editor.action.revealDefinition"]
    },
    {
        "before": ["g", "r"],
        "commands": ["editor.action.goToReferences"],
    },
    {
        "before": ["g", "i"],
        "commands": ["editor.action.goToImplementation"],
    },
    {
        "before": ["g", "b"],
        "commands": ["workbench.action.navigateBack"],
    },

    {
        "before": ["leader", "e", "f"],
        "commands": ["workbench.explorer.fileView.focus"],
    },

    {
        "before": [ "leader", "leader", "e", "f"],
        "commands": ["workbench.action.toggleActivityBarVisibility"]
    },

    {
        "before": ["leader", "g", "g"],
        "commands": ["workbench.action.quickOpen"],
    },
    {   // Global find
        "before": ["leader", "g", "f"],
        "commands": ["workbench.view.search"],
    },

    {
        "before": ["g", "[",], 
        "commands": ["editor.action.marker.prevInFiles"],
    },
    {
        "before": ["g", "]",],
        "commands": ["editor.action.marker.nextInFiles"],
    },
    {   // Source Control Git
        "before": [ "leader", "g", "i", "t" ],
        "commands": ["workbench.scm.focus"]
    },
    { // Start to debug
        "before": [ "leader", "d" ],
        "commands": ["workbench.action.debug.start"]
    },
    {
        "before": ["leader","w"],
        "commands": [":w!" ]
    },
    {
        "before": ["leader","q"],
        "commands": [":q" ]
    },
    {   // No highlight
        "before": ["leader", "n", "h"],
        "commands": [":nohl"]
    },
    {
        "before": ["H"],
        "after": ["^"]
    },
    {
        "before": ["L"], 
        "after": ["$"]
    },
    {   // Blockwise visual mode
        "before": ["\\"],
        "commands": ["extension.vim_ctrl+v"]
    },
    {
        "before": ["leader", "t"],
        "commands": [":terminal"] 
    }, 

    {
        "before": ["g", "t"],
        "commands": [":tabnext"]
    }, 
    {
        "before": ["g", "T"],
        "commands": [":tabprev"] 
    },
    {   // project-manager
        "before": ["leader", "p", "m"], 
        "commands": [{
            "command":"workbench.view.extension.project-manager",
            "when":"viewContainer.workbench.view.extension.project-manager.enabled"
        }]
    },
],
```

#### INSERT Mode

```json
"vim.insertModeKeyBindings": [
    {
        "before": ["j", "k"],
        "after": ["<Esc>"]
    }, 
],
```

#### VISUAL Mode

```json
"vim.visualModeKeyBindings": [
    {
        "before": ["H"],
        "after": ["^"] 
    },
    {
        "before": ["L"],
        "after": ["$"]
    },

    {
        "before": [">"],
        "commands": [ "editor.action.indentLines"]
    },
    {
        "before": ["<"],
        "commands": [ "editor.action.outdentLines"]
    },
],
```

#### COMMAND LINE Mode

```json
"vim.commandLineModeKeyBindingsNonRecursive": [
],
```

{{< admonition type=note title="`leader` 键注意事项" open=true >}}
`leader` 键只在代码编辑区域生效，它无法做到全 VSCode 生效
{{< /admonition >}}

### 3. 资源管理配置

`keybindings.json` 定义对于**非代码编辑区**的热键.

```json
[
    {
        "key": "cmd+h",
        "command": "workbench.action.focusLeftGroup"
    },
    {
        "key": "cmd+l",
        "command": "workbench.action.focusRightGroup"
    },

    {   // Rename file
        "key": "r",
        "command": "renameFile",
        "when": "explorerViewletVisible && filesExplorerFocus && !explorerResourceIsRoot && !explorerResourceReadonly && !inputFocus"
    },
    {   // New file
        "key": "a",
        "command": "explorer.newFile",
        "when": "explorerViewletVisible && filesExplorerFocus && !explorerResourceIsRoot && !explorerResourceReadonly && !inputFocus"
    },
    {   // New folder
        "key": "shift+a",
        "command": "explorer.newFolder",
        "when": "explorerViewletVisible && filesExplorerFocus && !explorerResourceIsRoot && !explorerResourceReadonly && !inputFocus"
    },
    {   // Delete file
        "key": "d",
        "command": "deleteFile",
        "when": "explorerViewletVisible && filesExplorerFocus && !explorerResourceIsRoot && !explorerResourceReadonly && !inputFocus"
    },
    {   // Copy
        "key": "y",
        "command": "filesExplorer.copy",
        "when": "explorerViewletVisible && filesExplorerFocus && !explorerResourceIsRoot && !explorerResourceReadonly && !inputFocus"
    },
    {   // Cut
        "key": "x",
        "command": "filesExplorer.cut",
        "when": "explorerViewletVisible && filesExplorerFocus && !explorerResourceIsRoot && !explorerResourceReadonly && !inputFocus"
    },
    {   // Paste
        "key": "p",
        "command": "filesExplorer.paste",
        "when": "explorerViewletVisible && filesExplorerFocus && !explorerResourceIsRoot && !explorerResourceReadonly && !inputFocus"
    },
    {   // 全局搜索后，在输入框按回车，即可聚焦到搜索结果列表
        "key": "enter",
        "command": "search.action.focusSearchList",
        "when": "inSearchEditor && inputBoxFocus && hasSearchResult || inputBoxFocus && searchViewletVisible && hasSearchResult"
    },
    {   // 在搜索结果列表，只需按 esc，就可回到搜索输入框
        "key": "escape",
        "command": "workbench.action.findInFiles",
        "when": "searchViewletVisible && hasSearchResult && searchViewletFocus"
    },
    {   // 在搜索输入框，只需按 esc，就可回到编辑器
        "key": "escape",
        "command": "workbench.action.focusFirstEditorGroup",
        "when": "inSearchEditor && inputBoxFocus|| inputBoxFocus && searchViewletVisible"
    },
    {   // 在文件浏览界面，只需按 esc，就可回到编辑器
        "key": "escape",
        "command": "workbench.action.focusFirstEditorGroup",
        "when": "explorerViewletVisible && filesExplorerFocus && !inputFocus"
    },
    {
        "key": "tab",
        "command": "acceptSelectedSuggestion",
        "when": "suggestWidgetVisible && textInputFocus"
    },
    { // Next Suggestion
        "key": "tab",
        "command": "selectNextSuggestion",
        "when": "editorTextFocus && suggestWidgetMultipleSuggestions && suggestWidgetVisible"
    },
    { // Prev Suggestion
        "key": "shift+tab",
        "command": "selectPrevSuggestion",
        "when": "editorTextFocus && suggestWidgetMultipleSuggestions && suggestWidgetVisible"
    },
    {
        "key": "cmd+k",
        "command": "workbench.action.focusActiveEditorGroup",
        "when": "terminalFocus"
    }
]
```

## Gitlens.pro

[gitlens-pro](https://github.com/chengazhen/gitlens-pro)

## Monokai Pro

License

- Email: `id@chinapyg.com`
- Password: `d055c-36b72-151ce-350f4-a8f69`

## Cursor 重置机器码

```bash
curl -sL dub.sh/cursorreset | python3
```

## Reference

- [Visual Studio Code](https://code.visualstudio.com/)
- [VSCodeVim](https://github.com/VSCodeVim/Vim/#key-remapping)
- [gitlens-pro](https://github.com/chengazhen/gitlens-pro)


---

> Author: [kyden](https://github.com/kydance)  
> URL: http://kydance.github.io/posts/vscode-vim/  

