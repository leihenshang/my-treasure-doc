<template>
    <div class="my-note-wrapper">
        <div class="left">
            <div class="title">
                <span>分组列表</span>
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
                                <SvgIcon icon-name="more" v-if="item.id !== 0"></SvgIcon>
                            </div>
                            <template #empty>
                                <ul class="handle-pop-select-options">
                                    <li @click="changeGroup('update', item)">编辑</li>
                                    <li @click="changeGroup('delete', item)">删除</li>
                                </ul>
                            </template>
                        </n-popselect>
                        <SvgIcon icon-name="add" v-if="item.id !== 0"
                            @click.stop="router.push('/Write?groupId=' + item.id)">
                        </SvgIcon>
                    </p>
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
const isFirst = ref<boolean>(false)

interface DocItemList {
    total: Number,
    list: Array<any>
}

const groupList = ref<Group[]>([{ title: '全部', id: 0 }]);


onMounted(() => {
    getGroupList(1, 10000)
    console.log('Note on mounted')
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
    id: number
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
const selectedCollectionId = ref(0);
const showModal = ref(false);
const newGroupName = ref('');
const updateGroup = ref<Group | null>(null);
const groupHandleType = ref('');

//选中收藏的某个分类，并保存分类的id
const viewCollectionList = (collectionId: number) => {
    selectedCollectionId.value = collectionId;
    docList(pagination.page, pagination.pageSize, true, Number(collectionId));
    console.log('viewCollectionList')
};
//左侧收藏分类的编辑删除操作
const changeGroup = (type: string, group: Group | null = null) => {
    groupHandleType.value = type
    if (type === 'update') {
        updateGroup.value = group;
        newGroupName.value = group?.title || '';
        showModal.value = true;
    } else if (type === 'delete') {
        if (!group) {
            return
        }
        //调用删除接口后，再次调用分组接口刷新页面的分组信息
        // message.info(JSON.stringify(group))
        doDelGroup(group)
        for (let i = 0; i < groupList.value.length; i++) {
            if (groupList.value[i].id === group.id) {
                groupList.value.splice(i, 1)
                i--
            }
        }
    } else if (type === 'add') {
        newGroupName.value = ''
        showModal.value = true;
    }
};
//更新分组名字
const updateGroupName = () => {
    //调用保存接口后，再次调用分组接口刷新页面的分组信息
    if (groupHandleType.value == 'update') {
        let g = updateGroup.value as Group
        if (g.id > 0) {
            g.title = newGroupName.value
            doUpdateGroup(g)
        }
    }

    if (groupHandleType.value == 'add') {
        doCreateGroup({
            title: newGroupName.value,
            id: 0
        })
    }


    //然后关闭弹窗
    showModal.value = false;
};

//获取右侧表格的标题
const selectedGroupName = computed(() => {
    return groupList.value.filter((item) => item.id === selectedCollectionId.value)[0]?.title;
});

function rowKey(rowData: any) {
    return rowData?.id
}


function handlePageChange(currentPage: any) {
    if (!loading.value) {
        loading.value = true
        docList(currentPage, 10, true);
        console.log('handlePageChange')
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


async function docList(page: number = 1, pageSize: number = 10, isDesc: boolean = true, groupId: number = 0) {
    await myHttp.post('api/doc/list', {
        page,
        pageSize,
        isDesc,
        groupId
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

function getGroupList(page: number = 1, pageSize: number = 10, refreshAll: boolean = false) {
    myHttp.post('/api/doc-group/list', {
        page,
        pageSize,
    }).then((response: any) => {
        if (!response) {
            message.error("获取分组数据失败！")
            return
        }

        if (response?.data?.code) {
            message.error(response?.data?.msg)
            return
        }
        if (refreshAll) {
            groupList.value = [groupList.value[0]]
        }
        groupList.value.push(...response.data?.data.list)


    }).catch((err: any) => {
        console.log(err)
    })
}

function doUpdateGroup(groupObj: Group) {
    myHttp.post('/api/doc-group/update', {
        ...groupObj
    }).then((response: any) => {
        if (!response) {
            message.error("更新分组数据失败！")
            return
        }

        if (response?.data?.code) {
            message.error(response?.data?.msg)
            return
        }

    }).catch((err: any) => {
        console.log(err)
    })
}

function doDelGroup(groupObj: Group) {
    myHttp.post('/api/doc-group/delete', {
        ...groupObj
    }).then((response: any) => {
        if (!response) {
            message.error("删除分组数据失败！")
            return
        }

        if (response?.data?.code) {
            message.error(response?.data?.msg)
            return
        }

    }).catch((err: any) => {
        console.log(err)
    })
}

function doCreateGroup(groupObj: Group) {
    myHttp.post('/api/doc-group/create', {
        ...groupObj
    }).then((response: any) => {
        if (!response) {
            message.error("创建分组数据失败！")
            return
        }

        if (response?.data?.code) {
            message.error(response?.data?.msg)
            return
        }
        getGroupList(1, 10000, true)
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

                .icon {
                    margin: 0 5px;
                }
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