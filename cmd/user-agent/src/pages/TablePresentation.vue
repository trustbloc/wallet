/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

<template>
  <div class="content">
    <div class="md-layout">
      <div
        class="md-layout-item md-medium-size-100 md-xsmall-size-100 md-size-100"
      >
        <md-card class="md-card-plain">
          <md-card-header data-background-color="green">
            <h4 class="title">Choose your credential</h4>
            <p class="category"> Generate button will create the signed presentation QR Code </p>
          </md-card-header>
          <md-card-content style="background-color: white;">
            <md-field>
            </md-field>
              <select id="selectVC" style="color: grey; width: 200px; height: 35px;">
                <option value="0">Select VC</option>
              </select>
            <md-field style="margin-top: -15px">
            </md-field>
            <md-button class="md-button md-info md-square md-theme-default md-large-size-100 md-size-100"
                       id="getVCBtn">Generate Presentation QR
            </md-button>
            <md-field>
            </md-field>
              <img src="" id="qr-result" style="width:25%" />
          </md-card-content>
        </md-card>
      </div>
    </div>
  </div>
</template>

<script>
  async function generateQRCode() {
    document.getElementById('getVCBtn').addEventListener('click', async () => {
      // Get the VC ID from UI selection
      let e = document.getElementById("selectVC");
      let vcID = e.options[e.selectedIndex].value;
      if (vcID == 0) {
        alert("please select vc")
        return
      }

      // Get the VC data
      let vcData
      await window.$aries.verifiable.getCredential({
        id: vcID
      }).then(resp => {
                vcData = JSON.stringify(JSON.parse(resp.verifiableCredential))
              }
      ).catch(err =>
              console.log('generateQRCode - get vc failed : errMsg=' + err)
      )

      // Generate QR code
      let QRCode = require('qrcode')
      QRCode.toDataURL(vcData, function (err, url) {
        let canvas = document.getElementById('qr-result')
        canvas.src = url
      })
    });
  }

  export default {
    beforeCreate: async function () {
      // Load the Credentials in the drop down
      let aries = await this.$arieslib
      await aries.verifiable.getCredentials()
              .then(resp => {
                        const data = resp.result
                        if (data && data.length !== 0) {
                          let dropdown = document.getElementById('selectVC');
                          let option;
                          for (let i = 0; i < data.length; i++) {
                            option = document.createElement('option');
                            option.text = data[i].name;
                            option.value = data[i].id;
                            dropdown.add(option);
                          }
                        } else {
                          console.log('no credentials exists')
                        }
                      }
              ).catch(err => {
                        console.log('get credentials failed : errMsg=' + err)
                      }
              )

      window.$webCredentialHandler = this.$webCredentialHandler
      window.$aries = aries

      generateQRCode()
    },
  }


</script>
