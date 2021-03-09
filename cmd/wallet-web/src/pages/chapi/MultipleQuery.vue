/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div v-if="loading" style="margin-left: 40%;margin-top: 20%;height: 200px;">
        <div class="md-layout">
            <md-progress-spinner :md-diameter="100" class="md-accent" :md-stroke="10"
                                 md-mode="indeterminate"></md-progress-spinner>
        </div>
    </div>
    <div v-else class="md-layout">
        <div class="md-layout-item md-medium-size-100 md-xsmall-size-100 md-size-100">
            <div class="md-layout md-alignment-center-center" style="margin-top: 20px;">
                <div class="md-headline">Credential Presentation Requested</div>
                <div class="md-subheading">A credential presentation is been requested:</div>
            </div>

            <div style="margin: 10px"></div>

            <div class="md-layout md-alignment-center-center md-subheading">
                <div>By</div>
            </div>

            <div class="md-layout md-alignment-center-center md-subheading" style="color: #025C8F;">
                <div><md-icon style="color: #00cc66;" class="md-size-1x">verified_user</md-icon>{{ requestOrigin }}</div>
            </div>

            <div style="margin: 20px"></div>

            <div v-if="errors.length" >
                <b>Failed with following error(s):</b>
                <md-field style="margin-top: -15px">
                    <ul>
                        <li v-for="error in errors" :key="error">{{ error }}</li>
                    </ul>
                </md-field>

                <md-button v-on:click="cancel" style="margin-left: 40%; background-color: #9d0006 !important;" class="md-cancel-text" id="cancelBtnNrc">
                    Cancel
                </md-button>
            </div>


            <md-card v-for="(record, key) in records" :key="key" md-with-hover>
                <md-card-header>
                    <md-avatar>
                        <md-icon class="md-size-2x"></md-icon>
                        <md-icon class="md-size-2x">{{getVCIcon( record.credential.type)}}</md-icon>
                    </md-avatar>

                    <div class="md-title">{{record.credential.name ? record.credential.name : getCredentialType(record.credential.type)}}</div>
                    <div class="md-subhead">{{record.credential.description}}</div>
                </md-card-header>

                <md-card-content>
                    <span class="md-caption">
                        <md-avatar>
                            <md-icon>info</md-icon>
                        </md-avatar>
                        <b>For the purpose of:</b> {{record.reason}}
                    </span>

                    <div v-if="record.output" class="md-caption" >
                        The verifier will only access below information from your <span class="md-body-1">{{record.credential.name ? record.credential.name  : getCredentialType(record.credential.type)}}</span>
                        <div style="margin: 10px"></div>
                        <div v-for="(subj, skey) in record.output.credentialSubject" :key="skey">
                            <div class="md-caption" style="padding-left: 35%" v-if="displayContent(skey)">
                                {{skey}} : {{subj}}
                            </div>
                        </div>
                    </div>
                </md-card-content>
            </md-card>


            <div v-if="showShareCredential" class="md-layout md-alignment-center-center" >
                <p class="md-body-1">By clicking Agree you will be sharing a unique identifier to <b style="color: #2E7D32">{{requestOrigin}}</b>, the Credential content, and your digital signature.
                    <a href="https://www.w3.org/TR/vc-data-model/#proofs-signatures" target="_blank">Learn more</a></p>

                 <md-button v-on:click="share" class="md-button md-info md-square md-theme-default md-large-size-100 md-size-100 col"
                            style="background-color: #29a329 !important;" id="share-credentials" >
                     Agree
                 </md-button>
                 <md-button v-on:click="cancel" style="margin-left: 5px; background-color: #9d0006 !important;" class="md-cancel-text" id="cancelBtn">
                     Cancel
                 </md-button>
            </div>

        </div>
    </div>
</template>
<script>

    import {MultipleQuery} from "./wallet"
    import {mapGetters} from 'vuex'

    const nonDisplayContent = ['id', 'type']
    const allIcons= ['account_box', 'contacts', 'person', 'person_outline', 'card_membership', 'portrait', 'bento']
    const vcIcons = {
        PermanentResidentCard:'perm_identity',
        UniversityDegreeCredential:'school',
        FlightBookingRef: 'flight',
        VaccinationCertificate: 'health_and_safety',
    }

    const getIcon = (type) =>  vcIcons[type] ? vcIcons[type] : allIcons[Math.floor(Math.random() * Math.floor(allIcons.length))]


    export default {
        errorCaptured(error) {
            console.log('An error has occurred!', error);
        },
        created: async function () {
            this.requestOrigin = this.$parent.credentialEvent.credentialRequestOrigin

            try {
                this.wallet = new MultipleQuery(this.getAgentInstance(), this.$parent.credentialEvent)
            } catch (e) {
                this.handleError(e)
                return
            }

            await this.loadCredentials()
            this.loading = false
        },
        data() {
            return {
                errors: [],
                requestOrigin: "",
                loading: true,
                credentialWarning: "",
                reason: "",
                allIcons: ['account_box', 'contacts', 'person', 'person_outline', 'card_membership', 'portrait', 'bento'],
                selectedFrames: [],
                decisionMade: []
            };
        },
        methods: {
            ...mapGetters('agent', {getAgentInstance: 'getInstance'}),
            ...mapGetters(['getCurrentUser']),
            handleError(e) {
                this.errors.push(e)
                this.loading = false
            },
            loadCredentials: async function () {
                try {
                    this.records = await this.wallet.queryCredentials()
                } catch (err) {
                    this.errors.push('failed to query your saved credentials.')
                    console.error('get credentials failed, error:', err)
                    return
                }

                if (this.records.length == 0) {
                    this.errors.push('No credentials found matching requested criteria')
                }

                console.log(`found ${this.records.length}`, this.records)
            },
            async share() {
                this.loading = true
                this.wallet.generatePresentation(this.getCurrentUser().username, this.records)
            },
            cancel() {
                this.wallet.cancel()
            },
            getVCIcon(types) {
                return getIcon(this.getCredentialType(types))
            },
            getCredentialType(types) {
                return types.filter(type => type != "VerifiableCredential")[0]
            },
            displayContent(k) {
                return !nonDisplayContent.includes(k)
            },
        },
        computed: {
            showShareCredential() {
                return this.records.length > 0;
            },
        },
    }
</script>
