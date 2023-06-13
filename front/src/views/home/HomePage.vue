<template>
  <div class="homePage-wrapper">
    <!-- <Header></Header> -->
    <n-layout has-sider class="menu-layout">
      <!-- left sidebar -->
      <n-layout-sider class="menu-sider" bordered collapse-mode="width" :collapsed-width="64" :width="280"
        :collapsed="collapsed" @collapse="collapsed = true" @expand="collapsed = false">
        <h3>treasure_doc</h3>
        <!-- user menu -->
        <n-menu v-model:value="activeKey" mode="horizontal" :options="horizontalMenuOptions" />
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

<script lang="ts">
import Header from '../../components/Header.vue';
import { h, ref, Component } from 'vue';
import type { MenuOption } from 'naive-ui';
import { useRoute, RouterLink } from 'vue-router';
import SvgIcon from '../../components/public/SvgIcon.vue';
import { NIcon } from 'naive-ui';
import {
  BookOutline as BookIcon,
  PersonOutline as PersonIcon,
  WineOutline as WineIcon,
  EllipsisHorizontalCircleOutline as EllipsisHorizontalCircle,
  Pencil as Pen,
  SearchSharp as Search,
  MailOpen
} from '@vicons/ionicons5'

function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const horizontalMenuOptions: MenuOption[] = [
  {
    label: "",
    key: 'hear-the-wind-sing',
    icon: renderIcon(MailOpen)
  },
  {
    label: '',
    key: 'pinball-1973',
    icon: renderIcon(Search),
  },
  {
    label: '',
    key: 'a-wild-sheep-chase',
    icon: renderIcon(Pen),
  },
  {
    label: '',
    key: 'dance-dance-dance',
    icon: renderIcon(EllipsisHorizontalCircle),
    children: [
      {
        type: 'group',
        label: '人物',
        key: 'people',
        children: [
          {
            label: '叙事者',
            key: 'narrator',
            icon: renderIcon(PersonIcon)
          },
          {
            label: '羊男',
            key: 'sheep-man',
            icon: renderIcon(PersonIcon)
          }
        ]
      },
      {
        label: '饮品',
        key: 'beverage',
        icon: renderIcon(WineIcon),
        children: [
          {
            label: '威士忌',
            key: 'whisky'
          }
        ]
      },
      {
        label: '食物',
        key: 'food',
        children: [
          {
            label: '三明治',
            key: 'sandwich'
          }
        ]
      },
      {
        label: '过去增多，未来减少',
        key: 'the-past-increases-the-future-recedes'
      }
    ]
  }
]
const menuOptions = [
  {
    label: '收藏',
    key: 'like',
    pathName: 'Collection',
    iconName: 'collect',
  }, {
    label: '我的笔记',
    key: 'notes',
    iconName: 'notes',
    children: [
      {
        label: '工作',
        key: 'work',
        pathName: 'Work',
        iconName: 'work',
      }, {
        label: '生活',
        pathName: 'Life',
        key: 'life',
        iconName: 'life',
      }, {
        label: '经验',
        key: 'experience',
        pathName: 'Experience',
        iconName: 'experience',
      }
    ]
  }, {
    label: '计划',
    key: 'plan',
    pathName: 'Plan',
    iconName: 'plan',
  }, {
    label: '我的日记本',
    key: 'diary',
    pathName: 'Diary',
    iconName: 'diary',
  },
];

export default {
  name: 'HomePage',
  components: { Header },
  setup() {
    const route = useRoute();
    return {
      activeKey: ref<string | null>(null),
      horizontalMenuOptions,
      route,
      collapsed: ref(false),
      menuOptions,
      renderMenuLabel(option: MenuOption) {
        if ('pathName' in option) {
          return h(RouterLink,
            {
              to: {
                name: option.pathName,
              }
            },
            { default: () => option.label }
          );
        }
        return option.label as string;
      },
      renderMenuIcon(option: MenuOption) {
        return option.iconName && h(SvgIcon, { iconName: option.iconName });
      },
    };
  }
};
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
       padding-left: 20px;
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