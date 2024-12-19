<template xmlns="http://www.w3.org/1999/html">
  <el-tabs v-model="activeName" class="demo-tabs"
           style="position: relative;left: 10px;display: inline-grid;justify-content: center;align-items: center;">
    <el-tab-pane label="Windows" name="Windows"/>
    <el-tab-pane label="MacOs" name="MacOs"/>
    <el-tab-pane label="Ios" name="Ios"/>
    <el-tab-pane label="Android" name="Android"/>
    <el-tab-pane label="注意事项" name="注意事项"/>
    <el-tab-pane label="源码" name="源码"/>
    <el-tab-pane label="开源协议" name="开源协议"/>
    <el-tab-pane label="打赏作者" name="打赏作者"/>
  </el-tabs>
  <div style="position: relative;height: calc(100% - 60px);width: 100%;overflow-y: auto; overflow-x: hidden">
    <el-container style="height: 100%;width: 100%">
      <el-container style="height: 100%;width: 100%">
        <el-main style="height: 100%">
          <div v-if="activeName==='Windows'">
            <Windows/>
          </div>
          <div v-if="activeName==='MacOs'">
            <MacOs/>
          </div>
          <div v-if="activeName==='Ios'">
            <IOS/>
          </div>
          <div v-if="activeName==='Android'" style="height: 100%;width: 100%">
            <Android1/>
          </div>
          <div v-if="activeName==='注意事项'">
            <attention/>
          </div>
          <div v-if="activeName==='源码'">
            <SDK/>
          </div>
          <div v-if="activeName==='开源协议'">
            <OpenSource/>
          </div>
          <div v-if="activeName==='打赏作者'">
            <img :src="ImagePay" style="width: 100%;max-width: 1200px;display: block;margin: auto;">
          </div>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>

import Windows from "./CertDoc/Windows.vue";
import Attention from "./CertDoc/attention.vue";
import MacOs from "./CertDoc/MacOs.vue";
import IOS from "./CertDoc/IOS.vue";
import SDK from "./CertDoc/SDK.vue";
import OpenSource from "./CertDoc/OpenSource.vue";
import Android1 from "./CertDoc/Android.vue";
import {GetDocImage} from "./image.js";

export default {
  components: {
    Android1, IOS, MacOs, Attention, Windows, SDK, OpenSource,
  },
  data() {
    return {
      activeName: "开源协议",
    }
  },
  mounted() {
  },
  beforeUnmount() {
  },
  computed: {
    ImagePay() {
      return GetDocImage("Imagepay")
    },
  },
  methods: {},

}
</script>