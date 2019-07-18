<template>
  <div class="card">
    <van-cell :title="invitation.code" :value="status" :label="label" value-class="{unused: true}" @click="copy" />
  </div>
</template>

<script>
export default {
  name: "InvitationCodeListItem",
  props: ["invitation"],
  computed: {
    status() {
      if (this.invitation.is_used) {
        if (this.invitation.invitee) {
          return "已使用"
        } else {
          return "待付费"
        }
      } else {
        return "未使用"
      }
    },

    label() {
      if (this.invitation.invitee) {
        return this.invitation.invitee.full_name + "已用"
      } else {
        return "点击复制"
      }
    }
  },
  methods: {
    copy() {
      this.$copyText(this.invitation.code).then(function (e) {
        alert('Copied')
      })
    }
  },
};
</script>

<style lang="scss" scoped>
.van-cell {
  padding: 0rem;
}

.card {
  background-color: #ffffff;
  border-radius: 4px;
  box-shadow: 1px 1px 1px 1px rgba(0, 0, 0, 0.1);
}

.unused {
  color: #52C41A;
}
</style>