<template>
  <div class="homePage-wrapper">
    <!-- <Header></Header> -->
    <n-layout has-sider class="menu-layout">
      <!-- left sidebar -->
      <n-layout-sider class="menu-sider" bordered collapse-mode="width" :collapsed-width="64" :width="280"
        :collapsed="collapsed" @collapse="collapsed = true" @expand="collapsed = false">
        <h3>treasure_doc</h3>
        <!-- user menu -->
        <n-menu v-model:value="activeKey" mode="horizontal" :options="horizontalMenuOptions" @update:value="topMenuUpdate"
          :icon-size="18" ref="topMenuRef" />
        <n-menu class="menu-menu" :collapsed="collapsed" :collapsed-width="64" :collapsed-icon-size="22"
          :options="menuOptions" :indent="24" :render-label="renderMenuLabel" :default-value="route.path"
          :render-icon="renderMenuIcon" />
          <n-divider />
        <n-tree block-line :data="data" :default-expanded-keys="defaultExpandedKeys" :checkable="false" expand-on-click
          :selectable="false" />
      </n-layout-sider>
      <!-- right sidebar -->
      <n-layout class="right">
        <router-view></router-view>
      </n-layout>
    </n-layout>
  </div>
</template>

<script lang="ts" setup>
import { h, ref, Component, onMounted } from 'vue';
import { MenuOption, TreeOption, useMessage, NButton, } from 'naive-ui';
import { useRoute, RouterLink } from 'vue-router';
import SvgIcon from '../../components/public/SvgIcon.vue';
import { NIcon } from 'naive-ui';
import {
  EllipsisHorizontalCircleOutline as EllipsisHorizontalCircle,
  Pencil as Pen,
  SearchSharp as Search,
  MailOpen,
  ArrowForwardCircleSharp,
  AppsSharp
} from '@vicons/ionicons5'
import { myHttp } from '@/api/myAxios';
import { router } from '@/router';
import { repeat } from 'seemly'



const route = useRoute();
const topMenuRef = ref(null)
const message = useMessage()
const data = ref<Array<any>>([])
const defaultExpandedKeys = ref(['40', '41'])

console.log(data.value)

function renderSuffix({ option }: { option: TreeOption }) {
  return h(
    NButton,
    { text: true, type: 'primary' },
    { default: () => `Suffix-${option.level}` }
  )
}

function renderPrefix({ option }: { option: TreeOption }) {
  return h(
    NButton,
    { text: true, type: 'primary' },
    { default: () => `Prefix-${option.level}` }
  )
}



function renderLabel({ option }: { option: TreeOption }) {
  return `${option.label} :)`
}

function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const horizontalMenuOptions: MenuOption[] = [
  {
    label: '',
    key: 'top-menu-message',
    icon: renderIcon(MailOpen)
  },
  {
    label: '',
    key: 'top-menu-search',
    icon: renderIcon(Search),
  },
  {
    label: () =>
      h(
        RouterLink,
        {
          to: {
            name: 'Write',
          }
        }
      )
    ,
    key: 'top-menu-write',
    icon: renderIcon(Pen),
  },
  {
    label: '',
    key: 'top-menu-my-center',
    icon: renderIcon(EllipsisHorizontalCircle),
    children: [
      {
        label: '退出登录',
        key: 'login-out',
        icon: renderIcon(ArrowForwardCircleSharp),

      },

    ]
  }
]
const menuOptions = [
  {
    label: '收藏',
    key: 'like',
    pathName: 'Collection',
    iconName: 'collect',
  },
  {
    label: '计划',
    key: 'plan',
    pathName: 'Plan',
    iconName: 'plan',
  },
  {
    label: '笔记',
    key: 'all-doc',
    pathName: 'Note',
    iconName: 'diary',
  },
];


function topMenuUpdate(key: string, item: MenuOption): void {
  console.log(key, item)
  if (key === 'top-menu-write') {
    console.log(menuOptions)
    console.log(topMenuRef)
  }
  if (key == 'login-out') {
    myHttp.post('/api/user/logout', {}).then(() => {
      router.push("/LogIn")
    })
  }
}

const activeKey = ref<string | null>(null)
const collapsed = ref(false)


function renderMenuLabel(option: MenuOption) {
  if ('pathName' in option) {
    return h(
      RouterLink as Component,
      {
        to: {
          name: option.pathName,
        }
      },
      { default: () => option.label }
    );
  }
  return option.label as string;
}

function renderMenuIcon(option: MenuOption) {
  return option.iconName && h(SvgIcon as Component, { iconName: option.iconName });
}


function getDocGroupTree(isAll: Number = 0, pId: Number = 0, withDoc: Number = 1) {
  myHttp.get('/api/doc-group/tree', {
    isAll,
    pId,
    withDoc
  }).then((response: any) => {
    if (!response) {
      message.error("获取分组数据失败！")
      return
    }

    if (response?.data?.code) {
      message.error(response?.data?.msg)
      return
    }

    console.log(response.data)
    data.value = [];
    for (const e of response.data?.data) {
      data.value.push({
        label: e.title,
        key: e.id,
        children: [],
        suffix: () =>
        h(
          NButton,
          { text: true, type: 'primary' },
          { default: () => 'Suffix' }
        ),
      })
    }

  }).catch((err: any) => {
    console.log(err)
  })
}

onMounted(() => {
  getDocGroupTree()
})

</script>

<style scoped lang='scss'>
@import "../../assets/style/helper";

.homePage-wrapper {
  height: 100%;

  >.menu-layout {
    height: 100%;

    .menu-sider {
      background: $menuBackground;

      h3 {
        text-align: left;
        padding-left: 10%;
        height: 40px;
        line-height: 40px;
        // color: rgb(115, 171, 231);
        font-size: 18px;
      }

      .menu-menu ::v-deep(.n-menu-item.n-menu-item--selected) {
        .n-menu-item-content {

          .n-menu-item-content__icon,
          .n-menu-item-content-header {
            color: darken($mainColor, 0.5);
          }
        }

        .n-menu-item>.n-menu-item-content:hover {

          .n-menu-item-content__icon,
          .n-menu-item-content-header {
            color: darken($mainColor, 0.5);
          }
        }

      }
      .n-tree {
        margin-left: 15px
      }
    }
  }
}
</style>