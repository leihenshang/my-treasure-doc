<template>
    <div id="markdown-container">
    </div>
</template>
  
<script lang="ts" setup>
import Cherry from 'cherry-markdown/dist/cherry-markdown.core'
import 'cherry-markdown/dist/cherry-markdown.min.css'
import { onMounted, ref } from 'vue'


let editor: any = ref(null)

onMounted(() => {
    editor = new Cherry({
        id: 'markdown-container',
        value: '# welcome to cherry editor!',
        callback: {
            afterChange(mb: any, htmlVal: any) {
                console.log(htmlVal)
                console.log(mb)
            }
        },
        toolbars: {
            theme: 'light'
        }
    });
})
</script>
  
<style scoped lang='scss'></style>