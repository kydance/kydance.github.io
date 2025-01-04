# 【最佳实践】VSCode &#43; Vim = 效率之神


{{&lt; admonition type=abstract title=&#34;导语&#34; open=true &gt;}}
在 Visual Studio Code 中引入 Vim 模式，无疑可以极大程度上提高个人的编码效率。
{{&lt; /admonition &gt;}}

&lt;!--more--&gt;

## I. VSCodeVim

&lt;!-- {{&lt; figure src=&#34;/posts/vscode-vim/vim.png&#34; title=&#34;&#34; height=0.5 width=0.5 &gt;}} --&gt;

[VSCodeVim](https://github.com/VSCodeVim/Vim/#key-remapping) 是一款 vim 模拟器，它将 vim 的大部分功能都集成在了 VSCode 中，即一个嵌入在 VSCode 中的 vim。

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

{{&lt; figure src=&#34;/posts/vscode-vim/vscode-extension-vim.png&#34; title=&#34;&#34; height=&#34;256&#34; width=&#34;256&#34; &gt;}}

{{&lt; admonition type=tip title=&#34;关闭 Mac 的重复键&#34; open=true &gt;}}

当使用 Mac 时，需要输入以下代码，用以关闭 Mac 的重复键

```bash
# For VS Code
$ defaults write com.microsoft.VSCode ApplePressAndHoldEnabled -bool false
# For VS Codium
$ defaults write com.vscodium ApplePressAndHoldEnabled -bool false
# If necessary, reset global default
$ defaults delete -g ApplePressAndHoldEnabled
```

{{&lt; /admonition &gt;}}

---

## III. 文件配置详解

VSCodeVim 的相关配置文件是放在 `settings.json` 中，而不是 `vimrc` 文件.

- 对于**非代码编辑区**的热键将其定义在 `keybindings.json` 中
- 对于**代码编辑区**且属于 vim 的热键将其定义在 `settings.json` 文件中

{{&lt; admonition type=tip title=&#34;个人配置参考&#34; open=true &gt;}}
[个人 vscode 配置文件参考：https://github.com/kydance/dotfiles](https://github.com/kydance/dotfiles)
{{&lt; /admonition &gt;}}

### 1. 基础配置

```json
// leader - prefix key
&#34;vim.leader&#34;: &#34;&lt;space&gt;&#34;,
// To improve performance
&#34;extensions.experimental.affinity&#34;: { 
    &#34;vscodevim.vim&#34;: 1 
},
// Easy motion
&#34;vim.easymotion&#34;: true,
// Use system clipboard
&#34;vim.useSystemClipboard&#34;: true,
// 由vim接管ctrl&#43;any的按键，而不是vscode
&#34;vim.useCtrlKeys&#34;: true,
&#34;vim.replaceWithRegister&#34;: true,
// 忽略大小写
&#34;vim.ignorecase&#34;: true,
&#34;vim.smartcase&#34;: true,
// 智能行号
&#34;vim.smartRelativeLine&#34;: true,
&#34;vim.foldfix&#34;: true,
// Highlight search
&#34;vim.hlsearch&#34;: true,
&#34;vim.highlightedyank.enable&#34;: true,
&#34;vim.highlightedyank.duration&#34;: 500,

// 由vscode进行处理，而不是vscode-vim插件
&#34;vim.handleKeys&#34;: {
    &#34;&lt;C-d&gt;&#34;: true, // 向下滚动半页
    &#34;&lt;C-f&gt;&#34;: true, // 向下滚动一页
    &#34;&lt;C-e&gt;&#34;: true, // 向下滚动一行
    &#34;&lt;C-s&gt;&#34;: true,
    &#34;&lt;C-z&gt;&#34;: false,
    &#34;&lt;C-a&gt;&#34;: true,
    &#34;&lt;C-c&gt;&#34;: true,
    &#34;&lt;C-v&gt;&#34;: true,
    &#34;&lt;C-x&gt;&#34;: true,
},
```

---

### 2. 快捷键配置

#### NORMAL Mode

```json
&#34;vim.normalModeKeyBindingsNonRecursive&#34;: [
    {   // 聚集 terminal
        &#34;before&#34;: [&#34;&lt;C-j&gt;&#34;], 
        &#34;commands&#34;: [&#34;workbench.action.terminal.focus&#34;]
    },
    {   // 语义级 重命名
        &#34;before&#34;: [&#34;leader&#34;, &#34;r&#34;],  
        &#34;commands&#34;: [ &#34;editor.action.rename&#34;] 
    },
    {   // 
        &#34;before&#34;: [&#34;g&#34;, &#34;h&#34;],  
        &#34;commands&#34;: [ &#34;editor.action.showHover&#34;]
    },
    {
        &#34;before&#34;: [&#34;g&#34;, &#34;d&#34;],
        &#34;commands&#34;: [&#34;editor.action.revealDefinition&#34;]
    },
    {
        &#34;before&#34;: [&#34;g&#34;, &#34;r&#34;],
        &#34;commands&#34;: [&#34;editor.action.goToReferences&#34;],
    },
    {
        &#34;before&#34;: [&#34;g&#34;, &#34;i&#34;],
        &#34;commands&#34;: [&#34;editor.action.goToImplementation&#34;],
    },
    {
        &#34;before&#34;: [&#34;g&#34;, &#34;b&#34;],
        &#34;commands&#34;: [&#34;workbench.action.navigateBack&#34;],
    },

    {
        &#34;before&#34;: [&#34;leader&#34;, &#34;e&#34;, &#34;f&#34;],
        &#34;commands&#34;: [&#34;workbench.explorer.fileView.focus&#34;],
    },

    {
        &#34;before&#34;: [ &#34;leader&#34;, &#34;leader&#34;, &#34;e&#34;, &#34;f&#34;],
        &#34;commands&#34;: [&#34;workbench.action.toggleActivityBarVisibility&#34;]
    },

    {
        &#34;before&#34;: [&#34;leader&#34;, &#34;g&#34;, &#34;g&#34;],
        &#34;commands&#34;: [&#34;workbench.action.quickOpen&#34;],
    },
    {   // Global find
        &#34;before&#34;: [&#34;leader&#34;, &#34;g&#34;, &#34;f&#34;],
        &#34;commands&#34;: [&#34;workbench.view.search&#34;],
    },

    {
        &#34;before&#34;: [&#34;g&#34;, &#34;[&#34;,], 
        &#34;commands&#34;: [&#34;editor.action.marker.prevInFiles&#34;],
    },
    {
        &#34;before&#34;: [&#34;g&#34;, &#34;]&#34;,],
        &#34;commands&#34;: [&#34;editor.action.marker.nextInFiles&#34;],
    },
    {   // Source Control Git
        &#34;before&#34;: [ &#34;leader&#34;, &#34;g&#34;, &#34;i&#34;, &#34;t&#34; ],
        &#34;commands&#34;: [&#34;workbench.scm.focus&#34;]
    },
    { // Start to debug
        &#34;before&#34;: [ &#34;leader&#34;, &#34;d&#34; ],
        &#34;commands&#34;: [&#34;workbench.action.debug.start&#34;]
    },
    {
        &#34;before&#34;: [&#34;leader&#34;,&#34;w&#34;],
        &#34;commands&#34;: [&#34;:w!&#34; ]
    },
    {
        &#34;before&#34;: [&#34;leader&#34;,&#34;q&#34;],
        &#34;commands&#34;: [&#34;:q&#34; ]
    },
    {   // No highlight
        &#34;before&#34;: [&#34;leader&#34;, &#34;n&#34;, &#34;h&#34;],
        &#34;commands&#34;: [&#34;:nohl&#34;]
    },
    {
        &#34;before&#34;: [&#34;H&#34;],
        &#34;after&#34;: [&#34;^&#34;]
    },
    {
        &#34;before&#34;: [&#34;L&#34;], 
        &#34;after&#34;: [&#34;$&#34;]
    },
    {   // Blockwise visual mode
        &#34;before&#34;: [&#34;\\&#34;],
        &#34;commands&#34;: [&#34;extension.vim_ctrl&#43;v&#34;]
    },
    {
        &#34;before&#34;: [&#34;leader&#34;, &#34;t&#34;],
        &#34;commands&#34;: [&#34;:terminal&#34;] 
    }, 

    {
        &#34;before&#34;: [&#34;g&#34;, &#34;t&#34;],
        &#34;commands&#34;: [&#34;:tabnext&#34;]
    }, 
    {
        &#34;before&#34;: [&#34;g&#34;, &#34;T&#34;],
        &#34;commands&#34;: [&#34;:tabprev&#34;] 
    },
    {   // project-manager
        &#34;before&#34;: [&#34;leader&#34;, &#34;p&#34;, &#34;m&#34;], 
        &#34;commands&#34;: [{
            &#34;command&#34;:&#34;workbench.view.extension.project-manager&#34;,
            &#34;when&#34;:&#34;viewContainer.workbench.view.extension.project-manager.enabled&#34;
        }]
    },
],
```

#### INSERT Mode

```json
&#34;vim.insertModeKeyBindings&#34;: [
    {
        &#34;before&#34;: [&#34;j&#34;, &#34;k&#34;],
        &#34;after&#34;: [&#34;&lt;Esc&gt;&#34;]
    }, 
],
```

#### VISUAL Mode

```json
&#34;vim.visualModeKeyBindings&#34;: [
    {
        &#34;before&#34;: [&#34;H&#34;],
        &#34;after&#34;: [&#34;^&#34;] 
    },
    {
        &#34;before&#34;: [&#34;L&#34;],
        &#34;after&#34;: [&#34;$&#34;]
    },

    {
        &#34;before&#34;: [&#34;&gt;&#34;],
        &#34;commands&#34;: [ &#34;editor.action.indentLines&#34;]
    },
    {
        &#34;before&#34;: [&#34;&lt;&#34;],
        &#34;commands&#34;: [ &#34;editor.action.outdentLines&#34;]
    },
],
```

#### COMMAND LINE Mode

```json
&#34;vim.commandLineModeKeyBindingsNonRecursive&#34;: [
],
```

{{&lt; admonition type=note title=&#34;`leader` 键注意事项&#34; open=true &gt;}}
`leader` 键只在代码编辑区域生效，它无法做到全 VSCode 生效
{{&lt; /admonition &gt;}}

### 3. 资源管理配置

`keybindings.json` 定义对于**非代码编辑区**的热键.

```json
[
    {
        &#34;key&#34;: &#34;cmd&#43;h&#34;,
        &#34;command&#34;: &#34;workbench.action.focusLeftGroup&#34;
    },
    {
        &#34;key&#34;: &#34;cmd&#43;l&#34;,
        &#34;command&#34;: &#34;workbench.action.focusRightGroup&#34;
    },

    {   // Rename file
        &#34;key&#34;: &#34;r&#34;,
        &#34;command&#34;: &#34;renameFile&#34;,
        &#34;when&#34;: &#34;explorerViewletVisible &amp;&amp; filesExplorerFocus &amp;&amp; !explorerResourceIsRoot &amp;&amp; !explorerResourceReadonly &amp;&amp; !inputFocus&#34;
    },
    {   // New file
        &#34;key&#34;: &#34;a&#34;,
        &#34;command&#34;: &#34;explorer.newFile&#34;,
        &#34;when&#34;: &#34;explorerViewletVisible &amp;&amp; filesExplorerFocus &amp;&amp; !explorerResourceIsRoot &amp;&amp; !explorerResourceReadonly &amp;&amp; !inputFocus&#34;
    },
    {   // New folder
        &#34;key&#34;: &#34;shift&#43;a&#34;,
        &#34;command&#34;: &#34;explorer.newFolder&#34;,
        &#34;when&#34;: &#34;explorerViewletVisible &amp;&amp; filesExplorerFocus &amp;&amp; !explorerResourceIsRoot &amp;&amp; !explorerResourceReadonly &amp;&amp; !inputFocus&#34;
    },
    {   // Delete file
        &#34;key&#34;: &#34;d&#34;,
        &#34;command&#34;: &#34;deleteFile&#34;,
        &#34;when&#34;: &#34;explorerViewletVisible &amp;&amp; filesExplorerFocus &amp;&amp; !explorerResourceIsRoot &amp;&amp; !explorerResourceReadonly &amp;&amp; !inputFocus&#34;
    },
    {   // Copy
        &#34;key&#34;: &#34;y&#34;,
        &#34;command&#34;: &#34;filesExplorer.copy&#34;,
        &#34;when&#34;: &#34;explorerViewletVisible &amp;&amp; filesExplorerFocus &amp;&amp; !explorerResourceIsRoot &amp;&amp; !explorerResourceReadonly &amp;&amp; !inputFocus&#34;
    },
    {   // Cut
        &#34;key&#34;: &#34;x&#34;,
        &#34;command&#34;: &#34;filesExplorer.cut&#34;,
        &#34;when&#34;: &#34;explorerViewletVisible &amp;&amp; filesExplorerFocus &amp;&amp; !explorerResourceIsRoot &amp;&amp; !explorerResourceReadonly &amp;&amp; !inputFocus&#34;
    },
    {   // Paste
        &#34;key&#34;: &#34;p&#34;,
        &#34;command&#34;: &#34;filesExplorer.paste&#34;,
        &#34;when&#34;: &#34;explorerViewletVisible &amp;&amp; filesExplorerFocus &amp;&amp; !explorerResourceIsRoot &amp;&amp; !explorerResourceReadonly &amp;&amp; !inputFocus&#34;
    },
    {   // 全局搜索后，在输入框按回车，即可聚焦到搜索结果列表
        &#34;key&#34;: &#34;enter&#34;,
        &#34;command&#34;: &#34;search.action.focusSearchList&#34;,
        &#34;when&#34;: &#34;inSearchEditor &amp;&amp; inputBoxFocus &amp;&amp; hasSearchResult || inputBoxFocus &amp;&amp; searchViewletVisible &amp;&amp; hasSearchResult&#34;
    },
    {   // 在搜索结果列表，只需按 esc，就可回到搜索输入框
        &#34;key&#34;: &#34;escape&#34;,
        &#34;command&#34;: &#34;workbench.action.findInFiles&#34;,
        &#34;when&#34;: &#34;searchViewletVisible &amp;&amp; hasSearchResult &amp;&amp; searchViewletFocus&#34;
    },
    {   // 在搜索输入框，只需按 esc，就可回到编辑器
        &#34;key&#34;: &#34;escape&#34;,
        &#34;command&#34;: &#34;workbench.action.focusFirstEditorGroup&#34;,
        &#34;when&#34;: &#34;inSearchEditor &amp;&amp; inputBoxFocus|| inputBoxFocus &amp;&amp; searchViewletVisible&#34;
    },
    {   // 在文件浏览界面，只需按 esc，就可回到编辑器
        &#34;key&#34;: &#34;escape&#34;,
        &#34;command&#34;: &#34;workbench.action.focusFirstEditorGroup&#34;,
        &#34;when&#34;: &#34;explorerViewletVisible &amp;&amp; filesExplorerFocus &amp;&amp; !inputFocus&#34;
    },
    {
        &#34;key&#34;: &#34;tab&#34;,
        &#34;command&#34;: &#34;acceptSelectedSuggestion&#34;,
        &#34;when&#34;: &#34;suggestWidgetVisible &amp;&amp; textInputFocus&#34;
    },
    { // Next Suggestion
        &#34;key&#34;: &#34;tab&#34;,
        &#34;command&#34;: &#34;selectNextSuggestion&#34;,
        &#34;when&#34;: &#34;editorTextFocus &amp;&amp; suggestWidgetMultipleSuggestions &amp;&amp; suggestWidgetVisible&#34;
    },
    { // Prev Suggestion
        &#34;key&#34;: &#34;shift&#43;tab&#34;,
        &#34;command&#34;: &#34;selectPrevSuggestion&#34;,
        &#34;when&#34;: &#34;editorTextFocus &amp;&amp; suggestWidgetMultipleSuggestions &amp;&amp; suggestWidgetVisible&#34;
    },
    {
        &#34;key&#34;: &#34;cmd&#43;k&#34;,
        &#34;command&#34;: &#34;workbench.action.focusActiveEditorGroup&#34;,
        &#34;when&#34;: &#34;terminalFocus&#34;
    }
]
```

## Reference

- [Visual Studio Code](https://code.visualstudio.com/)
- [VSCodeVim](https://github.com/VSCodeVim/Vim/#key-remapping)


---

> : [kyden](https:github.com/kydance)  
> URL: http://kyden.us.kg/posts/vscode-vim/  

