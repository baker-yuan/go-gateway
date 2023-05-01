<template>
  <div>
    <div class="container">

      <div class="handle-box">
        <el-select v-model="query.interfaceType" placeholder="接口类型" class="handle-select mr10">
          <el-option key="0" label="http" value="0"></el-option>
          <el-option key="1" label="https" value="1"></el-option>
          <el-option key="2" label="grpc" value="2"></el-option>
          <el-option key="3" label="double" value="3"></el-option>
        </el-select>
        <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
      </div>

      <el-button type="primary" :icon="Plus" @click="handleAdd" style="margin-bottom: 20px;">新增</el-button>

      <el-table :data="tableData" border class="table" ref="multipleTable" header-cell-class-name="table-header">
        <el-table-column prop="id" label="ID" min-width="60" align="center"></el-table-column>
        <el-table-column prop="gwUrl" label="网关路径" min-width="200" align="center"></el-table-column>
        <el-table-column prop="httpType" label="接口类型" align="center" min-width="90">
          <template #default="scope">
            <el-tag type="success">{{scope.row.httpType}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="接口状态" align="center" min-width="90">
          <template #default="scope">
            <el-tag type="success">{{scope.row.status}}</el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="application" label="应用" align="center"></el-table-column>
        <el-table-column prop="interfaceType" label="接口协议" min-width="100" align="center"></el-table-column>
        <el-table-column prop="interfaceUrl" label="接口路径" align="center"></el-table-column>

        <el-table-column label="操作" width="200" align="center">
          <template #default="scope">
            <el-button text :icon="Edit" @click="handleEdit(scope.$index, scope.row)" v-permiss="15">
              编辑
            </el-button>
            <el-button text :icon="Delete" class="red" @click="handleDelete(scope.$index)" v-permiss="16">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          background
          layout="total, prev, pager, next"
          :current-page="query.pageIndex"
          :page-size="query.pageSize"
          :total="pageTotal"
          @current-change="handlePageChange"
        ></el-pagination>
      </div>

    </div>

    <!-- 编辑弹出框 -->
    <el-dialog :title="form.id !== 0 ? '编辑': '新增'" v-model="addOrEditVisible" width="30%">
      <el-form ref="formRef" label-width="70px">
        <el-form-item prop="gwUrl" label="网关接口">
          <el-input v-model="form.gwUrl"></el-input>
        </el-form-item>
        <el-form-item prop="httpType" label="接口类型">
          <el-select v-model="form.httpType" placeholder="接口类型">
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="DELETE" value="DELETE" />
          </el-select>
        </el-form-item>
        <el-form-item prop="application" label="应用名称">
          <el-input v-model="form.application"></el-input>
        </el-form-item>
        <el-form-item prop="interfaceUrl" label="接口方法">
          <el-input v-model="form.interfaceUrl"></el-input>
        </el-form-item>
        <el-form-item prop="interfaceType" label="接口协议">
          <el-select v-model="form.interfaceType" placeholder="接口协议">
            <el-option label="HTTP" value="HTTP" />
            <el-option label="HTTP_S" value="HTTP_S" />
            <el-option label="G_RPC" value="G_RPC" />
            <el-option label="DOUBLE" value="DOUBLE" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
				<span class="dialog-footer">
			    <el-button @click="onReset(formRef)">取 消</el-button>
					<el-button type="primary" @click="saveOrUpdateEdit">确 定</el-button>
				</span>
      </template>
    </el-dialog>


  </div>
</template>

<script setup lang="ts" name="basetable">
import { ref, reactive } from 'vue';
import {ElMessage, ElMessageBox, FormInstance} from 'element-plus';
import { Delete, Edit, Search, Plus } from '@element-plus/icons-vue';
import {gatewayHttpRuleData} from '../api/index';

interface TableItem {
  id: number;
  gwUrl: string;
  httpType: string;
  application: string;
  interfaceType: string;
  interfaceUrl: string;
}

const formRef = ref<FormInstance>();

const query = reactive({
  interfaceType: '',
  pageIndex: 1,
  pageSize: 10
});
const tableData = ref<TableItem[]>([]);
const pageTotal = ref(0);
// 获取表格数据
const getData = () => {
  gatewayHttpRuleData(query).then(res => {
    tableData.value = res.data.rules;
    pageTotal.value = res.data.total || 50;
  });
};
getData();

// 查询操作
const handleSearch = () => {
  query.pageIndex = 1;
  getData();
};

// 分页导航
const handlePageChange = (val: number) => {
  query.pageIndex = val;
  getData();
};

// 删除操作
const handleDelete = (index: number) => {
  // 二次确认删除
  ElMessageBox.confirm('确定要删除吗？', '提示', {
    type: 'warning'
  })
    .then(() => {
      ElMessage.success('删除成功');
      tableData.value.splice(index, 1);
    })
    .catch(() => {});
};

// 表格编辑时弹窗和保存
const addOrEditVisible = ref(false);
let form = reactive({
  id : 0,
  gwUrl: '',
  httpType: '',
  application: '',
  interfaceType: '',
  interfaceUrl: ''
});
let idx: number = -1;
const handleEdit = (index: number, row: any) => {
  idx = index;
  form.id = row.id;
  form.gwUrl = row.gwUrl;
  form.httpType = row.httpType;
  form.application = row.application;
  form.interfaceType = row.interfaceType;
  form.interfaceUrl = row.interfaceUrl;
  addOrEditVisible.value = true;
};
const saveOrUpdateEdit = () => {
  addOrEditVisible.value = false;
  ElMessage.success(`修改第 ${idx + 1} 行成功`);
  tableData.value[idx].id = form.id;
  tableData.value[idx].gwUrl = form.gwUrl;
  tableData.value[idx].httpType = form.httpType;
  tableData.value[idx].application = form.application;
  tableData.value[idx].interfaceType = form.interfaceType;
  tableData.value[idx].interfaceUrl = form.interfaceUrl;
};

const handleAdd = (index: number, row: any) => {
  idx = index;
  tableData.value[idx] = {
    id: 0,
    gwUrl: "",
    httpType: "",
    application:"",
    interfaceType: "",
    interfaceUrl: ""
  }
  addOrEditVisible.value = true;
}

// 重置
const onReset = (formEl: FormInstance | undefined) => {
  console.log('onReset start...', formEl)
  addOrEditVisible.value = false
  if (!formEl) return;
  formEl.resetFields();
  console.log('onReset success...')
};

</script>

<style scoped>
.handle-box {
  margin-bottom: 20px;
}

.handle-select {
  width: 120px;
}

.handle-input {
  width: 300px;
}
.table {
  width: 100%;
  font-size: 14px;
}
.red {
  color: #F56C6C;
}
.mr10 {
  margin-right: 10px;
}
.table-td-thumb {
  display: block;
  margin: auto;
  width: 40px;
  height: 40px;
}
</style>
