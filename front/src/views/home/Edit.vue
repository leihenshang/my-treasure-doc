<template>
    <div style="width: 100%;height: 100%;">
    <p>{{ title }}</p>
    <div id="markdown-container">
    </div>
</div>
</template>
  
<script lang="ts" setup>
import Cherry from 'cherry-markdown/dist/cherry-markdown.core'
import 'cherry-markdown/dist/cherry-markdown.min.css'
import { onMounted, ref } from 'vue'
import { myHttp } from "../../api/myAxios";
import { useMessage } from 'naive-ui';
import { useRoute, RouterLink } from 'vue-router';

const route = useRoute()
const message = useMessage()
const editor: any = ref(null)
let docId: number = 0
let content: string = ''
let document: Object = ref(null)
let title = ref('')

console.log(route.query.id)

onMounted(() => {

    getDoc(Number(route.query?.id))


})


function getDoc(id: number) {
    myHttp.post('api/doc/detail', {
        id: Number(id)
    }).then((response: any) => {
        //todo save user information to vuex or state management?
        message.destroyAll()
        console.log(response)
        if (!response) {
            message.error("响应数据错误！")
            return
        }

        if (response?.data?.code) {
            message.error("msg:" + response?.data?.msg)
            return
        }

        document = response?.data?.data
        if (document) {
            docId = document.id
            title = document.title
            editor.value = new Cherry({
                id: 'markdown-container',
                value: document.content,
                callback: {
                    afterChange(mb: any, htmlVal: any) {
                        console.log(htmlVal)
                        console.log(mb)
                        // update content variable
                        content = mb
                        if (docId > 0) {
                            // update
                            updateDoc()
                        }

                    }
                },
                toolbars: {
                    theme: 'light'
                }
            });
        }

    }).catch((err: any) => {
        console.log(err)
    })
}


function updateDoc() {
    if (content.length === 0) {
        return
    }

    myHttp.post('api/doc/update', {
        // title: getTodayStr() + "速记",
        content: content,
        // groupId: 0,
        // isTop: 0,
        id: docId
    }).then((response: any) => {
        //todo save user information to vuex or state management?
        message.destroyAll()
        console.log(response)
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

</script>
  
<style scoped lang='scss'></style>