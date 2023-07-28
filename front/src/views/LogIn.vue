<template>
  <div class="log-in-wrap">
    <div class="img">
      <img src="../assets/images/logIn.png" alt="">
    </div>
    <div class="log-in">
      <h3>登陆</h3>
      <n-form :model="userInfo" ref="formRef" label-placement="left" class="long-in-form" size="large">
        <n-form-item path="username" :rule="getRules('用户名')">
          <n-input v-model:value="userInfo.username" autofocus clearable placeholder="用户名" ref="usernameInput" />
        </n-form-item>
        <n-form-item path="password" :rule="getRules('密码')">
          <n-input v-model:value="userInfo.password" type="password" clearable placeholder="密码" />
        </n-form-item>
        <n-form-item class="buttons-wrapper">
          <div class="buttons">
            <n-button attr-type="button" @click="longIn">登录</n-button>
            <n-button attr-type="button" @click="longIn">注册</n-button>
          </div>
        </n-form-item>
      </n-form>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { FormInst, } from 'naive-ui';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useMessage } from 'naive-ui';
import { useUserinfoStore } from "../stores/user/userinfo";
import { myHttp } from "../api/myAxios";


const message = useMessage()
const formRef = ref<FormInst | null>(null)
const router = useRouter()
const storeUserinfo = useUserinfoStore()

interface user {
  username: string,
  password: string
}

const userInfo = ref<user>({
  username: '',
  password: ''
});

const usernameInput = ref(null as any)

const longIn = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors) {
      message.loading("登录...")

      myHttp.post('api/user/login', {
        account: userInfo.value.username,
        password: userInfo.value.password,
        verifyCode: "123456"
      }).then((response: any) => {
        //todo save user information to vuex or state management?
        message.destroyAll()
        console.log(response)
        if (!response) {
          message.error("响应数据错误！")
          return
        }

        if (response?.data?.code) {
          message.error("登录失败:" + response?.data?.msg)
          localStorage.removeItem('userInfo')
          if (usernameInput.value) {
            (usernameInput.value.focus)()
          }

          return
        }

        // local storage
        localStorage.setItem('userInfo', JSON.stringify(response?.data?.data))
        storeUserinfo.updateUserinfo(response?.data?.data)
        message.success("登录成功")
        router.push({ name: 'HomePage' })
      }).catch((err: any) => {
        console.log(err)
      })
    }
  }).catch((err: any) => {
    console.log(err)
  })
}

const getRules = (name: string) => {
  return { required: true, trigger: ['blur', 'input'], message: '请输入' + name }
}

</script>

<style scoped lang='scss'>
.log-in-wrap {
  display: flex;
  width: 100%;
  height: 100%;

  >.img {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 70%;
  }

  >.log-in {
    width: 30%;
    background: #ececec;
    display: flex;
    align-items: center;
    flex-direction: column;
    justify-content: center;

    >h3 {
      font-size: 24px;
      text-align: center;
    }

    >.long-in-form {
      margin-top: 16px;

      > ::v-deep(.n-form-item) {
        .n-form-item-blank {
          .buttons {
            width: 100%;
            display: flex;
            justify-content: space-between;

            button {
              width: 150px;
              background-color: #21a497;
              color: #ffffff;

              .n-button__border {
                border: none;
              }

              &:first-child {
                margin-right: 8px;
              }
            }
          }
        }

        .n-input {
          width: 320px;

          .n-input__border {
            border: 1px solid #b8bcbf;
          }
        }
      }
    }
  }
}
</style>