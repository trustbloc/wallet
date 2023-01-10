<!--
 * Copyright Avast Software. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
-->

<script setup>
import { computed, onMounted, ref, watchEffect } from 'vue';
import { useRoute } from 'vue-router';
import { useStore } from 'vuex';
import { useI18n } from 'vue-i18n';
import { CollectionManager, CredentialManager, DIDManager, OpenID4CI } from '@trustbloc/wallet-sdk';
import { verifiableDataFormatCode } from '@/mixins';
import CustomSelectComponent from '@/components/CustomSelect/CustomSelectComponent.vue';
import StyledButtonComponent from '@/components/StyledButton/StyledButtonComponent.vue';
import CredentialOverviewComponent from '@/components/WACI/CredentialOverviewComponent.vue';
import WACIActionButtonsContainerComponent from '@/components/WACI/WACIActionButtonsContainerComponent.vue';
import WACIErrorComponent from '@/components/WACI/WACIErrorComponent.vue';
import WACILoadingComponent from '@/components/WACI/WACILoadingComponent.vue';
import WACISuccessComponent from '@/components/WACI/WACISuccessComponent.vue';
import CredentialDetailsTableComponent from '@/components/WACI/CredentialDetailsTableComponent.vue';

// Hooks
const store = useStore();
const route = useRoute();
const { t } = useI18n();

// Local Variables
const loading = ref(true);
const saving = ref(false);
const errors = ref([]);
const vaults = ref([]);
const selectedVault = ref('');
const processedCredentials = ref([]);
const savedSuccessfully = ref(false);
const token = ref(null);
const credentialManager = ref(null);
const collectionManager = ref(null);
const didManager = ref(null);
const openID4CI = ref(null);
const pinEntryRequired = ref(false);
const pin = ref(null);
const authorizeResp = ref();
const showMainState = computed(
  () =>
    !loading.value &&
    !errors.value.length &&
    !saving.value &&
    !savedSuccessfully.value &&
    !pinEntryRequired.value
);
const successButtonLabel = computed(() => t('WACI.Issue.viewCredential'));

// Store Getters
const currentUser = computed(() => store.getters['getCurrentUser']);
const agentInstance = computed(() => store.getters['agent/getInstance']);

// Methods
function setPinEntryRequired(value) {
  pinEntryRequired.value = value;
}

// Fetches all vaults and sets selected vault to default
async function fetchAllVaults() {
  const { contents } = await collectionManager.value.getAll(token.value);
  vaults.value = Object.values(contents).map((vault) => vault);
  // Default vault is selected vault by default, it is created on wallet setup and must be only one.
  selectedVault.value = vaults.value.find((vault) => vault.name === 'Default Vault').id;
}

onMounted(async () => {
  const { profile, preference } = currentUser.value;
  const { user } = profile;
  token.value = profile.token;

  const requestUrl = new URL(route.query.url);

  const issuer = requestUrl.searchParams.get('issuer');
  const credential_type = requestUrl.searchParams.get('credential_type');
  const preAuthCode = requestUrl.searchParams.get('pre-authorized_code');
  const op_state = requestUrl.searchParams.get('op_state');
  const user_pin_required = requestUrl.searchParams.get('user_pin_required') === 'true';
  pinEntryRequired.value = user_pin_required;

  credentialManager.value = new CredentialManager({ agent: agentInstance.value, user });
  collectionManager.value = new CollectionManager({ agent: agentInstance.value, user });
  didManager.value = new DIDManager({ agent: agentInstance.value, user });
  const { contents } = await didManager.value.getAllDIDs(token.value);
  const kid = Object.values(contents)[0].didDocument.verificationMethod[0].id;
  openID4CI.value = new OpenID4CI({
    agent: agentInstance.value,
    user,
    clientConfig: { userDID: kid, clientID: 'test-client-id' },
  });

  await fetchAllVaults();

  const req = {};
  req.issuer = issuer;
  req.credential_type = credential_type;
  req['pre-authorized_code'] = preAuthCode;
  req.user_pin_required = user_pin_required;
  req.format = verifiableDataFormatCode(preference.proofFormat);
  if (op_state) req.op_state = op_state;

  watchEffect(async () => {
    if (pin.value && !pinEntryRequired.value) {
      authorizeResp.value = await openID4CI.value.authorize(token.value, kid, req, pin.value);
    }
  });

  loading.value = false;
});
</script>

<template>
  <div v-if="!showMainState" class="flex h-full w-full grow flex-col items-center justify-center">
    <!-- Loading State -->
    <WACILoadingComponent v-if="loading" />

    <!-- Saving State -->
    <WACILoadingComponent v-else-if="saving" :message="t('WACI.Issue.savingCredential')" />

    <!-- Error State -->
    <WACIErrorComponent v-else-if="errors.length" @click="cancel" />

    <!-- Pin Prompt State -->
    <div
      v-else-if="pinEntryRequired"
      class="flex h-full w-full grow flex-col items-center justify-center"
    >
      <h4>Please, enter a PIN</h4>
      <div class="input-container mb-0 mt-5 w-auto pr-4">
        <input
          id="pin-input"
          v-model="pin"
          type="text"
          inputmode="decimal"
          required
          pattern="^\d{6}$"
          minlength="6"
          maxlength="6"
          autocomplete="off"
          placeholder="6-digit PIN"
          size="6"
        />
        <label for="pin-input" class="input-label">PIN</label>
      </div>
      <StyledButtonComponent
        id="enterPinButton"
        class="mt-5"
        type="btn-primary"
        @click="setPinEntryRequired(false)"
      >
        Submit
      </StyledButtonComponent>
    </div>

    <!-- Success State -->
    <WACISuccessComponent
      v-else-if="savedSuccessfully"
      id="issue-credentials-ok-btn"
      :heading="t('WACI.Issue.success', processedCredentials.length)"
      :message="
        t('WACI.Issue.message', {
          subject: processedCredentials[0].title,
        })
      "
      :button-label="successButtonLabel"
      @click="finish"
    />
  </div>

  <!-- Main State -->
  <div v-else class="flex h-full w-full grow flex-col items-center justify-between overflow-hidden">
    <div class="flex w-full justify-center overflow-auto">
      <div
        class="flex h-full w-full max-w-3xl grow flex-col items-start justify-start py-8 px-5 md:px-0"
      >
        <span class="mb-6 text-3xl font-bold">{{ t('WACI.Issue.saveCredential') }}</span>

        <div
          v-for="(credential, index) in processedCredentials"
          :key="index"
          class="flex w-full max-w-3xl flex-col justify-start"
        >
          <CredentialOverviewComponent :credential="credential">
            <template #bannerBottomContainer>
              <div
                class="mt-5 flex w-full grow flex-col items-start justify-start rounded-t-lg border-b border-neutrals-dark bg-neutrals-lilacSoft px-4"
              >
                <label for="select-key" class="mb-1 text-sm font-bold text-neutrals-dark">{{
                  t('Vaults.selectVault')
                }}</label>
                <CustomSelectComponent
                  id="waci-issue-select-vault"
                  :options="vaults"
                  default="Default Vault"
                  @selected="setSelectedVault"
                />
              </div>
            </template>
            <template #credentialDetails>
              <CredentialDetailsTableComponent
                :heading="t('WACI.Issue.verifiedInformation')"
                :credential="processedCredentials[0]"
                class="mt-8"
              />
            </template>
          </CredentialOverviewComponent>
        </div>
      </div>
    </div>

    <WACIActionButtonsContainerComponent>
      <template #leftButton>
        <StyledButtonComponent id="cancelBtn" type="btn-outline" @click="cancel">
          {{ t('CHAPI.Share.decline') }}
        </StyledButtonComponent>
      </template>
      <template #rightButton>
        <StyledButtonComponent id="storeVCBtn" type="btn-primary" @click="save">
          {{ t('WACI.Issue.save') }}
        </StyledButtonComponent>
      </template>
    </WACIActionButtonsContainerComponent>
  </div>
</template>

<style scoped>
.chapi-container {
  width: 28rem;
}

.credentialPreviewContainer:not(:focus-within):hover {
  box-shadow: 0px 4px 12px 0px rgba(25, 12, 33, 0.1);
}

.footerContainer {
  box-shadow: inset 0px 1px 0px 0px rgb(219, 215, 220);
}
</style>
