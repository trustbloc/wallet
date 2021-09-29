/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

"use strict";

const { chapi, wallet } = require("../helpers");

const vc = {
  "@context": [
    "https://www.w3.org/2018/credentials/v1",
    "https://w3id.org/citizenship/v1",
    "https://w3id.org/vc-revocation-list-2020/v1",
    "https://w3id.org/security/bbs/v1",
  ],
  credentialStatus: {
    id: "https://issuer-vcs.trustbloc.local/status/1#3",
    revocationListCredential: "https://issuer-vcs.trustbloc.local/status/1",
    revocationListIndex: "3",
    type: "RevocationList2020Status",
  },
  credentialSubject: {
    birthCountry: "Bahamas",
    birthDate: "1958-07-17",
    familyName: "Pasteur",
    gender: "Male",
    givenName: "Louis",
    id: "did:orb:AiMP4:EiA2Gtl-qHKjTouzu-Rd0cYOwQxJ-qN0DO0HNnhfXXCqCg",
    lprCategory: "C09",
    lprNumber: "999-999-999",
    residentSince: "2015-01-01",
    type: ["Person", "PermanentResident"],
  },
  description: "Permanent Resident Card of Mr.Louis Pasteur",
  id: "http://example.com/eb299a34-529e-4a84-a67c-573865db4aa7",
  issuanceDate: "2021-03-11T14:52:00.8492482Z",
  issuer: {
    id: "did:key:zUC72c7u4BYVmfYinDceXkNAwzPEyuEE23kUmJDjLy8495KH3pjLwFhae1Fww9qxxRdLnS2VNNwni6W3KbYZKsicDtiNNEp76fYWR6HCD8jAz6ihwmLRjcHH6kB294Xfg1SL1qQ",
    name: "didkey-bbsblssignature2020-bls12381g2",
  },
  name: "Permanent Resident Card",
  proof: {
    created: "2021-07-06T18:16:57.739627-04:00",
    proofPurpose: "assertionMethod",
    proofValue:
      "koyKGr8WwjCOUqm-HV_7SVtvIIM4EhnJJ_8P2k0RF3ElQP2ntQJMKKtpoQTqk5l3QI5jN0Zn8nHJm3gyFkKdYJpC4IseNTU98u9UTijHlABpAhGbDaKTHs-b1IDsHkx_DrR3BSktz1Va_cilRP2WqA",
    type: "BbsBlsSignature2020",
    verificationMethod:
      "did:key:zUC73A7EHiDAxxy29qox4hD5Dyc6fXqStkWjbW2V5uVtmdpAr33Lhtz2sb9m8WotP6WxvjWxGb4iVsPPM5EGkwq5NCNwb6sn9breK588SiEcBtQEPyK7wXzXBT9QcCZ3S5XWygm#zUC73A7EHiDAxxy29qox4hD5Dyc6fXqStkWjbW2V5uVtmdpAr33Lhtz2sb9m8WotP6WxvjWxGb4iVsPPM5EGkwq5NCNwb6sn9breK588SiEcBtQEPyK7wXzXBT9QcCZ3S5XWygm",
  },
  type: ["VerifiableCredential", "PermanentResidentCard"],
};

describe("TrustBloc Wallet - WACI Share flow", () => {
  const ctx = {
    email: `ui-aut-${new Date().getTime()}@test.com`,
  };

  // runs once before the first test in this block
  before(async () => {
    await browser.reloadSession();
    await browser.maximizeWindow();
  });

  beforeEach(function () {});

  afterEach(async function () {
    if (this.currentTest.state === "failed") {
      const logs = await browser.getLogs("browser");
      console.log(JSON.stringify(logs, null, 4));
    }
  });

  it(`User Sign up (${ctx.email})`, async function () {
    this.timeout(90000);

    // 1. Navigate to Wallet Website
    await browser.navigateTo(browser.config.walletURL);

    // 2. Initialize Wallet (register/sign-up/etc.)
    await wallet.init(ctx);

    // TODO - https://github.com/trustbloc/wallet/issues/1140 Dashboard loads before router connection is setup
    await new Promise((resolve) => setTimeout(resolve, 8000));
  });

  it(`User Stores Permanent Resident Card Credential`, async function () {
    // mock issuer (wallet page with sample requests)
    await browser.navigateTo(browser.config.walletURL + "/web-wallet");

    const vcSampleBtn = await $("#store-vc-sample-1");
    await vcSampleBtn.waitForExist();
    await vcSampleBtn.click();

    let vprs = {
      "@context": ["https://www.w3.org/2018/credentials/v1"],
      type: "VerifiablePresentation",
      verifiableCredential: [vc],
    };

    const sampleText = await $("#vcDataTextArea");
    await sampleText.clearValue();
    await sampleText.addValue(vprs);

    const storeButton = await $("button*=Store");
    await storeButton.waitForClickable();
    await storeButton.click();

    await chapi.chooseWallet({
      name: browser.config.walletName,
    });

    await wallet.storeCredentials(ctx);
    await browser.switchToFrame(null);

    const storeSuccessMsg = await $(
      "div*=Successfully stored verifiable presentation to wallet."
    );
    await storeSuccessMsg.waitForExist();
  });

  it(`User validates the saved credential from mock issuer`, async function () {
    this.timeout(90000);

    await browser.navigateTo(browser.config.walletURL);

    const vcName = await $("span*=" + vc.name);
    await vcName.waitForExist();
    await vcName.click();

    const credTab = await $("span*=Credentials");
    await credTab.waitForExist();
    await credTab.click();
  });

  it(`User presents credential through WACI-Share (Redirect) : already signed-in`, async function () {
    // demo verifier page
    await browser.navigateTo(browser.config.demoVerifierURL);

    const waciShareDemoBtn = await $("#waci-share-demo");
    await waciShareDemoBtn.waitForExist();
    await waciShareDemoBtn.click();

    const shareBtn = await $("#share-credentials");
    await shareBtn.waitForExist();
    await shareBtn.click();

    const getSuccessMsg = await $("b*=Successfully Received Presentation");
    await getSuccessMsg.waitForExist();
  });

  it(`User signs out`, async function () {
    // wallet
    await browser.navigateTo(browser.config.walletURL);

    await wallet.logout(ctx);
  });
});
