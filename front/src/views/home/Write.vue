<template>
    <div class="edit-box">
        <div class="edit-title">
            <n-input v-model:value="docObj.title" type="text" placeholder="标题" size="large" />
        </div>
        <div class="edit-content">
            <CherryMarkdown :content="docObj.content" @update="contentUpdate" :is-create="true"></CherryMarkdown>
        </div>
    </div>
</template>
  
<script lang="ts" setup>
import Cherry from 'cherry-markdown/dist/cherry-markdown.core'
import 'cherry-markdown/dist/cherry-markdown.min.css'
import { onMounted, ref, watch, onBeforeMount } from 'vue'
import { myHttp } from "../../api/myAxios";
import { useMessage } from 'naive-ui';
import { useRoute } from 'vue-router';
import CherryMarkdown from '@/components/CherryMarkdown.vue';

type DocumentObj = {
    id: number,
    title: string,
    content: string,
    groupId: number,
    isTop: number
}

const route = useRoute()

const docObj = ref<DocumentObj>({
    id: 0,
    title: getTodayStr() + " - 速记",
    content: '# 输入你的想法！',
    groupId: Number(route.query.groupId?.toString()),
    isTop: 0
})
const message = useMessage()

function contentUpdate(content: string) {
    docObj.value.content = content
}

onMounted(() => {
    let isCtrlPressed = false
    document.addEventListener('keydown', function (event) {
        if (event.key === 's') {
            if (isCtrlPressed) {
                message.info('你按下了 "Ctrl + s" 组合键！');
                event.preventDefault();
            }
        } else if (event.key === 'Control') {
            isCtrlPressed = true;
        } else if (event.key === 'Alt') {
            // 可以添加额外的逻辑来处理 Alt + 其他键的组合
        } else {
            // 清除 isCtrlPressed 标志，以便下次按下 Ctrl 时重新开始
            isCtrlPressed = false;
        }
    });
})



function createDoc(doc: any) {
    if (doc.title.length === 0) {
        return
    }

    myHttp.post('api/doc/create', {
        ...doc
    }).then((response: any) => {
        //todo save user information to vuex or state management?
        message.destroyAll()
        // console.log(response)
        if (!response) {
            message.error("响应数据错误！")
            return
        }

        if (response?.data?.code) {
            message.error("创建失败:" + response?.data?.msg)
            return
        }

        docObj.value.id = response.data.data.id


    }).catch((err: any) => {
        console.log(err)
    })
}

function updateDoc(doc: any) {
    if (doc.content.length === 0) {
        return
    }

    myHttp.post('api/doc/update', {
        ...doc
    }).then((response: any) => {
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

function getTodayStr() {
    const today = new Date()
    let todayStr = ''
    todayStr = todayStr.concat(today.getFullYear().toString(), (today.getMonth() + 1).toString().padStart(2, '0'), today.getDate().toString())
    return todayStr
}

watch(docObj, async (newD) => {
    if (newD.id > 0) {
        updateDoc(newD)
    } else {
        createDoc(newD)
    }
}, { deep: true, immediate: true })

</script>
  
<style scoped lang='scss'>
.edit-box {
    margin: 10px 10px;
    height: 100%;

    .edit-title {
        margin: 0 0 10px 0;
    }

    .edit-content {
        border: 1px dashed rgb(176, 170, 170);
        height: 100%;
    }
}
</style>