<template>
  <header>
    <section class="left">
      <img src="../assets/images/duck.png" class="favicon" alt="" width="32" height="32"/>
      <h3>my-treasure-doc</h3>
    </section>
    <section class="right">
      <div class="avatar-wrapper">
        <n-avatar class="avatar" round size="medium" src="https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg"/>
        <n-icon size="18">
          <chevron-down />
        </n-icon>
      </div>
    </section>
  </header>
</template>

<script lang="ts">
import { ChevronDown } from '@vicons/ionicons5'

export default {
  name: 'Header',
  components: {ChevronDown}

};
</script>

<style scoped lang='scss'>
@import "src/assets/style/helper.scss";

header {
  width: 100%;
  height: $headerHeight;
  background: $mainColor;
  color: #ffffff;
  display: flex;
  align-items: center;
  padding: 0 32px;
  > .left{
    display: flex;
    width: calc(#{$menuWidth} - 32px);
    > h3{
      font-size: 20px;
      margin-left: 8px;
    }
  }
  > .right{
    display: flex;
    width: calc(100% - #{$menuWidth});
    align-items: center;
    > .avatar-wrapper{
      margin-left: auto;
      display: flex;
      align-items: center;
      > .avatar{
        margin-right: 8px;
      }
    }
  }
}
</style>