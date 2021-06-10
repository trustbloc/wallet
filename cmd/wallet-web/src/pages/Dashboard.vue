/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
    <div class="content">
        <div class="content">
            <div class="md-layout">
                <div>
                    <md-label style="color: #1B5E20; font-size: 16px; margin: 5px">
                        <span v-if="loadingStatus === 'inprogress'">
                            <pulse-loader :color="'green'" :size="5"
                                          style="float: left; margin-right: 5px"></pulse-loader> Setting up your user for secured communication.
                        </span>
                        <span v-else-if="loadingStatus === 'success'" id="dashboard-success-msg" class="px-24">
                             <md-icon style="color: green" class="px-4">check_circle_outline</md-icon> Successfully setup your user for secured communication.
                        </span>
                        <span v-else-if="loadingStatus === 'failed'">
                            <md-icon style="color: red;" class="px-4">warning</md-icon>
                            <b>Warning:</b> Failed to connect to server. Your wallet can not participate in secured communication.
                        </span>
                    </md-label>
                    <div class="md-card px-24" md-with-hover v-if="verifiableCredentials.length">
                        <md-card-content>
                            <ul class="credential-list">
                                <li v-on:click="toggleCard(card)" v-for="(card, index) in cards" :key="index">
                                    <transition name="flip">
                                        <div v-if="!card.flipped" v-bind:key="card.flipped" class="card">
                                            <div class="cardContent">
                                                <div class="cardHeader">
                                                    {{credDisplayName(card)}}
                                                </div>

                                <university-card
                                        v-if="credDisplayName(card) === 'Bachelor Degree' || credDisplayName(card) === 'University Degree Credential'"
                                        :item="card"
                                />
                                <permanent-resident-card
                                        v-else-if="credDisplayName(card) === 'Permanent Resident Card'"
                                        :item="card"
                                />
                                <travel-card
                                        v-else-if="credDisplayName(card) === 'Travel Card'"
                                        :item="card"
                                />
                                <student-card
                                        v-else-if="credDisplayName(card) === 'Student Card'"
                                        :item="card"
                                />
                                <drivers-license
                                        v-else-if="credDisplayName(card) === 'Drivers License'"
                                        :item="card"
                                />
                                <crude-product-card
                                        v-else-if="credDisplayName(card) === 'Heavy Sour Dilbit' || credDisplayName(card) === 'Crude Product Credential'"
                                        :item="card"
                                />
                                <mill-test-card
                                        v-else-if="credDisplayName(card) ===  'Steel Inc. CMTR' || credDisplayName(card) === 'Certified Mill Test Report'"
                                        :item="card"
                                />
                                <general-card v-else :item="card"/>

                            </div>
                                <json-modal :item="verifiableCredentials[index]" />
                            </div>
                                <div v-else v-bind:key="card.flipped" class="card">
                                    <div class="cardContent cardBack">
                                        <p> Issuance Date: {{ card.credentialSubject.issue_date || card.credentialSubject.issuedate || card.issuanceDate || 'N/A' }} </p>
                                        <p> Expiration Date: {{ card.credentialSubject.expiry_date || card.credentialSubject.cardexpires || card.expirationDate || 'N/A' }} </p>
                                    </div>
                                </div>
                            </transition>
                            </li>
                            </ul>
                        </md-card-content>
                    </div>
                    <md-empty-state v-else
                                    md-icon="devices_other"
                                    :md-label=error
                                    :md-description=errorDescription>
                    </md-empty-state>
                </div>

            </div>
        </div>
    </div>

</template>

<script>
    import {filterCredentialsByType, getCredentialType} from "@/pages/chapi/wallet";
    import {mapActions, mapGetters} from 'vuex'
    import {PulseLoader} from "@saeris/vue-spinners";
    import PermanentResidentCard from "../components/CredentialCards/PermanentResidentCard";
    import UniversityCard from "../components/CredentialCards/UniversityCard";
    import TravelCard from "../components/CredentialCards/TravelCard";
    import StudentCard from "../components/CredentialCards/StudentCard";
    import DriversLicense from "../components/CredentialCards/DriversLicense";
    import CrudeProductCard from "../components/CredentialCards/CrudeProductCard";
    import MillTestCard from "../components/CredentialCards/MillTestCard";
    import GeneralCard from "../components/CredentialCards/GeneralCard";
    import JsonModal from "../components/CredentialCards/JsonModal";

    const manifestCredType = "IssuerManifestCredential"
    const governanceCredType = "GovernanceCredential"
    // TODO: issue-627 Add generic vue card for all the credentials to dynamically add support for all VC types.
    export default {
        components: {
            PulseLoader,
            PermanentResidentCard,
            UniversityCard,
            TravelCard,
            StudentCard,
            DriversLicense,
            CrudeProductCard,
            MillTestCard,
            GeneralCard,
            JsonModal
        },
        created: async function () {
            // Load the Credentials
            await this.getCredentials()
            await this.fetchAllCredentials()
            await this.refreshUserMetadata()

            this.username = this.getCurrentUser().username
        },
        methods: {
            ...mapGetters('agent', {getAgentInstance: 'getInstance'}),
            ...mapGetters(['getCurrentUser', 'allCredentials', 'getAgentOpts']),
            ...mapActions(['getCredentials', 'refreshUserMetadata']),
            fetchAllCredentials: async function () {
                this.verifiableCredentials = []
                try {
                    for (let c of filterCredentialsByType(this.allCredentials(), [manifestCredType, governanceCredType])) {
                        let resp = await this.getAgentInstance().verifiable.getCredential({
                            id: c.id
                        })
                        this.verifiableCredentials.push(JSON.parse(resp.verifiableCredential))
                    }
                } catch (e) {
                    console.error('failed to get all stored credentials', e)
                    this.error = 'Failed to get your stored credentials'
                    this.errorDescription = 'Unable to get stored credentials from your wallet, please try again later.'
                }
                this.cards = this.verifiableCredentials.map((credential) => {
                    return { ...credential, flipped: false };

                });

            },
            credDisplayName: function (vc) {
                return vc.name ? vc.name : getCredentialType(vc.type)
            },
            toggleCard: function(card) {
                card.flipped = !card.flipped;
            },
        },
        computed: {
            loadingStatus() {
                return this.getCurrentUser().setupStatus
            }
        },
        data() {
            return {
                verifiableCredentials: [],
                cards: [],
                username: '',
                agent: null,
                error: 'No stored credentials',
                errorDescription: 'Your wallet is empty, there aren\'t any stored credentials to show.',
            }
        }
    }
</script>

<style scoped>
    .md-card {
      display: inline-block !important;
      position: relative !important;
      width: 100% !important;
      margin: 25px 0 !important;
      overflow: unset !important;
      background: none !important;
      box-shadow:none !important;
      -webkit-box-shadow: none !important;
    }
    ul.credential-list {
        padding-left: 0;
        display: flex;
        flex-flow: row wrap;
    }

    li {
        list-style-type: none;
        padding: 10px 10px;
        transition: all 0.3s ease;
    }

    .card {
        display: block;
        width: 360px;
        height: 233px;
        padding: 10px;
        background-color: #FFFFFF ;
        border-radius: 7px;
        margin: 5px;
        text-align: center;
        line-height: 22px;
        cursor: pointer;
        position: relative;
        color: black;
        font-weight: 400;
        font-size: 16px;
        -webkit-box-shadow: 9px 10px 22px -8px rgba(209,193,209,.5);
        -moz-box-shadow: 9px 10px 22px -8px rgba(209,193,209,.5);
        box-shadow: 9px 10px 22px -8px rgba(209,193,209,.5);
        will-change: transform;
        user-select: none;
    }
    .card i {
        color: rgb(11, 151, 196) !important;
    }

    .cardContent {
        text-align: left;
    }
    .cardHeader {
        font-weight: 500;
        padding: 10px 15px;
    }
    .cardBack {
        padding-top:40px;
        color: rgba(0,0,0,.54);
    }
    li:hover{
        transform: scale(1.1);
    }

    .flip-enter-active {
        transition: all 0.4s ease;
    }

    .flip-leave-active {
        display: none;
    }

    .flip-enter, .flip-leave {
        transform: rotateY(180deg);
        opacity: 0;

    }
    .md-dialog-container {
        width: 100% !important;
    }

</style>


