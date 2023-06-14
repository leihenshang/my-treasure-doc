<template>
    <div id="markdown-container">
    </div>
</template>
  
<script lang="ts">
import Cherry from 'cherry-markdown/dist/cherry-markdown.core';
import 'cherry-markdown/dist/cherry-markdown.min.css'

export default {
    name: "Write",
    data() {
        return {
            editor: null
        }
    },
    mounted() {
        this.init()
    },
    methods: {
        init() {
            const cherryInstance = new Cherry({
                id: 'markdown-container',
                value: '# welcome to cherry editor!',
            });
        }
    }
}
</script>
  
<style scoped lang='scss'></style>