<template>
    <div class="my-note-wrapper">
        New111
        <ul>
        <li v-for="item in docListObj.list">
           {{ item.title }}
        </li>
    </ul>
    </div>
</template>
<script  lang="ts" setup>
import Cherry from 'cherry-markdown'
import { onMounted, ref } from 'vue'
import { myHttp } from "../../api/myAxios";
import { useMessage } from 'naive-ui';


const message = useMessage()

interface DocItemList {
    total: Number,
    list: Array<any>
}

let docListObj = ref<DocItemList>({
    total: 0,
    list: []
})

onMounted(() => {
    docList();
})

function docList() {
    myHttp.post('api/doc/list', {
        "page": 1,
        "pageSize": 10,
        "isDesc": true
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

        docListObj = response.data?.data
        console.log(docListObj)

    }).catch((err: any) => {
        console.log(err)
    })
}

</script>
  
<style scoped lang='scss'></style>