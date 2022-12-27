<template>
  <a-card :bordered="false">
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

      <!-- <span slot="action" slot-scope="text, record">
        <template>
          <a @click="handleEdit(record)">配置</a>
          <a-divider type="vertical" />
          <a @click="handleSub(record)">订阅报警</a>
        </template>
      </span> -->
    </s-table>
  </a-card>
</template>

<script>
import { STable, Ellipsis } from '@/components'
import { getProjectList } from '@/api/project'

const statusMap = {
  'active': {
    status: 'success',
    text: '激活'
  },
  'archive': {
    status: 'error',
    text: '归档'
  }
}

export default {
  name: 'ProjectList',
  data () {
    return {
      columns: [
        {
          title: '项目名',
          dataIndex: 'name'
        },
        {
          title: '状态',
          dataIndex: 'status',
          scopedSlots: { customRender: 'status' }
        },
        {
          title: '描述',
          dataIndex: 'desc',
          scopedSlots: { customRender: 'desc' }
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
        return getProjectList(requestParameters)
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
      return statusMap[type].text
    },
    statusTypeFilter (type) {
      return statusMap[type].status
    }
  },
  methods: {
    onSelectChange (selectedRowKeys, selectedRows) {
      this.selectedRowKeys = selectedRowKeys
      this.selectedRows = selectedRows
    }
  }
}
</script>
