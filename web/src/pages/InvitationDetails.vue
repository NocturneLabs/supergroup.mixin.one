<template>
  <div>
    <div class="layout header">
      <h2>我的邀请</h2>
      <p
        class="subtitle"
      >成功邀请好友：{{this.usedInvitations.length}} 待付费好友：{{this.pendingInvitations.length}}</p>
    </div>
    <van-tabs  v-model="activeName">
      <van-tab title="邀请码" name="codes">
        <van-list v-model="loading" :finished="finished">
          <InvitationCodeListItem class="layout" v-for="item in list" :key="item.code" :invitation="item" />
        </van-list>
        <div class="layout">
          <p class="subtitle">
            * 未付费：被邀请用户已使用邀请码但未成功付费
            <br />* 已使用：被邀请用户已成功使用邀请码并付费
            <br />* 当以上三个邀请码状态均为‘已使用’状态时才可申请新的一组邀请码
          </p>
          <van-button
            @disabled="applyDisabled"
            class="button"
            type="info"
            size="small"
            @click="apply"
          >申请</van-button>
        </div>
      </van-tab>
      <van-tab title="受邀者" name="invitees">
        <van-list v-model="loading" :finished="finished">
          <InvitationInviteeListItem class="layout" v-for="item in list" :key="item.code" :invitation="item" />
        </van-list>
      </van-tab>
    </van-tabs>
  </div>
</template>

<script>
import InvitationCodeListItem from '../components/Invitation/InvitationCodeListItem'
import InvitationInviteeListItem from '../components/Invitation/InvitationInviteeListItem'
export default {
  name: "InvitationDetails",

  data() {
    return {
      activeName: "codes",
      list: [
    {
      "type": "Invitation",
      "code": "bknia7a3q5626nglhs6g",
      "invitee": {
        "user_id": "9643f9b0-1c26-44cd-9785-a54c0585185c",
        "full_name": "Reason",
        "avatar_url": "https://images.mixin.one/vGj4xQirkdUg6C6zw6b49MmZPdSvxJc5N5geS1xaV-c7W5BblFq2qUq3aDAHCJqdDP_umRsMBoRiaIoBEesL0zM=s256",
        "state": "pending"
      },
      "is_used": true,
      "created_at": "2019-07-17T21:39:41.065295+08:00"
    },
    {
      "type": "Invitation",
      "code": "bknia7a3q5626nglhs60",
      "invitee": null,
      "is_used": false,
      "created_at": "2019-07-17T21:39:41.065279+08:00"
    },
    {
      "type": "Invitation",
      "code": "bknia7a3q5626nglhs5g",
      "invitee": null,
      "is_used": false,
      "created_at": "2019-07-17T21:39:41.065201+08:00"
    }
  ]
,
      loading: false,
      finished: true,
      invitationsHistory: [],
      invitationsCurrent: [],
    };
  },

  mounted() {
    this.GLOBAL.api.invitation.index(false).then((response) => {
      this.invitationsCurrent = response.data
    })
    this.GLOBAL.api.invitation.index(true).then((response) => {
      this.invitationsHistory = response.data
    })
  },

  computed: {
    usedInvitations() {
      return this.invitationsHistory
    },
    pendingInvitations() {
      return [];
    },
    unusedInvitations() {
      return [];
    },
    applyDisabled() {
      return false;
    }
  },

  components: {
    InvitationCodeListItem,
    InvitationInviteeListItem,
  },

  methods: {
    apply() {
      this.GLOBAL.api.invitation.create()
    }
  }
};
</script>

<style lang="scss" scoped>
.layout {
  margin: 1rem;
  padding: 1rem;
}

.header {
  margin-bottom: 0rem;
}

.button {
  width: 100%;
}

.subtitle {
  line-height: 22px;
  font-size: 14px;
  text-align: left;
  color: rgba(0, 0, 0, 0.25);
}
</style>

