<template>
    <div class="edit-box">
        <div class="edit-title">
            <n-input v-model:value="document.title" type="text" placeholder="标题" size="large" />
        </div>
        <div class="edit-content">
            <CherryMarkdown :content="document.content" @update="contentUpdate"></CherryMarkdown>
        </div>
    </div>
</template>
  
<script lang="ts" setup>
import 'cherry-markdown/dist/cherry-markdown.min.css'
import { onMounted, ref, watch } from 'vue'
import { myHttp } from "../../api/myAxios";
import { useMessage } from 'naive-ui';
import { useRoute, RouterLink } from 'vue-router';
import CherryMarkdown from '@/components/CherryMarkdown.vue';

type DocumentObj = {
    id: number,
    title: string,
    content: string
}


const route = useRoute()
const message = useMessage()
const document = ref<DocumentObj>({
    id: 0,
    title: '',
    content: ''
})

onMounted(() => {
    getDoc(Number(route.query?.id))
})

function contentUpdate(content: string) {
    document.value.content = content
}


function getDoc(id: number) {
    myHttp.post('api/doc/detail', {
        id: Number(id)
    }).then((response: any) => {
        //todo save user information to vuex or state management?
        message.destroyAll()
        // console.log(response)
        if (!response) {
            message.error("响应数据错误！")
            return
        }

        if (response?.data?.code) {
            message.error("msg:" + response?.data?.msg)
            return
        }

        document.value = response?.data?.data as DocumentObj
    }).catch((err: any) => {
        console.log(err)
    })
}


function updateDoc(doc: any) {
    if (!document.value || document.value?.content.length === 0) {
        return
    }

    myHttp.post('api/doc/update', { ...document.value }).then((response: any) => {
        //todo save user information to vuex or state management?
        message.destroyAll()
        // console.log(response)
        if (!response) {
            message.error("响应数据错误！")
            return
        }

        if (response?.data?.code) {
            message.error("更新失败:" + response?.data?.msg)
            return
        }

    }).catch((err: any) => {
        console.log(err)
    })
}

watch(document, async (newD, oldD) => {
    if (oldD.id > 0) {
        updateDoc(newD)
    }
}, { deep: true })

</script>
  
<style scoped lang='scss'>
.edit-box {
    margin: 5px 5px;
    height: 100%;

    .edit-title {
        margin: 0 0 5px 0;
    }

    .edit-content {
        height: 100%;
    }
}
</style>