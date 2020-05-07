/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <div class="content table-presentation">
    <div class="md-layout">
      <div
              class="md-layout-item md-medium-size-100 md-xsmall-size-100 md-size-100"
      >
        <md-card class="md-card-plain">
          <md-card-header data-background-color="green">
            <h4 class="title">Choose your options</h4>
            <p class="category"> Generate button will create the signed presentation </p>
          </md-card-header>
          <md-card-content style="background-color: white;">
            <div v-if="errors.length">
              <b>Failed with following error(s):</b>
              <ul>
                <li v-for="error in errors" :key="error">{{ error }}</li>
              </ul>
            </div>
            <md-field>
            </md-field>
            <label>
              <md-icon>how_to_reg</md-icon>
              Issuer</label><br>
            <select v-model="selectedIssuer" id="selectDID" style="color: grey; width: 300px; height: 35px;">
              <option v-for="issuer in issuers" :key="issuer.id" :value="issuer.id">
                {{issuer.name}}
              </option>
            </select><br><br>

            <label>
              <md-icon>fingerprint</md-icon>
              Credential</label><br>
            <select v-model="selectedVC" style="color: grey; width: 300px; height: 35px;">
              <option v-for="vc in savedVCs" :key="vc.id" :value="vc.id">
                {{vc.name}}
              </option>
            </select><br><br>

            <div class="md-layout-item md-size-100">
              <md-field maxlength="5">
                <label class="md-helper-text">Type presentation friendly name here</label>
                <md-input v-model="friendlyName" id="friendlyName" required></md-input>
              </md-field>
            </div>

            <md-button v-on:click="generatePresentation" class="md-button md-success md-square md-theme-default md-large-size-100 md-size-100"
                       id="getVPBtn">Generate Presentation
            </md-button>
            <md-field style="margin-top: 5px"></md-field>

            <div v-if="loading" style="margin-left: 40%;margin-top: 20%;height: 200px;">
            <div class="md-layout">
              <md-progress-spinner :md-diameter="100" class="md-primary" :md-stroke="10"
                                   md-mode="indeterminate"></md-progress-spinner>
            </div>
          </div>

            <md-tabs class="md-info md-ripple" md-alignment="left"  v-if="!isHidden">
              <md-tab id="tab-home" md-label="Source" md-icon="code">
                <md-card-content v-model="vpData" style="overflow: hidden">
                  <vue-json-pretty
                          :data="this.vpData"
                  >
                  </vue-json-pretty>
                </md-card-content>
              </md-tab>
              <md-tab id="tab-pages" md-label="QR Code" md-icon="rounded_corner">
                <md-content>
                  <img src="" id="qr-result" style="width:25%" />
                </md-content>
              </md-tab>

            </md-tabs>

          </md-card-content>
        </md-card>
      </div>
    </div>
  </div>
</template>

<script>
  import VueJsonPretty from 'vue-json-pretty';
  export default {
    components: {
      VueJsonPretty
    },
    beforeCreate: async function () {
      this.aries = await this.$arieslib
      await this.loadIssuers()
      // Load the Credentials in the drop down
      await this.aries.verifiable.getCredentials()
              .then(resp => {
                        const data = resp.result
                        if (data.length == 0) {
                          console.log("unable to get saved VCs")
                          return
                        }
                        this.savedVCs.length = 0
                        data.forEach((item) => {
                          this.savedVCs.push({id:item.id, name:item.name})
                        })
                        this.selectedVC = this.savedVCs[0].id
                      }
              ).catch(err => {
                        console.log('get credentials failed : errMsg=' + err)
                      }
              )
      window.$webCredentialHandler = this.$webCredentialHandler
      window.$aries = this.aries
    },
    data() {
      return {
        savedVCs: [{id: "", name: "Select VC"}],
        selectedVC: "",
        issuers: [{id: 0, name: "Select Identity"}],
        selectedIssuer: "",
        errors: [],
        vpData:"",
        isHidden: true,
        friendlyName:"",
        loading:false
      };
    },
    methods: {
      getDIDMetadata: function (id) {
        return new Promise(function(resolve) {
          var openDB = indexedDB.open("did-metadata", 1);

          openDB.onupgradeneeded = function () {
            var db = {}
            db.result = openDB.result;
            db.store = db.result.createObjectStore("metadata", {keyPath: "id"});
          };

          openDB.onsuccess = function () {
            var db = {};
            db.result = openDB.result;
            db.tx = db.result.transaction("metadata", "readonly");
            db.store = db.tx.objectStore("metadata");
            let getData = db.store.get(id);
            getData.onsuccess = function () {
              resolve(getData.result);
            };

            db.tx.oncomplete = function () {
              db.result.close();
            };
            console.log("got did metadata from db")
          }
        });
      },
      generatePresentation: async function () {
        const errorMsg = "friendly name required."
        if (this.friendlyName.length === 0) {
          if (!this.errors.includes(errorMsg)) {
            this.errors.push(errorMsg)
          }
          return
        }
        this.isHidden = false
        this.loading = true
        let didMetadata = await this.getDIDMetadata(this.issuers[this.selectedIssuer].key)

        // fetch the credential
        let data = await this.getSelectedCredentials()
        let QrData
        // generate presentation
        if (data.vc) {
          await window.$aries.verifiable.generatePresentation({
            verifiableCredential: data.vc,
            did: this.issuers[this.selectedIssuer].key,
            skipVerify: true,
            signatureType:didMetadata.signatureType,
            privateKey: didMetadata.privateKey,
            keyType: didMetadata.privateKeyType,
            verificationMethod: didMetadata.keyID
          }).then(resp => {
              this.vpData = resp.verifiablePresentation
                    QrData = JSON.stringify(resp.verifiablePresentation)
                  }
          ).catch(err =>
                  this.errors.push("failed to create presentation : errMsg="+ err)
          )

          console.log("Response presentation:", this.vpData)

         // save presentation
          this.aries.verifiable.savePresentation({
                   verifiablePresentation: this.vpData,
                   name: this.friendlyName
           }).then(resp => {
                console.log("Successfully saved presentation:", resp)
             }).catch( err =>{
                   this.vpData = err

                   console.log('failed to save presentation, errMsg:', err)
              }
           )

          // Generate QR code
          let QRCode = require('qrcode')
          QRCode.toDataURL(QrData, function (err, url) {
            let canvas = document.getElementById('qr-result')
            canvas.src = url
          })
        }

        this.loading = false
      },
      getSelectedCredentials: async function () {
        if (this.selectedVC.length === 0) {
          return {retry: "Please select at least one credential"}
        }

        try {
          let vc = []
          await window.$aries.verifiable.getCredential({
            id: this.selectedVC
          }).then(resp => {
                    vc.push(JSON.parse(resp.verifiableCredential))
                  }
          ).catch(err =>
                  console.log('get credential failed=' + err)
          )
          return {vc: vc}
        } catch (e) {
          return e
        }
      },
      loadIssuers: async function () {
        await this.aries.vdri.getDIDRecords().then(
                resp => {
                  const data = resp.result
                  if (!data || data.length == 0) {
                    this.errors.push("No issuers found to select, please create an issuer")
                    return
                  }

                  this.issuers = []
                  this.selectedIssuer = 0

                  data.forEach((item, id) => {
                    this.issuers.push({id: id, name: item.name, key: item.id})
                  })
                })
                .catch(err => {
                  this.errors.push(err)
                })
      }
    }
  }


</script>
<style>
  .table-presentation .md-ripple {
    margin-top: 10px;
  }
</style>
