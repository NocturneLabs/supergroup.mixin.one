<template>
  <loading :loading="loading" :fullscreen="true">
  <div class="pay-page" style="padding-top: 60px;">
    <nav-bar :title="$t('pay.title')" :hasTopRight="false"></nav-bar>
    <van-panel :title="$t('pay.welcome')" :desc="$t('pay.welcome_desc')">
    </van-panel>
    <br/>
    <van-panel :title="$t('pay.method_crypto')">
      <row-select
        :index="0"
        :title="$t('pay.select_assets')"
        :columns="assets"
        placeholder="Tap to Select"
        @change="onChangeAsset">
        <span slot="text">{{selectedAsset ? selectedAsset.text : 'Tap to Select'}}</span>
      </row-select>
      <van-cell
        :title="$t('pay.price_label', {price: currentCryptoPrice, unit: selectedAsset ? selectedAsset.text : '...'})"
        >
        <!-- <span>≈{{currencySymbol}}{{currentEstimatedPrice.toLocaleString()}}</span> -->
      </van-cell>
      <div slot="footer">
        <van-cell>
          <van-button style="width: 100%" type="info" :disabled="selectedAsset === null || loading" @click="payCrypto">{{$t('pay.pay_crypto')}}</van-button>
        </van-cell>
        <!-- <van-cell>
          <van-button style="width: 100%" type="warning" :disabled="selectedAsset === null" @click="payCrypto">{{$t('pay.pay_foxone')}}</van-button>
        </van-cell> -->
      </div>
    </van-panel>
    <br/>
    <van-panel v-if="acceptWechatPayment" :title="$t('pay.method_wechat')">
      <van-cell
        :title="$t('pay.price_label', {price: wechatPaymentAmount, unit: $t('currency.' + autoEstimateCurrency)})"
        >
      </van-cell>
      <div slot="footer">
        <van-cell>
          <van-button style="width: 100%" type="primary" @click="payWechatMobile">{{$t('pay.pay_wechat')}}</van-button>
        </van-cell>
      </div>
    </van-panel>
    <br/>
    <van-panel v-if="acceptCouponPayment" :title="$t('pay.method_coupon')">
      <van-cell>
        <van-field :placeholder="$t('pay.coupon_placeholder')" v-model="couponCode"></van-field>
      </van-cell>
      <div slot="footer">
        <van-cell>
          <van-button style="width: 100%" type="info" @click="payCoupon" :disabled="loading">{{$t('pay.pay_coupon')}}</van-button>
        </van-cell>
      </div>
    </van-panel>
  </div>
  </loading>
</template>

<script>
import NavBar from '@/components/Nav'
import RowSelect from '@/components/RowSelect'
import Loading from '../components/Loading'
import Row from '@/components/Nav'
import uuid from 'uuid'
import {Toast} from 'vant'
import { CLIENT_ID, WEB_ROOT } from '@/constants'

export default {
  name: 'Pay',
  props: {
    msg: String
  },
  data () {
    return {
      loading: true,
      config: null,
      meInfo: null,
      selectedAsset: null,
      autoEstimate: false,
      autoEstimateCurrency: 'usd',
      acceptWechatPayment: false,
      acceptCouponPayment: false,
      wechatPaymentAmount: '100',
      cryptoEsitmatedUsdMap: {},
      currencyTickers: [],
      cnyRatio: {},
      currentCryptoPrice: 0,
      currentEstimatedPrice: 0,
      couponCode: '',
      assets: []
    }
  },
  components: {
    NavBar, RowSelect, Loading
  },
  async mounted () {
    this.loading = true
    let config = await this.GLOBAL.api.website.config()
    this.assets = config.data.accept_asset_list.map((x) => {
      x.text = x.symbol
      return x
    })
    this.selectedAsset = this.assets[0]
    this.autoEstimate = config.data.auto_estimate
    this.autoEstimateCurrency = config.data.auto_estimate_currency
    this.autoEstimateBase = config.data.auto_estimate_base
    this.acceptWechatPayment = config.data.accept_wechat_payment
    this.wechatPaymentAmount = config.data.wechat_payment_amount
    this.acceptCouponPayment = config.data.accept_coupon_payment
    this.GLOBAL.api.fox.currency()
      .then((currencyInfo) => {
        this.currencyTickers = currencyInfo.data.cnyTickers.reduce((map, obj) => {
          map[obj.from] = obj.price;
          return map;
        }, {})
        this.cnyRatio = currencyInfo.data.currencies
        // console.log(this.currencyTickers)
      })
    this.meInfo = await this.GLOBAL.api.account.me()
    setTimeout(this.updatePrice, 2000)
    this.loading = false
  },
  computed: {
    currencySymbol() {
      if (this.autoEstimate) {
        if (this.autoEstimateCurrency === 'cny') return '¥'
        if (this.autoEstimateCurrency === 'usd') return '$'
      }
      return ''
    }
  },
  methods: {
    payCrypto () {
      this.loading = true
      let traceId = this.meInfo.data.trace_id
      setTimeout(async () => { await this.waitForPayment(); }, 1000)
      window.location.href = `mixin://pay?recipient=${CLIENT_ID}&asset=${this.selectedAsset.asset_id}&amount=${this.currentCryptoPrice}&trace=${traceId}&memo=PAY_TO_JOIN`
      // console.log(`mixin://pay?recipient=${CLIENT_ID}&asset=${this.selectedAsset.asset_id}&amount=${this.currentCryptoPrice}&trace=${traceId}&memo=PAY_TO_JOIN`);
    },
    async onChangeAsset (ix) {
      this.loading = true
      this.selectedAsset = this.assets[ix]
      await this.updatePrice()
      this.loading = false
    },
    async updatePrice () {
      if (this.selectedAsset.amount === 'auto') {
        let base = this.autoEstimateBase / parseFloat(this.cnyRatio.usdt)
        let priceUsdt = await this.getCryptoEsitmatedUsdt(this.selectedAsset.symbol)
        this.currentCryptoPrice = (base / priceUsdt).toFixed(8)
        if (this.autoEstimateCurrency === 'usd') {
          this.currentEstimatedPrice = base
        } else {
          this.currentEstimatedPrice = base * this.cnyRatio.usdt
        }
      } else {
        this.currentCryptoPrice = parseFloat(this.selectedAsset.amount).toFixed(8)
        this.currentEstimatedPrice = '-'
      }
    },
    async waitForPayment () {
      let meInfo = await this.GLOBAL.api.account.me()
      if (meInfo.error) {
        setTimeout(async () => { await this.waitForPayment(); }, 1500)
        return;
      }
      if (meInfo.data.state === 'paid') {
        Toast(this.$t('pay.success_toast'))
        this.$router.push('/');
        this.loading = false
        return;
      }
      setTimeout(async () => { await this.waitForPayment(); }, 1500)
    },
    async getCryptoEsitmatedUsdt (symbol) {
      if (this.cryptoEsitmatedUsdMap.hasOwnProperty(symbol)) {
        return this.cryptoEsitmatedUsdMap[symbol]
      }
      // only support fetching from big.one
      const pairName = symbol + '-' + 'USDT'
      let resp = await this.GLOBAL.api.fox.b1Ticker(pairName)
      if (resp && resp.data) {
        this.cryptoEsitmatedUsdMap[symbol] = resp.data.close
        return resp.data.close
      }
      return -1
    },
    payWechatMobile () {
      this.$router.push(`/pay/wxqr/?qr_url=${encodeURIComponent(WEB_ROOT + '/wechat/request/' + this.meInfo.data.user_id)}`)
    },
    async payCoupon () {
      this.loading = true
      try {
        let resp = await this.GLOBAL.api.coupon.occupy(this.couponCode)
        if (resp && resp.data) {
          Toast(this.$t('pay.correct_coupon_code_toast'))
          this.$router.push('/')
          this.loading = false
        } else {
          Toast(this.$t('pay.incorrect_coupon_code_toast'))
          this.loading = false
        }
      } catch (err) {
        Toast(this.$t('pay.incorrect_coupon_code_toast'))
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.pay-page {
  padding-top: 60px;
}
</style>
