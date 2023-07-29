<template>
  <div>
    <VueForm
      v-model="formData"
      :schema="schema"
      @submit="submitForm"
    >
    </VueForm>
  </div>
</template>
<!--
      <el-form-item>
        <el-button type="primary" @click="submitForm">Create</el-button>
        <el-button>Cancel</el-button>
      </el-form-item>
-->
<script>
//  使用
import VueForm from '@lljj/vue3-form-element';
import { defineComponent, ref } from 'vue';
export default defineComponent({
  components: {
    VueForm,
  },
  setup() {
    const schema = ref({
      type: 'object',
      required: [
        'userName',
        'age',
      ],
      properties: {
        userName: {
          type: 'string',
          title: '用户名',
          default: 'Liu.Jun',
        },
        age: {
          type: 'number',
          title: '年龄'
        },
        bio: {
          type: 'string',
          title: '签名',
          minLength: 10,
          default: '知道的越多、就知道的越少',
          'ui:options': {
            placeholder: '请输入你的签名',
            type: 'textarea',
            rows: 1
          }
        }
      }
    });
    const formData = ref({});

    const submitForm = () => {
      alert('ok')
      // 将formData发送到后端（Go）进行处理
      // 可以使用axios或其他HTTP库发送POST请求
    };

    return {
      schema,
      formData,
      submitForm,
    };
  },
});
</script>