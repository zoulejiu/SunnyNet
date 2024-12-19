<template>
  <div>
    <el-tag class="ml-2" type="danger">请尽量不要在WIFI 设置代理</el-tag>
    <br>
    <el-tag class="ml-2" type="danger">请尽量不要在WIFI 设置代理</el-tag>
    <br>
    <el-tag class="ml-2" type="danger">请尽量不要在WIFI 设置代理</el-tag>
    <br><br><br>
  </div>
  <div>
    <el-text class="mx-1" type="success">安卓推荐使用：</el-text>&nbsp;
    <el-tooltip
        class="box-item"
        effect="dark"
        content="点击去下载"
        placement="top"
    >
      <el-text class="mx-1" type="danger" @click="callUrl('https://wwxa.lanzouj.com/b02p44fnje','c1pe')"
               style="cursor: pointer; text-decoration: underline;">Kitsuneb APP
      </el-text>
    </el-tooltip>
    &nbsp;&nbsp;&nbsp;
    <el-text class="mx-1" type="success">设置Socks代理</el-text>
    <el-text class="mx-1" type="danger">注意：需设置全局代理</el-text>
    <br><br>
    <el-text class="mx-1" type="success">Ios推荐使用：</el-text>
    <el-tooltip
        class="box-item"
        effect="dark"
        content="自行想办法下载"
        placement="top"
    >
      <el-text class="mx-1" type="danger">Shadowrocket</el-text>&nbsp;
    </el-tooltip>&nbsp;&nbsp;、
    <el-tooltip
        class="box-item"
        effect="dark"
        content="自行想办法下载"
        placement="top"
    >
      <el-text class="mx-1" type="danger">Quantumult</el-text>&nbsp;
    </el-tooltip>&nbsp;&nbsp;
    <el-text class="mx-1" type="success">设置Socks5代理</el-text>
    <el-text class="mx-1" type="danger">注意：需设置全局代理</el-text>
    <br><br>
    <el-text class="mx-1" type="danger">在IOS中,若是Socks5代理,工作不正常,可尝试选择HTTP/HTTPS代理类型</el-text>
  </div>

  <div>
    <br>
    <el-text class="mx-1" type="danger">除此之外还需要注意电脑/手机、模拟器的系统时间是否正确</el-text>
    <br>
    <el-text class="mx-1" type="danger">除此之外还需要注意电脑/手机、模拟器的系统时间是否正确</el-text>
    <br>
    <el-text class="mx-1" type="danger">除此之外还需要注意电脑/手机、模拟器的系统时间是否正确</el-text>
  </div>
  <el-dialog
      v-model="dialogVisible"
      title="是否前去下载?"
      width="500"
  >
    <span>请记录下载密码:</span>
    <span
        style="background-color: #b69393;color: #216741;width: 100px;position: relative;display: block;text-align: center"> {{
        Pass
      }} </span>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="callUrl('do',this.Pass)">
          去下载
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script>


export default {
  data() {
    return {
      Port: "",
      Pass: "",
      URL: "",
      dialogVisible: false,
    }
  }, mounted() {

  }, methods: {
    callUrl(url, pass) {
      if (url === "do") {
        window.open(this.URL, '_blank');
        return
      }
      this.Pass = pass;
      this.URL = url;
      this.dialogVisible = true
    }
  }
}
</script>