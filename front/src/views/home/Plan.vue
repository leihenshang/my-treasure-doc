<template>
  <n-space vertical size="large">
    <!-- overall-start -->
    <n-layout has-sider>

      <!-- left-start -->
      <n-layout-sider content-style="padding: 24px;">
        <div>
          2002 <br>
          1月 <br>
          1月 <br>
          1月 <br>
          1月 <br>
          1月 <br>
          1月 <br>
          1月 <br>
          1月 <br>
          1月 <br>
        </div>

      </n-layout-sider>
      <!-- left-end -->

      <!-- right-start -->
      <n-layout>
        <n-layout-header>
          <p>计划</p>
        </n-layout-header>

        <n-layout-content content-style="padding: 24px;">
          2022年7月
          <n-radio-group v-model:value="size" name="left-size" style="margin-bottom: 12px">
            <n-radio-button value="small">
              1日
            </n-radio-button>
            <n-radio-button value="medium">
              2日
            </n-radio-button>
            <n-radio-button value="large">
              3日
            </n-radio-button>
            <n-radio-button value="large">
              ...日
            </n-radio-button>
          </n-radio-group>
        </n-layout-content>

        <n-layout-content content-style="padding: 24px;">
          <n-input v-model:value="value" type="text" placeholder="基本的 Input" />
          <n-date-picker v-model:value="range" type="daterange" clearable />
          <n-button type="primary">
            提交
          </n-button>

        </n-layout-content>

        <n-layout-footer>
          <p>今天要做的事情</p>
          <table>
            <th>
            <td>事项</td>
            <td>进度</td>
            </th>
            <tr>
              <td>sxxxxxxxxxxxx</td>
              <td>10%</td>
            </tr>
            <tr>
              <td>sxxxxxxxxxxxx</td>
              <td>10%</td>
            </tr>
            <tr>
              <td>sxxxxxxxxxxxx</td>
              <td>10%</td>
            </tr>
            <tr>
              <td>sxxxxxxxxxxxx</td>
              <td>10%</td>
            </tr>
          </table>
        </n-layout-footer>


      </n-layout>
      <!-- right-end -->

    </n-layout>
    <!-- overall-end -->
  </n-space>
</template>

<script lang="ts">
export default {
  name: "Plan"
}
</script>

<style scoped lang='scss'>
.n-layout-header,
.n-layout-footer {
  background: rgba(128, 128, 128, 0.2);
  padding: 24px;
}

.n-layout-sider {
  background: rgba(128, 128, 128, 0.3);
}

.n-layout-content {
  background: rgba(128, 128, 128, 0.4);
}
</style>