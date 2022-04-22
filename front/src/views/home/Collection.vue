<template>
  <div class="collect-wrapper">
    <div class="left">
      <div class="title">
        <span>收藏列表</span>
        <SvgIcon icon-name="add" @click="changeGroup('add')"></SvgIcon>
      </div>
      <ul class="list">
        <li v-for="item in groupList" :key="item.id"
            :class="{'selectedCollection':selectedCollectionId === item.id}"
            @click="viewCollectionList(item.id)">
          <p class="title">
            <span>{{ item.title }}</span>
            <n-popselect :options="[]" trigger="hover">
              <div>
                <SvgIcon icon-name="more" v-if="item.type!=='all'"></SvgIcon>
              </div>
              <template #empty>
                <ul class="handle-pop-select-options">
                  <li @click="changeGroup('update',item)">编辑</li>
                  <li @click="changeGroup('delete',item)">删除</li>
                </ul>
              </template>
            </n-popselect>
          </p>
          <p class="content">{{ item.count }}条内容</p>
        </li>
      </ul>
    </div>
    <div class="right">
      <div class="title">
        <h3>{{ selectedGroupName }}</h3>
        <div class="search">
          <n-input v-model:value="searchData.data1" type="text" size="small" placeholder="基本的 Input">
            <template #prefix>
              <SvgIcon icon-name="search"></SvgIcon>
            </template>
          </n-input>
          <SvgIcon icon-name="filter"></SvgIcon>
        </div>
      </div>
      <n-data-table :columns="columns" :data="data" :pagination="pagination" :bordered="false"/>
    </div>
  </div>
  <n-modal v-model:show="showModal" preset="dialog" title="Dialog" :show-icon="false" class="modal-dialog"
           :mask-closable=false style="position: fixed; left: 50%;transform: translateX(-50%);top: 100px">
    <template #header>
      <div>{{ groupHandleType === 'add'?'新增分组':'编辑分组' }}</div>
    </template>
    <div class="dialog-content">
      <label>分组名称</label>
      <n-input v-model:value="newGroupName" type="text" placeholder="分组名称"></n-input>
    </div>
    <template #action>
      <n-button type="primary" @click="updateGroupName">确定</n-button>
      <n-button @click="showModal = false">取消</n-button>
    </template>
  </n-modal>
</template>

<script lang="ts">
import SvgIcon from '../../components/public/SvgIcon.vue';
import {h, defineComponent, computed, ref,Component} from 'vue';
import {NButton, useMessage, DataTableColumns} from 'naive-ui';

type Song = {
  name: string
  belong: string
  time:string
}
type search = {
  data1?: string
  data2?: string
}
type Group = {
  title: string
  type: string
  id: string
  count: number | string
}

const createColumns = ({play}: { play: (row: Song) => void }): DataTableColumns<Song> => {
  return [
    {
      title: '名称',
      key: 'name'
    }, {
      title: '归属',
      key: 'belong'
    }, {
      title: '收藏时间',
      key: 'time'
    }, {
      title: '操作',
      key: 'actions',
      render(row) {
        return h(
            SvgIcon,
            {
              iconName:'collected',
              onClick: () => play(row)
            },
        );
      }
    }
  ];
};

const data: Song[] = [
  {name: 'Wonderwall', belong: 'Wonderwall',time:'2021-01-29 14:32'},
  {name: 'Don\'t Look Back in Anger', belong: 'Don\'t Look Back in Anger',time:'2021-01-29 14:32'},
  {name: 'Champagne Supernova', belong: 'Champagne Supernova',time:'2021-01-29 14:32'}
];

export default {
  name: 'Collection',
  components: {SvgIcon},
  setup() {
    const groupList: Group[] = [{title: '全部', type: 'all', count: 100, id: '1'}
      , {title: 'javaScript', type: '', count: 50, id: '2'}, {title: '好看的花花们', type: '', count: 50, id: '3'},];
    const searchData = ref<search>({});
    const selectedCollectionId = ref('1');
    const showModal = ref(false);
    const newGroupName = ref('');
    const updateGroup = ref({});
    const groupHandleType = ref('');
//选中收藏的某个分类，并保存分类的id
    const viewCollectionList = (collectionId: string) => {
      selectedCollectionId.value = collectionId;
    };
//左侧收藏分类的编辑删除操作
    const changeGroup = (type: string, group?: Group) => {
      groupHandleType.value = type
      if (type === 'update') {
        updateGroup.value = group || {};
        newGroupName.value = group?.title || '';
        showModal.value = true;
      } else if (type === 'delete')  {
        //调用删除接口后，再次调用分组接口刷新页面的分组信息
      }else if (type === 'add')  {
        newGroupName.value = ''
        showModal.value = true;
      }
    };
    //更新分组名字
    const updateGroupName = () => {
      //调用保存接口后，再次调用分组接口刷新页面的分组信息
      //然后关闭弹窗
      showModal.value = false;
    };

//获取右侧表格的标题
    const selectedGroupName = computed(() => {
      return groupList.filter((item) => item.id === selectedCollectionId.value)[0]?.title;
    });
//渲染右侧表格
    const message = useMessage();
    return {
      groupHandleType,
      groupList,
      viewCollectionList,
      changeGroup,
      updateGroupName,
      showModal,
      newGroupName,
      selectedCollectionId,
      selectedGroupName,
      data,
      searchData,
      columns: createColumns({
        play(row: Song) {
          message.info(`Play ${row.name}`);
        }
      }),
      pagination: false as const
    };
  }
};
</script>

<style scoped lang='scss'>
$grey-color: #8A8F8D;
$grey-background: #eff0f0;
.handle-pop-select-options {
  min-width: 104px;

  > li {
    padding: 8px 16px;
    cursor: pointer;

    &:hover {
      background: $grey-background;
    }
  }
}

.collect-wrapper {
  display: flex;
  height: 100%;

  .title {
    color: $grey-color;
    display: flex;
    justify-content: space-between;
    align-items: center;

    > .icon {
      color: #000000;
      cursor: pointer;
    }
  }

  > .left {
    width: 204px;
    padding: 20px 12px;
    border-right: 1px solid #f4f5f5;

    > .title {
      margin-bottom: 20px;

      > .icon {
        margin-right: 12px;
      }
    }

    > .list > li {
      padding: 10px;
      border-radius: 8px;
      margin-bottom: 4px;
      cursor: pointer;

      &.selectedCollection {
        background: $grey-background;
      }

      &:hover {
        background: $grey-background
      }

      > .title {
        color: #585A5A;
      }

      > .content {
        font-size: 12px;
        color: $grey-color;
        margin-top: 1em;
      }
    }
  }

  > .right {
    width: calc(100% - 204px);
    padding: 10px 24px 20px;

    > .title {
      display: flex;
      color: #262626;
      padding: 10px 0;
      align-items: center;

      h3 {
        font-size: 16px;
      }

      > .search {
        display: flex;
        align-items: center;

        > .n-input {
          background: #fafafa;

          &::v-deep(.n-input__border), &::v-deep(.n-input__state-border) {
            display: none;
          }

          &.n-input--focus {
            border: none;
          }
        }

        > .icon {
          color: #585A5A;
          font-size: 18px;
          margin-left: 4px;
        }
      }
    }
  }
}

.dialog-content {
  > label {
    display: inline-flex;
    margin-bottom: 4px;
  }
}
</style>