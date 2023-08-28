<template>
    <div class="my-note-wrapper">
        <n-data-table remote ref="table" :columns="columns" :data="data" :loading="loading" :pagination="pagination"
            :row-key="rowKey" @update:page="handlePageChange" />
    </div>
</template>
<script  lang="ts" setup>
import Cherry from 'cherry-markdown'
import { defineComponent, ref, reactive, onMounted, h } from 'vue'
import { myHttp } from "../../api/myAxios";
import { NButton, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'


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


function handlePageChange(currentPage) {
    if (!loadingRef.value) {
        loadingRef.value = true
        query(
            currentPage,
            paginationReactive.pageSize,
            column1Reactive.sortOrder,
            column2Reactive.filterOptionValues
        ).then((data) => {
            dataRef.value = data.data
            paginationReactive.page = currentPage
            paginationReactive.pageCount = data.pageCount
            paginationReactive.itemCount = data.total
            loadingRef.value = false
        })
    }
}


function docList(page: Number = 1, pageSize: Number = 10, isDesc: boolean = true) {
    myHttp.post('api/doc/list', {
        "page": page,
        "pageSize": pageSize,
        "isDesc": isDesc
    }, {
        headers: { "X-Token": '3b5b3d702a9637860ac351550859cd19' }
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
        data.value = response.data?.data.list
        loading.value = false

    }).catch((err: any) => {
        console.log(err)
    })
}

</script>
  
<style scoped lang='scss'></style>