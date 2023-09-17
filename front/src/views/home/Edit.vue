<template>
    <div class="edit-box">
        <div class="edit-title">
            <n-input v-model:value="document.title" type="text" placeholder="标题" size="large" />
        </div>
        <div class="edit-content">
            <div id="markdown-container">
            </div>
        </div>
    </div>
</template>
  
<script lang="ts" setup>
import Cherry from 'cherry-markdown/dist/cherry-markdown.core'
import 'cherry-markdown/dist/cherry-markdown.min.css'
import { onMounted, ref, watch } from 'vue'
import { myHttp } from "../../api/myAxios";
import { useMessage } from 'naive-ui';
import { useRoute, RouterLink } from 'vue-router';

type DocumentObj = {
    id: number,
    title: string,
    content: string
}


const route = useRoute()
const message = useMessage()
const editor: any = ref(null)
const document = ref<DocumentObj>({
    id: 0,
    title: '',
    content: ''
})

// console.log(route.query.id)

onMounted(() => {
    getDoc(Number(route.query?.id))
})


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
        if (document) {
            editor.value = new Cherry({
                id: 'markdown-container',
                value: document.value.content,
                callback: {
                    afterChange(mb: any, htmlVal: any) {
                        // console.log(htmlVal)
                        // console.log(mb)
                        // update content variable
                        if (document.value) {
                            document.value.content = mb
                        }
                    }
                },
                fileUpload(file: File, fCallback: any) {
                    // console.table(file, fCallback)
                    myHttp.postForm("api/file/upload", file).then((response: any) => {
                        if (!response) {
                            message.error("上传文件失败！")
                            return
                        }

                        if (response?.data?.code) {
                            message.error("上传失败:" + response?.data?.msg)
                            return
                        }

                        // console.log(response)
                        fCallback(response.data.data?.path)

                    }).catch((err: any) => {
                        console.log(err)
                    })
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