<template>
  <div ref="container" class="monaco-editor" style="width: 100%;height: 100%;"></div>
</template>
<script>
import * as monaco from 'monaco-editor'
import {ElMessage, ElMessageBox} from "element-plus";
import {HoverEvent, InputEvent} from "./Event.js"

const tmpCode = `
`


export default {
  data() {
    return {
      // 主要配置
      defaultOpts: {
        language: 'go',
        value: tmpCode, // 编辑器的值
        theme: 'vs', // 编辑器主题：vs, hc-black, or vs-dark，更多选择详见官网
        roundedSelection: true, // 右侧不显示编辑器预览框
        scrollBeyondLastLine: false, // 禁止滚动超过最后一行
        autoIndent: true, // 自动缩进
        automaticLayout: true,
        formatOnType: true,
        formatOnPaste: true,
        originalEditable: true,
        glyphMargin: false,
        diffViewport: false,
        wordWrap: 'on', // 设置自动换行
        validationOptions: {
          validate: false, // 禁用语法错误提示
        },
        minimap: {enabled: false},
        fontSize: 15,
        readOnly: false,
        hover: {
          enabled: true,
          delay: 300
        },
        suggest: {
          showInlineDetails: true,
          showStatusBar: true,
          showIcons: true, // 显示补全项图标
          preview: true,// 开启补全预览
          showDocs: true,  // 强制展示文档部分
          snippetsPreventQuickSuggestions: false,
          filterGraceful: true,
          insertMode: 'insert',
          showWords: true,
          localityBonus: true,
        },
        quickSuggestions: true, // 启用快速建议
        suggestOnTriggerCharacters: true, // 在输入字符时触发补全
        suggestSelection: 'first',  // 自动选择第一个建议
        wordBasedSuggestions: false,  // 禁用基于单词的补全
        parameterHints: {enabled: true},  // 启用参数提示
        tabCompletion: 'on',  // 使用 Tab 补全
        // 配置推荐超时
        suggestTimeout: 5000 // 延长超时

      },
      //内置命令列表
      wordBeforeDot: "",
      getEditor: null,
      getWsSocket: null,  // WebSocket 实例
    }
  },
  mounted() {

  },
  methods: {
    SendWebsocket(cmd, data) {
      window.SendWebsocket(cmd, data)
    },
    GetProvider(func, w1, w2, w3) {
      window.SendWsMessage(JSON.stringify({cmd: "Provider", func: func, w1: w1, w2: w2, w3: w3}))
    },
    //收到Websocket消息
    onWebsocket(cmd, data) {
      switch (cmd) {
        case "SetCodeInit":
          this.getEditor().setValue(data);
          return
        case "SetCode":
          const editor = this.getEditor();
          const currentPosition = editor.getPosition();
          const fullRange = editor.getModel().getFullModelRange();
          editor.executeEdits("my-source", [
            {
              range: fullRange,
              text: data,
              forceMoveMarkers: true
            }
          ]);
          editor.setPosition(currentPosition);
          editor.focus();
          return
        case "Message":
          ElMessage({
            message: data,
            type: 'success',
          })
          return
        case "Error":
          const match = data.match(/第(\d+)行/);
          if (match) {
            const lineNumber = parseInt(match[1]);
            if (lineNumber > 0) {
              const editor = this.getEditor();
              // 获取行的文本内容
              const lineContent = editor.getModel().getLineContent(lineNumber);

              // 获取行的开始和结束位置
              const startPosition = { lineNumber: lineNumber, column: 1 };
              const endPosition = { lineNumber: lineNumber, column: lineContent.length + 1 };

              // 选择这一行
              editor.setSelection(new monaco.Selection(
                  startPosition.lineNumber,
                  startPosition.column,
                  endPosition.lineNumber,
                  endPosition.column
              ));

              // 确保这一行在可视区域中
              editor.revealLineInCenter(lineNumber);
            }
          }
          ElMessageBox.alert(
              data,
              "载入代码错误",
              {
                dangerouslyUseHTMLString: true,
                confirmButtonText: '好的',
                closeOnClickModal: true, // 设置点击遮罩层关闭消息框
                closeOnPressEscape: true, // 设置按下 ESC 键关闭消息框
              }
          )
          return
        case "Provider":
          console.log(data)
          break
      }
    },
    Command(cmd) {
      this.SendWebsocket("LoadDefaultCode", cmd)
    },
    //解析所有单词
    ParsingWords() {
      let text = this.getEditor().getValue();
      text = text.replaceAll('\t', ' ').replaceAll('\n', ' ').replaceAll("\r", "").replaceAll("{", " ").replaceAll("}", " ").replaceAll("[", " ").replaceAll("]", " ").replaceAll(".", " ").replaceAll(".", " ").replaceAll(".", " ").replaceAll(",", " ").replaceAll("\\", " ").replaceAll("\u3000", " ")
      text = text.replaceAll(':', ' ').replaceAll("&", " ").replaceAll("*", " ").replaceAll("(", " ").replaceAll(")", " ").replaceAll("`", " ").replaceAll("/", " ").replaceAll(".", " ").replaceAll("'", " ").replaceAll("\"", " ").replaceAll("=", " ").replaceAll("：", " ")
      text = text.replaceAll('+', ' ').replaceAll("-", " ").replaceAll("*", " ").replaceAll("/", " ").replaceAll("_", " ").replaceAll("<", " ").replaceAll(">", " ").replaceAll("package", " ").replaceAll("main", " ")
      let m = text.length;
      while (true) {
        text = text.replaceAll('  ', ' ')
        if (text.length === m) {
          break
        }
        m = text.length;
      }
      return [...new Set(text.split(" "))]
    },
    //添加执行命令
    AddCommand(editor, title, cmd) {
      editor.addAction({
        id: title,
        label: '* 脚本代码 -> ' + title,
        contextMenuGroupId: '',
        contextMenuOrder: 0,
        run: () => {
          this.Command(cmd)
        }
      });
    },
    //获取当前位置的函数名
    getCurrentFuncName(editor, position) {
      const model = editor.getModel();
      if (position === null || position === undefined) {
        position = editor.getPosition();
      }
      // 获取全部文本
      const text = model.getValue();

      // 将文本分割成行
      const lines = text.split('\n');

      let currentLine = position.lineNumber - 1; // 因为Monaco Editor的行号从1开始
      let inMultiLineComment = false;
      let bracketCount = 0;

      while (currentLine >= 0) {
        const line = lines[currentLine].trim();

        // 处理多行注释
        if (line.endsWith('*/')) {
          inMultiLineComment = true;
        }
        if (line.startsWith('/*')) {
          inMultiLineComment = false;
          currentLine--;
          continue;
        }
        if (inMultiLineComment) {
          currentLine--;
          continue;
        }

        // 忽略单行注释
        if (line.startsWith('//')) {
          currentLine--;
          continue;
        }
        // 找到函数声明
        if (line.startsWith('func ')) {
          const match = line.match(/func\s+(\w+)/);
          if (match && match[1]) {
            return match[1]; // 返回方法名
          }
        }

        currentLine--;
      }

      return null; // 如果没有找到方法名，返回null
    },
    //获取倒数第二个单词
    getSecondLastWord(editor) {
      const model = editor.getModel();
      const position = editor.getPosition();
      // 获取当前行的内容
      const lineContent = model.getLineContent(position.lineNumber);
      // 获取光标前的内容
      const textBeforeCursor = lineContent.substring(0, position.column - 1);
      // 使用正则表达式匹配所有单词
      const words = textBeforeCursor.match(/\b\w+\b/g);
      // 如果找到至少两个单词，返回倒数第二个
      if (words && words.length >= 2) {
        return words[words.length - 2];
      }
      return '';
    },

    //获取倒数第一个单词
    getSecondWord(editor) {
      const model = editor.getModel();
      const position = editor.getPosition();
      // 获取当前行的内容
      const lineContent = model.getLineContent(position.lineNumber);
      // 获取光标前的内容
      const textBeforeCursor = lineContent.substring(0, position.column - 1);
      // 使用正则表达式匹配所有单词
      const words = textBeforeCursor.match(/\b\w+\b/g);
      // 如果找到至少两个单词，返回倒数第二个
      if (words && words.length >= 1) {
        return words[words.length - 1];
      }
      return '';
    },
    //获取鼠标悬停行处及前面的所有单词
    getHoverWords(model, position) {
      const Words = [];
      let co = 0;
      const max = position.column + 1;
      while (co < max) {
        position.column = co;
        const word = model.getWordAtPosition(position);
        if (word) {
          if (!Words.includes(word.word)) {
            Words.push(word.word)
          }
        }
        co++;
      }
      return Words
    }
    ,
    init() {

      this.$refs.container.innerHTML = ''
      monaco.languages.registerCompletionItemProvider('go', {
        provideCompletionItems: async (model, position) => {
          let suggestions = [];
          let wordInfo = "";
          const func = this.getCurrentFuncName(this.getEditor())
          const wos = this.ParsingWords()

          if (this.wordBeforeDot === "") {
            const obj = editor.getModel().getWordAtPosition(position);
            if (obj != null) {
              wordInfo = obj.word;
            }
            if (wordInfo.length < 1) {
              wordInfo = this.getSecondWord(this.getEditor())
              if (wordInfo.length < 1) {
                return {suggestions: []};
              }
              suggestions = await InputEvent(func, wos, [wordInfo], wordInfo, this.getEditor().getValue(), position)
            } else {
              const LastWord = this.getSecondLastWord(this.getEditor())
              suggestions = await InputEvent(func, wos, [LastWord], wordInfo, this.getEditor().getValue(), position)
            }
          } else {
            const LastWord = this.getSecondLastWord(this.getEditor())
            const mm = []
            if (LastWord !== "") {
              mm.push(LastWord)
            }
            mm.push(this.wordBeforeDot)
            suggestions = await InputEvent(func, wos, mm, ".", this.getEditor().getValue(), position)
          }
          let reversedArr = [...suggestions].reverse();
          return {suggestions: reversedArr};
        },
        triggerCharacters: ['.']
      });
      monaco.languages.registerHoverProvider('go', {
        provideHover: (model, position) => {
          const words = this.getHoverWords(model, position);
          const length = words.length;
          if (length < 1) {
            return null
          }
          const word = words[length - 1];
          const func = this.getCurrentFuncName(this.getEditor(), position)
          return HoverEvent(func, words, word);
        }
      });
      const editor = monaco.editor.create(this.$refs.container, this.defaultOpts);
      Array.prototype.add = function (label, kind, detail, documentation) {
        const completionItem = {
          label: label,
          kind: kind,
          insertText: label,
          insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
          detail: detail,
          documentation: {
            value: documentation,
            isTrusted: true
          }
        };
        this.push(completionItem);
      };
      Array.prototype.addFunc = function (label, insertText, kind, detail, documentation) {
        const completionItem = {
          label: label,
          kind: kind,
          insertText: insertText,
          insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
          detail: detail,
          documentation: {
            value: documentation,
            isTrusted: true
          }
        };
        this.push(completionItem);
      };
      const key = editor.createContextKey('wordWrapOn', true)
      editor.onDidChangeModelContent((event) => {
        {
          // 获取最近插入的文本
          const changes = event.changes;
          const lastChange = changes[changes.length - 1];
          if (lastChange && lastChange.text.length > 0) {
            const selectedText = lastChange.text.trim();
            // 检查是否是你定义的补全项之一
            if (selectedText === 'Conn.') {
              this.$nextTick(() => {
                editor.trigger('manual', 'editor.action.triggerSuggest', {});
              })
              return;
            }
          }
          if (lastChange && lastChange.text.length === 0) {
            const position = editor.getPosition(); // 获取光标当前位置
            if (position.column > 1) {
              // 创建一个范围，获取光标前的一个字符
              const range = new monaco.Range(position.lineNumber, position.column - 1, position.lineNumber, position.column);
              const previousChar = editor.getModel().getValueInRange(range); // 获取该范围内的字符
              if (/^[a-zA-Z]$/.test(previousChar)) {
                this.$nextTick(() => {
                  editor.trigger('manual', 'editor.action.triggerSuggest', {});
                })
              }
            }

            return;
          }
        }
        // 检查是否输入了 "."
        const changes = event.changes;
        const lastChange = changes[changes.length - 1];
        this.wordBeforeDot = "";
        if (lastChange.text === '.') {
          const model = editor.getModel();
          const position = editor.getPosition();

          // 获取当前行的内容
          const lineContent = model.getLineContent(position.lineNumber);

          // 获取光标前的内容
          const textBeforeCursor = lineContent.substring(0, position.column - 1);

          // 使用正则表达式匹配最后一个单词
          const match = textBeforeCursor.match(/(\w+)\.?$/);

          if (!(match && match[1])) {
            return
          }
          this.wordBeforeDot = match[1];
        }
      });
      editor.addAction({
        id: 'turnWordWrapOff',
        label: '关闭自动换行',
        contextMenuGroupId: 'my-commands',
        contextMenuOrder: Number.MAX_SAFE_INTEGER,
        precondition: 'wordWrapOn',
        run: () => {
          this.defaultOpts.wordWrap = 'off'
          editor.updateOptions({
            wordWrap: "off"
          });
          key.set(false)
        },
      })
      editor.addAction({
        id: 'turnWordWrapOn',
        label: '自动换行',
        contextMenuGroupId: 'my-commands',
        contextMenuOrder: Number.MAX_SAFE_INTEGER,
        precondition: '!wordWrapOn',
        run: () => {
          this.defaultOpts.wordWrap = 'on'
          editor.updateOptions({
            wordWrap: "on"
          });
          key.set(true)
        },
      })
      editor.addAction({
        id: 'qhtheme',
        label: '切换主题',
        contextMenuGroupId: '1_modification',
        contextMenuOrder: Number.MAX_SAFE_INTEGER,
        run: () => {
          if (this.defaultOpts.theme === "vs") {
            this.defaultOpts.theme = "vs-dark";
          } else {
            this.defaultOpts.theme = "vs";
          }
          monaco.editor.setTheme(this.defaultOpts.theme);
        },
      })
      editor.addAction({
        id: 'CodeLoadSave',
        label: '加载并且保存代码',
        contextMenuGroupId: 'navigation',
        keybindings: [monaco.KeyMod.chord(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS)],
        contextMenuOrder: Number.MAX_SAFE_INTEGER,
        run: () => {
          this.SendWebsocket("CodeLoadSave", editor.getValue());
        },
      })
      this.AddCommand(editor, "恢复到默认代码", "DefaultCode")
      this.AddCommand(editor, "加载 拦截修改HTTP/S模板", "httpDefaultCode")
      this.AddCommand(editor, "加载 拦截修改 TCP 模板", "tcpDefaultCode")
      this.AddCommand(editor, "加载 拦截修改 UDP 模板", "udpDefaultCode")
      this.AddCommand(editor, "加载 拦截修改 Websocket 模板", "WebsocketDefaultCode")
      window.onresize = function () {
        if (editor) {
          editor.layout();
        }
      };
      editor.onDidFocusEditorText(() => {
        if (window.vsFocus !== undefined && window.vsFocus != null) {
          window.vsFocus()
        }
      })

      this.getEditor = () => {
        return editor
      }
      window.openWebsocket()
    },
  },
  computed: {}
}
</script>

