<template>
    <div class="my-note-wrapper">
        <div class="left">
            <div class="title">
                <span>收藏列表</span>
                <SvgIcon icon-name="add" @click="changeGroup('add')"></SvgIcon>
            </div>
            <ul class="list">
                <li v-for="item in groupList" :key="item.id"
                    :class="{ 'selectedCollection': selectedCollectionId === item.id }"
                    @click="viewCollectionList(item.id)">
                    <p class="title">
                        <span>{{ item.title }}</span>
                        <n-popselect :options="[]" trigger="hover">
                            <div>
                                <SvgIcon icon-name="more" v-if="item.type !== 'all'"></SvgIcon>
                            </div>
                            <template #empty>
                                <ul class="handle-pop-select-options">
                                    <li @click="changeGroup('update', item)">编辑</li>
                                    <li @click="changeGroup('delete', item)">删除</li>
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
                    <n-input v-model:value="searchData.data1" type="text" size="small" placeholder="搜索">
                        <template #prefix>
                            <SvgIcon icon-name="search"></SvgIcon>
                        </template>
                    </n-input>
                    <SvgIcon icon-name="filter"></SvgIcon>
                </div>
            </div>
            <n-data-table remote ref="table" :columns="columns" :data="data" :loading="loading" :pagination="pagination"
                :row-key="rowKey" @update:page="handlePageChange" :row-props="rowProps" />
        </div>

    </div>
    <n-modal v-model:show="showModal" preset="dialog" title="Dialog" :show-icon="false" class="modal-dialog"
        :mask-closable=false style="position: fixed; left: 50%;transform: translateX(-50%);top: 100px">
        <template #header>
            <div>{{ groupHandleType === 'add' ? '新增分组' : '编辑分组' }}</div>
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
<script  lang="ts" setup>
import { ref, reactive, onMounted, h, nextTick, computed } from 'vue'
import { myHttp } from "@/api/myAxios";
import { NButton, useMessage } from 'naive-ui'
import { useRouter, RouterLink } from 'vue-router';
import SvgIcon from '@/components/public/SvgIcon.vue';

const message = useMessage()

interface DocItemList {
    total: Number,
    list: Array<any>
}

const groupList: Group[] = [{ title: '全部', type: 'all', count: 100, id: '1' }
    , { title: 'javaScript', type: '', count: 50, id: '2' }, { title: '好看的花花们', type: '', count: 50, id: '3' },];


onMounted(() => {
    docList(
        pagination.page,
        pagination.pageSize,
        true
    );
})

const columns = ref([
    {
        title: 'ID',
        key: 'id',
        sorter: true,
        sortOrder: false
    },
    {
        title: '名称',
        key: 'title'
    },
    {
        title: '置顶',
        key: 'isTop'
    },
    {
        title: '创建时间',
        key: 'createdAt'
    }
])

type Song = {
    name: string
    belong: string
    time: string
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


const router = useRouter()

const data = ref([])
const loading = ref(true)
const pagination = reactive({
    page: 1,
    pageCount: 1,
    pageSize: 10,
    prefix({ itemCount }: { itemCount: any }) {
        return `Total is ${itemCount}.`
    }
})
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
    } else if (type === 'delete') {
        //调用删除接口后，再次调用分组接口刷新页面的分组信息
    } else if (type === 'add') {
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

function rowKey(rowData: any) {
    return rowData?.id
}


function handlePageChange(currentPage: any) {
    if (!loading.value) {
        loading.value = true
        docList(currentPage, 10, true);
    }
}

type RowData = {
    key: number
    name: string
    age: string
    address: string
}

function rowProps(row: any) {
    return {
        style: 'cursor: pointer;',
        onClick: () => {
            // message.info(JSON.stringify(row))

            router.push({ path: '/Edit', query: { id: row.id } })
        }
    }
}


function docList(page: number = 1, pageSize: number = 10, isDesc: boolean = true) {
    myHttp.post('api/doc/list', {
        "page": page,
        "pageSize": pageSize,
        "isDesc": isDesc
    }).then((response: any) => {
        //todo save user information to vuex or state management?
        message.destroyAll()
        if (!response) {
            message.error("响应数据错误！")
            return
        }

        if (response?.data?.code) {
            message.error(response?.data?.msg)
            return
        }

        pagination.itemCount = response.data?.data.total
        pagination.page = page
        pagination.pageCount = response.data?.data.total / 10
        pagination.pageSize = pageSize
        data.value = response.data?.data.list
        loading.value = false
    }).catch((err: any) => {
        console.log(err)
    })
}

</script>
  
<style scoped lang='scss'>
$grey-color: #8A8F8D;
$grey-background: #eff0f0;

.handle-pop-select-options {
    min-width: 104px;

    >li {
        padding: 8px 16px;
        cursor: pointer;

        &:hover {
            background: $grey-background;
        }
    }
}

.my-note-wrapper {
    display: flex;
    height: 100%;

    .title {
        color: $grey-color;
        display: flex;
        justify-content: space-between;
        align-items: center;

        >.icon {
            color: #000000;
            cursor: pointer;
        }
    }

    >.left {
        width: 204px;
        padding: 20px 12px;
        border-right: 1px solid #f4f5f5;

        >.title {
            margin-bottom: 20px;

            >.icon {
                margin-right: 12px;
            }
        }

        >.list>li {
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

            >.title {
                color: #585A5A;
            }

            >.content {
                font-size: 12px;
                color: $grey-color;
                margin-top: 1em;
            }
        }
    }

    >.right {
        width: calc(100% - 204px);
        padding: 10px 24px 20px;

        >.title {
            display: flex;
            color: #262626;
            padding: 10px 0;
            align-items: center;

            h3 {
                font-size: 16px;
            }

            >.search {
                display: flex;
                align-items: center;

                >.n-input {
                    background: #fafafa;

                    &::v-deep(.n-input__border),
                    &::v-deep(.n-input__state-border) {
                        display: none;
                    }

                    &.n-input--focus {
                        border: none;
                    }
                }

                >.icon {
                    color: #585A5A;
                    font-size: 18px;
                    margin-left: 4px;
                }
            }
        }
    }
}

.dialog-content {
    >label {
        display: inline-flex;
        margin-bottom: 4px;
    }
}
</style>