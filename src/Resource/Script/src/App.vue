<script setup>
import HelloWorld from './components/HelloWorld.vue'</script>

<template>
  <cert ref="cert" style="width: 100%;height: 100%" v-if="IsCert"/>
  <div style="width: 100%; height: 100%; display: grid; place-items: center;" v-show="!ConnectionSuccessful"
       v-if="IsCert===false">
    <el-button type="danger" @click="ReconnectWebsocket">{{ ConnectionStatus }}</el-button>
  </div>
  <div ref="appMain" style="width: 100%;height: 100%" v-show="ConnectionSuccessful" v-if="IsCert===false">
    <el-container>
      <el-container>
        <div :style="getAppListStyle">
          <el-aside style="height: 100%; overflow-y: auto;overflow-x: hidden;width: 300px;">
            <div class="center-text" style="    background-color: gold;">SunnyNet内置函数列表</div>
            <div class="center-text2" style="    background-color: chartreuse;">除以下函数外</div>
            <div class="center-text2" style="    background-color: cyan;">您也可以使用Go语言函数</div>
            <el-tree
                style="max-width: 200px"
                :data="FuncList"
                :props="defaultProps"
                @node-click="handleNodeClick"
            />
          </el-aside>
        </div>
        <el-container>
          <el-main :style="getAppCodeStyle">
            <HelloWorld ref="vs" style="width: 100%;height: 100%" @focus="hideFooter"/>
          </el-main>
          <el-footer v-show="showFooterFlag" ref="footer" :style="getAppFooterStyle" v-html="FooterHtml">
          </el-footer>
        </el-container>
      </el-container>
    </el-container>
  </div>
</template>

<script>
import {marked} from 'marked';
import cert from './components/cert.vue'
import {builtCmdWords} from "./components/builtCmdWords.js";

export default {
  data() {
    return {
      FuncList: [],
      defaultProps: {
        children: 'children',
        label: 'label',
      },
      appMainWidth: 0,
      appMainHeight: 0,
      showFooterFlag: false,
      FooterHtml: "",
      ConnectionSuccessful: false,
      ConnectionStatus: "正在连接SunnyNet脚本服务",
      IsCert: false,
    }
  },
  computed: {
    getAppCodeStyle() {
      if (this.showFooterFlag) {
        return "width: 100%;height: " + (this.appMainHeight / 2) + "px"
      }
      return "width: 100%;height: " + (this.appMainHeight) + "px"
    },
    getAppListStyle() {
      return "width: 300px;height: " + (this.appMainHeight) + "px"
    },
    getAppFooterStyle() {
      return "width: 100%;height: " + (this.appMainHeight / 2) + "px; overflow-y: auto;border: 2px solid #b0c2e3;"
    },
  },
  methods: {
    //初始化连接Websocket
    openWebsocket() {
      let wsUrl = (document.location.toString().indexOf("https://") > -1 ? "wss" : "ws") + `://${window.location.host}${window.location.pathname}/WebSocketServer`;
      const socket = new WebSocket(wsUrl);
      socket.onopen = () => {
        this.ConnectionSuccessful = true
      };
      // 监听消息事件
      socket.onmessage = (event) => {
        try {
          const obj = JSON.parse(event.data)
          this.$refs.vs.onWebsocket(obj.cmd, obj.data)
        } catch (e) {
        }
      };
      // 监听连接关闭事件
      socket.onclose = () => {
        this.ConnectionStatus = "连接断开,点击重连"
        this.ConnectionSuccessful = false
      };
      window.getWsSocket = () => {
        return socket;
      }
      window.SendWebsocket = (cmd, data) => {
        this.SendWebsocket(cmd, data)
      }
      window.SendWsMessage = (msg) => {
        this.SendMessage(msg)
      }
    },
    //重新连接websocket
    ReconnectWebsocket() {
      if (this.ConnectionStatus === "正在连接SunnyNet脚本服务") {
        return;
      }
      this.ConnectionStatus = "正在连接SunnyNet脚本服务"
      this.openWebsocket()
    },
    showFooter() {
      this.showFooterFlag = true;
      this.$refs.footer.$el.scrollTop = 0;
    },
    SendWebsocket(cmd, data) {
      if (this.ConnectionSuccessful) {
        window.getWsSocket().send(JSON.stringify({cmd: cmd, data: data}))
      }
    },
    SendMessage(msg) {
      if (this.ConnectionSuccessful) {
        window.getWsSocket().send(msg)
      }
    },
    hideFooter() {
      this.showFooterFlag = false;
    },
    handleNodeClick(Tree) {
      const renderer = new marked.Renderer();
      renderer.code = (r) => {
        // 自定义代码块的 HTML 结构
        return `<pre style="background-color: #5e5d5d;color: #ffffff; border-radius: 5px; padding: 10px;"><code class="language-${r.lang}">${r.text}</code></pre>`;
      };

      let v = "## **" + Tree.label + "**\n\n---\n\n"
      v += "* 函数命令：**支持" + Tree.names.length + "种别名**" + "\n"
      for (let i = 0; i < Tree.names.length; i++) {
        v += "  * " + Tree.names[i] + "\n"
      }
      v += "\n---\n"
      for (let i = 1; i < Tree.contents.length; i++) {
        const val = Tree.contents[i].value + ""
        if (val.startsWith("**") || val.startsWith('```')) {
          v += "\n---\n"
        }
        if (val.endsWith('```')) {
          v += Tree.contents[i].value + "\n"
        } else {
          v += Tree.contents[i].value + '</br>' + "\n"
        }
      }
      //'# Hello, Vue!\n\nThis is my **Markdown** content.'
      this.FooterHtml = marked(v.replaceAll("\n---\n\n---\n", "\n---\n"), {renderer});
      this.showFooter()
    },
    init() {
      for (let index = 0; index < builtCmdWords.length; index++) {
        const obj = {
          label: builtCmdWords[index].detail,
          children: [],
          names: builtCmdWords[index].name,
          contents: builtCmdWords[index].contents
        }
        this.FuncList.push(obj)
      }
    }
  },
  mounted() {
    if (window.location.href.indexOf("install.html") > 0) {
      document.title = "SunnyNet证书安装文档";
      this.IsCert = true
      return
    }
    const elementRef = this.$refs.appMain; // 获取元素的引用
    // 创建 ResizeObserver 实例并监听元素尺寸变化
    const resizeObserver = new ResizeObserver(entries => {
      for (const entry of entries) {
        const {width, height} = entry.contentRect;
        this.appMainWidth = width;
        this.appMainHeight = height;
      }
    });
    resizeObserver.observe(elementRef); // 开始监听元素尺寸变化
    window.vsFocus = this.hideFooter
    window.openWebsocket = this.openWebsocket
    this.init()
    this.$refs.vs.init()
  }, components: {cert},
}
</script>
<style scoped>
.center-text {
  text-align: center;
  padding: 10px 0;
}

.center-text2 {
  text-align: center;
}

.el-aside:focus, .el-main:focus {
  outline: none;
}

.el-main {
  --el-main-padding: 0 !important;
}
</style>