<template>
  <div class="homePage-wrapper">
    <!-- <Header></Header> -->
    <n-layout has-sider class="menu-layout">
      <!-- left sidebar -->
      <n-layout-sider class="menu-sider" bordered collapse-mode="width" :collapsed-width="64" :width="280"
        :collapsed="collapsed" @collapse="collapsed = true" @expand="collapsed = false">
        <n-icon size="40">
         <svg t="1694998158700" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="6303" xmlns:xlink="http://www.w3.org/1999/xlink" width="200" height="200"><path d="M901.2 317.4c0-64.4-52.2-116.8-116.8-116.8H239.6c-64.4 0-116.8 52.2-116.8 116.8v447.6c0 32.2 26.2 58.4 58.4 58.4h661.6c32.2 0 58.4-26.2 58.4-58.4V317.4z" fill="#FCD64B" p-id="6304"></path><path d="M784.4 200.6h-39c64.4 0 116.8 52.2 116.8 116.8v447.6c0 32.2-26.2 58.4-58.4 58.4h39c32.2 0 58.4-26.2 58.4-58.4V317.4c0-64.4-52.4-116.8-116.8-116.8z" fill="#F5AB33" p-id="6305"></path><path d="M239.6 200.6h39c-64.4 0-116.8 52.2-116.8 116.8v447.6c0 32.2 26.2 58.4 58.4 58.4H181.2c-32.2 0-58.4-26.2-58.4-58.4V317.4c0-64.4 52.4-116.8 116.8-116.8z" fill="#FFECA6" p-id="6306"></path><path d="M161.8 434.2h739.4V512H161.8z" fill="#F5AB33" p-id="6307"></path><path d="M842.8 833H181.2c-37.6 0-68.2-30.6-68.2-68.2V317.2c0-69.8 56.8-126.4 126.4-126.4h544.8c69.8 0 126.4 56.8 126.4 126.4v447.6c0.2 37.8-30.2 68.2-67.8 68.2zM239.6 210.4c-59 0-107 48-107 107v447.6c0 26.8 21.8 48.6 48.6 48.6h661.6c26.8 0 48.6-21.8 48.6-48.6V317.4h9.8-9.8c0-59-48-107-107-107H239.6z" fill="#402943" p-id="6308"></path><path d="M122.8 473h778.4" fill="#FFFFFF" p-id="6309"></path><path d="M901.2 482.8H122.8c-5.4 0-9.8-4.4-9.8-9.8s4.4-9.8 9.8-9.8h778.4c5.4 0 9.8 4.4 9.8 9.8-0.2 5.4-4.4 9.8-9.8 9.8z" fill="#402943" p-id="6310"></path><path d="M200.6 551v136.2c0 32.2 26.2 58.4 58.4 58.4h506c32.2 0 58.4-26.2 58.4-58.4v-136.2H200.6z" fill="#FFFFFF" p-id="6311"></path><path d="M784.4 551v136.2c0 32.2-26.2 58.4-58.4 58.4h39c32.2 0 58.4-26.2 58.4-58.4v-136.2h-39zM239.6 551v136.2c0 32.2 26.2 58.4 58.4 58.4h-39c-32.2 0-58.4-26.2-58.4-58.4v-136.2h39z" fill="#E6E7E8" p-id="6312"></path><path d="M200.6 551h622.6V590H200.6z" fill="#E6E7E8" p-id="6313"></path><path d="M765 755.2h-506c-37.6 0-68.2-30.6-68.2-68.2v-136.2c0-2.6 1-5 2.8-6.8 1.8-1.8 4.2-2.8 6.8-2.8h622.6c5.4 0 9.8 4.4 9.8 9.8v136.2c0.2 37.4-30.2 68-67.8 68z m-554.6-194.6v126.4c0 26.8 21.8 48.6 48.6 48.6h506c26.8 0 48.6-21.8 48.6-48.6v-126.4H210.4z" fill="#402943" p-id="6314"></path><path d="M823.4 395.2v-77.8c0-64.4-52.2-116.8-116.8-116.8H317.4c-64.4 0-116.8 52.2-116.8 116.8v77.8h622.8z" fill="#FFFFFF" p-id="6315"></path><path d="M706.6 200.6h-39c64.4 0 116.8 52.2 116.8 116.8v77.8h39v-77.8c0-64.4-52.4-116.8-116.8-116.8zM317.4 200.6h39c-64.4 0-116.8 52.2-116.8 116.8v77.8h-39v-77.8c0-64.4 52.4-116.8 116.8-116.8z" fill="#E6E7E8" p-id="6316"></path><path d="M823.4 405H200.6c-5.4 0-9.8-4.4-9.8-9.8v-77.8c0-69.8 56.8-126.4 126.4-126.4h389.2c69.8 0 126.4 56.8 126.4 126.4v77.8c0.2 5.4-4 9.8-9.4 9.8z m-613-19.4h603.2v-68.2c0-59-48-107-107-107H317.4c-59 0-107 48-107 107v68.2z" fill="#402943" p-id="6317"></path><path d="M609.2 492.6c0 53.8-43.6 97.2-97.2 97.2s-97.2-43.6-97.2-97.2v-136.2h194.6l-0.2 136.2z" fill="#FCD64B" p-id="6318"></path><path d="M570.4 356.4v136.2c0 47-33.4 86.4-77.8 95.4 6.2 1.2 12.8 2 19.4 2 53.8 0 97.2-43.6 97.2-97.2v-136.2h-38.8z" fill="#F5AB33" p-id="6319"></path><path d="M453.6 356.4v136.2c0 47 33.4 86.4 77.8 95.4-6.2 1.2-12.8 2-19.4 2-53.8 0-97.2-43.6-97.2-97.2v-136.2h38.8z" fill="#FFECA6" p-id="6320"></path><path d="M512 599.6c-59 0-107-48-107-107v-136.2c0-5.4 4.4-9.8 9.8-9.8h194.6c2.6 0 5 1 6.8 2.8 1.8 1.8 2.8 4.2 2.8 6.8v136.2c0 59.2-48 107.2-107 107.2zM424.4 366v126.4c0 48.2 39.2 87.6 87.6 87.6 48.2 0 87.6-39.2 87.6-87.6V366h-175.2z" fill="#402943" p-id="6321"></path><path d="M551 400.2h-136.2c-2.6 0-4.8-2.2-4.8-4.8 0-2.6 2.2-4.8 4.8-4.8h136.2c2.6 0 4.8 2.2 4.8 4.8 0 2.6-2.2 4.8-4.8 4.8z" fill="#402943" p-id="6322"></path><path d="M609.2 400.2h-39c-2.6 0-4.8-2.2-4.8-4.8s2.2-4.8 4.8-4.8h39c2.6 0 4.8 2.2 4.8 4.8s-2 4.8-4.8 4.8z" fill="#402943" p-id="6323"></path><path d="M538 460.2c0-14.4-11.6-26-26-26s-26 11.6-26 26c0 9.6 5.2 18 13 22.4V512h26v-29.4c7.8-4.6 13-12.8 13-22.4z" fill="#59C3F0" p-id="6324"></path><path d="M525 521.8h-26c-5.4 0-9.8-4.4-9.8-9.8v-24.4c-8.2-6.8-13-16.8-13-27.4 0-19.6 16-35.6 35.6-35.6s35.6 16 35.6 35.6c0 10.6-4.8 20.8-13 27.4V512c0.4 5.4-4 9.8-9.4 9.8z m-16.2-19.6h6.4v-19.8c0-3.4 1.8-6.6 4.8-8.4 5-3 8.2-8.2 8.2-14 0-9-7.2-16.2-16.2-16.2s-16.2 7.2-16.2 16.2c0 5.8 3.2 11.2 8.2 14 3 1.8 4.8 5 4.8 8.4v19.8z" fill="#402943" p-id="6325"></path><path d="M200.6 560.6h-19.4c-5.4 0-9.8-4.4-9.8-9.8s4.4-9.8 9.8-9.8h19.4c5.4 0 9.8 4.4 9.8 9.8s-4.4 9.8-9.8 9.8zM842.8 560.6h-19.4c-5.4 0-9.8-4.4-9.8-9.8s4.4-9.8 9.8-9.8h19.4c5.4 0 9.8 4.4 9.8 9.8s-4.4 9.8-9.8 9.8z" fill="#402943" p-id="6326"></path><path d="M200.6 405h-19.4c-5.4 0-9.8-4.4-9.8-9.8s4.4-9.8 9.8-9.8h19.4c5.4 0 9.8 4.4 9.8 9.8s-4.4 9.8-9.8 9.8zM842.8 405h-19.4c-5.4 0-9.8-4.4-9.8-9.8s4.4-9.8 9.8-9.8h19.4c5.4 0 9.8 4.4 9.8 9.8s-4.4 9.8-9.8 9.8z" fill="#402943" p-id="6327"></path><path d="M825 755.2H200.6c-5.4 0-9.8-4.4-9.8-9.8s4.4-9.8 9.8-9.8h624.4c5.4 0 9.8 4.4 9.8 9.8s-4.4 9.8-9.8 9.8z" fill="#402943" p-id="6328"></path><path d="M317.4 346.6c-5.4 0-9.8-4.4-9.8-9.8V298c0-5.4 4.4-9.8 9.8-9.8s9.8 4.4 9.8 9.8v39c0 5.2-4.4 9.6-9.8 9.6z" fill="#402943" p-id="6329"></path><path d="M336.8 327.2H298c-5.4 0-9.8-4.4-9.8-9.8s4.4-9.8 9.8-9.8h39c5.4 0 9.8 4.4 9.8 9.8-0.2 5.4-4.6 9.8-10 9.8z" fill="#402943" p-id="6330"></path><path d="M395.2 307.6c-5.4 0-9.8-4.4-9.8-9.8v-39c0-5.4 4.4-9.8 9.8-9.8s9.8 4.4 9.8 9.8V298c0 5.4-4.4 9.6-9.8 9.6z" fill="#402943" p-id="6331"></path><path d="M414.8 288.2h-39c-5.4 0-9.8-4.4-9.8-9.8s4.4-9.8 9.8-9.8h39c5.4 0 9.8 4.4 9.8 9.8-0.2 5.4-4.6 9.8-9.8 9.8z" fill="#402943" p-id="6332"></path></svg>
        </n-icon>
        <h3>treasure_doc</h3>
        <!-- user menu -->
        <n-menu v-model:value="activeKey" mode="horizontal" :options="horizontalMenuOptions" @update:value="topMenuUpdate"
          :icon-size="18" ref="topMenuRef" />
        <n-menu class="menu-menu" :collapsed="collapsed" :collapsed-width="64" :collapsed-icon-size="22"
          :options="menuOptions" :indent="24" :render-label="renderMenuLabel" :default-value="route.path"
          :render-icon="renderMenuIcon" />
      </n-layout-sider>
      <!-- right sidebar -->
      <n-layout class="right">

        <router-view></router-view>
      </n-layout>
    </n-layout>
  </div>
</template>

<script lang="ts" setup>
import Header from '../../components/Header.vue';
import { h, ref, Component, computed, DefineComponent } from 'vue';
import { MenuOption, messageDark, useMessage } from 'naive-ui';
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



const route = useRoute();
const topMenuRef = ref(null)
const msg = useMessage()

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
        text-align: center;
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
    }
  }
}
</style>