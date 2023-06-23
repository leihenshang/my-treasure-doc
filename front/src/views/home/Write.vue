<template>
    <div id="markdown-container">
    </div>
</template>
  
<script lang="ts" setup>
import Cherry from 'cherry-markdown/dist/cherry-markdown.core'
import 'cherry-markdown/dist/cherry-markdown.min.css'
import { onMounted, ref } from 'vue'
import { myHttp } from "../../api/myAxios";
import { useMessage } from 'naive-ui';


const message = useMessage()
const editor: any = ref(null)
let docId: number = 0
let content: string = ''


onMounted(() => {
    editor.value = new Cherry({
        id: 'markdown-container',
        value: '# welcome to cherry editor!',
        callback: {
            afterChange(mb: any, htmlVal: any) {
                console.log(htmlVal)
                console.log(mb)
                // update content variable
                content = mb

                if (docId > 0) {
                    // update
                    updateDoc()
                } else {
                    // create
                    createDoc()
                }

            }
        },
        toolbars: {
            theme: 'light'
        }
    });
})



function createDoc() {
    if (content.length === 0) {
        return
    }

    myHttp.post('api/doc/create', {
        title: getTodayStr() + "速记" + Math.random().toString(),
        content: content,
        groupId: 0,
        isTop: 0
    }, {
        headers: { "X-Token": '3b5b3d702a9637860ac351550859cd19' }
    }).then((response: any) => {
        //todo save user information to vuex or state management?
        message.destroyAll()
        console.log(response)
        if (!response) {
            message.error("响应数据错误！")
            return
        }

        if (response?.data?.code) {
            message.error("创建失败:" + response?.data?.msg)
            return
        }

        docId = (Number)(response?.data?.data?.id)

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
    }, {
        headers: { "X-Token": '3b5b3d702a9637860ac351550859cd19' }
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

function getTodayStr() {
    const today = new Date()
    let todayStr = ''
    todayStr = todayStr.concat(today.getFullYear().toString(), (today.getMonth() + 1).toString().padStart(2, '0'), today.getDate().toString())
    return todayStr
}

</script>
  
<style scoped lang='scss'></style>