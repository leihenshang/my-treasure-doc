<template>
    <div class="my-note-wrapper">
        <n-data-table remote ref="table" :columns="columns" :data="data" :loading="loading" :pagination="pagination"
            :row-key="rowKey" @update:page="handlePageChange" :row-props="rowProps" />
    </div>
</template>
<script  lang="ts" setup>
import { ref, reactive, onMounted, h, nextTick } from 'vue'
import { myHttp } from "../../api/myAxios";
import { NButton, useMessage } from 'naive-ui'
import { useRouter, RouterLink } from 'vue-router';

const message = useMessage()

interface DocItemList {
    total: Number,
    list: Array<any>
}


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
  
<style scoped lang='scss'></style>