<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { getGroupList, createGroup, updateGroup, deleteGroup } from '@/api/group'
import type { TaskGroup, CreateGroupRequest } from '@/api/types'

// 列表数据
const tableData = ref<TaskGroup[]>([])
const loading = ref(false)
const total = ref(0)

// 查询参数
const queryParams = reactive({
  page: 1,
  page_size: 10,
  keyword: ''
})

// 表单相关
const dialogVisible = ref(false)
const dialogTitle = ref('新增任务组')
const formRef = ref<FormInstance>()
const formLoading = ref(false)

const formData = reactive<CreateGroupRequest & { id?: number }>({
  name: '',
  description: '',
  app_name: ''
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入任务组名称', trigger: 'blur' }],
  app_name: [
    { required: true, message: '请输入应用名称', trigger: 'blur' },
    { pattern: /^[a-zA-Z][a-zA-Z0-9_-]*$/, message: '应用名称只能包含字母、数字、下划线和中划线', trigger: 'blur' }
  ]
}

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const res = await getGroupList(queryParams)
    tableData.value = res.data.list || []
    total.value = res.data.total
  } catch (error) {
    console.error('加载数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  queryParams.page = 1
  loadData()
}

// 重置
const handleReset = () => {
  queryParams.keyword = ''
  queryParams.page = 1
  loadData()
}

// 分页
const handlePageChange = (page: number) => {
  queryParams.page = page
  loadData()
}

const handleSizeChange = (size: number) => {
  queryParams.page_size = size
  queryParams.page = 1
  loadData()
}

// 新增
const handleAdd = () => {
  dialogTitle.value = '新增任务组'
  formData.id = undefined
  formData.name = ''
  formData.description = ''
  formData.app_name = ''
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: TaskGroup) => {
  dialogTitle.value = '编辑任务组'
  formData.id = row.id
  formData.name = row.name
  formData.description = row.description
  formData.app_name = row.app_name
  dialogVisible.value = true
}

// 删除
const handleDelete = async (row: TaskGroup) => {
  try {
    await ElMessageBox.confirm(`确定要删除任务组"${row.name}"吗？`, '提示', {
      type: 'warning'
    })
    await deleteGroup(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    // 取消删除
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    formLoading.value = true
    
    if (formData.id) {
      await updateGroup(formData.id, formData)
      ElMessage.success('更新成功')
    } else {
      await createGroup(formData)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    loadData()
  } catch (error: any) {
    // 验证失败或请求失败
  } finally {
    formLoading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="group-page">
    <div class="page-card">
      <div class="page-title">任务组管理</div>
      
      <!-- 搜索栏 -->
      <div class="table-toolbar">
        <div class="table-toolbar-left">
          <el-input
            v-model="queryParams.keyword"
            placeholder="搜索任务组名称或应用名"
            style="width: 240px"
            clearable
            @keyup.enter="handleSearch"
          />
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </div>
        <div class="table-toolbar-right">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增任务组
          </el-button>
        </div>
      </div>

      <!-- 表格 -->
      <el-table v-loading="loading" :data="tableData" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="任务组名称" />
        <el-table-column prop="app_name" label="应用名称" />
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.page_size"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handlePageChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <!-- 表单对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="任务组名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入任务组名称" />
        </el-form-item>
        <el-form-item label="应用名称" prop="app_name">
          <el-input v-model="formData.app_name" placeholder="用于执行器注册，如: my-executor" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="formLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>

