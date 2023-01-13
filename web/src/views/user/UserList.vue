<template>
  <div>
    <a-card :bordered="false">
      <div class="table-operator">
        <a-button type="primary" icon="plus" @click="handleAddUser">新建</a-button>
        <a-button icon="reload" @click="refresh">刷新</a-button>
        <a-dropdown v-action:edit v-if="selectedRowKeys.length > 0" :trigger="['click']">
          <a-menu slot="overlay" @click="e => handleBatchOpClick(e)">
            <a-menu-item key="delete"><a-icon type="delete" />删除</a-menu-item>
            <!-- lock | unlock -->
            <a-menu-item key="disabled"><a-icon type="lock" />禁用</a-menu-item>
          </a-menu>
          <a-button style="margin-left: 8px">
            批量操作 <a-icon type="down" />
          </a-button>
        </a-dropdown>
      </div>
      <s-table
        ref="table"
        size="small"
        rowKey="id"
        :columns="columns"
        :data="loadData"
        :alert="true"
        :rowSelection="rowSelection"
        :showPagination="true"
      >
        <span slot="status" slot-scope="text">
          <a-badge :status="text | statusTypeFilter" :text="text | statusFilter" />
        </span>
        <span slot="desc" slot-scope="text">
          <ellipsis :length="30" tooltip>{{ text }}</ellipsis>
        </span>
        <span slot="isAdmin" slot-scope="text">
          <a-tag v-if="text" color="#87d068">是</a-tag>
          <a-tag v-else color="orange">否</a-tag>
        </span>
        <template slot="operation" slot-scope="record">
          <a-tooltip placement="right">
            <template slot="title">
              <span>编辑</span>
            </template>
            <a-button size="small" type="link" icon="edit" @click="handleEditUser(record)"></a-button>
          </a-tooltip>
        </template>
      </s-table>
    </a-card>
    <!-- BEGIN 创建账号 modal -->
    <a-modal
      :visible="addUserModalVisible"
      :maskCloseable="false"
      :confirmLoading="addUserLoading"
      width="50%"
      title="创建用户"
      @ok="createUser"
      @cancel="addUserModalVisible = false"
    >
      <a-form-model
        :model="userForm"
        :label-col="labelCol"
        :wrapper-col="wrapperCol"
        :rules="rules"
        layout="horizontal"
        style="width: 100%"
        ref="userForm"
      >
        <a-form-model-item label="用户名" prop="username">
          <a-input v-model="userForm.username"></a-input>
        </a-form-model-item>
        <a-form-model-item label="密码" prop="password">
          <a-input-password v-model="userForm.password"></a-input-password>
        </a-form-model-item>
        <a-form-model-item label="昵称" prop="nickname">
          <a-input v-model="userForm.nickname"></a-input>
        </a-form-model-item>
        <a-form-model-item label="邮箱" prop="email">
          <a-input v-model="userForm.email"></a-input>
        </a-form-model-item>
        <a-form-model-item label="管理员" prop="isAdmin">
          <a-radio-group v-model="userForm.is_admin">
            <a-radio :value="true">是</a-radio>
            <a-radio :value="false">否</a-radio>
          </a-radio-group>
        </a-form-model-item>
        <a-form-model-item label="状态" prop="enable">
          <a-radio-group v-model="userForm.enable">
            <a-radio :value="true">激活</a-radio>
            <a-radio :value="false">禁用</a-radio>
          </a-radio-group>
        </a-form-model-item>
      </a-form-model>
    </a-modal>
    <!-- END 创建账号 modal -->
  </div>
</template>

<script>
import { STable, Ellipsis } from '@/components'
import { getUserList, createUserList } from '@/api/user'
import { getPublicKey } from '@/api/common'
import JSEncrypt from '@/utils/jsencrypt'

const statusMap = {
  'true': {
    status: 'success',
    text: '激活'
  },
  'false': {
    status: 'error',
    text: '禁用'
  }
}

export default {
  name: 'UserList',
  data () {
    return {
      labelCol: { span: 4 },
      wrapperCol: { span: 20 },
      pubkey: undefined,
      columns: [
        {
          title: '用户名',
          dataIndex: 'username'
        },
        {
          title: '状态',
          dataIndex: 'enable',
          scopedSlots: { customRender: 'status' }
        },
        {
          title: '昵称',
          dataIndex: 'nickname'
        },
        {
          title: '管理员',
          dataIndex: 'is_admin',
          scopedSlots: { customRender: 'isAdmin' }
        },
        {
          title: '创建时间',
          dataIndex: 'createTime'
        },
        {
          title: '更新时间',
          dataIndex: 'updateTime'
        },
        {
          title: '操作',
          dataIndex: 'operation',
          scopedSlots: { customRender: 'operation' }
        }
      ],
      // 查询参数
      queryParam: {},
      loadData: parameter => {
        console.log(parameter)
        const requestParameters = Object.assign({}, parameter, { 'filters': this.queryParam })
        console.log('loadData request parameters:', requestParameters)
        return getUserList(requestParameters)
          .then(res => {
            console.log(res)
            return res.info
          })
      },
      projectList: [],
      selectedRowKeys: [],
      selectedRows: [],
      // 添加用户
      addUserModalVisible: false,
      addUserLoading: false,
      userForm: {
        username: undefined,
        nickname: undefined,
        password: undefined,
        email: undefined,
        enable: true,
        is_admin: false
      },
      rules: {
        username: [{ required: true, message: '用户名必填', trigger: 'blur' }],
        password: [{ required: true, message: '密码必填', trigger: 'blur' }]
      }
    }
  },
  components: {
    STable,
    Ellipsis
  },
  computed: {
    rowSelection () {
      return {
        selectedRowKeys: this.selectedRowKeys,
        onChange: this.onSelectChange
      }
    }
  },
  filters: {
    statusFilter (type) {
      return statusMap[String(type)].text
    },
    statusTypeFilter (type) {
      return statusMap[String(type)].status
    }
  },
  created () {
    getPublicKey().then(resp => {
      this.pubkey = resp.info
    }).catch(err => {
      this.$notification.error({
        message: '错误',
        description: `获取通信密钥失败: ${err}`
      })
    })
  },
  methods: {
    refresh () {
      this.$refs.table.refresh()
    },
    onSelectChange (selectedRowKeys, selectedRows) {
      this.selectedRowKeys = selectedRowKeys
      this.selectedRows = selectedRows
    },
    handleAddUser () {
      this.addUserModalVisible = true
      this.addUserLoading = false
    },
    handleBatchOpClick (e) {
      console.log(e)
    },
    handleEditUser (record) {

    },
    createUser () {
      this.$refs.userForm.validate(valid => {
        if (valid) {
          if (this.pubkey === undefined) {
            this.$message.error('没有通信密钥，请刷新后重试')
            return
          }
          this.addUserLoading = true
          const postData = Object.assign({}, this.userForm)
          var encrypt = new JSEncrypt()
          encrypt.setPublicKey(this.pubkey)
          var encrypted = encrypt.encrypt(this.userForm.password)
          console.log(encrypted)
          postData.password = encrypted
          createUserList(postData).then(resp => {
            console.log(resp)
            if (resp.code !== 0) {
              this.$message.error(`创建用户失败: ${resp.info}`)
            } else {
              this.$message.info('创建成功')
              this.addUserModalVisible = false
              this.refresh()
            }
            this.addUserLoading = false
          }).catch(err => {
            this.$message.error(`创建用户失败: ${err}`)
            this.addUserLoading = false
          })
        }
      })
    }
  }
}
</script>
