<template>
  <div>
    当前软件源码
    <br><br>
    <el-tag class="ml-2" type="success">https://github.com/qtgolang/SunnyNet</el-tag>
    <br><br>
    获取SDK及源码、文档,请访问 以下任意地址
    <br><br>
    <el-tag class="ml-2" type="danger">https://esunny.vip</el-tag>
    <br>
    <el-tag class="ml-2" type="danger">https://www.esunny.vip</el-tag>
    <br>
    <el-tag class="ml-2" type="danger">https://github.esunny.vip</el-tag>
    <br>
    <el-tag class="ml-2" type="danger">https://gitee.com/qtr/SunnyNet</el-tag>
    <br>
    <el-tag class="ml-2" type="danger">https://github.com/qtgolang/SunnyNet</el-tag>
    <br> <br>
    <el-text class="mx-1" type="success">QQ交流群 : </el-text>
    <el-tag class="ml-2" type="danger">751406884</el-tag>
    <el-tag class="ml-2" type="danger">545120699</el-tag>
    <el-tag class="ml-2" type="danger">170902713</el-tag>
    <br>
    <br>
    <el-text class="mx-1" type="success">QQ频道:</el-text>
    <el-tag class="ml-2" type="danger">https://pd.qq.com/g/SunnyNetV5</el-tag>
  </div>
</template>

<script>

export default {
  data() {
    return {
      Port: ""
    }
  }, mounted() {

  }, methods: {
    callUrl(a) {
      window.open(a, '_blank');
    }
  }
}
</script>