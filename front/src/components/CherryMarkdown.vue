<template>
    <div id="markdown-container">
    </div>
</template>
<script lang="ts" setup>
import Cherry from 'cherry-markdown/dist/cherry-markdown.core'
import 'cherry-markdown/dist/cherry-markdown.min.css'
import { onMounted, ref} from 'vue'
import { myHttp } from "@/api/myAxios";
import { useMessage } from 'naive-ui';


const message = useMessage()
const editor = ref<any>(null)
const props = defineProps({
    content: String
})
const emit = defineEmits<{
    (event: 'update', content: string): void
}>()

const docContent = ref<string | undefined>(props.content)

onMounted(() => {
    editor.value = new Cherry({
        id: 'markdown-container',
        value: docContent.value,
        callback: {
            afterChange(mb: any, htmlVal: any) {
                // console.log(htmlVal)
                // console.log(mb)
                // update content variable
                if (mb.length > 0) {
                    docContent.value = mb
                    emit('update', mb)
                }
            }
        },
        fileUpload(file: File, fCallback: any) {
            console.table(file, fCallback)
            myHttp.postForm("api/file/upload", file).then((response: any) => {
                if (!response) {
                    message.error("上传文件失败！")
                    return
                }

                if (response?.data?.code) {
                    message.error("上传失败:" + response?.data?.msg)
                    return
                }

                console.log(response)
                fCallback(response.data.data?.path)

            }).catch((err: any) => {
                console.log(err)
            })
        },
        toolbars: {
            theme: 'light'
        },
        engine: {
            syntax: {
                header: {
                    /**
                     * 标题的样式：
                     *  - default       默认样式，标题前面有锚点
                     *  - autonumber    标题前面有自增序号锚点
                     *  - none          标题没有锚点
                     */
                    anchorStyle: 'autonumber',
                },
            }
        }
    });

})



</script>
<style lang=""></style>