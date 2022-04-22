<template>
  <n-popselect :options="[]" trigger="click">
    <div>
      <SvgIcon icon-name="collected"></SvgIcon>
    </div>
    <template #empty>
      <div class="content-wrapper">
        <h3>选择分组</h3>
        <p><span>你可以选择分组或者直接</span> <n-button type="info" text>取消收藏</n-button></p>
        <n-checkbox-group v-model:value="groups" class="checkbox-group">
          <n-space item-style="display: flex;">
            <n-checkbox value="1" label="全部" />
            <n-checkbox value="2" label="javascript" />
            <n-checkbox value="Guangzhou" label="好看的花花们" />
          </n-space>
        </n-checkbox-group>
      </div>
    </template>
    <template #action>
      <p class="action">
        <SvgIcon icon-name="add"></SvgIcon>
        <span @click="addGroup">新建分组</span>
      </p>
    </template>
  </n-popselect>
</template>

<script>
import SvgIcon from '../public/SvgIcon.vue'
import {ref} from 'vue'

export default {
  name: 'CollectionHandle',
  components:{SvgIcon},
  setup(props,context){
    const groups = ref([])
    const addGroup = ()=>{
      context.emit('update:addGroup')
    }
    return{
      groups,
      addGroup
    }
  }
}
</script>

<style scoped lang='scss'>
$grey-color: #8A8F8D;
$grey-background: #eff0f0;
.content-wrapper{
  color: #262626;
  >h3{
    font-size: 16px;
    padding: 8px 16px 0;
  }
  >p{
    padding: 8px 16px 0;
    >span{
      color:$grey-color;
    }
  }
  .checkbox-group{
    padding: 16px;
    .n-space{
      display: flex;
      flex-direction: column;
    }
  }
}
.action{
  height: 36px;
  line-height: 36px;
  .icon{
    margin-right: 8px;
  }
}
</style>