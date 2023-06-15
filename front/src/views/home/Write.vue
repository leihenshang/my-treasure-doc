<template>
    <div id="markdown-container">
    </div>
</template>
  
<script lang="ts">
import Cherry from 'cherry-markdown/dist/cherry-markdown.core'
import 'cherry-markdown/dist/cherry-markdown.min.css'
import { setup } from 'naive-ui/es/radio/src/use-radio'
import { onMounted, getCurrentInstance, ComponentInternalInstance } from 'vue'

export default {
    name: "Write",
    data() {
        return {
            editor: null
        }
    },
    setup() {
        const instance = getCurrentInstance() as ComponentInternalInstance

        onMounted(() => {
            instance.data.editor = new Cherry({
                id: 'markdown-container',
                value: '# welcome to cherry editor!',
                callback: {
                    afterChange(mb: any, htmlVal: any) {
                        console.log(htmlVal)
                        console.log(mb)
                    }
                }
            });
        })
    }
}
</script>
  
<style scoped lang='scss'></style>