<template>
  <div>
    <h2>
      Hi, {{this.fullName}}
      <br />欢迎加入笑来课堂
    </h2>
    <div class="action">
      <van-field v-model="code" placeholder="邀请码" autosize />
      <van-button class="button" type="info" size="small" @click="apply">入群验证</van-button>
    </div>
  </div>
</template>

<script>
export default {
  name: "InvitationEntry",

  data() {
    return {
      meInfo: null,
      code: '',
    };
  },

  computed: {
    fullName() {
      if (this.meInfo) {
        return this.meInfo.data.full_name
      } else {
        return ""
      }
    }
  },

  async mounted() {
    this.meInfo = await this.GLOBAL.api.account.me()
  },

  methods: {
    apply() {
      this.GLOBAL.api.invitation.apply(this.code)
    }
  },
};
</script>

<style lang="scss" scoped>
h2 {
  padding-top: 10rem;
  text-align: center;
}

.action {
  margin-left: 2rem;
  margin-right: 2rem;
}

.button {
  margin-top: 2rem;
  width: 100%;
  border-radius: 4px;
}
</style>

