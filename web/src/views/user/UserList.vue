<template>
  <a-card :bordered="false">
    <!-- <a-form layout="inline">

    </a-form> -->

    <div class="table-operator">
      <a-button type="primary" icon="plus" @click="handleAddProject">新建</a-button>
      <a-button icon="reload" @click="refresh">刷新</a-button>
      <a-dropdown v-action:edit v-if="selectedRowKeys.length > 0" :trigger="['click']">
        <a-menu slot="overlay" @click="e => handleBatchOpClick(e)">
          <a-menu-item key="1"><a-icon type="delete" />删除</a-menu-item>
          <!-- lock | unlock -->
          <a-menu-item key="2"><a-icon type="lock" />锁定</a-menu-item>
        </a-menu>
        <a-button style="margin-left: 8px">
          批量操作 <a-icon type="down" />
        </a-button>
      </a-dropdown>
    </div>
    <s-table
      ref="table"
      size="default"
      rowKey="id"
      :columns="columns"
      :data="loadData"
      :alert="true"
      :rowSelection="rowSelection"
      showPagination="auto"
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
    </s-table>
  </a-card>
</template>

<script>
import { STable, Ellipsis } from '@/components'
import { getUserList } from '@/api/user'

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
      selectedRows: []
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
  methods: {
    refresh () {
      this.$refs.table.refresh()
    },
    onSelectChange (selectedRowKeys, selectedRows) {
      this.selectedRowKeys = selectedRowKeys
      this.selectedRows = selectedRows
    },
    handleAddProject () {

    },
    handleBatchOpClick (e) {
      console.log(e)
    }
  }
}
</script>
